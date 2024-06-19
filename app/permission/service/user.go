package service

import (
	"acmweb/app/errors"
	"acmweb/app/permission/model"
	"acmweb/system"
	"acmweb/system/common"
)

var User = new(UserSerivce)

type UserSerivce struct {
}

func (c *UserSerivce) CreateUser(u model.User) (*common.ZIResult, error) {
	u.UserId = system.Common.UUID.NextVal()
	zs := common.NewResult()
	var err error
	if bExists, err := u.ExistsUser(); err != nil {
		system.Log.Error("Query user failed ,err:", err.Error())
		return zs, &errors.ErrDBQry
	} else {
		if bExists {
			return zs, &errors.ErrUserExists
		}
	}

	zs.Count, err = u.InsertUser()
	if err != nil {
		system.Log.Error("Insert user failed ,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) DeleteUser(idsArry []int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	for _, id := range idsArry {
		m := &model.User{UserId: id}
		rows, err := m.DeleteUser()
		if rows == 0 || err != nil {
			continue
		}
		zs.Count++
	}
	return zs, nil
}

func (c *UserSerivce) ModifyUser(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.Update(""); err != nil {
		system.Log.Error("Update user failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) UpdateRole(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.Update("role_id=?", m.RuleId); err != nil {
		system.Log.Error("Update user failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) SetPassword(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	//查询user_code
	passwd := m.Passwd
	if users, err := m.Find(" and user_id = ?", m.UserId); err != nil {
		return zs, &errors.ErrDBexecute
	} else {
		if len(users) == 0 {
			return zs, &errors.ErrDBQry
		}
		passwd, _ = system.Secret.MD5.Password(passwd + system.Config.Application.SecretKey + users[0].UserCode)
		if m.Passwd == passwd {
			return zs, &errors.ErrSameUserPwd
		}
	}
	if zs.Count, err = m.Update(" passwd=?", passwd); err != nil {
		system.Log.Error("Update user password failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) ChangeRole(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	//判断角色是否存在
	//更新角色
	if zs.Count, err = m.Update("role_id=?", m.RuleId); err != nil {
		system.Log.Error("Update role failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) ChangeState(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.Update("state=?", m.State); err != nil {
		system.Log.Error("Update state failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) Rename(m model.User) (*common.ZIResult, error) {
	var err error
	zs := common.NewResult()
	if zs.Count, err = m.Update("user_name=?", m.UserName); err != nil {
		system.Log.Error("Update user name failed,err:", err.Error())
		return zs, &errors.ErrDBexecute
	}
	return zs, nil
}

func (c *UserSerivce) List(search *model.UserSearch) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.User
	var args = make([]any, 0)
	var cond string
	start := (search.PageNum - 1) * search.PageSize
	//根据用户查询
	if len(search.User) > 0 {
		if search.Vague {
			cond = " and (user_code like ? or user_name  like ?) "
			args = append(args, "%"+search.User+"%")
			args = append(args, "%"+search.User+"%")
		} else {
			cond = " and (user_code = ? or user_name  = ?) "
			args = append(args, search.User)
			args = append(args, search.User)
		}

	} else {
		cond = " and 1=1 "
	}
	//根据角色查询
	if search.RoleId > 0 {
		cond += " and a.role_id=? "
		args = append(args, search.RoleId)
	} else {
		cond += " and 1=1 "
	}
	//状态
	if search.State != -1 {
		cond += "and a.state = ? "
		args = append(args, search.State)
	} else {
		cond += " and (a.state = ? or a.state = ?) "
		args = append(args, 0, 1)

	}
	//是否有分页
	var users = make([]model.User, 0)
	var err error
	if search.PageSize > 0 {
		users, err = m.FindByIndex(start, search.PageSize, cond, args...)
	} else {
		users, err = m.Find(cond, args...)

	}
	if err != nil {
		return zs, &errors.ErrDBexecute
	}

	for _, v := range users {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}

func (c *UserSerivce) GetUserById(id int64) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.User
	users, err := m.Find("user_id=?", id)
	if err != nil {
		return zs, &errors.ErrDBexecute
	}
	for _, v := range users {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}
func (c *UserSerivce) GetUserByCodeOrName(code, name string) (*common.ZIResult, error) {
	zs := common.NewResult()
	var m model.User
	users, err := m.Find("user_code=? or user_name = ?", code, name)
	if err != nil {
		return zs, &errors.ErrDBexecute
	}
	for _, v := range users {
		zs.Data = append(zs.Data, v)
		zs.Count++
	}
	return zs, nil
}
