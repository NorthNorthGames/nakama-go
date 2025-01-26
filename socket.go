package nakama

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type PromiseExecutor struct {
	Resolve func(value interface{})
	Reject  func(reason error)
}

type Presence struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
	Node      string `json:"node"`
}

type Channel struct {
	ID        string     `json:"id"`
	Presences []Presence `json:"presences"`
	Self      Presence   `json:"self"`
}

type ChannelJoin struct {
	ChannelJoin struct {
		Target      string `json:"target"`
		Type        int    `json:"type"`
		Persistence bool   `json:"persistence"`
		Hidden      bool   `json:"hidden"`
	} `json:"channel_join"`
}

type ChannelLeave struct {
	ChannelLeave struct {
		ChannelID string `json:"channel_id"`
	} `json:"channel_leave"`
}

type ChannelMessageAck struct {
	ChannelID   string `json:"channel_id"`
	MessageID   string `json:"message_id"`
	Code        int    `json:"code"`
	Username    string `json:"username"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Persistence bool   `json:"persistence"`
}

type ChannelMessageSend struct {
	ChannelMessageSend struct {
		ChannelID string      `json:"channel_id"`
		Content   interface{} `json:"content"`
	} `json:"channel_message_send"`
}

type ChannelMessageUpdate struct {
	ChannelMessageUpdate struct {
		ChannelID string      `json:"channel_id"`
		MessageID string      `json:"message_id"`
		Content   interface{} `json:"content"`
	} `json:"channel_message_update"`
}

type ChannelMessageRemove struct {
	ChannelMessageRemove struct {
		ChannelID string `json:"channel_id"`
		MessageID string `json:"message_id"`
	} `json:"channel_message_remove"`
}

type ChannelPresenceEvent struct {
	ChannelID string     `json:"channel_id"`
	Joins     []Presence `json:"joins"`
	Leaves    []Presence `json:"leaves"`
}

type StreamId struct {
	Mode       int    `json:"mode"`
	Subject    string `json:"subject"`
	Subcontext string `json:"subcontext"`
	Label      string `json:"label"`
}

type StreamData struct {
	Stream   StreamId  `json:"stream"`
	Sender   *Presence `json:"sender,omitempty"`
	Data     string    `json:"data"`
	Reliable *bool     `json:"reliable,omitempty"`
}

type StreamPresenceEvent struct {
	Stream StreamId   `json:"stream"`
	Joins  []Presence `json:"joins"`
	Leaves []Presence `json:"leaves"`
}

type MatchPresenceEvent struct {
	MatchID string     `json:"match_id"`
	Joins   []Presence `json:"joins"`
	Leaves  []Presence `json:"leaves"`
}

type MatchmakerAdd struct {
	MatchmakerAdd struct {
		MinCount          int                `json:"min_count"`
		MaxCount          int                `json:"max_count"`
		Query             string             `json:"query"`
		StringProperties  map[string]string  `json:"string_properties,omitempty"`
		NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
	} `json:"matchmaker_add"`
}

type MatchmakerTicket struct {
	Ticket string `json:"ticket"`
}

type MatchmakerRemove struct {
	MatchmakerRemove struct {
		Ticket string `json:"ticket"`
	} `json:"matchmaker_remove"`
}

type MatchmakerUser struct {
	Presence          Presence           `json:"presence"`
	PartyID           string             `json:"party_id"`
	StringProperties  map[string]string  `json:"string_properties,omitempty"`
	NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
}

type MatchmakerMatched struct {
	Ticket  string           `json:"ticket"`
	MatchID string           `json:"match_id"`
	Token   string           `json:"token"`
	Users   []MatchmakerUser `json:"users"`
	Self    MatchmakerUser   `json:"self"`
}

type Match struct {
	MatchID       string     `json:"match_id"`
	Authoritative bool       `json:"authoritative"`
	Label         *string    `json:"label,omitempty"`
	Size          int        `json:"size"`
	Presences     []Presence `json:"presences"`
	Self          Presence   `json:"self"`
}

type CreateMatch struct {
	MatchCreate struct {
		Name string `json:"name,omitempty"`
	} `json:"match_create"`
}

type JoinMatch struct {
	MatchJoin struct {
		MatchID  *string                `json:"match_id,omitempty"`
		Token    *string                `json:"token,omitempty"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	} `json:"match_join"`
}

