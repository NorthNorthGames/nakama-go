package nakama

import (
	"fmt"
	"log"
	"time"
)

// Default configuration values
const (
	DefaultHost              = "127.0.0.1"
	DefaultPort              = "7350"
	DefaultServerKey         = "defaultkey"
	DefaultTimeoutMs         = 7000
	DefaultExpiredTimespanMs = 5 * 60 * 1000 // 5 minutes in milliseconds
)

// RpcResponse defines the response for an RPC function executed on the server.
type RpcResponse struct {
	// ID is the identifier of the function.
	ID string

	// Payload is the payload of the function, which must be a JSON object.
	Payload map[string]interface{}
}

type LeaderboardRecord struct {
	CreateTime    *string
	ExpiryTime    *string
	LeaderboardID *string
	Metadata      map[string]interface{}
	NumScore      *int
	OwnerID       *string
	Rank          *int
	Score         *int
	SubScore      *int
	UpdateTime    *string
	Username      *string
	MaxNumScore   *int
}

type LeaderboardRecordList struct {
	NextCursor   *string
	OwnerRecords []LeaderboardRecord
	PrevCursor   *string
	RankCount    *int
	Records      []LeaderboardRecord
}

type Tournament struct {
	Authoritative *bool
	ID            *string
	Title         *string
	Description   *string
	Duration      *int
	Category      *int
	SortOrder     *int
	Size          *int
	MaxSize       *int
	MaxNumScore   *int
	CanEnter      *bool
	EndActive     *int
	NextReset     *int
	Metadata      map[string]interface{}
	CreateTime    *string
	StartTime     *string
	EndTime       *string
	StartActive   *int
}

type TournamentList struct {
	Tournaments []Tournament
	Cursor      *string
}

type TournamentRecordList struct {
	NextCursor   *string
	OwnerRecords []LeaderboardRecord
	PrevCursor   *string
	Records      []LeaderboardRecord
}

type WriteTournamentRecord struct {
	Metadata map[string]interface{}
	Score    *string
	SubScore *string
}

type WriteLeaderboardRecord struct {
	Metadata map[string]interface{}
	Score    *string
	SubScore *string
}

type WriteStorageObject struct {
	Collection      *string
	Key             *string
	PermissionRead  *int
	PermissionWrite *int
	Value           map[string]interface{}
	Version         *string
}

type StorageObject struct {
	Collection      *string
	CreateTime      *string
	Key             *string
	PermissionRead  *int
	PermissionWrite *int
	UpdateTime      *string
	UserID          *string
	Value           map[string]interface{}
	Version         *string
}

type StorageObjectList struct {
	Cursor  *string
	Objects []StorageObject
}

type StorageObjects struct {
	Objects []StorageObject
}

type ChannelMessage struct {
	ChannelID   *string
	Code        *int
	Content     map[string]interface{}
	CreateTime  *string
	GroupID     *string
	MessageID   *string
	Persistent  *bool
	RoomName    *string
	ReferenceID *string
	SenderID    *string
	UpdateTime  *string
	UserIDOne   *string
	UserIDTwo   *string
	Username    *string
}

type ChannelMessageList struct {
	CacheableCursor *string
	Messages        []ChannelMessage
	NextCursor      *string
	PrevCursor      *string
}

type User struct {
	AvatarURL             *string
	CreateTime            *string
	DisplayName           *string
	EdgeCount             *int
	FacebookID            *string
	FacebookInstantGameID *string
	GamecenterID          *string
	GoogleID              *string
	ID                    *string
	LangTag               *string
	Location              *string
	Metadata              map[string]interface{}
	Online                *bool
	SteamID               *string
	Timezone              *string
	UpdateTime            *string
	Username              *string
}

type Users struct {
	Users []User
}

type Friend struct {
	State *int
	User  *User
}

type Friends struct {
	Friends []Friend
	Cursor  *string
}

type FriendOfFriend struct {
	Referrer *string
	User     *User
}

type FriendsOfFriends struct {
	Cursor           *string
	FriendsOfFriends []FriendOfFriend
}

type GroupUser struct {
	User  *User
	State *int
}

type GroupUserList struct {
	GroupUsers []GroupUser
	Cursor     *string
}

type Group struct {
	AvatarURL   *string
	CreateTime  *string
	CreatorID   *string
	Description *string
	EdgeCount   *int
	ID          *string
	LangTag     *string
	MaxCount    *int
	Metadata    map[string]interface{}
	Name        *string
	Open        *bool
	UpdateTime  *string
}

type GroupList struct {
	Cursor *string
	Groups []Group
}

type UserGroup struct {
	Group *Group
	State *int
}

