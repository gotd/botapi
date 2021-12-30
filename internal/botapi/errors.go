package botapi

// NotImplementedError is stub error for not implemented methods.
type NotImplementedError struct {
	Message string
}

// Error implements error.
func (n *NotImplementedError) Error() string {
	if n.Message == "" {
		return "method not implemented yet"
	}
	return n.Message
}

// BadRequestError reports bad request.
type BadRequestError struct {
	Message string
}

// Error implements error.
func (p *BadRequestError) Error() string {
	return p.Message
}

func chatNotFound() *BadRequestError {
	return &BadRequestError{Message: "Bad Request: chat not found"}
}
