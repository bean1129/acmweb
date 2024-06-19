package model

import (
	"acmweb/system"
	"acmweb/system/text"
	"database/sql"

	"github.com/pkg/errors"
)

type AuthSearch struct {
	PageNum      int    `mapstructure:"page_num" json:"page_num,omitempty"`
	PageSize     int    `mapstructure:"page_size" json:"page_size,omitempty"`
	InstName     string `mapstructure:"inst_name" json:"inst_name,omitempty"`
	AuthState    int    `mapstructure:"auth_state" json:"auth_state,omitempty"`
	ApproveState int    `mapstructure:"approve_state" json:"approve_state,omitempty"`
	NextApprover string `mapstructure:"next_approver" json:"next_approver"`
	CurApprover  string `mapstructure:"cur_approver" json:"cur_approver"`
	UnitId       int64  `mapstructure:"unit_id" json:"unit_id"`
	Submitter    string `mapstructure:"proposer" json:"proposer"`
}

// 授权申请信息
type AuthApply struct {
	AuthId     int64    `mapstructure:"auth_id" json:"auth_id"`         // 授权申请标识
	Proposer   string   `mapstructure:"proposer" json:"proposer"`       // 申请人姓名
	ProposerId int64    `mapstructure:"proposer_id" json:"proposer_id"` // 申请人标识
	InstId     int64    `mapstructure:"inst_id" json:"inst_id"`         // 场景ID
	Host       []string `mapstructure:"host_ip" json:"host_ip"`         // 客户端主机地址
	Mac        []string `mapstructure:"host_mac" json:"host_mac"`       // 客户端设备mac地址
	CrtTime    string   `mapstructure:"create_time" json:"create_time"` // 创建时间
	AppvState  int      `mapstructure:"appv_state" json:"appv_state"`   // 审批状态：0-待审,1-审批通过,2-审批中,3-作废
	AppvTime   string   `mapstructure:"appv_time" json:"appv_time"`     // 审批时间
	AuthState  int      `mapstructure:"state" json:"auth_state"`        // 授权状态：0-待授权，1-已授权
	WFInstId   int64    `mapstructure:"wf_inst_id" json:"wf_inst_id"`   // 工作流Id
	Cert       string   `json:"cert,omitempty"`                         // 证书
	CertMD5    string   `json:"cert_md5,omitempty"`                     // 证书
	Reason     string   `mapstructure:"reason" json:"reason"`           // 申请理由
	AppList    string   `mapstructure:"app_list" json:"app_list,omitempty"`
	TrialDays  int      `mapstructure:"trial_days" json:"trial_days,omitempty"`
	AuthMode   int      `mapstructure:"auth_mode" json:"auth_mode,omitempty"`
	UnitId     int64    `mapstructure:"unit_id" json:"unit_id,omitempty"`
	Scene
}

