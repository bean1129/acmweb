package service

import (
	"acmweb/app/authorization/model"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/config"
	"io/ioutil"
	"path/filepath"
)

var (
	AuthApp = new(AuthAppService)
)

type AuthAppService struct {
}

// CreateApp 创建应用
func (c *AuthAppService) CreateApp(m model.App) (*common.ZIResult, error) {
	//定义应用的属性
	zs := common.NewResult()
	var err error
	if bExist, err := m.ExistsApp(); err != nil {
		system.Log.Error("Query app information failed,err:", err.Error())
		return zs, &errors.ErrDBQry
	} else {
		if bExist {
			return zs, &errors.ErrNameExists
		}
	}

	m.AppId = system.Common.UUID.NextVal()
	if zs.Count, err = m.InsertApp(); err != nil {
		system.Log.Error("Insert app information failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	return zs, nil
}

func (c *AuthAppService) RenameApp(m model.App) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if bExist, err := m.ExistsApp(); err != nil {
		system.Log.Error("Query app information failed,err:", err.Error())
		return zs, &errors.ErrDBQry
	} else {
		if bExist {
			return zs, &errors.ErrNameExists
		}
	}
	if zs.Count, err = m.UpdateApp(); err != nil {
		system.Log.Error("Rename app failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *AuthAppService) DeleteApp(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.App{AppId: id}
		rows, err := m.DeleteApp()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}
func (c *AuthAppService) ModifyApp(m model.App) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.ModifyApp(); err != nil {
		system.Log.Error("Modify app failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}

	return zs, nil
}

func (c *AuthAppService) ChangeModule(m model.AppModule) (*common.ZIResult, error) {
	//定义模块列表
	zs := common.NewResult()
	//插入
	var err error
	if m.ModuleId == 0 {
		m.ModuleId = system.Common.UUID.NextVal()
		if zs.Count, err = m.InsertAppModule(); err != nil {
			system.Log.Error("Insert app information failed,err:", err.Error())
			return zs, &errors.ErrDBexecute
		}
	} else {
		if zs.Count, err = m.UpdateAppModule(); err != nil {
			system.Log.Error("Update app information failed,err:", err.Error())
			return zs, &errors.ErrDBexecute
		}
	}

	return zs, nil
}

func (c *AuthAppService) List(search model.AppSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.App
	start := (search.PageNum - 1) * search.PageSize
	//是否有分页
	var apps = make([]model.AppRst, 0)
	var err error
	apps, err = m.FindByIndex(start, search.PageSize, search.AppData.AppName)

	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range apps {
		zs.Data = append(zs.Data, v)
	}
	zs.Count, _ = m.FindTotalNum(search.AppData.AppName)
	return zs, nil
}

func (c *AuthAppService) ListModules(appId int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	m := model.AppModule{AppId: appId}
	modules, err := m.QueryModuels()
	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range modules {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}

func (c *AuthAppService) DeleteModule(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.AppModule{ModuleId: id}
		rows, err := m.DeleteModules()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}

func (c *AuthAppService) ListAppType() (*common.ZIResult, error) {
	zs := common.NewResult()
	sql := "select type_id,type_code,type_name from app_type "
	ds, err := system.Data.DB.QueryResultRows(sql)
	if err != nil {
		return nil, err
	}

	for _, v := range ds {
		m := model.AppType{
			TypeId:   system.Text.Conv.Int(v["type_id"]),
			TypeCode: v["type_code"].(string),
			TypeName: v["type_name"].(string),
		}
		zs.Data = append(zs.Data, m)
		zs.Count++
	}
	return zs, nil
}
func (c *AuthAppService) ListAppGrp() (*common.ZIResult, error) {
	zs := common.NewResult()
	sql := "select grp_id,grp_name,remark from app_group "
	ds, err := system.Data.DB.QueryResultRows(sql)
	if err != nil {
		return nil, err
	}

	for _, v := range ds {
		m := model.AppGroup{
			GroupId:   system.Text.Conv.Int64(v["grp_id"]),
			GroupName: v["grp_name"].(string),
		}
		zs.Data = append(zs.Data, m)
		zs.Count++
	}
	return zs, nil
}

func (c *AuthAppService) writeAppFile(filename string, byteFile []byte) (string, error) {
	addPath := "/app/"
	fullFilePath := config.CONFIG.Attachment.FilePath + addPath
	newFileName := addPath + filename
	if !system.Text.File.Exists(fullFilePath) {
		system.Text.File.Mkdir(fullFilePath)
	}
	return newFileName, ioutil.WriteFile(fullFilePath+filename, byteFile, 0777)
}
func (c *AuthAppService) DownloadApp(appId int64) (srcFile, destFile string, err error) {
	m := model.App{AppId: appId}
	appfile, err := m.FindAppFile()
	if err != nil {
		return "", "", &errors.ErrDBexecute
	}
	srcFile = config.CONFIG.Attachment.FilePath + appfile
	destFile = filepath.Base(srcFile)
	return srcFile, destFile, err
}
