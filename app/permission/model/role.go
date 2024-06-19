package model

// Role 角色
type Role struct {
	Id     int64  `json:"id"`     // 角色ID
	Name   string `json:"name"`   // 角色名称
	PId    int64  `json:"pid"`    // 父级角色ID
	PName  string `json:"p_name"` // 角色名称
	MId    int64  `json:"mid"`    //管理员
	MName  string `json:"m_name"` // 角色名称
	State  int    `json:"state"`  // 状态
	Remark string `json:"remark"` // 描述
}

// 角色菜单权限
type RoleFuncReq struct {
	RoleId         int64     `mapstructure:"role_id" json:"role_id,omitempty"`
	RoleFuncDetail []FuncGrp `mapstructure:"role_func_grp" json:"role_func_grp,omitempty"`
}
type FuncGrp struct {
	ID   int          `mapstructure:"id" json:"id,omitempty"`
	Name string       `mapstructure:"name" json:"name,omitempty"`
	Func []FuncDetail `mapstructure:"func" json:"func,omitempty"`
}

type FuncDetail struct {
	ID    int    `mapstructure:"id" json:"id,omitempty"`
	Name  string `mapstructure:"name" json:"name,omitempty"`
	State int    `mapstructure:"state" json:"state"`
}

type RoleFuncResult struct {
	RoleId   int64  `mapstructure:"role_id"`
	FuncId   int    `mapstructure:"func_id"`
	FuncName string `mapstructure:"func_name"`
	GrpId    int    `mapstructure:"grp_id"`
	GrpName  string `mapstructure:"grp_name"`
	State    int    `mapstructure:"state"`
}
