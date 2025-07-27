package error

type Status struct {
	Code           int         `json:"code"`
	HttpStatusCode int         `json:"httpStatusCode"`
	Message        string      `json:"message"`
	Detail         []string    `json:"detail"`
	Data           interface{} `json:"data"`
}

func (curStatus *Status) Error() string {
	if curStatus == nil {
		return ""
	}

	return curStatus.Message
}

func (curStatus *Status) AddDetail(detailMsgList ...string) *Status {
	if curStatus.Detail == nil {
		curStatus.Detail = make([]string, 0, 3)
	}

	curStatus.Detail = append(curStatus.Detail, detailMsgList...)

	return curStatus
}

func (curStatus *Status) Equal(comparedStatus *Status) bool {
	return curStatus.Code == comparedStatus.Code
}
