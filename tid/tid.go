// Copyright © 2023 Yoshiki Shibata, Inc. All rights reserved.

package tid

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	courierTestingIDHeader = "courier-testing-id"
)

// New returns a TID(testing ID) and a context which contains the TID.
func New(ctx context.Context) (string, context.Context) {
	tid := uuid.NewString()
	return tid, metadata.AppendToOutgoingContext(ctx, courierTestingIDHeader, tid)
}

// Extract extracts an tid from the context.
func Extract(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("meta data is not found in the incomming context")
		return ""
	}

	values := md.Get(courierTestingIDHeader)
	if len(values) == 0 {
		log.Printf("No values for %v found", courierTestingIDHeader)
		return ""
	}
	return values[0]
}
