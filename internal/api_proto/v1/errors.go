package v1

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

func (x *Error) Error() string {
	return fmt.Sprintf("%d: %s", x.Code, x.Message)
}

/*
var (

	// ErrorInternal means server error, if returned this is a signal
	//something went wrong in the service itself
	ErrorInternal = &Error{
		Code:    100,
		Message: "internal service error",
	}

	ErrorBadRequest = &Error{
		Code:    107,
		Message: "bad request",
	}

	//TODO: Add more service-specific errors

)
*/
func ErrorInvalidArgument(msg string) *Error {
	return &Error{
		Code:    uint32(codes.InvalidArgument),
		Message: msg,
	}
}

func ErrorResourceNotFound(msg string) *Error {
	return &Error{
		Code:    uint32(codes.NotFound),
		Message: msg,
	}
}

func ErrorInternal(msg string) *Error {
	return &Error{
		Code:    uint32(codes.Internal),
		Message: msg,
	}
}

func ErrorUnauthenticated(msg string) *Error {
	return &Error{
		Code:    uint32(codes.Unauthenticated),
		Message: msg,
	}
}