type LeaveMatch struct {
	MatchLeave struct {
		MatchID string `json:"match_id"`
	} `json:"match_leave"`
}

type MatchData struct {
	MatchID  string    `json:"match_id"`
	OpCode   int       `json:"op_code"`
	Data     []byte    `json:"data"`
	Presence *Presence `json:"presence,omitempty"`
	Reliable *bool     `json:"reliable,omitempty"`
}

type MatchDataSend struct {
	MatchDataSend struct {
		MatchID   string      `json:"match_id"`
		OpCode    int         `json:"op_code"`
		Data      interface{} `json:"data"`
		Presences []Presence  `json:"presences"`
		Reliable  *bool       `json:"reliable,omitempty"`
	} `json:"match_data_send"`
}

type Party struct {
	PartyID   string     `json:"party_id"`
	Open      bool       `json:"open"`
	MaxSize   int        `json:"max_size"`
	Self      Presence   `json:"self"`
	Leader    Presence   `json:"leader"`
	Presences []Presence `json:"presences"`
}

type PartyCreate struct {
	PartyCreate struct {
		Open    bool `json:"open"`
		MaxSize int  `json:"max_size"`
	} `json:"party_create"`
}

type PartyJoin struct {
	PartyJoin struct {
		PartyID string `json:"party_id"`
	} `json:"party_join"`
}

type PartyLeave struct {
	PartyLeave struct {
		PartyID string `json:"party_id"`
	} `json:"party_leave"`
}

type PartyPromote struct {
	PartyPromote struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_promote"`
}

type PartyLeader struct {
	PartyID  string   `json:"party_id"`
	Presence Presence `json:"presence"`
}

type PartyAccept struct {
	PartyAccept struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_accept"`
}

type PartyClose struct {
	PartyClose struct {
		PartyID string `json:"party_id"`
	} `json:"party_close"`
}

type PartyData struct {
	PartyID  string   `json:"party_id"`
	Presence Presence `json:"presence"`
	OpCode   int      `json:"op_code"`
	Data     []byte   `json:"data"`
}

type PartyDataSend struct {
	PartyDataSend struct {
		PartyID string      `json:"party_id"`
		OpCode  int         `json:"op_code"`
		Data    interface{} `json:"data"`
	} `json:"party_data_send"`
}

type PartyJoinRequest struct {
	PartyID   string     `json:"party_id"`
	Presences []Presence `json:"presences"`
}

type PartyJoinRequestList struct {
	PartyJoinRequestList struct {
		PartyID string `json:"party_id"`
	} `json:"party_join_request_list"`
}

type PartyMatchmakerAdd struct {
	PartyMatchmakerAdd struct {
		PartyID           string             `json:"party_id"`
		MinCount          int                `json:"min_count"`
		MaxCount          int                `json:"max_count"`
		Query             string             `json:"query"`
		StringProperties  map[string]string  `json:"string_properties,omitempty"`
		NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
	} `json:"party_matchmaker_add"`
}

type PartyMatchmakerRemove struct {
	PartyMatchmakerRemove struct {
		PartyID string `json:"party_id"`
		Ticket  string `json:"ticket"`
	} `json:"party_matchmaker_remove"`
}

type PartyMatchmakerTicket struct {
	PartyID string `json:"party_id"`
	Ticket  string `json:"ticket"`
}

type PartyPresenceEvent struct {
	PartyID string     `json:"party_id"`
	Joins   []Presence `json:"joins"`
	Leaves  []Presence `json:"leaves"`
}

type PartyRemove struct {
	PartyRemove struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_remove"`
}

type Rpc struct {
	Rpc ApiRpc `json:"rpc"`
}

type Ping struct {
	// No fields needed for Ping
}

type Status struct {
	Presences []Presence `json:"presences"`
}

type StatusFollow struct {
	StatusFollow struct {
		UserIDs []string `json:"user_ids"`
	} `json:"status_follow"`
}

type StatusPresenceEvent struct {
	Joins  []Presence `json:"joins"`
	Leaves []Presence `json:"leaves"`
}

type StatusUnfollow struct {
	StatusUnfollow struct {
		UserIDs []string `json:"user_ids"`
	} `json:"status_unfollow"`
}

type StatusUpdate struct {
	StatusUpdate struct {
		Status *string `json:"status,omitempty"`
	} `json:"status_update"`
}

