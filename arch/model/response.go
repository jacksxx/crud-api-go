package model

type Response struct {
	Message string
}

func ResponseMessage(str string) Response {
	return Response{
		Message: str,
	}
}

type PaginatedResponse[T interface{}] struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	Data       []T `json:"data"`
}

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}
