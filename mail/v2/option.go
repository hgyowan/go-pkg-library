package v2

import "text/template"

type EmailSenderOption func(*emailSender)
type options struct {
	duplicateChecker func(r *Recipient) bool
	templateFuncMap  template.FuncMap
}

// WithDuplicateChecker
// 중복 체크 로직 (true 일 경우 중복)
func WithDuplicateChecker(f func(*Recipient) bool) EmailSenderOption {
	return func(e *emailSender) {
		e.options.duplicateChecker = f
	}
}

func WithTemplateFuncMap(tf template.FuncMap) EmailSenderOption {
	return func(e *emailSender) {
		e.options.templateFuncMap = tf
	}
}