// Socket defines the Go struct with corresponding methods.
type Socket interface {
	OnDisconnect(err error)
	OnError(err error)
	OnHeartbeatTimeout()
}

// SocketError represents an error received from a socket message.
type SocketError struct {
	Code    int    `json:"code"`    // The error code
	Message string `json:"message"` // A message in English to help developers debug the response
}

type Message struct {
	Cid           *string         `json:"cid"`
	Error         *error          `json:"error"`
	Notifications *[]Notification `json:"notifications"`
	Payload       interface{}     `json:"payload"`
}

// DefaultSocket constants
const (
	DefaultHeartbeatTimeoutMs = 10000
	DefaultSendTimeoutMs      = 10000
	DefaultConnectTimeoutMs   = 30000
)

// DefaultSocket represents a WebSocket connection to the Nakama server
type DefaultSocket struct {
	Host               string
	Port               string
	UseSSL             bool
	Verbose            bool
	Adapter            WebSocketAdapter
	SendTimeoutMs      int
	HeartbeatTimeoutMs int
	cIds               map[string]*PromiseExecutor
	nextCid            int
}

// NewDefaultSocket creates an instance of DefaultSocket.
func NewDefaultSocket(host, port string, useSSL, verbose bool, adapter *WebSocketAdapter, sendTimeoutMs *int) DefaultSocket {
	if adapter == nil {
		adapter = NewWebSocketAdapterText()
	}
	if sendTimeoutMs == nil {
		defaultTimeout := DefaultSendTimeoutMs
		sendTimeoutMs = &defaultTimeout
	}

	return DefaultSocket{
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Verbose:            verbose,
		Adapter:            *adapter,
		SendTimeoutMs:      *sendTimeoutMs,
		HeartbeatTimeoutMs: DefaultHeartbeatTimeoutMs,
		cIds:               make(map[string]*PromiseExecutor),
		nextCid:            1,
	}
}

// GenerateCID generates a unique client ID for requests.
func (socket *DefaultSocket) GenerateCID() string {
	cid := fmt.Sprintf("%d", socket.nextCid)
	socket.nextCid++
	return cid
}

// Connect establishes the WebSocket connection with optional timeouts.
func (socket *DefaultSocket) Connect(session Session, createStatus *bool, timeoutMs *int) (*Session, error) {
	if createStatus == nil {
		defaultStatus := false
		createStatus = &defaultStatus
	}

	if timeoutMs == nil {
		defaultTimeout := DefaultConnectTimeoutMs
		timeoutMs = &defaultTimeout
	}

	if socket.Adapter.IsOpen() {
		return &session, nil
	}

	scheme := "ws://"
	if socket.UseSSL {
		scheme = "wss://"
	}

	err := socket.Adapter.Connect(scheme, socket.Host, socket.Port, *createStatus, session.Token)
	if err != nil {
		return nil, err
	}

	socket.Adapter.onClose = func(err error) {
		socket.OnDisconnect(err)
	}

	socket.Adapter.onError = func(err error) {
		socket.OnError(err)
	}

	socket.Adapter.onMessage = func(message []byte) {
		if socket.Verbose == true {
			fmt.Println("Received message:", string(message))
		}

		var messageObject *Message
		if err := json.Unmarshal(message, &messageObject); err != nil {
			if socket.Verbose {
				fmt.Println("Failed to unmarshal message into custom object:", err)
			}
			return
		}

		if messageObject == nil {
			if socket.Verbose {
				fmt.Println("Received empty message")
			}
		}

		if messageObject.Cid == nil {
			if messageObject.Notifications != nil {

			}
		} else {
			executor := socket.cIds[*messageObject.Cid]
			if executor == nil {
				if socket.Verbose {
					log.Printf("No promise executor for message: %v\n", messageObject)
				}
				return
			}

			delete(socket.cIds, *messageObject.Cid)

			if messageObject.Error != nil {
				executor.Reject(*messageObject.Error)
			} else {
				executor.Resolve(messageObject)
			}
		}
	}

	go func() {
		socket.Adapter.onOpen = func(event interface{}) error {
			log.Printf("Socket opened: %v\n", event)

			socket.pingPong()

			// Set a timeout for the connection process
			resChan := make(chan error, 1)
			go func() {
				time.Sleep(time.Duration(*timeoutMs) * time.Millisecond)
				resChan <- errors.New("socket connection timed out")
			}()

			select {
			case err := <-resChan:
				if err != nil {
					socket.Adapter.Close()
					return err
				}
			}

			return nil
		}
	}()

	return &session, nil
}

