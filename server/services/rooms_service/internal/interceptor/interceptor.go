package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"rooms_service/internal/models"
)

type UserContextKeyType string

const (
	TokenMetadataKey                    = "token"
	UserContextKey   UserContextKeyType = "user"
)

type AuthServiceClient interface {
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
}

// TokenMiddleware is func that returns a middleware that extracts the token from the metadata of incoming gRPC requests,
// sends it to another service, and saves the user in the context
func TokenMiddleware(cl AuthServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		token, ok := extractToken(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "No token found in the request metadata")
		}

		// Send the token to another service and save the user
		user, err := cl.GetUserByToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Invalid token")
		}

		// Save the user in the context
		ctx = context.WithValue(ctx, UserContextKey, user)

		// Invoke the actual RPC method
		resp, err := handler(ctx, req)

		return resp, err
	}
}

// extractToken extracts the token from the metadata of the incoming gRPC request
// and returns it along with a boolean indicating whether the token was found
// name in metadata must be TokenMetadataKey
func extractToken(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	tokens := md.Get(TokenMetadataKey)
	if len(tokens) == 0 {
		return "", false
	}

	return tokens[0], true
}
