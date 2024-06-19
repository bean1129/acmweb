package service

import (
	"acmweb/app/authorization/model"
	"acmweb/app/constants"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/config"
	"database/sql"
	"os"
	"path/filepath"
	"time"
)

type AuthService struct {
	dbTx        *sql.Tx
	ApplyData   model.AuthApply
	ApproveData model.AuthApprove
}

func NewAuthService() *AuthService {
	return new(AuthService)
}

func init() {
	//定时检测申请的证书内容是否变化
	if system.Config.Application.CheckTime > 0 {
		go func() {
			authSer := NewAuthService()
			for {
				time.Sleep(time.Duration(system.Config.Application.CheckTime) * time.Second)
				total, change, err := authSer.CheckLicense()
				if err != nil {
					continue
				}
				system.Log.Infof("Check license data total[%d],change[%d]...", total, change)
			}
		}()
	}
}

func (c *AuthService) CheckLicense() (total, change int, err error) {
	search := model.AuthSearch{}
	//是否有分页
	var authIds = make(map[int64]string, 0)
	if authIds, err = search.GetVaildAuth(); err != nil {
		return
	}
	total = len(authIds)
	for authId, certMd5 := range authIds {
		licData, licMD5, filename, err := c.creatLicFileData(authId)
		if err != nil {
			system.Log.Error("Create license file failed,err:%s", err.Error())
			continue
		}
		//更新数据库cert
		if certMd5 != licMD5 {
			filename, _ := c.writeLicFile(authId, filename, licData)
			ap := model.AuthApply{AuthId: authId, Cert: filename, CertMD5: licMD5}
			if _, err := ap.UpdateCert(); err != nil {
				system.Log.Error("Upadte auth req failed,err:", err.Error())
			}
			system.Log.Infof("License data have change,auth[%d]", authId)
			change += 1
		}
	}
	return
}

func (c *AuthService) beginWork() error {
	var err error
	if c.dbTx, err = system.Data.DB.StartTrans(); err != nil {
		return err
	}
	return nil
}

func (c *AuthService) commit() bool {
	return system.Data.DB.CommitTrans(c.dbTx)
}

func (c *AuthService) rollback() {
	system.Data.DB.Rollback(c.dbTx)
}

