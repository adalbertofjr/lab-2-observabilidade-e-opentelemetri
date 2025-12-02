package internalerror

type InternalError struct {
	MSG  string
	Code int
}

func CEPInvalidError() *InternalError {
	return &InternalError{
		MSG:  "Invalid zipcode",
		Code: 422,
	}
}

func CEPNotFoundError() *InternalError {
	return &InternalError{
		MSG:  "Can not find zipcode",
		Code: 404,
	}
}
