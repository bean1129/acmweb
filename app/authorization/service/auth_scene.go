package service

import (
	"acmweb/app/authorization/model"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
)

var (
	Scene = new(SceneService)
)

type SceneService struct {
}

func (a *SceneService) CreateScene(m model.Scene) (*common.ZIResult, error) {
	//定义应用的属性
	zs := common.NewResult()
	if m.InstId == 0 {
		if bExist, err := m.ExistsScene(); err != nil {
			system.Log.Error("Query scene information failed,err:", err.Error())
			return zs, &errors.ErrDBQry
		} else {
			if bExist {
				return zs, &errors.ErrNameExists
			}
		}
	}
	//插入
	var err error
	if err = m.Insert(); err != nil {
		system.Log.Error("Insert Scene information failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	zs.Count = 1
	system.Log.Infof("Create scene success,id[%d] name[%s]", m.InstId, m.InstName)
	return zs, nil
}

func (a *SceneService) DeleteScene(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.Scene{InstId: id}
		row, _ := m.Delete()
		zs.Count += row
	}
	return zs, nil
}
func (a *SceneService) DeleteGroup(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.SceneGroup{GroupId: id}
		if err := m.Delete(); err != nil {
			zs.Count += 1
		}
	}
	return zs, nil
}
func (a *SceneService) CreateGrp(m model.SceneGroup) (*common.ZIResult, error) {
	//定义应用的属性
	zs := common.NewResult()
	//插入
	var err error
	if err = m.Insert(); err != nil {
		system.Log.Error("Insert scene group information failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	zs.Count = 1
	system.Log.Infof("Create scene group success,id[%d] name[%s]", m.GroupId, m.GrpName)
	return zs, nil
}
func (a *SceneService) ListScene(search *model.SceneSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	//是否有分页
	var apps = make([]model.SceneShow, 0)
	var err error
	apps, err = search.FindByIndex()

	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range apps {
		zs.Data = append(zs.Data, v)
	}
	zs.Count, _ = search.FindTotalNum()
	return zs, nil
}
func (a *SceneService) ListGroup() (*common.ZIResult, error) {
	zs := common.NewResult()
	//是否有分页
	var err error
	search := model.SceneGroup{}
	grps, err := search.List()
	if err != nil {
		return zs, &errors.ErrDBexecute
	}
	for _, v := range grps {
		zs.Data = append(zs.Data, v)
	}
	zs.Count = int64(len(grps))
	return zs, nil
}

func (a *SceneService) ListSceneApp(id []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	//是否有分页
	var err error
	search := model.Scene{InstId: id[0]}
	results, err := search.ListApp()
	if err != nil {
		return zs, &errors.ErrDBexecute
	}
	for _, v := range results {
		zs.Data = append(zs.Data, v)
	}
	zs.Count = int64(len(results))
	return zs, nil
}
