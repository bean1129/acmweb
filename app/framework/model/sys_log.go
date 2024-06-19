package model

import (
	"acmweb/system"
	"acmweb/system/text"
)

type SysLog struct {
	UnitId  int64  `mapstructure:"unit_id"`
	IpAddr  string `mapstructure:"ip_addr"`
	Mac     string `mapstructure:"mac"`
	AppCode string `mapstructure:"app_code"`
	LogType int    `mapstructure:"log_type"`
	LogTime string `mapstructure:"log_time"`
}

func (s *SysLog) Insert() (affect int64, err error) {
	sql := "insert into syslog(log_id,unit_id,ip_addr,mac,app_code,log_type,log_time) values(?,?,?,?,?,?)"
	affect, err = system.Data.DB.Exec(sql, system.Common.UUID.NextVal(), s.UnitId, s.IpAddr, s.Mac, s.AppCode, s.LogType, s.LogTime)
	if err != nil {
		return 0, err
	}
	return affect, nil
}

type LOGTYPE int

const (
	AuthCert LOGTYPE = iota + 1
)

func GetLogTypeDes(val LOGTYPE) string {
	switch val {
	case AuthCert:
		return "授权认证"
	}

	return "unknow"
}

type SysLogShow struct {
	UnitId    int64  `json:"unit_id"`
	UnitName  string `json:"unit_name"`
	AppCode   string `json:"app_code"`
	AppName   string `json:"app_name"`
	IpAddr    string `json:"ip_addr"`
	Mac       string `json:"mac"`
	LogType   string `json:"log_type"`
	LogTypeId int    `json:"log_type_id"`
	LogTime   string `json:"log_time"`
}

type SysLogSearch struct {
	PageNum  int    `mapstructure:"page_num"`
	PageSize int    `mapstructure:"page_size"`
	Search   string `mapstructure:"search"`
	LogType  int    `mapstructure:"log_type"`
}

func (ss *SysLogSearch) FindTotalNum() (int64, error) {
	var args []any
	sqlQry := "select count(*) total from syslog a ,unit b where a.unit_id = b.unit_id "
	if ss.LogType > 0 {
		sqlQry += " and a.log_type= ?"
		args = append(args, ss.LogType)
	}
	if len(ss.Search) > 0 {
		sqlQry += " and (b.unit_code = ? or b.unit_name = ?)"
		args = append(args, ss.Search, ss.Search)
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return 0, err
	}
	result := system.Text.Conv.Int64(ds[0]["total"])
	return result, nil
}

func (ss *SysLogSearch) FindByIndex() (result []SysLogShow, err error) {
	var args []any
	start := (ss.PageNum - 1) * ss.PageSize
	size := ss.PageSize
	sqlQry := "select a.unit_id,b.unit_name,a.app_code,c.app_name,a.ip_addr,a.mac,a.log_type,date_format(a.log_time,'%Y-%m-%d %H:%i:%s') log_time from syslog a,unit b,app c where a.unit_id = b.unit_id and a.app_code = c.app_code"
	if ss.LogType > 0 {
		sqlQry += " and log_type= ?"
		args = append(args, ss.LogType)
	}
	if len(ss.Search) > 0 {
		sqlQry += " and (b.unit_code = ? or b.unit_name = ?)"
		args = append(args, ss.Search, ss.Search)
	}
	if size > 0 {
		sqlQry += " limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return result, err
	}
	for _, rs := range ds {
		var ret SysLogShow
		ret.UnitName = rs["unit_name"]
		ret.AppCode = rs["app_code"]
		ret.AppName = rs["app_name"]
		ret.IpAddr = rs["ip_addr"]
		ret.UnitId = system.Text.Conv.Int64(rs["unit_id"])
		ret.Mac = rs["mac"]
		ret.LogTypeId = system.Text.Conv.Int(rs["log_type"])
		ret.LogType = GetLogTypeDes(LOGTYPE(ret.LogTypeId))
		ret.LogTime = rs["log_time"]
		result = append(result, ret)
	}
	return result, nil
}
