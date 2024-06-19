package constants

//审批操作
const (
	WF_OPR_START  = 1 //发起审批
	WF_OPR_NEXT   = 2 //提交下一级
	WF_OPR_PRE    = 3 //返回上一级
	WF_OPR_BACK   = 4 //返回修改
	WF_OPR_END    = 5 //终止
	WF_OPR_Finish = 6 //流程结束
)

//审批状态
const (
	WF_APPROVE_TODO   = 0 //待审
	WF_APPROVE_PASS   = 1 //通过
	WF_APPROVE_DOING  = 2 //审批中
	WF_APPROVE_NOPASS = 3 //作废
	WF_APPROVE_Modify = 4 //修改
)
