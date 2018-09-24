package enums

type JsonResultCode int

const (
	JRCodeSucc JsonResultCode = 1
	JRCodeFailed JsonResultCode = 0
	JRCode302 = 302 //跳转至地址
	JRCode401 = 401 //未授权访问
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)
