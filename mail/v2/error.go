package v2

type EmailError string

const (
	EmailErrorConnection        EmailError = "connection"
	EmailErrorGetTemplate       EmailError = "getTemplate"
	EmailErrorTemplateNotExists EmailError = "templateNotExists"
	EmailErrorTemplateExecute   EmailError = "templateExecute"
	EmailErrorSetFrom           EmailError = "setFrom"
	EmailErrorSetRcpt           EmailError = "setRcpt"
	EmailErrorSendStart         EmailError = "sendStart"
	EmailErrorSendWrite         EmailError = "sendWrite"
	EmailErrorSendClose         EmailError = "sendClose"
	EmailErrorDuplicate         EmailError = "duplicate"
)