// 申请授权
func (c *AuthService) Apply() (*common.ZIResult, error) {

	system.Log.Infof("Apply request:wf_inst_id[%d],inst_id[%d],trial_days[%d],host[%s] mac[%s]",
		c.ApplyData.WFInstId, c.ApplyData.InstId, c.ApplyData.TrialDays, c.ApplyData.Host, c.ApplyData.Mac)

	zs := common.NewResult()
	//生成申请信息
	var err error
	if err = c.beginWork(); err != nil {
		system.Log.Error("Start database transaction failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	//查询组
	ap := model.Scene{InstId: c.ApplyData.InstId}
	if err = ap.GetSceneInfo(); err != nil {
		return zs, &errors.ErrDBQry
	}

	grpId := ap.GetSceneGrp()
	if grpId == 0 {
		return zs, &errors.ErrNopeAppGRP
	}

	//查询流程信息
	wf := model.WF{GroupId: grpId}
	if err := wf.GetWFByGrp(); err != nil {
		return zs, &errors.ErrDBexecute
	}

	if wf.WFId == 0 {
		return zs, &errors.ErrNopeAppWF
	}
	//
	//生成流程实例信息
	_ = c.ApplyData.GetSubmitterId()
	wfInst := model.WFInst{WFInstId: c.ApplyData.WFInstId, Submitter: c.ApplyData.ProposerId, AppvCid: wf.ApproveOne, AppNid: wf.ApproveTwo, WFId: wf.WFId, WFLevel: 1}
	if wfInst.WFInstId == 0 {
		wfInst.WFInstId = system.Common.UUID.NextVal()
		if zs.Count, err = wfInst.Insert(c.dbTx, c.ApplyData.Proposer); err != nil {
			system.Log.Error("Insert work flow inst failed,err:", err.Error())
			c.rollback()
			return zs, &errors.ErrDBexecute
		}
	} else {
		//重新申请
		if zs.Count, err = wfInst.Modify(c.dbTx); err != nil {
			system.Log.Error("Insert work flow inst failed,err:", err.Error())
			c.rollback()
			return zs, &errors.ErrDBexecute
		}
	}

	//生成申请
	c.ApplyData.WFInstId = wfInst.WFInstId
	c.ApplyData.Scene = ap
	if zs.Count, err = c.ApplyData.Insert(c.dbTx); err != nil {
		system.Log.Error("Insert apply information failed,err:", err.Error())
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//记录日志
	wfLog := model.WFLog{InstId: wfInst.WFInstId,
		AppvId:     wfInst.Submitter,
		WfoState:   constants.WF_APPROVE_TODO,
		WfnState:   constants.WF_APPROVE_TODO,
		NextAppvId: wfInst.AppvCid,
		OprId:      constants.WF_OPR_START,
	}
	if err = c.addWFLog(wfLog); err != nil {
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//提交事务
	if !c.commit() {
		system.Log.Error("Trans commit failed,err:", err.Error())
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	return zs, nil
}

func (c *AuthService) Approve() (*common.ZIResult, error) {
	system.Log.Infof("Next approve request:wf_inst_id[%d],approver[%s],Oper[%d]",
		c.ApproveData.WfInstId, c.ApproveData.Approver, c.ApproveData.WfOprState)

	zs := common.NewResult()
	var err error
	//查询时实例信息
	wfInst := model.WFInst{WFInstId: c.ApproveData.WfInstId}
	if _, err = wfInst.GetWfInstInfo(); err != nil {
		system.Log.Error("Get work flow inst failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	//查询流程
	wf := model.WF{WFId: wfInst.WFId}
	if err = wf.GetWFById(); err != nil {
		system.Log.Error("Get work flow by id[%d] failed,err:", wf.WFId, err.Error())
		return zs, &errors.ErrDBexecute
	}
	//判断操作是否允许
	if wfInst.WFLevel == 1 && c.ApproveData.WfOprState == constants.WF_OPR_PRE {
		system.Log.Errorf("First level ,no allow to return")
		return zs, &errors.ErrInvaildApprover
	}

	wfLog := model.WFLog{
		InstId:     wfInst.WFInstId,
		AppvId:     wfInst.AppvCid,
		NextAppvId: wfInst.AppNid,
		Opinion:    c.ApproveData.Opinion,
		WfoState:   wfInst.WfState,
		OprId:      c.ApproveData.WfOprState,
	}

	//设置流程下一级
	c.setNextWorkFlow(&wfInst, wf)

	wfLog.WfnState = wfInst.WfState
	if wfInst.WfState == 1 {
		wfLog.OprId = constants.WF_OPR_Finish
	} else if wfInst.WfState == 3 {
		wfLog.NextAppvId = 0
	}

	//以下事务处理
	if err = c.beginWork(); err != nil {
		system.Log.Error("Start database transaction failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	//更新流程实例信息
	if _, err := wfInst.Update(c.dbTx); err != nil {
		system.Log.Errorf("Update work flow inst by id[%d] failed,err:%s", wf.WFId, err.Error())
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//更新授权申请状态
	ap := model.AuthApply{WFInstId: wfInst.WFInstId}
	if _, err := ap.UpdateApproveState(c.dbTx, wfInst.WfState); err != nil {
		system.Log.Errorf("Update auth_req approve state by id[%d] failed,err:%s", wf.WFId, err.Error())
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//记录日志
	if c.ApproveData.Approver == 1 {
		//admin记录备注
	}
	if err = c.addWFLog(wfLog); err != nil {
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//提交事务
	if !c.commit() {
		system.Log.Error("Trans commit failed,err:", err.Error())
		c.rollback()
		return zs, &errors.ErrDBexecute
	}

	//审批通过授权
	if wfInst.WfState == 1 {
		//审批通过授权
		ap.AuthState = 1
		if err := ap.GetAuthReqInfo(); err == nil {
			//异步生成
			go func() {
				ap.Cert, ap.CertMD5, err = c.CreatLicFile(ap.AuthId)
				if err != nil {
					system.Log.Error("Create license file failed,err:%s", err.Error())
				}
				//更新数据库cert
				if _, err := ap.UpdateCert(); err != nil {
					system.Log.Error("Upadte auth req failed,err:", err.Error())
				}
			}()
		}
	}

	return zs, nil
}

func (c *AuthService) CreatLicFile(AuthId int64) (string, string, error) {
	//生成
	LicData, LicMD5, curfilename, err := c.creatLicFileData(AuthId)
	if err != nil {
		return "", "", err
	}

	filename, err := c.writeLicFile(AuthId, curfilename, LicData)
	if err != nil {
		return "", "", err
	}
	return filename, LicMD5, nil

	//zipWriter := zip.NewWriter(fZip)
	//defer zipWriter.Close()

	//生成license文件，文件名：{appcode}-{authid}.lic
	//filename := + "-" + licDesc.Id + ".lic"
	//if file, err := zipWriter.Create(filename); err != nil {
	//	system.Log.Error("Create file failed,filename[%s],err:", fullfilename, err.Error())
	//	return "", err
	//} else {
	//	_, err = file.Write([]byte(genLic.GetLicData()))
	//	if err != nil {
	//		return "", err
	//	}
	//}
	//
	//licDesc.FileName = filename
	////生成license desc xml文件，文件名：desc.xml
	//xmlDescFile := "desc.xml"
	//xmldata, _ := xml.MarshalIndent(&licDesc, " ", " ")
	//if xmlfile, err := zipWriter.Create(xmlDescFile); err != nil {
	//	return "", err
	//} else {
	//	_, err = xmlfile.Write(xmldata)
	//}
	//
	//return "license/" + fzipfilename, nil
}

func (c *AuthService) creatLicFileData(Id int64) (fileContent, fileContentMd5, filename string, err error) {
	//生成
	genLic := NewGenLic(Id)
	var filenamekey string
	if filenamekey, err = genLic.GenLicData(); err != nil {
		system.Log.Error("Generator license file failed.err:%s", err.Error())
		return "", "", "", err
	}
	if err := genLic.GenSign(); err != nil {
		system.Log.Error("Generator license sign failed,err:%s", err.Error())
		return "", "", "", err
	}

	fileContent = genLic.GetLicData()
	fileContentMd5 = genLic.GetLicMd5()
	filename = filenamekey + ".lic"
	return
}

func (c *AuthService) writeLicFile(AuthId int64, curfilename string, LicData string) (filename string, err error) {
	//存储路径是否存在
	filePath := config.CONFIG.Attachment.FilePath + "/license/"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			system.Log.Errorf("mkdir path[%s] faield ,err:%s", filePath, err.Error())
			return "", err
		}
	}

	//生成license文件，文件名：{authid}-{md5}.lic
	//[auth_id]+[unit_id]+[1个IP地址].lic
	//filename = system.Text.Conv.String(AuthId) + "-" + LicMD5 + ".lic"
	fullfilename := filePath + curfilename
	fp, err := os.Create(fullfilename)
	if err != nil {
		system.Log.Errorf("Create file[%s] failed,err:", fullfilename, err.Error())
	}
	defer fp.Close()
	_, err = fp.WriteString(LicData)
	if err != nil {
		system.Log.Errorf("Write license data to file failed.err:", err.Error())
	}
	return "license/" + curfilename, nil
}
func (c *AuthService) addWFLog(m model.WFLog) error {
	if err := m.InsertWfLog(c.dbTx); err != nil {
		system.Log.Error("Insert work flow log  failed,err:", err.Error())
		return &errors.ErrDBexecute
	}
	return nil
}

func (c *AuthService) setNextWorkFlow(wfInst *model.WFInst, wf model.WF) {
	//提交下一级审批人

	switch c.ApproveData.WfOprState {
	case constants.WF_OPR_NEXT:
		wfInst.AppvCid = wfInst.AppNid
		wfInst.SetNextApprover(wf)
		if wfInst.AppvCid == 0 {
			//表示审批结束
			wfInst.WfState = constants.WF_APPROVE_PASS
		} else if wfInst.WFLevel == 0 {
			wfInst.WfState = constants.WF_APPROVE_TODO
		} else {
			wfInst.WfState = constants.WF_APPROVE_DOING
		}
		wfInst.WFLevel += 1
	case constants.WF_OPR_PRE:
		wfInst.AppNid = wfInst.AppvCid
		wfInst.SetPreApprover(wf)
		if wfInst.WFLevel != 1 {
			wfInst.WFLevel -= 1
		}
	case constants.WF_OPR_END:
		//终止流程
		wfInst.AppNid = 0
		wfInst.AppvCid = 0
		wfInst.WfState = constants.WF_APPROVE_NOPASS
	case constants.WF_OPR_BACK:
		//返回修改
		wfInst.AppvCid = wfInst.Submitter
		wfInst.WfState = constants.WF_APPROVE_Modify
		//表示重新提交
		wfInst.WFLevel = 0
	}
	system.Log.Infof("Set next approver[%d],level[%d]", wfInst.AppvCid, wfInst.WFLevel)
}

func (c *AuthService) List(search model.AuthSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	var args = make([]any, 0)
	var cond string
	start := (search.PageNum - 1) * search.PageSize

	//根据用户查询
	if len(search.InstName) > 0 {
		cond = " and app_name like ? "
		args = append(args, "%"+search.InstName+"%")
	} else {
		cond = " and 1=1 "
	}

	//根据组查询
	if search.AuthState != -1 {
		cond += " and a.auth_state=? "
		args = append(args, search.AuthState)
	} else {
		cond += " and 1=1 "
	}

	//审批状态
	if search.ApproveState != -1 {
		cond += " and a.appv_state=? "
		args = append(args, search.ApproveState)
	} else {
		cond += " and 1=1 "
	}
	//下一级审批人
	if len(search.CurApprover) > 0 {
		cond += " and e.user_code = ?"
		args = append(args, search.CurApprover)
	} else {
		cond += " and 1=1"
	}
	//组织
	if search.UnitId > 0 {
		cond += " and c.unit_id = ?"
		args = append(args, search.UnitId)
	} else {
		cond += " and 1=1"
	}
	if len(search.Submitter) > 0 {
		cond += " and a.proposer like ?"
		args = append(args, "%"+search.Submitter+"%")
	} else {
		cond += " and 1=1"
	}
	//是否有分页
	var authInfos = make([]model.AuthRes, 0)
	var err error
	authInfos, err = search.FindByIndex(start, search.PageSize, cond, args...)

	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range authInfos {
		zs.Data = append(zs.Data, v)
	}
	//总数
	zs.Count, _ = search.FindTotalNum(cond, args...)
	return zs, nil
}

func (c *AuthService) ListNum(search model.AuthSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	authInfos, err := search.FindAuthNum()

	if err != nil {
		return zs, &errors.ErrDBexecute
	}
	zs.Data = append(zs.Data, authInfos)
	zs.Count = 1
	return zs, nil
}

func (c *AuthService) ListAuthMode() (*common.ZIResult, error) {
	zs := common.NewResult()
	sql := "SELECT * from auth_mode a"
	ds, err := system.Data.DB.QueryRows(sql)
	if err != nil {
		return zs, err
	}
	if len(ds) == 0 {
		return zs, nil
	}
	for _, rs := range ds {
		info := model.AuthMode{
			TypeId:   system.Text.Conv.Int(rs["type_id"]),
			TypeName: rs["type_name"],
			TypeCode: rs["type_code"],
		}
		zs.Data = append(zs.Data, info)
		zs.Count++
	}
	return zs, nil
}

func (c *AuthService) AuthStatis(appId []int64, tierLevel int) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.AuthStatis

	retLevel, _ := m.GetTierLevel(tierLevel)

	if resultApp, err := m.GetHeader(appId, tierLevel); err != nil {
		return nil, &errors.ErrDBexecute
	} else {
		if err := m.GetStatisData(resultApp, retLevel); err != nil {
			return nil, &errors.ErrDBexecute
		}
	}

	zs.Data = append(zs.Data, m)
	zs.Count = int64(len(m.ResultData))
	return zs, nil
}

func (c *AuthService) PersistentLicFile(authId int64) (filename string, err error) {
	var licMd5 string
	if filename, licMd5, err = c.CreatLicFile(authId); err != nil {
		return "", err
	}
	fullname := system.Config.Attachment.FilePath + "/" + filename
	system.Log.Infof("Create license file[%s] success.", fullname)
	//入库
	ap := model.AuthApply{AuthId: authId, Cert: filename, CertMD5: licMd5}
	//更新数据库cert
	if _, err := ap.UpdateCert(); err != nil {
		system.Log.Error("Upadte auth req failed,err:", err.Error())
	}
	return fullname, nil
}

func (c *AuthService) DownloadLicFile(authId int64) (srcFile, destFile string, err error) {
	system.Log.Infof("Received download file request,auth id[%d]", authId)
	genLic := NewGenLic(authId)
	cert, err := genLic.FindLicData()
	if err != nil {
		system.Log.Errorf("Get license data failed,err:%s", err.Error())
		return srcFile, destFile, err
	} else {
		system.Log.Infof("Get license data success,len:%d", len(cert))
	}
	if len(cert) == 0 {
		system.Log.Warnf("License file not exists,to create")
		srcFile, err = c.PersistentLicFile(authId)
		if err != nil {
			return srcFile, destFile, err
		}
	} else {
		system.Log.Infof("Get license file[%s] success.", cert)
		//判断文件是否存在
		if _, err := os.Stat(cert); !os.IsExist(err) {
			srcFile, err = c.PersistentLicFile(authId)
			if err != nil {
				return srcFile, destFile, err
			}
		} else {
			srcFile = cert
		}
	}
	destFile = filepath.Base(srcFile)
	return srcFile, destFile, err
}

func (c *AuthService) AuthStatisByAxis(appId []int64, tierLevel int) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.AuthAxisStatis

	var result = make(map[string][]model.AuthAxisStatis, 0)
	if appStatis, err := m.GetUnitStatisByApp(appId, tierLevel); err != nil {
		return nil, &errors.ErrDBexecute
	} else {
		result["app_statis"] = append(result["app_statis"], appStatis...)
	}
	if unitStatis, err := m.GetAppStatisByUnit(appId, tierLevel); err != nil {
		return nil, &errors.ErrDBexecute
	} else {
		result["unit_statis"] = append(result["unit_statis"], unitStatis...)
	}

	zs.Data = append(zs.Data, result)
	zs.Count = int64(len(result))
	return zs, nil
}

func (c *AuthService) DeleteAuth(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.AuthApply{AuthId: id}
		rows, err := m.DeleteAuth()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}
