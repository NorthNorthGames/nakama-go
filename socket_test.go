package nakama

import (
	"testing"
)

func TestCreateMatch(t *testing.T) {
	deviceId := "376C007D-260F-579B-BD75-A3CBBFC2EF99"

	client := NewClient("defaultkey", "127.0.0.1", "7350", false, nil, nil)

	create := true

	session, _ := client.AuthenticateDevice(deviceId, &create, nil, nil)

	timeout := 1000
	socket := client.CreateSocket(false, true, nil, &timeout)
	connect, err := socket.Connect(*session, nil, &timeout)
	if err != nil {
		t.Error(err)
	}

	if connect == nil {
		t.Error("Connect is nil")
	}

	//match, err := socket.CreateMatch(nil)
	//
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//if match == nil {
	//	t.Error("Match is nil")
	//}
}
