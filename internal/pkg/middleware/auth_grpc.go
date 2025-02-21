package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not provided")
	}
	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	//tokenString := authHeaders[0]
	//username, err := validateJWT(tokenString)
	username := "Zack"

	/*s
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	*/

	newCtx := context.WithValue(ctx, "username", username)
	return handler(newCtx, req)
}
