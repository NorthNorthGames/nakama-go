package nakama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSocket(t *testing.T) (Client, Session) {
	client := NewClient("defaultkey", "127.0.0.1", "7350", false, nil, nil)

	deviceId := "376C007D-260F-579B-BD75-A3CBBFC2EF99"
	create := true
	session, _ := client.AuthenticateDevice(deviceId, &create, nil, nil)

	return *client, *session
}

func TestCreateMatch(t *testing.T) {
	client, session := setupSocket(t)

	socket, connect := createSocket(t, client, session)
	session = connect

	matchName := "Test"
	match, err := socket.CreateMatch(&matchName)

	assert.NoError(t, err)
	assert.NotNil(t, match)
	assert.IsType(t, &Match{}, match)
}

func createSocket(t *testing.T, client Client, session Session) (DefaultSocket, Session) {
	timeout := 1000
	socket := client.CreateSocket(false, true, nil, &timeout)

	assert.IsType(t, DefaultSocket{}, socket)

	connect, err := socket.Connect(session, nil, &timeout)

	assert.NoError(t, err)
	assert.NotNil(t, connect)
	assert.IsType(t, &Session{}, connect)

	return socket, *connect
}
