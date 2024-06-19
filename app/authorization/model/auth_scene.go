package model

import (
	"acmweb/system"
	"acmweb/system/text"
	sql "database/sql"
)

type Scene struct {
	InstId     int64   `mapstructure:"inst_id" json:"inst_id,omitempty"`
	InstName   string  `mapstructure:"inst_name" json:"inst_name,omitempty"`
	UnitId     int64   `mapstructure:"unit_id" json:"unit_id,omitempty"`
	GrpId      int64   `mapstructure:"grp_id" json:"grp_id,omitempty"`
	AuthMode   int     `mapstructure:"auth_mode" json:"auth_mode,omitempty"`
	TrialDays  int     `mapstructure:"trial_days" json:"trial_days,omitempty"`
	WaterMark  string  `mapstructure:"watermark" json:"watermark"`
	App        []int64 `mapstructure:"app" json:"app,omitempty"`
	State      int     `json:"state,omitempty"`
	CreateTime string  `mapstructure:"create_time" json:"create_time"`
	LastedTime string  `mapstructure:"lasted_time" json:"lasted_time"`
	IsVisable  int     `mapstructure:"is_visable" json:"is_visable"`
}

type SceneShow struct {
	Scene
	AppName  string `json:"app_name"`
	UnitName string `json:"unit_name"`
	GrpName  string `json:"grp_name"`
	TypeName string `json:"type_name"`
	TypeCode string `json:"type_code"`
}

type SceneApp struct {
	AppCode string `json:"app_code"`
	AppName string `json:"app_name"`
	AppDesc string `json:"app_desc"`
}

func (s *Scene) Insert() error {
	var tx *sql.Tx
	var err error
	if tx, err = system.Data.DB.StartTrans(); err != nil {
		return err
	}
	if 0 == s.InstId {
		s.InstId = system.Common.UUID.NextVal()
	}
	for {
		sqlQry := "insert into scene(inst_id,inst_name,unit_id,grp_id,auth_mode,trial_days,watermark,create_time,lasted_time,state)" +
			" values(?,?,?,?,?,?,?,now(),NOW(),1)" +
			" ON DUPLICATE KEY UPDATE inst_name=values(inst_name)" +
			",unit_id=values(unit_id),grp_id=values(grp_id),auth_mode=values(auth_mode),trial_days=values(trial_days),watermark=values(watermark),lasted_time=values(lasted_time)"
		_, err := system.Data.DB.ExecTx(tx, sqlQry, s.InstId, s.InstName, s.UnitId, s.GrpId, s.AuthMode, s.TrialDays, s.WaterMark)
		if err != nil {
			break
		}

		delAppSql := "delete from scene_app where inst_id=?"
		_, err = system.Data.DB.ExecTx(tx, delAppSql, s.InstId)
		if err != nil {
			break
		}
		sqlQry = ""
		var args []any
		for idx, value := range s.App {
			if idx == 0 {
				sqlQry = "insert into scene_app values(?,?,?)"
			} else {
				sqlQry += ",(?,?,?)"
			}
			args = append(args, system.Common.UUID.NextVal(), s.InstId, value)
		}
		if len(sqlQry) > 0 {
			_, err = system.Data.DB.ExecTx(tx, sqlQry, args...)
		}
		break
	}
	if err == nil {
		if !system.Data.DB.CommitTrans(tx) {
			system.Log.Errorf("Commit trans failed")
		}
	} else {
		system.Log.Errorf("Insert failed,err:%s", err.Error())
		system.Data.DB.Rollback(tx)
	}
	return err
}

