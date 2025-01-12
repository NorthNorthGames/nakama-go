package nakama

import (
	"testing"
)

func TestAuthenticateWithDeviceId(t *testing.T) {
	deviceId := "376C007D-260F-579B-BD75-A3CBBFC2EF99"

	client := NewClient("defaultkey", "127.0.0.1", "7350", false, nil, nil)

	create := true

	session, err := client.AuthenticateDevice(deviceId, &create, nil, nil)

	if err != nil {
		t.Error(err)
	}

	if session == nil {
		t.Error("Session is nil")
	}
}
