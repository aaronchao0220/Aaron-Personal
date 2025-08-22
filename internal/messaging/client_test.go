package messaging

import (
	"context"
	"testing"
	"time"

	gskJTW "github.com/qlik-trial/go-service-kit/v29/jwt"
	"github.com/stretchr/testify/assert"
)

type TokenGeneratorMock struct {
}

func (tg TokenGeneratorMock) GetServiceToServiceJWT(ctx context.Context, audience string, tenantID string) (string, error) {
	return "token", nil

}
func (tg TokenGeneratorMock) GetServiceIdentityJWT(audience string, expiry time.Duration, privateClaims ...gskJTW.PrivateClaim) (string, error) {
	return "token", nil
}

func TestCreateClient(t *testing.T) {
	_, err := CreateClient(context.TODO(), TokenGeneratorMock{}, "clientId")
	assert.NoError(t, err)
}
