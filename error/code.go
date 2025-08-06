package error

type Code int

const (
	None     Code = 0
	StatusOk Code = 1

	DBQuery Code = 10

	Create Code = 10001
	Update Code = 10002
	Delete Code = 10003
	Get    Code = 10004
	Tx     Code = 10005
	Upsert Code = 10006

	Email Code = 30001
	Kakao Code = 30002

	WrongParam Code = 40001
	Duplicate  Code = 40002
	Expired    Code = 40003
	NotFound   Code = 40004
	CryptData  Code = 40005

	PasswordMisMatch         Code = 50001
	AgreeRequired            Code = 50002
	UnsupportedOAuthProvider Code = 50003
	InvalidSSOAccount        Code = 50004
	InvalidPassword          Code = 50005
)
