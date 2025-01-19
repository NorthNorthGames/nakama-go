package nakama

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/coder/websocket"
)

// WebSocketAdapter is a text-based WebSocket adapter for transmitting payloads over UTF-8.
type WebSocketAdapter struct {
	socket    *websocket.Conn
	onClose   func(err error)
	onError   func(err error)
	onMessage func(message []byte)
	onOpen    func(event interface{}) error
	mu        sync.Mutex // To guard websocket connection reference
}

// NewWebSocketAdapterText creates a new instance of WebSocketAdapter.
func NewWebSocketAdapterText() *WebSocketAdapter {
	return &WebSocketAdapter{}
}

// IsOpen determines if the WebSocket connection is open.
func (w *WebSocketAdapter) IsOpen() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.socket != nil
}

// Close closes the WebSocket connection.
func (w *WebSocketAdapter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.socket != nil {
		_ = w.socket.Close(websocket.StatusNormalClosure, "Client closed connection")
		w.socket = nil
	}
}

// Connect connects to the WebSocket using the specified arguments.
func (w *WebSocketAdapter) Connect(scheme, host, port string, createStatus bool, token string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	urlStr := fmt.Sprintf("%s%s:%s/ws?lang=en&status=%s&token=%s",
		scheme,
		host,
		port,
		url.QueryEscape(fmt.Sprintf("%v", createStatus)),
		url.QueryEscape(token),
	)

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.socket, _, err = websocket.Dial(ctx, urlStr, nil)
	if err != nil {
		return err
	}

	go w.listen()

	return nil
}

// Send sends a message through the WebSocket connection.
func (w *WebSocketAdapter) Send(message interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.socket == nil {
		return fmt.Errorf("WebSocket is not connected")
	}

	// Handle specific cases of match_data_send and party_data_send
	if msgMap, ok := message.(map[string]interface{}); ok {
		handleEncodedData(msgMap, "match_data_send")
		handleEncodedData(msgMap, "party_data_send")
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return w.socket.Write(ctx, websocket.MessageText, msgBytes)
}

// listen listens for messages or errors from the WebSocket server.
func (w *WebSocketAdapter) listen() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		_, message, err := w.socket.Read(ctx)
		if err != nil {
			w.mu.Lock()
			socket := w.socket
			w.mu.Unlock()

			if socket != nil {
				closeStatus := websocket.CloseStatus(err)
				fmt.Printf("WebSocket closed with status: %d\n", closeStatus)

				w.Close()
			}
			break
		}

		var decodedMessage map[string]interface{}
		if err := json.Unmarshal(message, &decodedMessage); err != nil {
			fmt.Printf("Error unmarshalling WebSocket message: %v\n", err)
			continue
		}

		// Handle specific decoding logic for match_data and party_data
		decodeReceivedData(decodedMessage, "match_data")
		decodeReceivedData(decodedMessage, "party_data")

		messageBytes, err := json.Marshal(decodedMessage)
		if err == nil {
			w.onMessage(messageBytes)
		} else if w.onError != nil {
			w.onError(err)
		}
	}
}

// handleEncodedData handles encoding of match_data_send and party_data_send fields.
func handleEncodedData(msg map[string]interface{}, field string) {
	if sendData, exists := msg[field]; exists {
		if sendMap, ok := sendData.(map[string]interface{}); ok {
			// Convert op_code to string
			if opCode, ok := sendMap["op_code"]; ok {
				sendMap["op_code"] = fmt.Sprintf("%v", opCode)
			}

			// Encode data
			if payload, exists := sendMap["data"]; exists {
				switch v := payload.(type) {
				case []byte:
					sendMap["data"] = base64.StdEncoding.EncodeToString(v)
				case string:
					sendMap["data"] = base64.StdEncoding.EncodeToString([]byte(v))
				}
			}
		}
	}
}

// decodeReceivedData decodes the match_data and party_data fields in messages received from the server.
func decodeReceivedData(msg map[string]interface{}, field string) {
	if data, exists := msg[field]; exists {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if encoded, exists := dataMap["data"]; exists {
				if encodedStr, ok := encoded.(string); ok {
					decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
					if err == nil {
						dataMap["data"] = decodedBytes
					}
				}
			}
		}
	}
}
