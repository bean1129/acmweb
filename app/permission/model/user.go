package model

import (
	"acmweb/system"
	"acmweb/system/text"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type UserSearch struct {
	User     string `mapstructure:"user"`
	RoleId   int64  `mapstructure:"role_id"`
	State    int64  `mapstructure:"state"`
	Vague    bool   `mapstructure:"vague"`
	PageNum  int    `mapstructure:"page_num"`
	PageSize int    `mapstructure:"page_size"`
}

// User 用户
type User struct {
	UserId          int64  `mapstructure:"user_id" json:"user_id"`
	UserCode        string `mapstructure:"user_code" json:"user_code"`
	UserName        string `mapstructure:"user_name" json:"user_name"`
	RuleId          int64  `mapstructure:"role_id" json:"role_id"`
	RuleName        string `mapstructure:"role_name" json:"role_name"`
	Passwd          string `mapstructure:"password" json:"passwd"`
	State           int    `mapstructure:"state" json:"state"`
	Phone           string `mapstructure:"phone" json:"phone"`
	Address         string `mapstructure:"address" json:"address"`
	LastedIp        string `mapstructure:"lasted_ip" json:"lasted_ip"`
	CreatTime       string `mapstructure:"create_time" json:"create_time"`
	LastedLoginTime string `mapstructure:"lasted_login_time" json:"lasted_login_time"`
	LastedTime      string `mapstructure:"lasted_time" json:"lasted_time"`
	Remark          string `mapstructure:"remark" json:"remark"`
}

func (u *User) ExistsUser() (bool, error) {
	bExists := false
	sql := "select 1 from user where user_code =?"
	ds, err := system.Data.DB.QueryRows(sql, u.UserCode)
	if err != nil {
		return false, err
	}
	if len(ds) > 0 {
		bExists = true
	}
	return bExists, nil
}

func (u *User) InsertUser() (int64, error) {
	sql := "insert into user(user_id,user_code,user_name,role_id,passwd,state,phone,address,create_time,remark)" +
		" values(?,?,?,?,?,?,?,?,now(),?)"
	rc, err := system.Data.DB.Exec(sql, u.UserId, u.UserCode, u.UserName, u.RuleId, u.Passwd, u.State, u.Phone, u.Address, u.Remark)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (u *User) DeleteUser() (int64, error) {
	sql := "delete from  user where user_id = ?"
	rc, err := system.Data.DB.Exec(sql, u.UserId)
	if err != nil {
		return 0, err
	}
	return rc, nil
}

func (u *User) Update(fields string, values ...any) (int64, error) {
	if len(fields) == 0 {
		sql := "update user set user_code=?,user_name=?,address=?,phone=?,lasted_time=now(),remark=? where user_id = ?"
		rc, err := system.Data.DB.Exec(sql, u.UserCode, u.UserName, u.Address, u.Phone, u.Remark, u.UserId)
		if err != nil {
			return 0, err
		}
		return rc, nil
	} else {
		sql := "update user set " + fields + " where user_id = ?"
		values = append(values, u.UserId)
		rc, err := system.Data.DB.Exec(sql, values...)
		if err != nil {
			return 0, err
		}
		return rc, nil
	}
}

func (u *User) Find(cond string, args ...any) ([]User, error) {
	users := make([]User, 0)
	sql := "select user_id,user_code,user_name,a.role_id,a.state,phone,address,lasted_ip,date_format(create_time,'%Y-%m-%d %H:%i:%s') create_time" +
		",date_format(lasted_login_time,'%Y-%m-%d %H:%i:%s') lasted_login_time,date_format(lasted_time,'%Y-%m-%d %H:%i:%s') lasted_time,passwd,b.role_name,a.remark from user a,user_role b where a.role_id=b.role_id"
	if cond != "" {
		sql += cond
	}
	ds, err := system.Data.DB.QueryResultRows(sql, args...)
	if err != nil {
		return users, err
	}
	for _, rs := range ds {
		var result User
		if err := mapstructure.Decode(rs, &result); err != nil {
			continue
		}
		users = append(users, result)
	}
	return users, nil
}

func (u *User) FindByIndex(start, size int, cond string, args ...any) ([]User, error) {
	users := make([]User, 0)
	sql := "select a.user_id,a.user_code,user_name,a.role_id,a.state,phone,address,lasted_ip,lasted_login_time,date_format(create_time,'%Y-%m-%d %H:%i:%s') create_time" +
		",date_format(lasted_login_time,'%Y-%m-%d %H:%i:%s') lasted_login_time,date_format(lasted_time,'%Y-%m-%d %H:%i:%s') lasted_time,b.role_name,a.remark from user a,user_role b where a.role_id = b.role_id "
	if cond != "" {
		sql += cond
	}
	sql += " order by user_id limit " + text.Conv.String(start) + "," + text.Conv.String(size)
	ds, err := system.Data.DB.QueryResultRows(sql, args...)
	if err != nil {
		return users, err
	}
	for _, rs := range ds {
		var result User
		if err := mapstructure.Decode(rs, &result); err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, result)
	}
	return users, nil
}