// 权限请求信息响应
type AuthRes struct {
	AuthId       int64    `mapstructure:"auth_id" json:"auth_id"`         //授权申请标识
	UnitName     string   `mapstructure:"unit_name" json:"unit_name"`     //申请组织标识
	Proposer     string   `mapstructure:"proposer" json:"proposer"`       // 申请人姓名
	AppList      string   `mapstructure:"app_list" json:"app_list"`       // 申请人姓名
	InstName     string   `mapstructure:"inst_name" json:"inst_name"`     // 场景
	ModuleList   string   `mapstructure:"module_list" json:"module_list"` // 应用的模块
	TrialDays    int      `mapstructure:"trial_days" json:"trial_days"`
	HostIp       []string `mapstructure:"host_ip" json:"host_ip"`   //客户端主机地址
	HostId       string   `mapstructure:"host_id" json:"host_id"`   //客户端主机地址
	HostMac      []string `mapstructure:"host_mac" json:"host_mac"` //客户端设备mac地址
	CreatedTime  string   `mapstructure:"created_time" json:"created_time"`
	AppvState    int      `mapstructure:"appv_state" json:"appv_state"`        //审批状态：0-待审,1-审批通过,2-审批中,3-作废
	AppvTime     string   `mapstructure:"appv_time" json:"appv_time"`          //审批时间
	AuthState    int      `mapstructure:"state" json:"auth_state"`             //授权状态：0-待授权，1-已授权
	Cert         string   `mapstructure:"cert" json:"cert,omitempty" `         //证书
	CertMD5      string   `mapstructure:"cert_md5" json:"cert_md5,omitempty" ` //证书
	CurApprover  string   `mapstructure:"cur_approver" json:"cur_approver"`
	NextApprover string   `mapstructure:"next_approver" json:"next_approver"`
	WfInstId     int64    `mapstructure:"wf_inst_id" json:"wf_inst_id"`
	Reason       string   `mapstructure:"reason" json:"reason"`
	InstId       int64    `mapstructure:"inst_id" json:"inst_id"`
	AuthMode     int      `mapstructure:"auth_mode" json:"auth_mode"`
	AuthModeCode string   `mapstructure:"auth_mode_code" json:"auth_mode_code"`
	AuthModeName string   `mapstructure:"auth_mode_name" json:"auth_mode_name"`
}
type AuthResNum struct {
	WfstateIszero  int `json:"wfstate_iszero"`
	WfstateIstwo   int `json:"wfstate_istwo"`
	AuthstateIsone int `json:"authstate_isone"`
}

// 授权审批信息
type AuthApprove struct {
	WfInstId   int64  `mapstructure:"wf_inst_id"`
	Approver   int64  `mapstructure:"approver_id"`
	WfOprState int    `mapstructure:"opr_state"`
	Opinion    string `mapstructure:"opinion"`
}
type AuthMode struct {
	TypeId   int    `json:"type_id"`
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
}

// 查询流程定义
func (a *AuthApply) QueryFlow() {

}

// 单位和应用的申请数
func (a *AuthApply) GetApplyCount() (int, error) {
	sql := "select sum(a.auth_id) rstCount from auth_req a,scene b " +
		"where a.inst_id = b.inst_id and b.inst_id = ? "
	args := []any{
		a.InstId,
	}
	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return -1, err
	}
	return system.Text.Conv.Int(ds[0]["rstCount"]), nil
}

// 生成流程实例
func (a *AuthApply) Insert(tx *sql.Tx) (int64, error) {
	if a.AuthId == 0 {
		a.AuthId = system.Common.UUID.NextVal()
	}
	//查询
	if 1 != a.InstId {
		sqlOtherQry := "SELECT b.trial_days trial_days,GROUP_CONCAT(d.app_code) app_list,b.auth_mode,b.unit_id from scene b,scene_app c,app d " +
			"where  b.inst_id=c.inst_id and c.app_id=d.app_id and b.inst_id=?"
		ds, err := system.Data.DB.QueryRows(sqlOtherQry, a.InstId)
		if err != nil || len(ds) == 0 {
			return 0, err
		}
		a.TrialDays = system.Text.Conv.Int(ds[0]["trial_days"])
		a.AppList = ds[0]["app_list"]
		a.AuthMode = system.Text.Conv.Int(ds[0]["auth_mode"])
		a.UnitId = system.Text.Conv.Int64(ds[0]["unit_id"])
	}

	sql := "insert into auth_req(auth_id,proposer,inst_id,host_ip,host_mac,created_time,appv_state,auth_state,wf_inst_id,reason,trial_days,app_list,auth_mode,unit_id)" +
		" values(?,?,?,?,?,now(),0,0,?,?,?,?,?,?)" +
		" ON DUPLICATE KEY UPDATE inst_id=values(inst_id)" +
		",host_ip=values(host_ip),host_mac=values(host_mac),appv_state=values(appv_state),reason=values(reason)" +
		",trial_days=values(trial_days),app_list=values(app_list),auth_mode=values(auth_mode),unit_id=values(unit_id)"
	var hosts string
	var macs string
	for idx, host := range a.Host {
		if idx == 0 {
			hosts = host
		} else {
			hosts += "," + host
		}
	}
	for idx, mac := range a.Mac {
		if idx == 0 {
			macs = mac
		} else {
			macs += "," + mac
		}
	}
	rc, err := system.Data.DB.ExecTx(tx, sql, a.AuthId, a.Proposer, a.InstId, hosts, macs, a.WFInstId, a.Reason, a.TrialDays, a.AppList, a.AuthMode, a.UnitId)
	if err != nil {
		return 0, err
	}

	return rc, nil
}

