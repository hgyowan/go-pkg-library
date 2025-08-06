package error

import "net/http"

var businessCodeMap = map[Code]Status{
	None:                     {int(None), http.StatusInternalServerError, "not exists error", nil, nil},
	DBQuery:                  {int(DBQuery), http.StatusInternalServerError, "DB query error", nil, nil},
	Create:                   {int(Create), http.StatusInternalServerError, "fail to create data", nil, nil},
	Update:                   {int(Update), http.StatusInternalServerError, "fail to update data", nil, nil},
	Delete:                   {int(Delete), http.StatusInternalServerError, "fail to delete data", nil, nil},
	Get:                      {int(Get), http.StatusInternalServerError, "fail to get data", nil, nil},
	Tx:                       {int(Tx), http.StatusInternalServerError, "fail to start a db transaction", nil, nil},
	Upsert:                   {int(Upsert), http.StatusInternalServerError, "fail to upsert data", nil, nil},
	Email:                    {int(Email), http.StatusInternalServerError, "fail to send email", nil, nil},
	Kakao:                    {int(Kakao), http.StatusBadGateway, "fail to request kakao", nil, nil},
	PasswordMisMatch:         {int(PasswordMisMatch), http.StatusBadRequest, "password mismatch", nil, nil},
	AgreeRequired:            {int(AgreeRequired), http.StatusBadRequest, "terms agree required", nil, nil},
	UnsupportedOAuthProvider: {int(UnsupportedOAuthProvider), http.StatusBadRequest, "unsupported oauth provider", nil, nil},
	InvalidSSOAccount:        {int(InvalidSSOAccount), http.StatusBadRequest, "invalid sso account", nil, nil},
	InvalidPassword:          {int(InvalidPassword), http.StatusBadRequest, "invalid password", nil, nil},
	WrongParam:               {int(WrongParam), http.StatusBadRequest, "wrong parameters", nil, nil},
	Duplicate:                {int(Duplicate), http.StatusInternalServerError, "duplicate data", nil, nil},
	Expired:                  {int(Expired), http.StatusInternalServerError, "expired data", nil, nil},
	NotFound:                 {int(NotFound), http.StatusNotFound, "not found data", nil, nil},
}
