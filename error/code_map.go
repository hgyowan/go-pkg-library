package error

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

var businessCodeMap = map[Code]Status{
	None:                     {int(None), http.StatusInternalServerError, "not exists error", nil, nil, codes.Internal},
	DBQuery:                  {int(DBQuery), http.StatusInternalServerError, "DB query error", nil, nil, codes.Internal},
	Create:                   {int(Create), http.StatusInternalServerError, "fail to create data", nil, nil, codes.Internal},
	Update:                   {int(Update), http.StatusInternalServerError, "fail to update data", nil, nil, codes.Internal},
	Delete:                   {int(Delete), http.StatusInternalServerError, "fail to delete data", nil, nil, codes.Internal},
	Get:                      {int(Get), http.StatusInternalServerError, "fail to get data", nil, nil, codes.Internal},
	Tx:                       {int(Tx), http.StatusInternalServerError, "fail to start a db transaction", nil, nil, codes.Internal},
	Upsert:                   {int(Upsert), http.StatusInternalServerError, "fail to upsert data", nil, nil, codes.Internal},
	Email:                    {int(Email), http.StatusInternalServerError, "fail to send email", nil, nil, codes.Internal},
	Kakao:                    {int(Kakao), http.StatusBadGateway, "fail to request kakao", nil, nil, codes.Unavailable},
	OdCloud:                  {int(OdCloud), http.StatusBadGateway, "fail to request odcloud", nil, nil, codes.Unavailable},
	Google:                   {int(Google), http.StatusBadGateway, "fail to request google", nil, nil, codes.Unavailable},
	PasswordMisMatch:         {int(PasswordMisMatch), http.StatusBadRequest, "password mismatch", nil, nil, codes.InvalidArgument},
	AgreeRequired:            {int(AgreeRequired), http.StatusBadRequest, "terms agree required", nil, nil, codes.InvalidArgument},
	UnsupportedOAuthProvider: {int(UnsupportedOAuthProvider), http.StatusBadRequest, "unsupported oauth provider", nil, nil, codes.InvalidArgument},
	InvalidSSOAccount:        {int(InvalidSSOAccount), http.StatusBadRequest, "invalid sso account", nil, nil, codes.InvalidArgument},
	InvalidPassword:          {int(InvalidPassword), http.StatusBadRequest, "invalid password", nil, nil, codes.InvalidArgument},
	AlreadyExistsEmail:       {int(AlreadyExistsEmail), http.StatusConflict, "already exists email", nil, nil, codes.AlreadyExists},
	WithdrawalPending:        {int(WithdrawalPending), http.StatusForbidden, "withdrawal pending", nil, nil, codes.FailedPrecondition},
	WrongParam:               {int(WrongParam), http.StatusBadRequest, "wrong parameters", nil, nil, codes.InvalidArgument},
	Duplicate:                {int(Duplicate), http.StatusConflict, "duplicate data", nil, nil, codes.AlreadyExists},
	Expired:                  {int(Expired), http.StatusGone, "expired data", nil, nil, codes.NotFound},
	NotFound:                 {int(NotFound), http.StatusNotFound, "not found data", nil, nil, codes.NotFound},
	CryptData:                {int(CryptData), http.StatusInternalServerError, "fail to crypto", nil, nil, codes.Internal},
	Permission:               {int(Permission), http.StatusForbidden, "access denied", nil, nil, codes.PermissionDenied},
	CanNotEmpty:              {int(CanNotEmpty), http.StatusInternalServerError, "can not empty", nil, nil, codes.Internal},
	AlreadyExists:            {int(AlreadyExists), http.StatusConflict, "already exists", nil, nil, codes.AlreadyExists},
	InvalidFileType:          {int(InvalidFileType), http.StatusBadRequest, "invalid file type", nil, nil, codes.InvalidArgument},
	FailReadFile:             {int(FailReadFile), http.StatusInternalServerError, "fail read file", nil, nil, codes.Internal},
	FailWriteFile:            {int(FailWriteFile), http.StatusInternalServerError, "fail write file", nil, nil, codes.Internal},
	FileSizeExceeded:         {int(FileSizeExceeded), http.StatusRequestEntityTooLarge, "file size exceeded", nil, nil, codes.Internal},
	MisMatchOwnerName:        {int(MisMatchOwnerName), http.StatusForbidden, "mismatch owner name", nil, nil, codes.PermissionDenied},
	NotAllowedOwner:          {int(NotAllowedOwner), http.StatusForbidden, "not allowed owner", nil, nil, codes.PermissionDenied},
	SettlementInProgress:     {int(SettlementInProgress), http.StatusLocked, "settlement in progress", nil, nil, codes.FailedPrecondition},
}