func (a *AuthApply) UpdateApproveState(tx *sql.Tx, wfState int) (int64, error) {
	sql := "update auth_req set appv_state= ?,appv_time=now() where wf_inst_id=?"
	rc, err := system.Data.DB.ExecTx(tx, sql, wfState, a.WFInstId)
	if err != nil {
		return 0, err
	}
	if rc == 0 {
		return 0, errors.New("No find auth req formation")
	}
	return rc, nil
}

func (a *AuthApply) UpdateCert() (int64, error) {
	sql := "update auth_req set cert=?,cert_md5=?,auth_state=1 where "
	var args int64
	if a.AuthId > 0 {
		args = a.AuthId
		sql += " auth_id = ?"
	} else if a.WFInstId > 0 {
		args = a.WFInstId
		sql += " wf_inst_id = ?"
	} else {
		return 0, errors.New("cond err")
	}
	rc, err := system.Data.DB.Exec(sql, a.Cert, a.CertMD5, args)
	if err != nil {
		return 0, err
	}
	if rc == 0 {
		//return 0, errors.New("Update effective empty.")
	}
	return rc, nil
}

func (a *AuthApply) GetAuthReqInfo() error {
	sql := "select * from auth_req where auth_id = ? or wf_inst_id = ?"
	ds, err := system.Data.DB.QueryRows(sql, a.AuthId, a.WFInstId)
	if err != nil {
		return err
	}
	for _, rs := range ds {
		a.AuthId = system.Text.Conv.Int64(rs["auth_id"])
		a.WFInstId = system.Text.Conv.Int64(rs["wf_inst_id"])
		break
	}
	return nil
}

func (a *AuthApply) GetSubmitterId() error {
	sql := "select user_id from user where user_code = ?"
	ds, err := system.Data.DB.QueryRows(sql, a.Proposer)
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return errors.New("empty")
	}
	a.ProposerId = system.Text.Conv.Int64(ds[0]["user_id"])
	return nil
}

