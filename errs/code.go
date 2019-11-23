package errs

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK", Hint: "成功"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error", Hint: "服务器错误，请稍后重试"}
)