// Disconnect terminates the WebSocket connection.
func (socket *DefaultSocket) Disconnect(fireDisconnectEvent bool) {
	if socket.Adapter.IsOpen() {
		socket.Adapter.Close()
	}
	if fireDisconnectEvent {
		socket.OnDisconnect(fmt.Errorf("socket disconnected"))
	}
}

// SetHeartbeatTimeoutMs sets the timeout for heartbeat pings.
func (socket *DefaultSocket) SetHeartbeatTimeoutMs(ms int) {
	socket.HeartbeatTimeoutMs = ms
}

// GetHeartbeatTimeoutMs gets the timeout for heartbeat pings.
func (socket *DefaultSocket) GetHeartbeatTimeoutMs() int {
	return socket.HeartbeatTimeoutMs
}

// OnDisconnect handles WebSocket disconnections.
func (socket *DefaultSocket) OnDisconnect(evt error) {
	if socket.Verbose {
		fmt.Println("OnDisconnect:", evt)
	}
}

// OnError handles WebSocket errors.
func (socket *DefaultSocket) OnError(evt error) {
	if socket.Verbose {
		fmt.Println("OnError:", evt)
	}
}

// HandleMessage processes incoming WebSocket messages.
func (socket *DefaultSocket) HandleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		if socket.Verbose {
			fmt.Println("Failed to parse message:", err)
		}
		return
	}

	if cid, ok := msg["cid"].(string); ok {
		executor, exists := socket.cIds[cid]
		if exists {
			delete(socket.cIds, cid)
			if _, hasError := msg["error"]; hasError {
				executor.Reject(errors.New(msg["error"].(string)))
			} else {
				executor.Resolve(msg)
			}
		} else {
			if socket.Verbose {
				fmt.Println("No promise executor for message CID:", cid)
			}
		}
	} else {
		// Handle different message types here (notifications, match data, etc.)
		if socket.Verbose {
			fmt.Println("Message received:", string(message))
		}
	}
}

// Send sends a message to the WebSocket server with optional timeout.
func (socket *DefaultSocket) Send(message interface{}, sendTimeout *int) error {
	if sendTimeout == nil {
		sendTimeout = new(int)
		*sendTimeout = DefaultTimeoutMs
	}

	if !socket.Adapter.IsOpen() {
		return errors.New("socket connection is not established")
	}

	cid := socket.GenerateCID()
	socket.cIds[cid] = &PromiseExecutor{
		Resolve: func(result interface{}) {
			if socket.Verbose {
				fmt.Println("Message sent successfully")
			}
		},
		Reject: func(e error) {
			if socket.Verbose {
				fmt.Println("Message failed:", e)
			}
		},
	}

	err := socket.Adapter.Send(message)
	if err != nil {
		log.Print(err)
		return err
	}

	// Set a timeout for the send operation
	go func(cid string) {
		time.Sleep(time.Duration(*sendTimeout) * time.Millisecond)
		delete(socket.cIds, cid)
	}(cid)

	return nil
}

// ReadResponse reads and parses the next response from the WebSocket connection.
func (socket *DefaultSocket) Read() (map[string]interface{}, error) {
	if !socket.Adapter.IsOpen() {
		return nil, errors.New("socket connection is not established")
	}

	fmt.Println("Reading message...")
	message, err := socket.Adapter.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read message from socket: %w", err)
	}

	fmt.Println("Message:", string(message))

	var response map[string]interface{}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("failed to parse socket response: %w", err)
	}

	return response, nil
}

