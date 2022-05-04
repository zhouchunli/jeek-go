package microerr

var HttpErrMessages = map[ErrType]string{
	Success:           "Success",
	Failure:           "FAILED",
	ClientCommonError: "Illegal request",
	ServerCommonError: "Server internal microerr",
}

func GetHttpErrMessage(httpErrCode ErrType) string {
	return HttpErrMessages[httpErrCode]
}
