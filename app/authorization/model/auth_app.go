package model

import (
	"acmweb/system"
	"acmweb/system/text"
	"github.com/pkg/errors"
)

type AppSearch struct {
	PageNum  int `mapstructure:"page_num"`
	PageSize int `mapstructure:"page_size"`
	AppData  App `mapstructure:"app_data"`
}

type App struct {
	AppId       int64  `mapstructure:"app_id" json:"app_id"`
	AppName     string `mapstructure:"app_name" json:"app_name,omitempty"`
	ReleaseDate string `mmapstructure:"release_date" json:"release_date"`
	AppDesc     string `mapstructure:"app_desc" json:"app_desc,omitempty"`
	GroupId     int64  `mapstructure:"group_id" json:"group_id,omitempty"`
	TypeId      int    `mapstructure:"type_id" json:"type_id,omitempty"`
	AppCode     string `mapstructure:"app_code" json:"app_code"`
	Water       string `mapstructure:"water" json:"water"`
	Title       string `mapstructure:"title" json:"title"`
	About       string `mapstructure:"about" json:"about"`
	Features    string `mapstructure:"features" json:"features"`
	Details     string `mapstructure:"details" json:"details"`
	TrialDays   int    `mapstructure:"trial_days" json:"trial_days"`
	FileBytes   []byte //应用包的二进制
	FileName    string `json:"app_file"`
}

type AppModule struct {
	ModuleId   int64  `mapstructure:"module_id" json:"module_id,omitempty"`
	AppId      int64  `mapstructure:"app_id" json:"app_id,omitempty"`
	ModuleName string `mapstructure:"module_name" json:"module_name,omitempty"`
}

type AppGroup struct {
	GroupId   int64  `mapstructure:"group_id" json:"group_id"`
	GroupName string `mapstructure:"group_name" json:"group_name"`
}

type AppType struct {
	TypeId   int    `mapstructure:"type_id" json:"type_id"`
	TypeCode string `mapstructure:"type_code" json:"type_code"`
	TypeName string `mapstructure:"type_name" json:"type_name"`
}

type AppRst struct {
	App
	Type     string `mapstructure:"type" json:"type,omitempty"`
	TypeCode string `mapstructure:"type_code" json:"type_code,omitempty"`
	Grp      string `mapstructure:"group" json:"grp,omitempty"`
	//Modules  []string `mapstructure:"modules" json:"modules,omitempty"`
	Modules string `mapstructure:"modules" json:"modules"`
}

func (c *App) ExistsApp() (bool, error) {
	bExists := false
	sql := "select 1 from app where app_name=?"
	ds, err := system.Data.DB.QueryRows(sql, c.AppName)
	if err != nil {
		return false, err
	}
	if len(ds) > 0 {
		bExists = true
	}
	return bExists, nil
}

func (c *App) FindAppFile() (string, error) {
	sql := "select app_file from app where app_id=?"
	ds, err := system.Data.DB.QueryRows(sql, c.AppId)
	if err != nil {
		return "", err
	}
	if len(ds) == 0 {
		return "", errors.New("empty")
	}
	return ds[0]["app_file"], nil
}

func (c *App) InsertApp() (int64, error) {
	sql := "insert into app(app_id,app_name,app_desc,grp_id,type_id,water,title,about,app_code,details,features,trial_days,app_file)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	rc, err := system.Data.DB.Exec(sql, c.AppId, c.AppName, c.AppDesc, c.GroupId, c.TypeId, c.Water, c.Title, c.About, c.AppCode, c.Details, c.Features, c.TrialDays, c.FileName)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (c *App) UpdateApp() (int64, error) {
	sql := "update app set app_name = ? where app_id= ?"
	rc, err := system.Data.DB.Exec(sql, c.AppName, c.AppId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (c *App) DeleteApp() (int64, error) {
	sql := "delete from  app where app_id = ?"
	rc, err := system.Data.DB.Exec(sql, c.AppId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (c *App) ModifyApp() (int64, error) {
	var args []any
	sql := "update app set grp_id= ?,type_id=?,app_desc=? ,water=?,title=?,about=?,details=?,features=?,trial_days = ? "
	args = append(args, c.GroupId, c.TypeId, c.AppDesc, c.Water, c.Title, c.About, c.Details, c.Features, c.TrialDays)
	if len(c.FileName) > 0 {
		sql += ",app_file=?"
		args = append(args, c.FileName)
	}
	sql += " where app_id = ? "
	args = append(args, c.AppId)
	rc, err := system.Data.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return rc, nil
}
func (c *App) FindTotalNum(cond string) (int64, error) {
	var args []any
	sqlQry := "select count(*) total from app a where 1=1"
	if cond != "" {
		sqlQry += " and (app_name like ? or app_code like ?)"
		args = append(args, "%"+cond+"%", "%"+cond+"%")
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return 0, err
	}
	result := system.Text.Conv.Int64(ds[0]["total"])
	return result, nil
}

func (c *App) FindByIndex(start, size int, cond string) ([]AppRst, error) {
	apps := make([]AppRst, 0)
	var args []any
	sqlQry := "select app_id ,app_name,app_code,app_desc,app_file from app a where 1=1"
	if cond != "" {
		sqlQry += " and (app_name like ? or app_code like ?)"
		args = append(args, "%"+cond+"%", "%"+cond+"%")
	}
	sqlQry += " order by a.app_id "
	if size > 0 {
		sqlQry += " limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return apps, err
	}
	for _, rs := range ds {
		var appret AppRst
		appret.FileName = rs["app_file"]
		appret.AppName = rs["app_name"]
		appret.AppCode = rs["app_code"]
		appret.AppId = system.Text.Conv.Int64(rs["app_id"])
		appret.AppDesc = rs["app_desc"]
		apps = append(apps, appret)
	}
	return apps, nil
}

func (m *AppModule) InsertAppModule() (int64, error) {
	if m.ModuleId == 0 {
		m.ModuleId = system.Common.UUID.NextVal()
	}
	sql := "insert into app_module(module_id,app_id,module_name) values(?,?,?)"
	rc, err := system.Data.DB.Exec(sql, m.ModuleId, m.AppId, m.ModuleName)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (m *AppModule) UpdateAppModule() (int64, error) {
	sql := "update app_module set module_name=? where module_id = ?"
	rc, err := system.Data.DB.Exec(sql, m.ModuleName, m.ModuleId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (m *AppModule) QueryModuels() ([]AppModule, error) {
	var appModules = make([]AppModule, 0)
	sql := "select * from app_module where app_id=?"
	ds, err := system.Data.DB.QueryRows(sql, m.AppId)
	if err != nil {
		return nil, err
	}

	for _, v := range ds {
		var moduel AppModule
		moduel.AppId = m.AppId
		moduel.ModuleName = v["module_name"]
		moduel.ModuleId = system.Text.Conv.Int64(v["module_id"])
		appModules = append(appModules, moduel)
	}
	return appModules, nil
}

func (m *AppModule) DeleteModules() (int64, error) {
	sql := "delete from app_module where module_id = ? or app_id=?"
	rc, err := system.Data.DB.Exec(sql, m.ModuleId, m.AppId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}
