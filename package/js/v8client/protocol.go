package v8client

type input struct {
	Type string // "script" or "expr"
	Path string
	Code string
}

type output struct {
	Result string
	Error  string
}
