package error

import "net/http"

var businessCodeMap = map[Code]Status{
	None:    {int(None), http.StatusInternalServerError, "not exists error", nil, nil},
	DBQuery: {int(DBQuery), http.StatusInternalServerError, "DB query error", nil, nil},
	Create:  {int(Create), http.StatusInternalServerError, "fail to create data", nil, nil},
	Update:  {int(Update), http.StatusInternalServerError, "fail to update data", nil, nil},
	Delete:  {int(Delete), http.StatusInternalServerError, "fail to delete data", nil, nil},
	Get:     {int(Get), http.StatusInternalServerError, "fail to get data", nil, nil},
	Tx:      {int(Tx), http.StatusInternalServerError, "fail to start a db transaction", nil, nil},
	Upsert:  {int(Upsert), http.StatusInternalServerError, "fail to upsert data", nil, nil},
}
