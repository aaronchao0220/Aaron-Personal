package auth

import (
	"context"
	"time"

	gskJTW "github.com/qlik-trial/go-service-kit/v29/jwt"
)

// TokenGenerator is the interface for generating tokens
type TokenGenerator interface {
	GetServiceToServiceJWT(ctx context.Context, audience string, tenantID string) (string, error)
	GetServiceIdentityJWT(audience string, expiry time.Duration, privateClaims ...gskJTW.PrivateClaim) (string, error)
}
