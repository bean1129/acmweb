package model

import (
	"acmweb/system"
	"acmweb/system/text"
	"fmt"
	"github.com/pkg/errors"
)

type Unit struct {
	UnitId       int64  `mapstructure:"unit_id" json:"unit_id"`
	UnitCode     string `mapstructure:"unit_code" json:"unit_code"`
	UnitName     string `mapstructure:"unit_name" json:"unit_name"`
	UnitFullName string `mapstructure:"unit_fullname" json:"unit_fullname"`
	ParentUnitId int64  `mapstructure:"p_unit_id" json:"p_unit_id"`
	TierId       int64  `mapstructure:"tier_id" json:"tier_id"`
	State        int    `mapstructure:"state" json:"state"`
	Principal    string `mapstructure:"principal" json:"principal"`
}

type UnitSearch struct {
	Unit     int64  `mapstructure:"unit" json:"unit"`
	State    int64  `mapstructure:"state" json:"state"`
	TireName string `mapstructure:"tire" json:"tire"`
	PUnit    string `mapstructure:"parent_unit" json:"parent_unit"`
	UnitName string `mapstructure:"unit_name" json:"unit_name"`
}

type UnitRsp struct {
	UnitId       int64  `mapstructure:"unit_id" json:"unit_id"`
	UnitCode     string `mapstructure:"unit_code" json:"unit_code"`
	UnitName     string `mapstructure:"unit_name" json:"unit_name"`
	UnitFullName string `mapstructure:"unit_fullname" json:"unit_fullname"`
	PunitCode    string `mapstructure:"p_unit_code" json:"p_unit_code"`
	PunitName    string `mapstructure:"p_unit_name" json:"p_unit_name"`
	PunitId      int64  `mapstructure:"p_unit_id" json:"p_unit_id"`
	TierName     string `mapstructure:"tier_name" json:"tier_name"`
	State        int    `mapstructure:"state" json:"state"`
	Principal    string `mapstructure:"principal" json:"principal"`
}

type UnitTier struct {
	TierId    int    `json:"tier_id"`
	TierName  string `json:"tier_name"`
	TierLevel int    `json:"tier_level"`
	Remark    string `json:"remark"`
}

func (u *Unit) Exists() (bool, error) {
	bExists := false
	sql := "select 1 from unit where unit_name =?"
	ds, err := system.Data.DB.QueryRows(sql, u.UnitName)
	if err != nil {
		return false, err
	}
	if len(ds) > 0 {
		bExists = true
	}
	return bExists, nil
}

func (u *Unit) SetUnitLevel() error {
	sql := "SELECT unit_code,ifnull(unit_fullname,unit_name) unit_name from unit  where  unit_id= ? union ALL select  unit_code,unit_name from unit  where parent_unit_id = ?"
	ds, err := system.Data.DB.QueryRows(sql, u.ParentUnitId, u.ParentUnitId)
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return errors.New("lack parent unit")
	} else if len(ds) == 1 {
		u.UnitCode = ds[0]["unit_code"] + "01"
		u.UnitFullName = ds[0]["unit_name"] + "/" + u.UnitName
	} else {
		u.UnitCode = ds[0]["unit_code"] + fmt.Sprintf("%02d", len(ds))
		u.UnitFullName = ds[0]["unit_name"] + "/" + u.UnitName
	}
	return nil
}

func (u *Unit) Insert() (int64, error) {
	sql := "insert into unit(unit_id,unit_code,unit_name,unit_fullname,parent_unit_id,tier_id,state,principal)" +
		" values(?,?,?,?,?,?,?,?)"
	rc, err := system.Data.DB.Exec(sql, u.UnitId, u.UnitCode, u.UnitName, u.UnitFullName, u.ParentUnitId, u.TierId, u.State, u.Principal)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (u *Unit) Delete() (int64, error) {
	sql := "delete from  unit where unit_id = ?"
	rc, err := system.Data.DB.Exec(sql, u.UnitId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (u *Unit) Update(fields string, values ...any) (int64, error) {

	sql := "update unit set " + fields + " where unit_id = ?"
	values = append(values, u.UnitId)
	rc, err := system.Data.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (u *Unit) UpdateUnitName() (int64, error) {
	tx, _ := system.Data.DB.StartTrans()
	sqlUpt := "update unit a,unit b set a.unit_fullname=REPLACE(a.unit_fullname,concat(\"/\",b.unit_name,\"/\"),concat(\"/\",?,\"/\")) where b.unit_id=?"
	rc, err := system.Data.DB.ExecTx(tx, sqlUpt, u.UnitName, u.UnitId)
	if err != nil {
		system.Data.DB.Rollback(tx)
		return 0, err
	}
	sqlUpt = "update unit set unit_name = ? where unit_id = ?"
	rc, err = system.Data.DB.ExecTx(tx, sqlUpt, u.UnitName, u.UnitId)
	if err != nil {
		system.Data.DB.Rollback(tx)
		return 0, err
	}
	system.Data.DB.CommitTrans(tx)
	return rc, nil
}

func (u *Unit) Find(cond string, args ...any) ([]UnitRsp, error) {
	unitRet := make([]UnitRsp, 0)
	sql := "SELECT a.unit_id,a.unit_name,a.unit_code,a.unit_fullname,b.tier_name,a.principal,a.state,c.unit_name as p_unit_name,c.unit_code as p_unit_code,c.unit_id as p_unit_id " +
		"from unit a left join unit_tier b on a.tier_id = b.tier_id left join unit c on a.parent_unit_id = c.unit_id where 1=1 "
	if cond != "" {
		sql += cond
	}
	ds, err := system.Data.DB.QueryRows(sql, args...)
	if err != nil {
		return unitRet, err
	}
	for _, rs := range ds {
		var us UnitRsp
		us.mapToStruct(rs)
		unitRet = append(unitRet, us)
	}
	return unitRet, nil
}

func (us *UnitRsp) mapToStruct(rs map[string]string) {
	us.UnitId = text.Conv.Int64(rs["unit_id"])
	us.UnitCode = rs["unit_code"]
	us.UnitName = rs["unit_name"]
	us.UnitFullName = rs["unit_fullname"]
	us.PunitCode = rs["p_unit_code"]
	us.PunitName = rs["p_unit_name"]
	us.TierName = rs["tier_name"]
	us.State = text.Conv.Int(rs["state"])
	us.Principal = rs["principal"]
	us.PunitId = text.Conv.Int64(rs["p_unit_id"])
}