type UserGroupList struct {
	UserGroups []UserGroup
	Cursor     *string
}

type Notification struct {
	Code       *int
	Content    map[string]interface{}
	CreateTime *string
	ID         *string
	Persistent *bool
	SenderID   *string
	Subject    *string
}

type NotificationList struct {
	CacheableCursor *string
	Notifications   []Notification
}

type ValidatedSubscription struct {
	Active                *bool
	CreateTime            *string
	Environment           *string
	ExpiryTime            *string
	OriginalTransactionID *string
	ProductID             *string
	ProviderNotification  *string
	ProviderResponse      *string
	PurchaseTime          *string
	RefundTime            *string
	Store                 *string
	UpdateTime            *string
	UserID                *string
}

type SubscriptionList struct {
	Cursor                 *string
	PrevCursor             *string
	ValidatedSubscriptions []ValidatedSubscription
}

// Client represents a client for the Nakama server.
type Client struct {
	ExpiredTimespanMs  int64      // The expired timespan used to check session lifetime.
	ApiClient          *NakamaApi // The low-level API client for Nakama server.
	ServerKey          string
	Host               string
	Port               string
	UseSSL             bool
	Timeout            int
	AutoRefreshSession bool
}

// NewClient creates a new instance of Client with the specified configuration.
func NewClient(
	serverKey string,
	host string,
	port string,
	useSSL bool,
	timeout int,
	autoRefreshSession bool,
) *Client {
	// Default values if not provided
	if serverKey == "" {
		serverKey = DefaultServerKey
	}
	if host == "" {
		host = DefaultHost
	}
	if port == "" {
		port = DefaultPort
	}
	if timeout == 0 {
		timeout = DefaultTimeoutMs
	}

	scheme := "http://"
	if useSSL {
		scheme = "https://"
	}
	basePath := scheme + host + ":" + port

	return &Client{
		ExpiredTimespanMs:  DefaultExpiredTimespanMs,
		ApiClient:          &NakamaApi{serverKey, basePath, timeout},
		ServerKey:          serverKey,
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Timeout:            timeout,
		AutoRefreshSession: autoRefreshSession,
	}
}

// AddGroupUsers adds users to a group, or accepts their join requests.
func (c *Client) AddGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().UnixMilli()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.RefreshSession(session, nil); err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.AddGroupUsers(session.Token, groupId, ids, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// AddFriends adds friends by ID or username to a user's account.
func (c *Client) AddFriends(session *Session, ids []string, usernames []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().UnixMilli()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.RefreshSession(session, nil); err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.AddFriends(session.Token, ids, usernames, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// AuthenticateApple authenticates a user with an Apple ID against the server.
func (c *Client) AuthenticateApple(token string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountApple{
		Token: token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Apple
	apiSession, err := c.ApiClient.AuthenticateApple(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateCustom authenticates a user with a custom ID against the server.
func (c *Client) AuthenticateCustom(id string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountCustom{
		ID:   id,
		Vars: vars,
	}

	// Call the API client to authenticate with a custom ID
	apiSession, err := c.ApiClient.AuthenticateCustom(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateDevice authenticates a user with a device ID against the server.
func (c *Client) AuthenticateDevice(id string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountDevice{
		ID:   id,
		Vars: vars,
	}

	// Call the API client to authenticate with a device ID
	apiSession, err := c.ApiClient.AuthenticateDevice(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// RefreshSession refreshes a user's session using a refresh token retrieved from a previous authentication request.
func (c *Client) RefreshSession(session *Session, vars map[string]string) (*Session, error) {
	if session == nil {
		return nil, fmt.Errorf("cannot refresh a null session")
	}

	if session.ExpiresAt != nil && *session.ExpiresAt-session.CreatedAt < 70 {
		log.Println("Session lifetime too short, please set '--session.token_expiry_sec' option. See the documentation for more info: https://heroiclabs.com/docs/nakama/getting-started/configuration/#session")
	}

	if session.RefreshExpiresAt != nil && *session.RefreshExpiresAt-session.CreatedAt < 3700 {
		log.Println("Session refresh lifetime too short, please set '--session.refresh_token_expiry_sec' option. See the documentation for more info: https://heroiclabs.com/docs/nakama/getting-started/configuration/#session")
	}

	apiSession, err := c.ApiClient.SessionRefresh(c.ServerKey, "", ApiSessionRefreshRequest{
		Token: session.RefreshToken,
		Vars:  vars,
	}, make(map[string]string))

	if err != nil {
		return nil, err
	}

	session.Update(apiSession.Token, apiSession.RefreshToken)
	return session, nil
}
