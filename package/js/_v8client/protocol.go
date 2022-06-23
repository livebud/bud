package v8client

type Input struct {
	Type string // "script" or "expr"
	Path string
	Code string
}

type Output struct {
	Result string
	Error  string
}