func (s *Scene) Delete() (int64, error) {
	var tx *sql.Tx
	var err error
	var rc int64
	if tx, err = system.Data.DB.StartTrans(); err != nil {
		return 0, err
	}
	for {
		sqlQry := "delete from  scene where inst_id = ?"
		rc, err = system.Data.DB.ExecTx(tx, sqlQry, s.InstId)
		if err != nil {
			break
		}
		sqlQry = "delete from scene_app where inst_id = ?"
		_, err = system.Data.DB.ExecTx(tx, sqlQry, s.InstId)
		break
	}
	if err == nil {
		system.Data.DB.CommitTrans(tx)
	} else {
		system.Log.Errorf("Delete scene faield,err:", err.Error())
		system.Data.DB.Rollback(tx)
	}
	return rc, err
}

func (s *Scene) ExistsScene() (bFind bool, err error) {
	sqlQry := "select 1 from scene where inst_name = ?"
	rc, err := system.Data.DB.QueryRows(sqlQry, s.InstName)
	if err != nil {
		return false, err
	}
	if len(rc) > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Scene) ListApp() (result []SceneApp, err error) {
	sqlQry := "SELECT b.app_code,b.app_name,b.app_desc from scene_app a,app b where a.app_id = b.app_id and a.inst_id =?"
	ds, err := system.Data.DB.QueryRows(sqlQry, s.InstId)
	if err != nil {
		return
	}
	for _, rc := range ds {
		result = append(result, SceneApp{
			rc["app_code"],
			rc["app_name"],
			rc["app_desc"],
		})
	}

	return result, nil
}

// 场景信息
func (s *Scene) GetSceneInfo() error {
	sqlQry := "select * from scene where inst_id = ? "
	ds, err := system.Data.DB.QueryRows(sqlQry, s.InstId)
	if err != nil {
		return err
	}
	for _, rs := range ds {
		if 1 == s.InstId {
			s.GrpId = system.Text.Conv.Int64(rs["grp_id"])
		} else {
			s.UnitId = system.Text.Conv.Int64(rs["unit_id"])
			s.AuthMode = system.Text.Conv.Int(rs["auth_mode"])
			s.InstId = system.Text.Conv.Int64(rs["inst_id"])
			s.GrpId = system.Text.Conv.Int64(rs["grp_id"])
			s.WaterMark = rs["watermark"]
			s.TrialDays = system.Text.Conv.Int(rs["trial_days"])
		}

	}
	return nil
}

// 查询流应用所属组
func (s *Scene) GetSceneGrp() int64 {
	return s.GrpId
}

type SceneSearch struct {
	PageNum   int    `mapstructure:"page_num"`
	PageSize  int    `mapstructure:"page_size"`
	State     int    `mapstructure:"state"`
	SearchBox string `mapstructure:"search_box"`
}

func NewScenseSearch() *SceneSearch {
	return &SceneSearch{
		PageNum:   0,
		PageSize:  0,
		State:     -1,
		SearchBox: "",
	}
}

func (o *SceneSearch) FindTotalNum() (int64, error) {
	sqlQry := "select count(*) total from (SELECT a.inst_id from scene a left join scene_group b on a.grp_id = b.grp_id " +
		"LEFT JOIN scene_app c on a.inst_id = c.inst_id " +
		"LEFT join app d on c.app_id = d.app_id " +
		"LEFT join unit e on a.unit_id=e.unit_id where a.is_visable=1"
	var args []any
	if o.SearchBox != "" {
		sqlQry += "and a.inst_name=? "
		args = append(args, o.SearchBox)
	}
	if o.State >= 0 {
		sqlQry += "and a.state = ?"
		args = append(args, o.State)
	}
	sqlQry += " GROUP BY a.inst_id) t"
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return 0, err
	}
	result := system.Text.Conv.Int64(ds[0]["total"])
	return result, nil
}

