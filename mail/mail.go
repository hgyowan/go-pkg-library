package mail

import (
	"bytes"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"mime"
	"net/smtp"
	"text/template"
)

type EmailTemplateKey int

const (
	TemplateKeyVerifyEmail EmailTemplateKey = iota + 1
	_mime                                   = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type EmailConfig struct {
	ServerHost string
	ServerPort string
	SenderAddr string
	Username   string
	Password   string
}

type sendFuncType func(string, smtp.Auth, string, []string, []byte) error

type emailSender struct {
	conf        *EmailConfig
	templateMap map[EmailTemplateKey]*template.Template
	sendFunc    func(string, smtp.Auth, string, []string, []byte) error
}

type EmailSender interface {
	SendMail(from string, to []string, body []byte) error
	SendMailWithTemplate(to, subject string, templateType EmailTemplateKey, templateData interface{}) error
}

// MustNewEmailSender
// inviteTemplate, err := template.ParseFiles(templateDir + "inviteCoworkerTemplate.html")
//
//	if err != nil {
//		pkgLogger.ZapLogger.Logger.Sugar().Fatal(err)
//	}
//	templateMap[TemplateKeyCoworkerInvite] = inviteTemplate
func MustNewEmailSender(conf *EmailConfig, templateMap map[EmailTemplateKey]*template.Template, sendFunc ...sendFuncType) EmailSender {
	if len(sendFunc) != 0 {
		return &emailSender{
			conf:        conf,
			sendFunc:    sendFunc[0],
			templateMap: templateMap,
		}
	}

	return &emailSender{
		conf:        conf,
		sendFunc:    smtp.SendMail,
		templateMap: templateMap,
	}
}

func (e *emailSender) SendMail(from string, to []string, body []byte) error {
	emailServerAddr := e.conf.ServerHost + ":" + e.conf.ServerPort
	auth := smtp.PlainAuth("", e.conf.Username, e.conf.Password, e.conf.ServerHost)
	return pkgError.Wrap(e.sendFunc(emailServerAddr, auth, from, to, body))
}

func (e *emailSender) SendMailWithTemplate(to, subject string, templType EmailTemplateKey, templateData interface{}) error {
	tmpl := e.getTemplate(templType)
	if tmpl == nil {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.Get, "no template")
	}
	encodedSubject := mime.QEncoding.Encode("utf-8", subject)
	generatedSubj := "Subject: " + encodedSubject + "\n"

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, templateData); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}
	fromAddress := envs.SMTPAccount
	generatedFrom := "From: " + fromAddress + "\n"
	generatedBody := []byte(generatedFrom + generatedSubj + _mime + "\n" + buf.String())
	return pkgError.WrapWithCode(e.SendMail(fromAddress, []string{to}, generatedBody), pkgError.Email)
}

func (e *emailSender) getTemplate(key EmailTemplateKey) *template.Template {
	return e.templateMap[key]
}
