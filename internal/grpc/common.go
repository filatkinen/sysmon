package grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMessage(err error, message string, field string) error {
	errorStatus := status.New(codes.Internal, message)
	ds, er := errorStatus.WithDetails(&errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	})
	if er != nil {
		return errorStatus.Err()
	}
	return ds.Err()
}
