package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	AlertMsgGeneric = "Something went wrong. Please try again."
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Yield interface{}
}
