package xutil

/*
服务级错误码 模块级错误码 子模块错误码 具体错误码

- 服务级错误码：1 位数进行表示，比如 1 为系统级错误；2 为普通错误。
- 模块级错误码：2 位数进行表示。
- 子模块错误码：2 位数进行表示。
- 具体的错误码：2 位数进行表示。
*/

const (
	NoError = 0
)

type Error struct {
	Code int
	Msg  string
}

// system error
const (
	RequestMethodError            = 1000001
	RequestAppIdError             = 1000002
	RequestSignatureError         = 1000003
	RequestTimestampValidateError = 1000004
	RequestTimestampExpiredError  = 1000005
	RequestNonceValidateError     = 1000006
	RequestNonceDuplicatedError   = 1000007
	RequestIpNotAllowedError      = 1000008
	RequestHeaderError            = 1000009
	RequestUserError              = 1000010

	RequestBodyParseError    = 1000030
	RequestBodyValidateError = 1000031

	SystemError = 1000098
	OtherError  = 1000099
)
