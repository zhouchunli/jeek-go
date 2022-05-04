package errcode

/*
错误码设计
*/

// 系统错误码
var (
	Success       = &Errno{Code: 0, Message: "Success"}
	ApiSuccess    = &Errno{Code: 0, Message: "Success"}
	SysSuccess    = &Errno{Code: 0, Message: "Success"}
	ApiFailure    = &Errno{Code: 1, Message: "FAILED"}
	SystemFailure = &Errno{Code: 1, Message: "FAILED"}
	Failed        = &Errno{Code: 1, Message: "Request failed"}
	SystemError   = &Errno{Code: 1, Message: "System internal microerr,please try again"}
	Unknown       = &Errno{Code: 1, Message: "unknown microerr"}
)
