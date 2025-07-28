package error

import "net/http"

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

	PasswordMisMatch Code = 20001
	AgreeRequired    Code = 20002

	WrongParam Code = 40001
)

var (
	OK = Status{1, http.StatusOK, "", nil, nil}
)
