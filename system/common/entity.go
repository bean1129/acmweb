package common

type ZIResult struct {
	Code   int           `json:"code"`    // 执行结果，0表示成功，其他值表示错误码
	Count  int64         `json:"count"`   // 影响行数
	Msg    string        `json:"msg"`     // 当Code不为0时，为错误原因描述
	Data   []interface{} `json:"data"`    // 数据结果
	SumAll []interface{} `json:"sum_all"` // 所有数据汇总结果
	SumCur []interface{} `json:"sum_cur"` // 当前页数据汇总结果
}

func NewResult() *ZIResult {
	return &ZIResult{
		Code:   OK,
		Count:  0,
		Msg:    "Success",
		Data:   make([]interface{}, 0),
		SumCur: make([]interface{}, 0),
		SumAll: make([]interface{}, 0),
	}
}

func NewErrResult(code int, msg string) *ZIResult {
	return &ZIResult{
		Code:   code,
		Count:  0,
		Msg:    msg,
		Data:   make([]interface{}, 0),
		SumCur: make([]interface{}, 0),
		SumAll: make([]interface{}, 0),
	}
}
