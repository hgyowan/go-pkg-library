package v2

type EmailTemplateType string

const (
	EmailTemplateTypeInviteSend  EmailTemplateType = "invite_send"
	EmailTemplateTypeJoinConfirm EmailTemplateType = "join_confirm"
	EmailTemplateTypeJoinMessage EmailTemplateType = "join_message"
	EmailTemplateTypeVerifyEmail EmailTemplateType = "verify_email"
)