func (a *AuthApply) DeleteAuth() (int64, error) {
	sql := "delete from  auth_req where auth_id = ?"
	rc, err := system.Data.DB.Exec(sql, a.AuthId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (a *AuthSearch) FindTotalNum(cond string, args ...any) (int64, error) {
	sql := "SELECT count(a.auth_id) total from auth_req a,unit b,scene c,wf_inst d LEFT JOIN user f on d.appv_nid = f.user_id " +
		" LEFT JOIN user e on e.user_id = d.appv_cid " +
		"where c.unit_id = b.unit_id  and a.inst_id= c.inst_id  and a.wf_inst_id = d.inst_id and 1=1 "
	if cond != "" {
		sql += cond
	}
	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return 0, err
	}
	result := system.Text.Conv.Int64(ds[0]["total"])
	return result, nil
}

func (a *AuthSearch) GetVaildAuth() (result map[int64]string, err error) {
	sqlQry := "SELECT auth_id ,cert_md5 from auth_req where auth_state=1"
	ds, err := system.Data.DB.QueryRows(sqlQry)
	if err != nil {
		return
	}
	result = make(map[int64]string)
	for _, rs := range ds {
		result[system.Text.Conv.Int64(rs["auth_id"])] = rs["cert_md5"]
	}
	return
}

func (a *AuthSearch) FindByIndex(start, size int, cond string, args ...any) ([]AuthRes, error) {
	auths := make([]AuthRes, 0)
	sql := "SELECT a.auth_mode,g.type_code ,g.type_name,a.app_list,a.inst_id,a.unit_id,a.reason,a.wf_inst_id,a.auth_id,date_format(a.created_time,'%Y-%m-%d %H:%i:%s') created_time" +
		",a.proposer,c.inst_name,a.trial_days,a.host_ip,a.host_mac,a.appv_state,date_format(a.appv_time,'%Y-%m-%d %H:%i:%s') appv_time" +
		",a.auth_state,e.user_code,e.user_name,f.user_code next_user,f.user_name next_name,b.unit_name,a.cert,a.cert_md5" +
		" from auth_req a left join auth_mode g on a.auth_mode=g.type_id left join unit b on a.unit_id = b.unit_id ,scene c,wf_inst d LEFT JOIN user f on d.appv_nid = f.user_id  " +
		"LEFT JOIN user e on e.user_id = d.appv_cid " +
		"where a.inst_id= c.inst_id  and a.wf_inst_id = d.inst_id  and 1=1 "
	if cond != "" {
		sql += cond
	}
	if size > 0 {
		sql += " order by a.created_time desc,a.inst_id limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	}
	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return auths, err
	}
	for _, rs := range ds {
		var authResp AuthRes
		if len(rs["user_code"]) > 0 {
			authResp.CurApprover = rs["user_code"] + "|" + rs["user_name"]
		}
		if len(rs["next_user"]) > 0 {
			authResp.NextApprover = rs["next_user"] + "|" + rs["next_name"]
		}
		authResp.AuthMode = system.Text.Conv.Int(rs["auth_mode"])
		authResp.AuthModeCode = rs["type_code"]
		authResp.AuthModeName = rs["type_name"]
		authResp.UnitName = rs["unit_name"]
		authResp.InstName = rs["inst_name"]
		authResp.AppList = rs["app_list"]
		authResp.Cert = rs["cert"]
		authResp.CertMD5 = rs["cert_md5"]
		authResp.TrialDays = system.Text.Conv.Int(rs["trial_days"])
		Macs := system.Text.String.Split(rs["host_mac"], ",")
		for _, mac := range Macs {
			authResp.HostMac = append(authResp.HostMac, mac)
		}
		Ips := system.Text.String.Split(rs["host_ip"], ",")
		for _, ip := range Ips {
			authResp.HostIp = append(authResp.HostIp, ip)
		}
		if len(authResp.HostMac) == 1 && len(authResp.HostIp) == 1 && authResp.HostMac[0] == "" && authResp.HostIp[0] == "" {
			authResp.HostIp = authResp.HostIp[:0]
			authResp.HostMac = authResp.HostMac[:0]
		}
		authResp.Proposer = rs["proposer"]
		authResp.AppvTime = rs["appv_time"]
		authResp.AppvState = system.Text.Conv.Int(rs["appv_state"])
		authResp.AuthState = system.Text.Conv.Int(rs["auth_state"])
		authResp.CreatedTime = rs["created_time"]
		authResp.AuthId = system.Text.Conv.Int64(rs["auth_id"])
		authResp.WfInstId = system.Text.Conv.Int64(rs["wf_inst_id"])
		authResp.InstId = system.Text.Conv.Int64(rs["inst_id"])
		authResp.Reason = rs["reason"]
		auths = append(auths, authResp)
	}
	return auths, nil
}
func (a *AuthSearch) FindAuthNum() (AuthResNum, error) {
	var auths AuthResNum
	sql := "SELECT COUNT(CASE WHEN auth_state=1 THEN 1 ELSE null END) as authstate_isone,COUNT(CASE WHEN appv_state=0 THEN 1 ELSE null END) as wfstate_iszero,COUNT(CASE WHEN appv_state=2 THEN 1 ELSE null END) as wfstate_istwo from auth_req"
	ds, err := system.Data.DB.QueryRows(sql)
	if err != nil {
		return auths, err
	}
	if len(ds) == 0 {
		return auths, errors.New("empty")
	}
	auths.WfstateIstwo = system.Text.Conv.Int(ds[0]["wfstate_istwo"])
	auths.WfstateIszero = system.Text.Conv.Int(ds[0]["wfstate_iszero"])
	auths.AuthstateIsone = system.Text.Conv.Int(ds[0]["authstate_isone"])
	return auths, nil
}

type AuthStatis struct {
	ResultData []map[string]interface{} `mapstructure:"result_data" json:"result_data"`
	HeaderData []Header                 `mapstructure:"header_data" json:"header_data"`
}

type Header struct {
	Label string `mapstructure:"label" json:"label"`
	Prop  string `mapstructure:"prop" json:"prop"`
}

func (a *AuthStatis) GetTierLevel(tierLevel int) (ret []string, err error) {
	sql := "select tier_name,tier_level from unit_tier where tier_level <= ? ORDER BY tier_level"
	ds, err := system.Data.DB.QueryRows(sql, tierLevel)
	if err != nil {
		return ret, err
	}
	for _, v := range ds {
		propVal := "tier_level" + v["tier_level"]
		a.HeaderData = append(a.HeaderData, Header{Label: v["tier_name"], Prop: propVal})
		ret = append(ret, propVal)
	}
	return ret, err
}

func (a *AuthStatis) GetHeader(appId []int64, tierLevel int) (result []int64, err error) {

	//应用
	var field string
	var fieldVal []any
	for i, v := range appId {
		if i == 0 {
			field = "?"
		} else {
			field += ",?"
		}
		fieldVal = append(fieldVal, v)
	}
	var sql string
	if len(appId) == 0 {
		sql = "select app_id prop,app_name label from app"
	} else {
		sql = "select app_id prop,app_name label from app where app_id in(" + field + ")"
	}
	ds, err := system.Data.DB.QueryRows(sql, fieldVal...)
	if err != nil {
		return nil, err
	}

	if len(appId) != 0 && len(appId) != len(ds) {
		return nil, errors.New("Query app not match app table data")
	}
	for i, rs := range ds {
		result = append(result, system.Text.Conv.Int64(rs["prop"]))
		qryAppColumn := "app" + system.Text.Conv.String(i)
		a.HeaderData = append(a.HeaderData, Header{Label: rs["label"], Prop: qryAppColumn})
	}
	a.HeaderData = append(a.HeaderData, Header{Label: "授权总数", Prop: "auth_total"})
	a.HeaderData = append(a.HeaderData, Header{Label: "应用总数", Prop: "app_total"})
	return result, nil
}

func (a *AuthStatis) GetStatisData(appId []int64, tierLevel []string) (err error) {
	unitColumnUp1 := "(SELECT unit_name from unit t1 where t1.unit_id = a.parent_unit_id) "
	unitColumnUp2 := "(SELECT t2.unit_name from unit t1,unit t2 where t1.parent_unit_id = t2.unit_id and t1.unit_id = a.parent_unit_id) "
	var columnUnit string
	if len(tierLevel) == 3 {
		columnUnit = unitColumnUp1 + tierLevel[1] + "," + unitColumnUp2 + tierLevel[0] + ","
	} else if len(tierLevel) == 2 {
		columnUnit = unitColumnUp1 + tierLevel[0] + ","
	}
	columnUnit += "a.unit_name " + tierLevel[len(tierLevel)-1]
	//应用
	var qryVal []any
	for i, v := range appId {
		qryColumn := "app" + text.Conv.String(i)
		columnUnit += ",COUNT( CASE WHEN d.app_id=? THEN 1 ELSE NULL END ) AS " + qryColumn
		qryVal = append(qryVal, v)
	}
	//固定列
	columnUnit += ",COUNT(CASE WHEN auth_id >0 THEN 1 ELSE NULL END) AS auth_total,COUNT(CASE WHEN d.app_id >0 THEN 1 ELSE NULL END) AS app_total"
	//查询
	sql := "select " + columnUnit + " from  unit a INNER JOIN unit_tier b on a.tier_id = b.tier_id and b.tier_level=? " +
		"LEFT JOIN auth_req c on a.unit_id=c.unit_id and c.auth_state = 1 " +
		"LEFT JOIN app d on  c.app_id = d.app_id GROUP BY a.unit_code"
	qryVal = append(qryVal, len(tierLevel))
	a.ResultData, err = system.Data.DB.QueryResultRows(sql, qryVal...)
	return err
}

type AuthAxisStatis struct {
	XData string      `json:"x_data"`
	YData []AixsYData `json:"y_data"`
}

type AixsYData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (aa *AuthAxisStatis) GetUnitStatisByApp(appId []int64, tier int) ([]AuthAxisStatis, error) {
	appQryCond := ""
	var args []any
	args = append(args, tier)
	if len(appId) > 0 {
		for index, app := range appId {
			if 0 == index {
				appQryCond = " and a.app_id in(?"
			} else {
				appQryCond += ",?"
			}
			args = append(args, app)
		}
		appQryCond += ") "
	}
	sql := "SELECT b.app_name,c.unit_name,count(1) as statis_num from auth_req a,app b,unit c,unit d " +
		"where a.app_id = b.app_id and a.unit_id = c.unit_id and d.tier_id = ? " +
		appQryCond +
		" and a.auth_state=1 and LEFT(c.unit_code,4)=d.unit_code group by a.app_id,a.unit_id order by a.app_id"

	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		return nil, errors.New("empty!")
	}
	var retdata []AuthAxisStatis
	var data AuthAxisStatis
	for _, d := range ds {
		xname := d["app_name"]
		if xname != data.XData {
			if len(data.XData) > 0 {
				retdata = append(retdata, data)
			}
			data.XData = xname
			data.YData = []AixsYData{}
		}
		oneYData := AixsYData{Name: d["unit_name"], Count: system.Text.Conv.Int(d["statis_num"])}
		data.YData = append(data.YData, oneYData)
	}
	if len(data.XData) > 0 {
		retdata = append(retdata, data)
	}
	return retdata, nil
}
func (aa *AuthAxisStatis) GetAppStatisByUnit(appId []int64, tier int) ([]AuthAxisStatis, error) {
	appQryCond := ""
	var args []any
	args = append(args, tier)
	if len(appId) > 0 {
		for index, app := range appId {
			if 0 == index {
				appQryCond = " and a.app_id in(?"
			} else {
				appQryCond += ",?"
			}
			args = append(args, app)
		}
		appQryCond += ") "
	}
	sql := "SELECT b.app_name,c.unit_name,count(1) as statis_num from auth_req a,app b,unit c,unit d " +
		"where a.app_id = b.app_id and a.unit_id = c.unit_id and d.tier_id = ? " +
		appQryCond +
		" and a.auth_state=1 and LEFT(c.unit_code,4)=d.unit_code group by a.unit_id,a.app_id order by a.unit_id"

	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(ds) == 0 {
		return nil, errors.New("empty!")
	}
	var retdata []AuthAxisStatis
	var data AuthAxisStatis
	for _, d := range ds {
		xname := d["unit_name"]
		if xname != data.XData {
			if len(data.XData) > 0 {
				retdata = append(retdata, data)
			}
			data.XData = xname
			data.YData = []AixsYData{}
		}
		oneYData := AixsYData{Name: d["app_name"], Count: system.Text.Conv.Int(d["statis_num"])}
		data.YData = append(data.YData, oneYData)
	}
	if len(data.XData) > 0 {
		retdata = append(retdata, data)
	}
	return retdata, nil
}
