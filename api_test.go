package nakama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupApi(t *testing.T) Client {
	client := NewClient("defaultkey", "127.0.0.1", "7350", false, nil, nil)

	return *client
}

func TestAuthenticateWithDeviceId(t *testing.T) {
	client := setupApi(t)

	deviceId := "376C007D-260F-579B-BD75-A3CBBFC2EF99"
	create := true
	session, err := client.AuthenticateDevice(deviceId, &create, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.IsType(t, &Session{}, session)
}
