package lib

//通用返回数据
type Result struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//分页返回数据
type Page struct {
	PageNo     int64       `json:"current"`
	PageSize   int64       `json:"size"`
	TotalCount int64       `json:"total"`
	Records    interface{} `json:"records"`
}

func Ok() Result {
	return Result{Code: Success, Msg: "操作成功"}
}
func OkMsg(msg string) Result {
	return Result{Code: Success, Msg: msg}
}
func OkData(data interface{}) Result {
	return Result{Code: Success, Msg: "操作成功", Data: data}
}

func ErrMsg(msg string) Result {
	return Result{Code: Error, Msg: msg}
}

func Err() Result {
	return Result{Code: Error, Msg: "操作失败"}
}

func LoginErr() Result {
	return Result{Code: NotLogin, Msg: "登录过期"}
}
