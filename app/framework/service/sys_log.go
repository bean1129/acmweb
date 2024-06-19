package service

import (
	"acmweb/app/errors"
	"acmweb/app/framework/model"
	"acmweb/system"
	"acmweb/system/common"
)

var (
	SysLog = new(SysLogService)
)

type SysLogService struct {
}

func (c *SysLogService) CreateSysLog(syslog model.SysLog) (*common.ZIResult, error) {
	//插入
	zs := common.NewResult()
	var err error
	if zs.Count, err = syslog.Insert(); err != nil {
		system.Log.Error("Insert syslog information failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}
func (c *SysLogService) List(search model.SysLogSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	//是否有分页
	var err error
	results, err := search.FindByIndex()
	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range results {
		zs.Data = append(zs.Data, v)
	}
	zs.Count, _ = search.FindTotalNum()
	return zs, nil
}
