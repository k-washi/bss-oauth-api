package grpc_errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//InvalidArgument invalid argument grpc error
func InvalidArgument(msg string) error {
	return status.Errorf(codes.InvalidArgument, msg)
}

//UnImplemented dose not imp
func UnImplemented(msg string) error {
	return status.Errorf(codes.Unimplemented, msg)
}

//InternalError internal error
func InternalError(msg string) error {
	return status.Errorf(codes.Internal, msg)
}
