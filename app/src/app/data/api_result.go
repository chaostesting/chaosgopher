package data

type APIResult struct {
	Success      bool
	ErrorMessage string
	// TODO: Payload interface{} ?
}

func SuccessResult() *APIResult {
	return &APIResult{true, ""}
}

func NotFoundResult() *APIResult {
	return &APIResult{false, "not found."}
}

func ErrorResult(e error) *APIResult {
	return &APIResult{false, e.Error()}
}

func (r *APIResult) Error() string {
	return r.ErrorMessage
}
