package model

import (
	"acmweb/system"
	"acmweb/system/text"
	"database/sql"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type WF struct {
	WFId         int64  `mapstructure:"wf_id"`
	WFName       string `mapstructure:"wf_name"`
	ApproveOne   int64  `mapstructure:"approver1"`
	ApproveTwo   int64  `mapstructure:"approver2"`
	ApproveThree int64  `mapstructure:"approver3"`
	ApproveFour  int64  `mapstructure:"approver4"`
	ApproveFive  int64  `mapstructure:"approver5"`
	ApproveSix   int64  `mapstructure:"approver6"`
	ApproveSeven int64  `mapstructure:"approver7"`
	ApproveEight int64  `mapstructure:"approver8"`
	ApproveNine  int64  `mapstructure:"approver9"`
	GroupId      int64  `mapstructure:"group_id"`
	CrtDate      string `mmapstructure:"create_time"`
	AuthId       int64  `mapstructure:"auth_id"`
}

type WFInst struct {
	WFInstId   int64  `mapstructure:"wf_inst_id"`
	Title      string `mapstructure:"title"`
	Submitter  int64  `mapstructure:"submitter"`
	AppvCid    int64  `mapstructure:"appv_cid"`
	AppNid     int64  `mapstructure:"appv_nid"`
	WfState    int    `mapstructure:"wf_state"`
	CrtTime    string `mapstructure:"crt_time"`
	LastedTime string `mapstructure:"lasted_time"`
	WFId       int64  `mapstructure:"wf_id"` //正在使用的流程
	WFLevel    int    `mapstructure:"wf_level"`
}

type WFLog struct {
	LogId      int64  `mapstructure:"log_id" json:"log_id"`
	InstId     int64  `mapstructure:"inst_id" json:"inst_id"`
	AppvId     int64  `mapstructure:"appv_id" json:"appv_id"`
	NextAppvId int64  `mapstructure:"next_appv_id" json:"next_appv_id"`
	AppTime    int64  `mapstructure:"app_time" json:"app_time"`
	OprId      int    `mapstructure:"opr_id" json:"opr_id"`
	Opinion    string `mapstructure:"opinion" json:"opinion"`
	WfoState   int    `mapstructure:"wfo_state" json:"wfo_state"`
	WfnState   int    `mapstructure:"wfn_state" json:"wfn_state"`
}

type WFLogShow struct {
	LogId     int64  `mapstructure:"log_id" json:"log_id"`
	InstId    int64  `mapstructure:"inst_id" json:"inst_id"`
	Appv      string `mapstructure:"approver" json:"approver"`
	Operator  string `mapstructure:"next_approver" json:"next_approver"`
	Opinion   string `mapstructure:"opinion" json:"opinion"`
	WfoState  int    `mapstructure:"wfo_state" json:"wfo_state"`
	WfnState  int    `mapstructure:"wfn_state" json:"wfn_state"`
	OprName   string `mapstructure:"opr_name" json:"opr_name"`
	AppvTimeT string `mapstructure:"approve_time" json:"approve_time"`
}

type WFInfo struct {
	Id          int64  `mapstructure:"wf_id" json:"wf_id,omitempty"`
	WFName      string `mapstructure:"wf_name" json:"wf_name,omitempty"`
	GrpName     string `mapstructure:"grp_name" json:"grp_name,omitempty"`
	Approver1   string `mapstructure:"approver1" json:"approver1"`
	Approver2   string `mapstructure:"approver2" json:"approver2"`
	Approver3   string `mapstructure:"approver3" json:"approver3"`
	Approver4   string `mapstructure:"approver4" json:"approver4"`
	Approver5   string `mapstructure:"approver5" json:"approver5"`
	Approver6   string `mapstructure:"approver6" json:"approver6"`
	Approver7   string `mapstructure:"approver7" json:"approver7"`
	Approver8   string `mapstructure:"approver8" json:"approver8"`
	Approver9   string `mapstructure:"approver9" json:"approver9"`
	GrpId       int64  `mapstructure:"grp_id" json:"grp_id,omitempty"`
	Approver1Id int64  `mapstructure:"approver1_id" json:"approver1_id"`
	Approver2Id int64  `mapstructure:"approver2_id" json:"approver2_id"`
	Approver3Id int64  `mapstructure:"approver3_id" json:"approver3_id"`
	Approver4Id int64  `mapstructure:"approver4_id" json:"approver4_id"`
	Approver5Id int64  `mapstructure:"approver5_id" json:"approver5_id"`
	Approver6Id int64  `mapstructure:"approver6_id" json:"approver6_id"`
	Approver7Id int64  `mapstructure:"approver7_id" json:"approver7_id"`
	Approver8Id int64  `mapstructure:"approver8_id" json:"approver8_id"`
	Approver9Id int64  `mapstructure:"approver9_id" json:"approver9_id"`
}

type WFHisSearch struct {
	PageNum    int   `mapstructure:"page_num" json:"page_num,omitempty"`
	PageSize   int   `mapstructure:"page_size" json:"page_size,omitempty"`
	ApproverId int64 `mapstructure:"approver_id" json:"approver_id"`
	AppvState  int   `mapstructure:"appv_state" json:"appv_state"`
}

// 按组创建流程
func (w *WF) CreateWf() error {
	sql := "insert into wf(wf_id,wf_name,approver1,approver2,approver3,approver4,approver5,approver6,approver7,approver8,approver9,group_id)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?) on DUPLICATE key update wf_name = VALUES(wf_name),approver1 = VALUES(approver1),approver2 = VALUES(approver2)," +
		" approver3 = VALUES(approver3),approver4 = VALUES(approver4),approver5 = VALUES(approver5),approver6 = VALUES(approver6)," +
		" approver7 = VALUES(approver7),approver8 = VALUES(approver8),approver9 = VALUES(approver9),group_id = VALUES(group_id)"
	_, err := system.Data.DB.Exec(sql, w.WFId, w.WFName, w.ApproveOne, w.ApproveTwo, w.ApproveThree, w.ApproveFour, w.ApproveFive, w.ApproveSix, w.ApproveSeven, w.ApproveEight, w.ApproveNine, w.GroupId)
	if err != nil {
		return err
	}
	return nil
}
func (w *WF) DeleteWf() (int64, error) {
	sql := "delete from  wf where wf_id = ?"
	rc, err := system.Data.DB.Exec(sql, w.WFId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}
func (w *WF) UpdateWf(fields string, values ...any) (int64, error) {
	sql := "update wf set " + fields + " where wf_id = ?"
	values = append(values, w.WFId)
	rc, err := system.Data.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

// 按条件查询所有流程信息
func (w *WF) GetWFByCond(start, size int, grpID int64) (result []WFInfo, err error) {

	sqlQry := "SELECT wf_id,wf_name,b.grp_name,b.grp_id,c1.user_id approver1_id,c1.user_name approver1,c2.user_id approver2_id,c2.user_name approver2,c3.user_id approver3_id,c3.user_name approver3," +
		"c4.user_id approver4_id,c4.user_name approver4,c5.user_id approver5_id ,c5.user_name approver5,c6.user_id approver6_id,c6.user_name approver6" +
		",c7.user_id approver7_id,c7.user_name approver7,c8.user_id approver8_id,c8.user_name approver8,c9.user_id approver9_id,c9.user_name approver9"
	sqlQry += " from wf a INNER JOIN scene_group b on a.group_id = b.grp_id " +
		"LEFT JOIN user c1 on a.approver1 = c1.user_id " +
		"LEFT JOIN user c2 on a.approver2 = c2.user_id " +
		"LEFT JOIN user c3 on a.approver3 = c3.user_id " +
		"LEFT JOIN user c4 on a.approver4 = c4.user_id " +
		"LEFT JOIN user c5 on a.approver5 = c5.user_id " +
		"LEFT JOIN user c6 on a.approver6 = c6.user_id " +
		"LEFT JOIN user c7 on a.approver7 = c7.user_id " +
		"LEFT JOIN user c8 on a.approver8 = c8.user_id " +
		"LEFT JOIN user c9 on a.approver9 = c9.user_id where 1 = 1 "
	var args []any
	if grpID > 0 {
		sqlQry += " and a.group_id = ?"
		args = append(args, grpID)
	} else {
		//不加这个返回的是字符串不是int64
		sqlQry += "and a.group_id <>?"
		args = append(args, grpID)
	}
	if size > 0 {
		sqlQry += " order by wf_id limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	}
	system.Log.Info(sqlQry)
	ds, err := system.Data.DB.QueryResultRows(sqlQry, args...)
	if err != nil {
		return result, err
	}
	if len(ds) == 0 {
		return result, err
	}
	for _, rc := range ds {
		var wfRsp WFInfo
		if err := mapstructure.Decode(rc, &wfRsp); err != nil {
			continue
		}
		result = append(result, wfRsp)
	}
	return result, err
}

func (w *WF) FindTotalNum(grpID int64) (int64, error) {
	sqlQry := " select count(*)  total from wf a INNER JOIN app_group b on a.group_id = b.grp_id " +
		"LEFT JOIN user c1 on a.approver1 = c1.user_id " +
		"LEFT JOIN user c2 on a.approver2 = c2.user_id " +
		"LEFT JOIN user c3 on a.approver3 = c3.user_id " +
		"LEFT JOIN user c4 on a.approver4 = c4.user_id " +
		"LEFT JOIN user c5 on a.approver5 = c5.user_id " +
		"LEFT JOIN user c6 on a.approver6 = c6.user_id " +
		"LEFT JOIN user c7 on a.approver7 = c7.user_id " +
		"LEFT JOIN user c8 on a.approver8 = c8.user_id " +
		"LEFT JOIN user c9 on a.approver9 = c9.user_id where 1 = 1 "
	var args []any
	if grpID > 0 {
		sqlQry += " and a.group_id = ?"
		args = append(args, grpID)
	} else {
		//不加这个返回的是字符串不是int64
		sqlQry += "and a.group_id <>?"
		args = append(args, grpID)
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return 0, err
	}
	result := system.Text.Conv.Int64(ds[0]["total"])
	return result, nil
}

// 根据分组获取流程信息
func (w *WF) GetWFByGrp() error {
	sqlQry := "select * from wf where group_id=?"
	ds, err := system.Data.DB.QueryRows(sqlQry, w.GroupId)
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return nil
	}
	w.WFId = system.Text.Conv.Int64(ds[0]["wf_id"])
	w.WFName = ds[0]["wf_name"]
	w.ApproveOne = system.Text.Conv.Int64(ds[0]["approver1"])
	w.ApproveTwo = system.Text.Conv.Int64(ds[0]["approver2"])
	w.ApproveThree = system.Text.Conv.Int64(ds[0]["approver3"])
	w.ApproveFour = system.Text.Conv.Int64(ds[0]["approver4"])
	w.ApproveFive = system.Text.Conv.Int64(ds[0]["approver5"])
	w.ApproveSix = system.Text.Conv.Int64(ds[0]["approver6"])
	w.ApproveSeven = system.Text.Conv.Int64(ds[0]["approver7"])
	w.ApproveEight = system.Text.Conv.Int64(ds[0]["approver8"])
	w.ApproveNine = system.Text.Conv.Int64(ds[0]["approver9"])
	return nil
}

// 根据实例id
func (w *WF) GetWFById() error {
	sqlQry := "select a.*,b.auth_id from wf a,auth_req b where wf_id=? and a.wf_id = b.wf_inst_id"
	ds, err := system.Data.DB.QueryRows(sqlQry, w.WFId)
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return nil
	}
	w.WFName = ds[0]["wf_name"]
	w.AuthId = system.Text.Conv.Int64(ds[0]["auth_id"])
	w.ApproveOne = system.Text.Conv.Int64(ds[0]["approver1"])
	w.ApproveTwo = system.Text.Conv.Int64(ds[0]["approver2"])
	w.ApproveThree = system.Text.Conv.Int64(ds[0]["approver3"])
	w.ApproveFour = system.Text.Conv.Int64(ds[0]["approver4"])
	w.ApproveFive = system.Text.Conv.Int64(ds[0]["approver5"])
	w.ApproveSix = system.Text.Conv.Int64(ds[0]["approver6"])
	w.ApproveSeven = system.Text.Conv.Int64(ds[0]["approver7"])
	w.ApproveEight = system.Text.Conv.Int64(ds[0]["approver8"])
	w.ApproveNine = system.Text.Conv.Int64(ds[0]["approver9"])
	return nil
}

// 生成审批信息
func (wi *WFInst) Insert(tx *sql.Tx, submitCode string) (int64, error) {
	sql := "insert into wf_inst (inst_id,title,submitter,appv_cid,appv_nid,wf_state,crt_time,lasted_time,wf_id,wf_level)" +
		" select ?,?,user_id,?,?,0,now(),now(),?,? from user where user_code = ?"
	rc, err := system.Data.DB.ExecTx(tx, sql, wi.WFInstId, wi.Title, wi.AppvCid, wi.AppNid, wi.WFId, wi.WFLevel, submitCode)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

// 修改审批信息
func (wi *WFInst) Modify(tx *sql.Tx) (int64, error) {
	sql := "update wf_inst set title=?,appv_cid=?,appv_nid=?,wf_state=0,lasted_time=now(),wf_id=?,wf_level=? where inst_id=?"
	rc, err := system.Data.DB.ExecTx(tx, sql, wi.Title, wi.AppvCid, wi.AppNid, wi.WFId, wi.WFLevel, wi.WFInstId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (wi *WFInst) Update(tx *sql.Tx) (int64, error) {
	sql := "update wf_inst set appv_cid = ?,appv_nid=?,wf_state=?,wf_level=? ,lasted_time= now() where inst_id = ?"
	rc, err := system.Data.DB.Exec(sql, wi.AppvCid, wi.AppNid, wi.WfState, wi.WFLevel, wi.WFInstId)
	if err != nil {
		return 0, err
	}
	if rc == 0 {
		return 0, errors.New("No find inst formation")
	}
	return rc, nil
}

func (wi *WFInst) GetWfInstInfo() (int, error) {
	sql := "select * from wf_inst where inst_id = ?"
	ds, err := system.Data.DB.QueryRows(sql, wi.WFInstId)
	if err != nil {
		return 0, err
	}
	for _, rs := range ds {
		wi.WFLevel = system.Text.Conv.Int(rs["wf_level"])
		wi.WfState = system.Text.Conv.Int(rs["wf_state"])
		wi.AppvCid = system.Text.Conv.Int64(rs["appv_cid"])
		wi.AppNid = system.Text.Conv.Int64(rs["appv_nid"])
		wi.WFId = system.Text.Conv.Int64(rs["wf_id"])
		wi.Submitter = system.Text.Conv.Int64(rs["submitter"])
	}
	return len(ds), nil
}

func (wi *WFInst) SetNextApprover(wf WF) {
	switch wi.WFLevel {
	case 0:
		//表示重新修改提交
		wi.AppNid = wf.ApproveTwo
	case 1:
		wi.AppNid = wf.ApproveThree
		break
	case 2:
		wi.AppNid = wf.ApproveFour
		break
	case 3:
		wi.AppNid = wf.ApproveFive
		break
	case 4:
		wi.AppNid = wf.ApproveSix
		break
	case 5:
		wi.AppNid = wf.ApproveSeven
		break
	case 6:
		wi.AppNid = wf.ApproveEight
		break
	case 7:
		wi.AppNid = wf.ApproveNine
		break
	case 8:
		wi.AppNid = 0
		break
	default:
		break
	}
}

func (wi *WFInst) SetPreApprover(wf WF) {

	switch wi.WFLevel {
	case 1:
		wi.WFLevel = 1
		break
	case 2:
		wi.AppvCid = wf.ApproveOne
		break
	case 3:
		wi.AppvCid = wf.ApproveTwo
		break
	case 4:
		wi.AppvCid = wf.ApproveThree
		break
	case 5:
		wi.AppvCid = wf.ApproveFour
		break
	case 6:
		wi.AppvCid = wf.ApproveFive
		break
	case 7:
		wi.AppvCid = wf.ApproveSix
		break
	case 8:
		wi.AppvCid = wf.ApproveSeven
		break
	default:
		break
	}
}

func (wo *WFLog) InsertWfLog(tx *sql.Tx) error {
	id := system.Common.UUID.NextVal()
	sql := "insert into wf_log(log_id,inst_id,appv_id,appv_time,opr_id,opinion,wfo_state,wfn_state,next_appv_id) values(?,?,?,now(),?,?,?,?,?)"
	_, err := system.Data.DB.ExecTx(tx, sql, id, wo.InstId, wo.AppvId, wo.OprId, wo.Opinion, wo.WfoState, wo.WfnState, wo.NextAppvId)
	if err != nil {
		return err
	}
	return nil
}
func (wo *WFLog) GetWFLog() ([]WFLogShow, error) {
	var result []WFLogShow
	sql := "SELECT d.opr_name,date_format(a.appv_time,'%Y-%m-%d %H:%i:%s') approve_time, a.log_id,a.inst_id,a.appv_time,a.opinion,a.wfn_state,a.wfo_state,CONCAT(b.user_code ,\"|\",b.user_name) approver,CONCAT(c.user_code ,\"|\",c.user_name) next_approver from wf_log a " +
		"INNER JOIN wf_opr d on a.opr_id = d.opr_id LEFT JOIN user b ON a.appv_id = b.user_id LEFT JOIN user c on a.next_appv_id = c.user_id where inst_id = ? order by log_id "
	ds, err := system.Data.DB.QueryResultRows(sql, wo.InstId)
	if err != nil {
		return result, err
	}
	if len(ds) == 0 {
		return result, errors.New("empty")
	}
	for _, rc := range ds {
		var log WFLogShow
		if err := mapstructure.Decode(rc, &log); err != nil {
			system.Log.Error("Parase work flow log errr:", err.Error())
			continue
		}
		//log.AppvTime = log.AppvTimeT.Format("2006-01-02 15:04:05")
		result = append(result, log)
	}
	return result, nil
}
func (his *WFHisSearch) FindByIndex(start, size int, cond string, args ...any) ([]AuthRes, error) {
	auths := make([]AuthRes, 0)
	sql := "SELECT a.wf_inst_id,a.auth_id,date_format(a.created_time,'%Y-%m-%d %H:%i:%s') created_time,a.proposer,c.inst_name,a.reason,a.host_ip" +
		",a.host_mac,a.appv_state,date_format(a.appv_time,'%Y-%m-%d %H:%i:%s') appv_time,a.auth_state,e.user_code,e.user_name,b.unit_name" +
		" from auth_req a,unit b,scene c,user e,wf_log d  where c.unit_id = b.unit_id  and a.inst_id= c.inst_id  and a.wf_inst_id = d.inst_id and e.user_id = d.appv_id  "
	if cond != "" {
		sql += cond
	}

	sql += " group by a.auth_id "
	if size > 0 {
		sql += " order by a.created_time desc, a.inst_id limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	}
	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return auths, err
	}
	for _, rs := range ds {
		var authResp AuthRes
		authResp.CurApprover = rs["user_code"] + "|" + rs["user_name"]
		authResp.UnitName = rs["unit_name"]
		authResp.InstName = rs["inst_name"]
		hosts := system.Text.String.Split(rs["host_ip"], ",")
		for _, host := range hosts {
			authResp.HostIp = append(authResp.HostIp, host)
		}
		macs := system.Text.String.Split(rs["host_mac"], ",")
		for _, mac := range macs {
			authResp.HostMac = append(authResp.HostMac, mac)
		}
		authResp.Reason = rs["reason"]
		authResp.Proposer = rs["proposer"]
		authResp.AppvTime = rs["appv_time"]
		authResp.AppvState = system.Text.Conv.Int(rs["appv_state"])
		authResp.AuthState = system.Text.Conv.Int(rs["auth_state"])
		authResp.CreatedTime = rs["created_time"]
		authResp.AuthId = system.Text.Conv.Int64(rs["auth_id"])
		authResp.WfInstId = system.Text.Conv.Int64(rs["wf_inst_id"])
		auths = append(auths, authResp)
	}
	return auths, nil
}