func (o *SceneSearch) FindByIndex() ([]SceneShow, error) {
	start := (o.PageNum - 1) * o.PageSize
	var result []SceneShow
	sqlQry := "SELECT a.inst_id,a.inst_name,a.unit_id,a.grp_id,a.auth_mode,a.trial_days,a.watermark,a.state," +
		"DATE_FORMAT(create_time,'%Y-%m-%d %H:%s:%i') create_time,DATE_FORMAT(lasted_time,'%Y-%m-%d %H:%s:%i') lasted_time" +
		",grp_name,group_concat(d.app_name) app_name,group_concat(d.app_id) app_id,f.type_name,e.unit_name from scene a left join scene_group b on a.grp_id = b.grp_id " +
		"LEFT JOIN scene_app c on a.inst_id = c.inst_id " +
		"LEFT join app d on c.app_id = d.app_id " +
		"LEFT join unit e on a.unit_id=e.unit_id " +
		"LEFT join auth_mode f on a.auth_mode = f.type_id where a.is_visable=1 "
	var args []any
	if o.SearchBox != "" {
		sqlQry += "and a.inst_name like ? "
		args = append(args, "%"+o.SearchBox+"%")
	}
	if o.State >= 0 {
		sqlQry += "and a.state = ? "
		args = append(args, o.State)
	}
	sqlQry += " group by a.inst_id order by a.inst_id "
	if start > 0 {
		sqlQry += " limit " + text.Conv.String(start) + "," + text.Conv.String(o.PageSize)
	}
	ds, err := system.Data.DB.QueryRows(sqlQry, args...)
	if err != nil {
		return result, err
	}
	for _, rs := range ds {
		var ret SceneShow
		ret.InstId = system.Text.Conv.Int64(rs["inst_id"])
		ret.AppName = rs["app_name"]
		ret.InstName = rs["inst_name"]
		ret.State = system.Text.Conv.Int(rs["state"])
		ret.WaterMark = rs["watermark"]
		ret.UnitName = rs["unit_name"]
		ret.UnitId = system.Text.Conv.Int64(rs["unit_id"])
		ret.GrpName = rs["grp_name"]
		ret.GrpId = system.Text.Conv.Int64(rs["grp_id"])
		ret.TypeName = rs["type_name"]
		ret.AuthMode = system.Text.Conv.Int(rs["auth_mode"])
		ret.TrialDays = system.Text.Conv.Int(rs["trial_days"])
		ret.CreateTime = rs["create_time"]
		ret.LastedTime = rs["lasted_time"]
		appidarry := system.Text.String.Split(rs["app_id"], ",")
		for _, appid := range appidarry {
			ret.App = append(ret.App, system.Text.Conv.Int64(appid))
		}
		result = append(result, ret)
	}
	return result, nil
}

type SceneGroup struct {
	GroupId   int64  `mapstructure:"grp_id" json:"grp_id,omitempty"`
	GrpName   string `mapstructure:"grp_name" json:"grp_name,omitempty"`
	Remark    string `mapstructure:"remark" json:"remark"`
	IsVisable int    `mapstructure:"is_visable" json:"is_visable"`
}

func (g *SceneGroup) Insert() error {
	if 0 == g.GroupId {
		g.GroupId = system.Common.UUID.NextVal()
	}

	sqlInsert := "Insert into scene_group values(?,?,?) ON DUPLICATE KEY UPDATE grp_name=values(grp_name),remark=values(remark)"
	_, err := system.Data.DB.Exec(sqlInsert, g.GroupId, g.GrpName, g.Remark)
	if err != nil {
		return err
	}
	return nil
}
func (g *SceneGroup) Delete() error {
	sqlQry := "delete from scene_group where grp_id = ?"
	_, err := system.Data.DB.Exec(sqlQry, g.GroupId)
	if err != nil {
		return err
	}
	return nil
}
func (g *SceneGroup) List() (result []SceneGroup, err error) {
	sqlQry := "select * from scene_group where is_visable = 1"
	ds, err := system.Data.DB.QueryRows(sqlQry)
	if err != nil {
		return
	}
	for _, rs := range ds {
		grp := SceneGroup{
			GroupId: system.Text.Conv.Int64(rs["grp_id"]),
			GrpName: rs["grp_name"],
			Remark:  rs["remark"],
		}
		result = append(result, grp)
	}
	return
}
