package response

type ErrorBody struct {
	Error ErrorItem `json:"error"`
}
type ErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(code, message string) ErrorBody {
	return ErrorBody{
		Error: ErrorItem{Code: code, Message: message},
	}
}
