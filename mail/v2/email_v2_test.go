package v2

import (
	"github.com/hgyowan/go-pkg-library/envs"
	"testing"
	"time"
)

func Test_SendEmail(t *testing.T) {
	sender := MustNewEmailSender(&EmailConfig{
		SMTPHost:         envs.SMTPServer,
		SMTPPort:         envs.SMTPPort,
		SMTPSender:       envs.SMTPSender,
		Username:         envs.SMTPAccount,
		Password:         envs.SMTPPassword,
		DelayBetweenMsg:  500 * time.Microsecond,
		MaxRetries:       3,
		EmailWorkerCount: 3,
	}, "")

	type BestContents struct {
		ThumbnailURL string `json:"thumbnailURL"`
		SpaceTitle   string `json:"spaceTitle"`
		SpaceInfo    string `json:"spaceInfo"`
	}

	type Test struct {
		MailOpenCheckURL string
		BestContents     []*BestContents
	}

	test := &Test{
		MailOpenCheckURL: "https://xromeda-back-office-dev.sandbox-olimzone.com/v1/email/send",
		BestContents: []*BestContents{
			{
				SpaceInfo:    "SpaceInfo1",
				SpaceTitle:   "SpaceTitle1",
				ThumbnailURL: "ThumbnailURL1",
			},
			{
				SpaceInfo:    "SpaceInfo2",
				SpaceTitle:   "SpaceTitle2",
				ThumbnailURL: "ThumbnailURL2",
			},
		},
	}

	res := sender.SendMailWithTemplateV2Parallel([]*Recipient{
		{
			TemplateType:     EmailTemplateTypeInviteSend,
			ToEmails:         []string{"rydhkstptkd@naver.com"},
			Subject:          "제목 테스트",
			TemplateMetaData: test,
			LangCode:         "KR",
		},
		{
			TemplateType:     EmailTemplateTypeInviteSend,
			ToEmails:         []string{"gwhwang@olimplanet.com"},
			Subject:          "제목 테스트",
			TemplateMetaData: test,
			LangCode:         "KR",
		},
	})

	for r := range res {
		t.Log(r)
	}
}