// CreateMatch sends a request to create a match and returns the created Match.
func (socket *DefaultSocket) CreateMatch(name *string) (*Match, error) {
	request := map[string]interface{}{
		"match_create": map[string]interface{}{},
	}
	if name != nil {
		request["match_create"].(map[string]interface{})["name"] = name
	}

	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	response, err := socket.Read()

	if socket.Verbose {
		fmt.Println("Response received:", response)
	}

	if match, ok := response["match"].(*Match); ok {
		return match, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid match field")
}

// CreateParty Example methods for handling specific socket calls
func (socket *DefaultSocket) CreateParty(open bool, maxSize int) (*Party, error) {
	request := map[string]interface{}{
		"party_create": map[string]interface{}{
			"open":     open,
			"max_size": maxSize,
		},
	}

	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	return &Party{Open: open, MaxSize: maxSize}, nil
}

// FollowUsers sends a request to follow a list of user IDs and returns the status.
func (socket *DefaultSocket) FollowUsers(userIds []string) (*Status, error) {
	request := map[string]interface{}{
		"status_follow": map[string]interface{}{
			"user_ids": userIds,
		},
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if respStatus, ok := response["status"].(*Status); ok {
		return respStatus, nil
	}

	return nil, fmt.Errorf("invalid response format")
}

// JoinChat sends a request to join a chat and returns the joined Channel.
func (socket *DefaultSocket) JoinChat(target string, chatType int, persistence, hidden bool) (*Channel, error) {
	request := map[string]interface{}{
		"channel_join": map[string]interface{}{
			"target":      target,
			"type":        chatType,
			"persistence": persistence,
			"hidden":      hidden,
		},
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if channel, ok := response["channel"].(*Channel); ok {
		return channel, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid channel field")
}

// JoinMatch sends a request to join a match and returns the joined Match.
func (socket *DefaultSocket) JoinMatch(matchID, token *string, metadata *map[string]interface{}) (*Match, error) {
	request := map[string]interface{}{
		"match_join": map[string]interface{}{
			"metadata": metadata,
		},
	}

	if token != nil && *token != "" {
		request["match_join"].(map[string]interface{})["token"] = token
	} else {
		request["match_join"].(map[string]interface{})["match_id"] = matchID
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if match, ok := response["match"].(*Match); ok {
		return match, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid match field")
}

// JoinParty sends a request to join a party.
func (socket *DefaultSocket) JoinParty(partyID string) error {
	request := map[string]interface{}{
		"party_join": map[string]interface{}{
			"party_id": partyID,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// LeaveChat sends a request to leave a chat channel.
func (socket *DefaultSocket) LeaveChat(channelID string) error {
	request := map[string]interface{}{
		"channel_leave": map[string]interface{}{
			"channel_id": channelID,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// LeaveMatch sends a request to leave a match.
func (socket *DefaultSocket) LeaveMatch(matchID string) error {
	request := map[string]interface{}{
		"match_leave": map[string]interface{}{
			"match_id": matchID,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// LeaveParty sends a request to leave a party.
func (socket *DefaultSocket) LeaveParty(partyID string) error {
	request := map[string]interface{}{
		"party_leave": map[string]interface{}{
			"party_id": partyID,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// ListPartyJoinRequests fetches the list of join requests for a given party ID.
func (socket *DefaultSocket) ListPartyJoinRequests(partyID string) (*PartyJoinRequest, error) {
	request := map[string]interface{}{
		"party_join_request_list": map[string]interface{}{
			"party_id": partyID,
		},
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if joinRequest, ok := response["party_join_request"].(*PartyJoinRequest); ok {
		return joinRequest, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid party_join_request field")
}

// RemoveChatMessage sends a request to remove a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) RemoveChatMessage(channelID, messageID string) (*ChannelMessageAck, error) {
	request := map[string]interface{}{
		"channel_message_remove": map[string]interface{}{
			"channel_id": channelID,
			"message_id": messageID,
		},
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if messageAck, ok := response["channel_message_ack"].(*ChannelMessageAck); ok {
		return messageAck, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid channel_message_ack field")
}

// PromotePartyMember promotes a party member to party leader and returns the new PartyLeader.
func (socket *DefaultSocket) PromotePartyMember(partyID string, partyMember Presence) (*PartyLeader, error) {
	request := map[string]interface{}{
		"party_promote": map[string]interface{}{
			"party_id": partyID,
			"presence": partyMember,
		},
	}

	var response map[string]interface{}
	err := socket.Send(request, nil)
	if err != nil {
		return nil, err
	}

	if partyLeader, ok := response["party_leader"].(*PartyLeader); ok {
		return partyLeader, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid party_leader field")
}

// RemoveMatchmaker sends a request to remove a matchmaker ticket.
func (socket *DefaultSocket) RemoveMatchmaker(ticket string) error {
	request := map[string]interface{}{
		"matchmaker_remove": map[string]interface{}{
			"ticket": ticket,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// RemoveMatchmakerParty sends a request to remove a matchmaker ticket from a party.
func (socket *DefaultSocket) RemoveMatchmakerParty(partyID, ticket string) error {
	request := map[string]interface{}{
		"party_matchmaker_remove": map[string]interface{}{
			"party_id": partyID,
			"ticket":   ticket,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// RemovePartyMember sends a request to remove a member from a party.
func (socket *DefaultSocket) RemovePartyMember(partyID string, member Presence) error {
	request := map[string]interface{}{
		"party_remove": map[string]interface{}{
			"party_id": partyID,
			"presence": member,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// Rpc sends an RPC request and returns an ApiRpc response.
func (socket *DefaultSocket) Rpc(id, payload, httpKey string) (*ApiRpc, error) {
	request := map[string]interface{}{
		"rpc": map[string]interface{}{
			"id":       id,
			"payload":  payload,
			"http_key": httpKey,
		},
	}

	var response map[string]interface{}
	if err := socket.Send(request, nil); err != nil {
		return nil, err
	}

	if rpc, ok := response["rpc"].(*ApiRpc); ok {
		return rpc, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid rpc field")
}

// SendMatchState sends match state updates to the server.
func (socket *DefaultSocket) SendMatchState(matchID string, opCode int, data interface{}, presences []Presence, reliable bool) error {
	request := map[string]interface{}{
		"match_data_send": map[string]interface{}{
			"match_id":  matchID,
			"op_code":   opCode,
			"data":      data,
			"presences": presences,
			"reliable":  reliable,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// SendPartyData sends party data updates to the server.
func (socket *DefaultSocket) SendPartyData(partyID string, opCode int, data interface{}) error {
	request := map[string]interface{}{
		"party_data_send": map[string]interface{}{
			"party_id": partyID,
			"op_code":  opCode,
			"data":     data,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// UnfollowUsers sends a request to unfollow the specified users.
func (socket *DefaultSocket) UnfollowUsers(userIDs []string) error {
	request := map[string]interface{}{
		"status_unfollow": map[string]interface{}{
			"user_ids": userIDs,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// UpdateChatMessage sends a request to update a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) UpdateChatMessage(channelID, messageID string, content interface{}) (*ChannelMessageAck, error) {
	request := map[string]interface{}{
		"channel_message_update": map[string]interface{}{
			"channel_id": channelID,
			"message_id": messageID,
			"content":    content,
		},
	}

	var response map[string]interface{}
	if err := socket.Send(request, nil); err != nil {
		return nil, err
	}

	if messageAck, ok := response["channel_message_ack"].(*ChannelMessageAck); ok {
		return messageAck, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid channel_message_ack field")
}

// UpdateStatus sends a status update to the server.
func (socket *DefaultSocket) UpdateStatus(status *string) error {
	request := map[string]interface{}{
		"status_update": map[string]interface{}{
			"status": status,
		},
	}

	if err := socket.Send(request, nil); err != nil {
		return err
	}

	return nil
}

// WriteChatMessage sends a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) WriteChatMessage(channelID string, content interface{}) (*ChannelMessageAck, error) {
	request := map[string]interface{}{
		"channel_message_send": map[string]interface{}{
			"channel_id": channelID,
			"content":    content,
		},
	}

	var response map[string]interface{}
	if err := socket.Send(request, nil); err != nil {
		return nil, err
	}

	if messageAck, ok := response["channel_message_ack"].(*ChannelMessageAck); ok {
		return messageAck, nil
	}

	return nil, fmt.Errorf("invalid response format: missing or invalid channel_message_ack field")
}

// pingPong does a periodic ping-pong check with the server.
func (socket *DefaultSocket) pingPong() {
	ticker := time.NewTicker(time.Duration(socket.HeartbeatTimeoutMs) * time.Millisecond)
	defer ticker.Stop()
	log.Println("before pingpong socket is nil:", socket.Adapter.socket == nil)

	for {
		select {
		case <-ticker.C:
			ping := map[string]interface{}{"ping": struct{}{}}
			if err := socket.Send(ping, &socket.HeartbeatTimeoutMs); err != nil {
				log.Println("after pingpong socket is nil:", socket.Adapter.socket == nil)
				log.Println("Failed to send ping:", err)
				if socket.Adapter.IsOpen() {
					socket.OnHeartbeatTimeout()
					socket.Adapter.Close()
				}
				return
			}
		}
	}
}

// OnHeartbeatTimeout handles heartbeat timeouts.
func (socket *DefaultSocket) OnHeartbeatTimeout() {
	if socket.Verbose {
		fmt.Println("Heartbeat timeout")
	}
}
