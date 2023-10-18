// Copyright © 2023 Yoshiki Shibata, Inc. All rights reserved.

package tid

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewGRPCHeaderPropagator returns an UnaryServerInterceptor that propagtes
// TID header from the incoming context to the outgoin context.
func NewGRPCHeaderPropagator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for k, v := range md {
				k = strings.ToLower(k)

				if strings.HasPrefix(k, courierTestingIDHeader) && len(v) > 0 {
					ctx = metadata.AppendToOutgoingContext(ctx, k, v[0])
					return handler(ctx, req)
				}
			}
		}

		return handler(ctx, req)
	}
}
