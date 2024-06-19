package service

import (
	"acmweb/app/authorization/model"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
)

var (
	AuthWF = new(AuthWFService)
)

type AuthWFService struct {
}

// CreateApp 创建应用
func (c *AuthWFService) CreateWorkFlow(m model.WF) (*common.ZIResult, error) {
	zs := common.NewResult()
	//判断组是否已存在流程
	//生成
	if m.WFId == 0 {
		m.WFId = system.Common.UUID.NextVal()
	}
	if err := m.CreateWf(); err != nil {
		system.Log.Error("Insert work flow failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	zs.Count = 1

	return zs, nil
}
func (c *AuthWFService) Delete(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.WF{WFId: id}
		rows, err := m.DeleteWf()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}
func (c *AuthWFService) RenameWorkFlow(m model.WF) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.UpdateWf("wf_name=?", m.WFName); err != nil {
		system.Log.Error("Rename app failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}
func (c *AuthWFService) GetWorkFlow(start, pageSize int, grpID int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	var wfResult []model.WFInfo
	var m model.WF
	var err error
	wfResult, err = m.GetWFByCond(start, pageSize, grpID)
	if err != nil {
		system.Log.Error("Get work flow failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	for _, v := range wfResult {
		zs.Data = append(zs.Data, v)
		zs.Count += 1
	}
	zs.Count, _ = m.FindTotalNum(grpID)
	return zs, nil
}
func (c *AuthWFService) GetWFLog(instId int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m = model.WFLog{InstId: instId}
	var err error
	logs, err := m.GetWFLog()
	if err != nil {
		system.Log.Error("Get work flow log failed,err:", err.Error())
		return zs, &errors.ErrDBQry
	}
	for _, v := range logs {
		zs.Data = append(zs.Data, v)
		zs.Count += 1
	}
	return zs, nil
}

func (c *AuthWFService) GetHisAppv(search model.WFHisSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	var args = make([]any, 0)
	var cond string
	cond = " and (d.appv_id = ? or d.next_appv_id = ?)"
	args = append(args, search.ApproverId, search.ApproverId)
	if search.AppvState != -1 {
		cond += " and d.wfo_state = ?"
		args = append(args, search.AppvState)
	} else {
		cond += " and 1 =1 "
	}
	start := (search.PageNum - 1) * search.PageSize
	//是否有分页
	var authInfos = make([]model.AuthRes, 0)
	var err error
	authInfos, err = search.FindByIndex(start, search.PageSize, cond, args...)

	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range authInfos {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}
