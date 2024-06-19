package service

import (
	"acmweb/app/permission/model"
	"acmweb/constants"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/text"
	"github.com/mitchellh/mapstructure"
)

var (
	Role = new(RoleService)
)

type RoleService struct {
}

// CreateRole 创建角色
func (c *RoleService) CreateRole(m model.Role) *common.ZIResult {
	zs := c.GetByName(m.Name)
	if zs.Code != common.OK {
		return zs
	} else if len(zs.Data) > 0 {
		zs.Code = constants.ErrNameExists
		zs.Msg = "Name repeat"
		zs.Data = nil
		return zs
	}
	tx, _ := system.Data.DB.StartTrans()
	sql := "insert into user_role(role_id,role_name,parent_role_id,remark,manage_user_id) values(?,?,?,?,?)"
	rc, err := system.Data.DB.ExecTx(tx, sql, m.Id, m.Name, m.PId, m.Remark, m.MId)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	//绑定权限菜单
	sql = "insert into role_func (role_id,func_id,state) select ?,func_id,0 from func"
	rc, err = system.Data.DB.ExecTx(tx, sql, m.Id)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		system.Data.DB.Rollback(tx)
	}
	system.Data.DB.CommitTrans(tx)
	return zs
}

// DeleteRole 删除角色
func (c *RoleService) DeleteRole(id int64) *common.ZIResult {
	zs := common.NewResult()
	sql := "delete from user_role where role_id = ?"
	rc, err := system.Data.DB.Exec(sql, id)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	return zs
}

// UpdateRole 更新角色信息
func (c *RoleService) UpdateRole(id int64, state int) *common.ZIResult {
	zs := common.NewResult()
	sql := "update user_role set state=? where role_id=?"
	rc, err := system.Data.DB.Exec(sql, state, id)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	return zs
}

// Rename 重命名角色
func (c *RoleService) Rename(id int64, name string) *common.ZIResult {
	zs := common.NewResult()
	sql := "update user_role set role_name=? where role_id=?"
	rc, err := system.Data.DB.Exec(sql, name, id)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	return zs
}

// SetPermission 设置权限
func (c *RoleService) SetPermission(idRole, idManage, idParent int64) *common.ZIResult {
	zs := common.NewResult()
	sql := "update user_role set manage_user_id=?,parent_role_id=? where role_id=?"
	rc, err := system.Data.DB.Exec(sql, idManage, idParent, idRole)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	return zs
}

/*func (c *RoleService) SetRoleFunc(roleFunc model.RoleFuncReq) *common.ZIResult {
	zs := common.NewResult()
	tx, _ := system.Data.DB.StartTrans()
	sql := "delete from role_func  where role_id=?"
	rc, err := system.Data.DB.ExecTx(tx, sql, roleFunc.RoleId)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		system.Data.DB.Rollback(tx)
		return zs
	} else {
		zs.Count = rc
	}
	//插入
	var fieldVal []any
	first := true
	sqlInsert := "insert into role_func (sys_id,role_id,func_id,state) values"
	for _, grpVal := range roleFunc.RoleFuncDetail {
		for _, funcVal := range grpVal.Group.Func {
			if first {
				sqlInsert += "(?,?,?,1)"
				first = false
			} else {
				sqlInsert += ",(?,?,?,1)"
			}
			fieldVal = append(fieldVal, system.Common.UUID.NextVal(), roleFunc.RoleId, funcVal.ID)
		}
	}
	rc, err = system.Data.DB.ExecTx(tx, sqlInsert, fieldVal...)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
	} else {
		zs.Count = rc
	}
	system.Data.DB.CommitTrans(tx)
	return zs
}*/

func (c *RoleService) SetRoleFunc(roleFunc model.RoleFuncReq) *common.ZIResult {
	zs := common.NewResult()
	tx, _ := system.Data.DB.StartTrans()
	for _, v := range roleFunc.RoleFuncDetail {
		for _, f := range v.Func {
			sql := "update role_func  set state = ? where role_id=? and func_id = ?"
			rc, err := system.Data.DB.ExecTx(tx, sql, f.State, roleFunc.RoleId, f.ID)
			if err != nil {
				zs.Code = constants.ErrDBExecute
				zs.Msg = err.Error()
				system.Data.DB.Rollback(tx)
				return zs
			} else {
				zs.Count += rc
			}
		}
	}
	system.Data.DB.CommitTrans(tx)
	return zs
}

// ListAll 获取角色列表
func (c *RoleService) ListAll() *common.ZIResult {
	return c.Find("", make([]any, 0))
}

