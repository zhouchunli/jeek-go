package microerr

const (
	NoTypeErr = ErrType(-1)

	/*
	 * 外部错误
	 */

	//HTTP错误码
	Success           = ErrType(0)
	HttpSuccess       = ErrType(200)
	Failure           = ErrType(1)
	ClientCommonError = ErrType(440) //通用客户端错误
	ServerCommonError = ErrType(550) //通用服务端错误


)
