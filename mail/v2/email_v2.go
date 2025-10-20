package v2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
)

type EmailConfig struct {
	SMTPHost         string
	SMTPPort         string
	SMTPSender       string
	Username         string
	Password         string
	DelayBetweenMsg  time.Duration
	MaxRetries       int
	EmailWorkerCount int
}

type Recipient struct {
	ID               string // 해당 메일에 대한 고유 키 (중복 발송 방지)
	LangCode         string
	TemplateType     EmailTemplateType
	ToEmails         []string // 변하지 않는 템플릿 메일일경우 여러명에게 전송 가능
	Subject          string
	TemplateMetaData interface{}
}

type SendMailResponse struct {
	MailID       string
	Status       string
	LangCode     string
	TemplateType EmailTemplateType
	TemplateData string
	Emails       []string
	FailReason   EmailError
	SendingDate  time.Time
}

type emailSender struct {
	tempDirectory string
	conf          *EmailConfig
	client        *smtp.Client
	tmplCache     map[string]*template.Template
	options       options
	mu            sync.RWMutex
}

type EmailSenderV2 interface {
	SendMailWithTemplateV2Parallel(list []*Recipient) <-chan *SendMailResponse
	GetTemplate(tmplType EmailTemplateType, langCode string) (*template.Template, EmailError)
	SetOptions(opts ...EmailSenderOption)
}

func MustNewEmailSender(conf *EmailConfig, tempDirectory string) EmailSenderV2 {
	e := &emailSender{
		conf:          conf,
		tempDirectory: tempDirectory,
	}

	smtpClient, err := e.connect(e.conf)
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatal(err)
	}

	e.client = smtpClient

	// option default setting
	e.options.duplicateChecker = func(r *Recipient) bool {
		return false
	}

	return e
}

func (e *emailSender) SetOptions(opts ...EmailSenderOption) {
	for _, opt := range opts {
		opt(e)
	}
}

func (e *emailSender) SendMailWithTemplateV2Parallel(list []*Recipient) <-chan *SendMailResponse {
	responseCh := make(chan *SendMailResponse, len(list))
	if len(list) == 0 {
		close(responseCh)
		return responseCh
	}

	jobs := make(chan *Recipient, len(list))

	worker := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()

		// worker 전용 SMTP client 생성
		client, err := e.connect(e.conf)
		if err != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
			return
		}
		defer func() {
			_ = client.Quit()
			_ = client.Close()
		}()

		for r := range jobs {
			b, _ := json.Marshal(r.TemplateMetaData)
			resp := &SendMailResponse{
				MailID:       r.ID,
				Status:       "FAIL",
				LangCode:     r.LangCode,
				TemplateType: r.TemplateType,
				TemplateData: string(b),
				Emails:       r.ToEmails,
				SendingDate:  time.Now().UTC(),
			}

			for attempt := 0; attempt < e.conf.MaxRetries; attempt++ {
				emailErr := e.sendMailWithTemplateV2(r, client)
				if emailErr == "" {
					resp.Status = "SUCCESS"
					break
				}

				switch emailErr {
				case EmailErrorSetFrom, EmailErrorSetRcpt, EmailErrorSendStart: // 세션 끊김 재연결
					_ = client.Close()
					client, err = e.connect(e.conf)
					if err != nil {
						resp.FailReason = EmailErrorConnection
						break
					}
					continue
				case EmailErrorDuplicate: // 중복일 경우 중복 처리
					resp.Status = "DUPLICATE"
					break
				}

				resp.FailReason = emailErr
				break
			}

			responseCh <- resp

			if e.conf.DelayBetweenMsg > 0 {
				time.Sleep(e.conf.DelayBetweenMsg)
			}
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < e.conf.EmailWorkerCount; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	go func() {
		for _, r := range list {
			jobs <- r
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(responseCh)
	}()

	return responseCh
}

func (e *emailSender) connect(cfg *EmailConfig) (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort), 30*time.Second)
	if err != nil {
		return nil, pkgError.Wrap(err)
	}
	c, err := smtp.NewClient(conn, cfg.SMTPHost)
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: envs.NASTLSHost}
		if err = c.StartTLS(tlsConfig); err != nil {
			_ = c.Close()
			return nil, pkgError.Wrap(err)
		}
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	if err = c.Auth(auth); err != nil {
		_ = c.Close()
		return nil, pkgError.Wrap(err)
	}
	return c, nil
}

func (e *emailSender) sendMailWithTemplateV2(r *Recipient, client *smtp.Client) EmailError {
	tmpl, emailErr := e.GetTemplate(r.TemplateType, r.LangCode)
	if emailErr != "" {
		return emailErr
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, r.TemplateMetaData); err != nil {
		return EmailErrorTemplateExecute
	}

	headers := "From: " + e.conf.SMTPSender + "\r\n" +
		"To: " + strings.Join(r.ToEmails, ", ") + "\r\n" +
		"Subject: " + r.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n"

	body := []byte(headers + buf.String())

	// 중복 메일 처리 하지 않음
	if e.options.duplicateChecker(r) {
		return EmailErrorDuplicate
	}

	if err := client.Mail(e.conf.SMTPSender); err != nil {
		return EmailErrorSetFrom
	}

	for _, rcpt := range r.ToEmails {
		if err := client.Rcpt(rcpt); err != nil {
			return EmailErrorSetRcpt
		}
	}

	w, err := client.Data()
	if err != nil {
		return EmailErrorSendStart
	}
	if _, err = w.Write(body); err != nil {
		_ = w.Close()
		return EmailErrorSendWrite
	}

	if err = w.Close(); err != nil {
		return EmailErrorSendClose
	}

	return ""
}

func (e *emailSender) GetTemplate(tmplType EmailTemplateType, langCode string) (*template.Template, EmailError) {
	if langCode == "" {
		langCode = "ko"
	}

	e.mu.RLock()
	if tmpl, ok := e.tmplCache[fmt.Sprintf("%s-%s", tmplType, strings.ToLower(langCode))]; ok {
		e.mu.RUnlock()
		return tmpl, ""
	}
	e.mu.RUnlock()

	tmpl := template.New(fmt.Sprintf("%s_%s.html", tmplType, strings.ToLower(langCode)))
	if len(e.options.templateFuncMap) > 0 {
		tmpl = tmpl.Funcs(e.options.templateFuncMap)
	}

	var err error
	tmpl, err = tmpl.ParseFiles(e.tempDirectory + fmt.Sprintf("%s_%s.html", string(tmplType), strings.ToLower(langCode)))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, EmailErrorGetTemplate
		}

		tmpl, err = template.ParseFiles(e.tempDirectory + fmt.Sprintf("%s_ko.html", string(tmplType)))
		if err != nil {
			return nil, EmailErrorTemplateNotExists
		}
	}

	e.mu.Lock()
	if e.tmplCache == nil {
		e.tmplCache = make(map[string]*template.Template)
	}
	e.tmplCache[fmt.Sprintf("%s-%s", tmplType, langCode)] = tmpl
	e.mu.Unlock()

	return tmpl, ""
}
