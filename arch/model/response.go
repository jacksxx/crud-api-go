package model

type Response struct {
	Message string
}

func ResponseMessage(str string) Response {
	return Response{
		Message: str,
	}
}
