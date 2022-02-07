package codes

const (
	OK = iota + 10000
	DEFAULT
	WRONG_PARAM

)

var CodeMap = map[int]string{
	OK: "ok",
	WRONG_PARAM: "wrong param",
	DEFAULT: "default message",
}

func GetMessage(code int) (int, string) {
	msg, ok := CodeMap[code]
	if !ok {
		msg = CodeMap[DEFAULT]
		code = DEFAULT
	}
	return code, msg
}
