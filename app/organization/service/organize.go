package service

import (
	"acmweb/app/errors"
	"acmweb/app/organization/model"
	"acmweb/system"
	"acmweb/system/common"
)

var Org = new(OrganizeSerivce)

type OrganizeSerivce struct {
}

func (c *OrganizeSerivce) Create(u *model.Unit) (*common.ZIResult, error) {
	u.UnitId = system.Common.UUID.NextVal()
	zs := common.NewResult()
	var err error
	if bExists, err := u.Exists(); err != nil {
		system.Log.Error("Query unit failed ,err:", err.Error())
		return zs, &errors.ErrDBQry
	} else {
		if bExists {
			return zs, &errors.ErrNameExists
		}
	}

	//生成编号
	err = u.SetUnitLevel()
	if err != nil {
		system.Log.Error("Invaild parent unit,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	zs.Count, err = u.Insert()
	if err != nil {
		system.Log.Error("Insert unit failed ,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *OrganizeSerivce) Delete(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.Unit{UnitId: id}
		rows, err := m.Delete()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}

func (c *OrganizeSerivce) Rename(m model.Unit) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if bExists, err := m.Exists(); err != nil {
		system.Log.Error("Query unit failed ,err:", err.Error())
		return zs, &errors.ErrDBQry
	} else {
		if bExists {
			return zs, &errors.ErrNameExists
		}
	}
	//更新部门名称同时更子级的全称
	if zs.Count, err = m.UpdateUnitName(); err != nil {
		system.Log.Error("Update unit name failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *OrganizeSerivce) List(search *model.UnitSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.Unit
	var cond string
	var args []any
	if len(search.UnitName) > 0 {
		cond += " and a.unit_name like ? "
		args = append(args, "%"+search.UnitName+"%")
	}
	units, err := m.Find(cond, args...)
	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range units {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}
func (c *OrganizeSerivce) ListTier() (*common.ZIResult, error) {
	zs := common.NewResult()
	sql := "select tier_id,tier_level,tier_name,remark from unit_tier"
	ds, err := system.Data.DB.QueryRows(sql)
	if err != nil {
		return nil, err
	}
	for _, v := range ds {
		m := model.UnitTier{
			TierId:    system.Text.Conv.Int(v["tier_id"]),
			TierName:  v["tier_name"],
			TierLevel: system.Text.Conv.Int(v["tier_level"]),
			Remark:    v["remark"],
		}
		zs.Data = append(zs.Data, m)
		zs.Count++
	}
	return zs, nil
}
