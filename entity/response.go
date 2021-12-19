package entity

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BadRequestResponse() *DefaultResponse {
	return &DefaultResponse{
		"Failed",
		"Bad request",
		"",
	}
}

func InternalServerErrorResponse(err string) *DefaultResponse {
	return &DefaultResponse{
		"Failed",
		"Bad request",
		err,
	}
}

func NewSuccessResponse(data interface{}) *DefaultResponse {
	return &DefaultResponse{
		"Success",
		"OK",
		data,
	}
}