// Find 获取角色列表
func (c *RoleService) Find(cond string, condVal []any) *common.ZIResult {
	zs := common.NewResult()
	sql := "select a.role_id,a.role_name,a.parent_role_id,b.role_name parent_role_name,a.state,a.remark,a.manage_user_id,c.user_name from user_role a " +
		"LEFT JOIN user_role b on a.parent_role_id = b.role_id LEFT  JOIN user c  on b.manage_user_id = c.user_id where 1=1  "
	if cond != "" {
		sql += " and " + cond
	}
	ds, err := system.Data.DB.QueryRows(sql, condVal...)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		return zs
	}
	if ds == nil || len(ds) == 0 {
		return zs
	}

	for _, rs := range ds {
		zs.Data = append(zs.Data, model.Role{
			Id:     text.Conv.Int64(rs["role_id"]),
			Name:   rs["role_name"],
			PId:    text.Conv.Int64(rs["parent_role_id"]),
			PName:  rs["parent_role_name"],
			MId:    text.Conv.Int64(rs["manage_user_id"]),
			MName:  rs["user_name"],
			State:  text.Conv.Int(rs["state"]),
			Remark: rs["remark"],
		})
	}
	return zs
}
func (c *RoleService) ListFunc(cond string) *common.ZIResult {
	zs := common.NewResult()
	sql := "select a.func_id,a.func_name,a.grp_id,b.grp_name from func a,func_group b where a.grp_id = b.grp_id and 1 =1  "
	if cond != "" {
		sql += " and " + cond
	}
	sql += " order by b.sort_id "
	ds, err := system.Data.DB.QueryRows(sql)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		return zs
	}
	if ds == nil || len(ds) == 0 {
		return zs
	}

	for _, rs := range ds {
		zs.Data = append(zs.Data, model.FuncDetail{
			ID:   text.Conv.Int(rs["func_id"]),
			Name: rs["func_name"],
		})
	}
	return zs
}
func (c *RoleService) GetById(id int64) *common.ZIResult {
	zs := common.NewResult()
	sql := "select role_id,role_name,parent_role_id,state,remark,manage_user_id from user_role where role_id=?"
	ds, err := system.Data.DB.QueryRows(sql, id)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		return zs
	}
	if ds == nil || len(ds) == 0 {
		return zs
	}
	rs := ds[0]
	zs.Data[0] = &model.Role{
		Id:     text.Conv.Int64(rs["role_id"]),
		Name:   rs["role_name"],
		PId:    text.Conv.Int64(rs["parent_role_id"]),
		State:  text.Conv.Int(rs["state"]),
		Remark: rs["remark"],
	}
	return zs
}

func (c *RoleService) GetByName(name string) *common.ZIResult {
	zs := common.NewResult()
	sql := "select role_id,role_name,parent_role_id,state,remark,manage_user_id from user_role where role_name=?"
	ds, err := system.Data.DB.QueryRows(sql, name)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		return zs
	}
	if ds == nil || len(ds) == 0 {
		return zs
	}
	rs := ds[0]
	zs.Data = append(zs.Data, model.Role{
		Id:     text.Conv.Int64(rs["role_id"]),
		Name:   rs["role_name"],
		PId:    text.Conv.Int64(rs["parent_role_id"]),
		State:  text.Conv.Int(rs["state"]),
		Remark: rs["remark"],
	})
	return zs
}

func (c *RoleService) GetRoleFunc(id int64) *common.ZIResult {
	zs := common.NewResult()
	var args []any
	sql := "SELECT a.role_id,b.func_id,b.func_name,c.grp_id,c.grp_name,a.state from role_func a,func b,func_group c " +
		"where a.func_id = b.func_id and b.grp_id = c.grp_id and 1=1 "
	if id > 0 {
		sql += " and a.role_id = ?"
		args = append(args, id)
	}
	sql += " order by c.sort_id,b.func_id"
	ds, err := system.Data.DB.QueryResultRows(sql, args...)
	if err != nil {
		zs.Code = constants.ErrDBExecute
		zs.Msg = err.Error()
		return zs
	}
	if ds == nil || len(ds) == 0 {
		return zs
	}

	var retData model.RoleFuncReq
	var mSort = make([]int, 0)
	retData.RoleId = id
	var grpMap = make(map[int]model.FuncGrp, 0)
	for _, rs := range ds {
		var result model.RoleFuncResult
		if err := mapstructure.Decode(rs, &result); err != nil {
			continue
		}
		funcVal := model.FuncDetail{ID: result.FuncId, Name: result.FuncName, State: result.State}
		if v, ok := grpMap[result.GrpId]; ok {
			v.Func = append(v.Func, funcVal)
			grpMap[result.GrpId] = v
		} else {
			funcGrp := model.FuncGrp{ID: result.GrpId, Name: result.GrpName}
			funcGrp.Func = append(funcGrp.Func, funcVal)
			grpMap[result.GrpId] = funcGrp
			mSort = append(mSort, result.GrpId)
		}
	}
	//var mSort = make([]int, 0)
	//for key, _ := range grpMap {
	//	mSort = append(mSort, key)
	//}
	//sort.Ints(mSort)
	for _, v := range mSort {
		retData.RoleFuncDetail = append(retData.RoleFuncDetail, grpMap[v])
	}
	zs.Data = append(zs.Data, retData)
	zs.Count = 1
	return zs
}
