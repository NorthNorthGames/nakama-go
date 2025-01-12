package nakama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	CreateTime    *string                `json:"create_time,omitempty"`
	ExpiryTime    *string                `json:"expiry_time,omitempty"`
	LeaderboardID *string                `json:"leaderboard_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	NumScore      *int                   `json:"num_score,omitempty"`
	OwnerID       *string                `json:"owner_id,omitempty"`
	Rank          *int                   `json:"rank,omitempty"`
	Score         *int                   `json:"score,omitempty"`
	SubScore      *int                   `json:"subscore,omitempty"`
	UpdateTime    *string                `json:"update_time,omitempty"`
	Username      *string                `json:"username,omitempty"`
	MaxNumScore   *int                   `json:"max_num_score,omitempty"`
}

type LeaderboardRecordList struct {
	NextCursor   *string             `json:"next_cursor,omitempty"`
	OwnerRecords []LeaderboardRecord `json:"owner_records,omitempty"`
	PrevCursor   *string             `json:"prev_cursor,omitempty"`
	RankCount    *int                `json:"rank_count,omitempty"`
	Records      []LeaderboardRecord `json:"records,omitempty"`
}

type Tournament struct {
	Authoritative *bool                  `json:"authoritative,omitempty"`
	ID            *string                `json:"id,omitempty"`
	Title         *string                `json:"title,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Duration      *int                   `json:"duration,omitempty"`
	Category      *int                   `json:"category,omitempty"`
	SortOrder     *int                   `json:"sort_order,omitempty"`
	Size          *int                   `json:"size,omitempty"`
	MaxSize       *int                   `json:"max_size,omitempty"`
	MaxNumScore   *int                   `json:"max_num_score,omitempty"`
	CanEnter      *bool                  `json:"can_enter,omitempty"`
	EndActive     *int                   `json:"end_active,omitempty"`
	NextReset     *int                   `json:"next_reset,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreateTime    *string                `json:"create_time,omitempty"`
	StartTime     *string                `json:"start_time,omitempty"`
	EndTime       *string                `json:"end_time,omitempty"`
	StartActive   *int                   `json:"start_active,omitempty"`
}

type TournamentList struct {
	Tournaments []Tournament `json:"tournaments,omitempty"`
	Cursor      *string      `json:"cursor,omitempty"`
}

type TournamentRecordList struct {
	NextCursor   *string             `json:"next_cursor,omitempty"`
	OwnerRecords []LeaderboardRecord `json:"owner_records,omitempty"`
	PrevCursor   *string             `json:"prev_cursor,omitempty"`
	Records      []LeaderboardRecord `json:"records,omitempty"`
}

type WriteTournamentRecord struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    *string                `json:"score,omitempty"`
	SubScore *string                `json:"subscore,omitempty"`
}

type WriteLeaderboardRecord struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    *string                `json:"score,omitempty"`
	SubScore *string                `json:"subscore,omitempty"`
}

type WriteStorageObject struct {
	Collection      *string                `json:"collection,omitempty"`
	Key             *string                `json:"key,omitempty"`
	PermissionRead  *int                   `json:"permission_read,omitempty"`
	PermissionWrite *int                   `json:"permission_write,omitempty"`
	Value           map[string]interface{} `json:"value,omitempty"`
	Version         *string                `json:"version,omitempty"`
}

type StorageObject struct {
	Collection      *string                `json:"collection,omitempty"`
	CreateTime      *string                `json:"create_time,omitempty"`
	Key             *string                `json:"key,omitempty"`
	PermissionRead  *int                   `json:"permission_read,omitempty"`
	PermissionWrite *int                   `json:"permission_write,omitempty"`
	UpdateTime      *string                `json:"update_time,omitempty"`
	UserID          *string                `json:"user_id,omitempty"`
	Value           map[string]interface{} `json:"value,omitempty"`
	Version         *string                `json:"version,omitempty"`
}

type StorageObjectList struct {
	Cursor  *string         `json:"cursor,omitempty"`
	Objects []StorageObject `json:"objects"`
}

type StorageObjects struct {
	Objects []StorageObject `json:"objects"`
}

type ChannelMessage struct {
	ChannelID   *string                `json:"channel_id,omitempty"`
	Code        *int                   `json:"code,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	CreateTime  *string                `json:"create_time,omitempty"`
	GroupID     *string                `json:"group_id,omitempty"`
	MessageID   *string                `json:"message_id,omitempty"`
	Persistent  *bool                  `json:"persistent,omitempty"`
	RoomName    *string                `json:"room_name,omitempty"`
	ReferenceID *string                `json:"reference_id,omitempty"`
	SenderID    *string                `json:"sender_id,omitempty"`
	UpdateTime  *string                `json:"update_time,omitempty"`
	UserIDOne   *string                `json:"user_id_one,omitempty"`
	UserIDTwo   *string                `json:"user_id_two,omitempty"`
	Username    *string                `json:"username,omitempty"`
}

type ChannelMessageList struct {
	CacheableCursor *string          `json:"cacheable_cursor,omitempty"`
	Messages        []ChannelMessage `json:"messages,omitempty"`
	NextCursor      *string          `json:"next_cursor,omitempty"`
	PrevCursor      *string          `json:"prev_cursor,omitempty"`
}

type User struct {
	AvatarURL             *string                `json:"avatar_url,omitempty"`
	CreateTime            *string                `json:"create_time,omitempty"`
	DisplayName           *string                `json:"display_name,omitempty"`
	EdgeCount             *int                   `json:"edge_count,omitempty"`
	FacebookID            *string                `json:"facebook_id,omitempty"`
	FacebookInstantGameID *string                `json:"facebook_instant_game_id,omitempty"`
	GameCenterID          *string                `json:"gamecenter_id,omitempty"`
	GoogleID              *string                `json:"google_id,omitempty"`
	ID                    *string                `json:"id,omitempty"`
	LangTag               *string                `json:"lang_tag,omitempty"`
	Location              *string                `json:"location,omitempty"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
	Online                *bool                  `json:"online,omitempty"`
	SteamID               *string                `json:"steam_id,omitempty"`
	Timezone              *string                `json:"timezone,omitempty"`
	UpdateTime            *string                `json:"update_time,omitempty"`
	Username              *string                `json:"username,omitempty"`
}

type Users struct {
	Users []User `json:"users,omitempty"`
}

type Friend struct {
	State *int  `json:"state,omitempty"`
	User  *User `json:"user,omitempty"`
}

type Friends struct {
	Friends []Friend `json:"friends,omitempty"`
	Cursor  *string  `json:"cursor,omitempty"`
}

type FriendOfFriend struct {
	Referrer *string `json:"referrer,omitempty"`
	User     *User   `json:"user,omitempty"`
}

type FriendsOfFriends struct {
	Cursor           *string          `json:"cursor,omitempty"`
	FriendsOfFriends []FriendOfFriend `json:"friends_of_friends,omitempty"`
}

type GroupUser struct {
	User  *User `json:"user,omitempty"`
	State *int  `json:"state,omitempty"`
}

type GroupUserList struct {
	GroupUsers []GroupUser `json:"group_users,omitempty"`
	Cursor     *string     `json:"cursor,omitempty"`
}

type Group struct {
	AvatarURL   *string                `json:"avatar_url,omitempty"`
	CreateTime  *string                `json:"create_time,omitempty"`
	CreatorID   *string                `json:"creator_id,omitempty"`
	Description *string                `json:"description,omitempty"`
	EdgeCount   *int                   `json:"edge_count,omitempty"`
	ID          *string                `json:"id,omitempty"`
	LangTag     *string                `json:"lang_tag,omitempty"`
	MaxCount    *int                   `json:"max_count,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Open        *bool                  `json:"open,omitempty"`
	UpdateTime  *string                `json:"update_time,omitempty"`
}

type GroupList struct {
	Cursor *string `json:"cursor,omitempty"`
	Groups []Group `json:"groups,omitempty"`
}

type UserGroup struct {
	Group *Group `json:"group,omitempty"`
	State *int   `json:"state,omitempty"`
}

type UserGroupList struct {
	UserGroups []UserGroup `json:"user_groups,omitempty"`
	Cursor     *string     `json:"cursor,omitempty"`
}

type Notification struct {
	Code       *int                   `json:"code,omitempty"`
	Content    map[string]interface{} `json:"content,omitempty"`
	CreateTime *string                `json:"create_time,omitempty"`
	ID         *string                `json:"id,omitempty"`
	Persistent *bool                  `json:"persistent,omitempty"`
	SenderID   *string                `json:"sender_id,omitempty"`
	Subject    *string                `json:"subject,omitempty"`
}

type NotificationList struct {
	CacheableCursor *string        `json:"cacheable_cursor,omitempty"`
	Notifications   []Notification `json:"notifications,omitempty"`
}

type ValidatedSubscription struct {
	Active                *bool   `json:"active,omitempty"`
	CreateTime            *string `json:"create_time,omitempty"`
	Environment           *string `json:"environment,omitempty"`
	ExpiryTime            *string `json:"expiry_time,omitempty"`
	OriginalTransactionID *string `json:"original_transaction_id,omitempty"`
	ProductID             *string `json:"product_id,omitempty"`
	ProviderNotification  *string `json:"provider_notification,omitempty"`
	ProviderResponse      *string `json:"provider_response,omitempty"`
	PurchaseTime          *string `json:"purchase_time,omitempty"`
	RefundTime            *string `json:"refund_time,omitempty"`
	Store                 *string `json:"store,omitempty"`
	UpdateTime            *string `json:"update_time,omitempty"`
	UserID                *string `json:"user_id,omitempty"`
}

type SubscriptionList struct {
	Cursor                 *string                 `json:"cursor,omitempty"`
	PrevCursor             *string                 `json:"prev_cursor,omitempty"`
	ValidatedSubscriptions []ValidatedSubscription `json:"validated_subscriptions,omitempty"`
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
	timeout *int,
	autoRefreshSession *bool,
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
	if timeout == nil {
		timeout = new(int)
		*timeout = DefaultTimeoutMs
	}
	if autoRefreshSession == nil {
		autoRefreshSession = new(bool)
		*autoRefreshSession = true
	}

	scheme := "http://"
	if useSSL {
		scheme = "https://"
	}
	basePath := scheme + host + ":" + port

	return &Client{
		ExpiredTimespanMs:  DefaultExpiredTimespanMs,
		ApiClient:          &NakamaApi{serverKey, basePath, *timeout},
		ServerKey:          serverKey,
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Timeout:            *timeout,
		AutoRefreshSession: *autoRefreshSession,
	}
}

// AddGroupUsers adds users to a group, or accepts their join requests.
func (c *Client) AddGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().UnixMilli()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
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
		if _, err := c.SessionRefresh(session, nil); err != nil {
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
		Token: &token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Apple
	apiSession, err := c.ApiClient.AuthenticateApple(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateCustom authenticates a user with a custom ID against the server.
func (c *Client) AuthenticateCustom(id string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountCustom{
		ID:   &id,
		Vars: vars,
	}

	// Call the API client to authenticate with a custom ID
	apiSession, err := c.ApiClient.AuthenticateCustom(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateDevice authenticates a user with a device ID against the server.
func (c *Client) AuthenticateDevice(id string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountDevice{
		ID:   &id,
		Vars: vars,
	}

	// Call the API client to authenticate with a device ID
	apiSession, err := c.ApiClient.AuthenticateDevice(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	created := false
	if apiSession.Created != nil {
		created = *apiSession.Created
	} else if create != nil {
		created = *create
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      created,
	}, nil
}

// AuthenticateEmail authenticates a user with an email and password against the server.
func (c *Client) AuthenticateEmail(email string, password string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountEmail{
		Email:    &email,
		Password: &password,
		Vars:     vars,
	}

	// Call the API client to authenticate with email and password
	apiSession, err := c.ApiClient.AuthenticateEmail(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateFacebookInstantGame authenticates a user with a Facebook Instant Game token against the server.
func (c *Client) AuthenticateFacebookInstantGame(signedPlayerInfo string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountFacebookInstantGame{
		SignedPlayerInfo: &signedPlayerInfo,
		Vars:             vars,
	}

	// Call the API client to authenticate with Facebook Instant Game
	apiSession, err := c.ApiClient.AuthenticateFacebookInstantGame(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateFacebook authenticates a user with a Facebook OAuth token against the server.
func (c *Client) AuthenticateFacebook(token string, create *bool, username *string, sync *bool, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountFacebook{
		Token: &token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Facebook
	apiSession, err := c.ApiClient.AuthenticateFacebook(c.ServerKey, "", request, create, username, sync, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateGoogle authenticates a user with a Google token against the server.
func (c *Client) AuthenticateGoogle(token string, create *bool, username *string, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountGoogle{
		Token: &token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Google
	apiSession, err := c.ApiClient.AuthenticateGoogle(c.ServerKey, "", request, create, username, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateGameCenter authenticates a user with GameCenter against the server.
func (c *Client) AuthenticateGameCenter(bundleId string, playerId string, publicKeyUrl string, salt string, signature string, timestamp string, create *bool, username *string, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountGameCenter{
		BundleID:     &bundleId,
		PlayerID:     &playerId,
		PublicKeyURL: &publicKeyUrl,
		Salt:         &salt,
		Signature:    &signature,
		Timestamp:    &timestamp,
		Vars:         vars,
	}

	// Call the API client to authenticate with GameCenter
	apiSession, err := c.ApiClient.AuthenticateGameCenter(c.ServerKey, "", request, create, username, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// AuthenticateSteam authenticates a user with a Steam token against the server.
func (c *Client) AuthenticateSteam(token string, create *bool, username *string, sync *bool, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := ApiAccountSteam{
		Token: &token,
		Vars:  vars,
		Sync:  sync,
	}

	// Call the API client to authenticate with Steam
	apiSession, err := c.ApiClient.AuthenticateSteam(c.ServerKey, "", request, create, username, nil, make(map[string]string))

	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        *apiSession.Token,
		RefreshToken: *apiSession.RefreshToken,
		Created:      *apiSession.Created,
	}, nil
}

// BanGroupUsers bans users from a group.
func (c *Client) BanGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.BanGroupUsers(session.Token, groupId, ids, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// BlockFriends blocks one or more users by ID or username.
func (c *Client) BlockFriends(session *Session, ids []string, usernames []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.BlockFriends(session.Token, ids, usernames, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// CreateGroup creates a new group with the current user as the creator and superadmin.
func (c *Client) CreateGroup(session *Session, request ApiCreateGroupRequest) (*Group, error) {
	// Check if the session requires refresh
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	// Call the API client to create the group
	apiGroup, err := c.ApiClient.CreateGroup(session.Token, request, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Map the response to the Group struct
	return &Group{
		AvatarURL:   apiGroup.AvatarURL,
		CreateTime:  timeToStringPointer(*apiGroup.CreateTime, time.RFC3339),
		CreatorID:   apiGroup.CreatorID,
		Description: apiGroup.Description,
		EdgeCount:   apiGroup.EdgeCount,
		ID:          apiGroup.ID,
		LangTag:     apiGroup.LangTag,
		MaxCount:    apiGroup.MaxCount,
		Metadata: func() map[string]interface{} {
			var metadata map[string]interface{}
			json.Unmarshal([]byte(*apiGroup.Metadata), &metadata)
			return metadata
		}(),
		Name:       apiGroup.Name,
		Open:       apiGroup.Open,
		UpdateTime: timeToStringPointer(*apiGroup.UpdateTime, time.RFC3339),
	}, nil
}

// CreateSocket creates a socket using the client's configuration.
func (c *Client) CreateSocket(useSSL bool, verbose bool, adapter *WebSocketAdapter, sendTimeoutMs *int) DefaultSocket {
	return NewDefaultSocket(c.Host, c.Port, useSSL, verbose, *adapter, sendTimeoutMs)
}

// DeleteAccount deletes the current user's account.
func (c *Client) DeleteAccount(session *Session) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteAccount(session.Token, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DeleteFriends deletes one or more users by ID or username.
func (c *Client) DeleteFriends(session *Session, ids []string, usernames []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteFriends(session.Token, ids, usernames, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DeleteGroup deletes a group the user is part of and has permissions to delete.
func (c *Client) DeleteGroup(session *Session, groupId string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteGroup(session.Token, groupId, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DeleteNotifications deletes one or more notifications.
func (c *Client) DeleteNotifications(session *Session, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteNotifications(session.Token, ids, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DeleteStorageObjects deletes one or more storage objects.
func (c *Client) DeleteStorageObjects(session *Session, request ApiDeleteStorageObjectsRequest) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteStorageObjects(session.Token, request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DeleteTournamentRecord deletes a tournament record.
func (c *Client) DeleteTournamentRecord(session *Session, tournamentId string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DeleteTournamentRecord(session.Token, tournamentId, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// DemoteGroupUsers demotes a set of users in a group to the next role down.
func (c *Client) DemoteGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.DemoteGroupUsers(session.Token, groupId, ids, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// EmitEvent submits an event for processing in the server's registered runtime custom events handler.
func (c *Client) EmitEvent(session *Session, request ApiEvent) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.Event(session.Token, request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// GetAccount fetches the current user's account.
func (c *Client) GetAccount(session *Session) (*ApiAccount, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	account, err := c.ApiClient.GetAccount(session.Token, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetSubscription fetches a subscription by product ID.
func (c *Client) GetSubscription(session *Session, productId string) (*ApiValidatedSubscription, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	subscription, err := c.ApiClient.GetSubscription(session.Token, productId, make(map[string]string))
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// ImportFacebookFriends imports Facebook friends and adds them to a user's account.
func (c *Client) ImportFacebookFriends(session *Session, request ApiAccountFacebook) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.ImportFacebookFriends(session.Token, request, false, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// ImportSteamFriends imports Steam friends and adds them to a user's account.
func (c *Client) ImportSteamFriends(session *Session, request ApiAccountSteam, reset bool) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.ImportSteamFriends(session.Token, request, reset, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// FetchUsers fetches zero or more users by ID and/or username.
func (c *Client) FetchUsers(session *Session, ids []string, usernames []string, facebookIds []string) (*Users, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.ApiClient.GetUsers(session.Token, ids, usernames, facebookIds, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &Users{
		Users: []User{},
	}

	if apiResponse.Users == nil {
		return result, nil
	}

	for _, u := range *apiResponse.Users {
		user := User{
			AvatarURL:    u.AvatarURL,
			CreateTime:   timeToStringPointer(*u.CreateTime, time.RFC3339),
			DisplayName:  u.DisplayName,
			EdgeCount:    u.EdgeCount,
			FacebookID:   u.FacebookID,
			GameCenterID: u.GameCenterID,
			GoogleID:     u.GoogleID,
			ID:           u.ID,
			LangTag:      u.LangTag,
			Location:     u.Location,
			Online:       u.Online,
			SteamID:      u.SteamID,
			Timezone:     u.Timezone,
			UpdateTime:   timeToStringPointer(*u.UpdateTime, time.RFC3339),
			Username:     u.Username,
			Metadata:     nil,
		}
		if u.Metadata != nil {
			if err := json.Unmarshal([]byte(*u.Metadata), &user.Metadata); err != nil {
				return nil, err
			}
		}
		result.Users = append(result.Users, user)
	}

	return result, nil
}

// JoinGroup either joins a group that's open or sends a request to join a group that's closed.
func (c *Client) JoinGroup(session *Session, groupId string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.JoinGroup(session.Token, groupId, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// JoinTournament allows a user to join a tournament by its ID.
func (c *Client) JoinTournament(session *Session, tournamentId string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.JoinTournament(session.Token, tournamentId, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// KickGroupUsers kicks users from a group or declines their join requests.
func (c *Client) KickGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.KickGroupUsers(session.Token, groupId, ids, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LeaveGroup allows a user to leave a group they are part of.
func (c *Client) LeaveGroup(session *Session, groupId string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LeaveGroup(session.Token, groupId, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// ListChannelMessages retrieves a channel's message history.
func (c *Client) ListChannelMessages(session *Session, channelId string, limit *int, forward *bool, cursor *string) (*ChannelMessageList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.ApiClient.ListChannelMessages(session.Token, channelId, limit, forward, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &ChannelMessageList{
		Messages:        []ChannelMessage{},
		NextCursor:      apiResponse.NextCursor,
		PrevCursor:      apiResponse.PrevCursor,
		CacheableCursor: apiResponse.CacheableCursor,
	}

	if apiResponse.Messages == nil {
		return result, nil
	}

	for _, m := range apiResponse.Messages {
		message := ChannelMessage{
			ChannelID:  m.ChannelID,
			Code:       m.Code,
			CreateTime: timeToStringPointer(*m.CreateTime, time.RFC3339),
			MessageID:  m.MessageID,
			Persistent: m.Persistent,
			SenderID:   m.SenderID,
			UpdateTime: timeToStringPointer(*m.UpdateTime, time.RFC3339),
			Username:   m.Username,
			Content:    nil,
			GroupID:    m.GroupID,
			RoomName:   m.RoomName,
			UserIDOne:  m.UserIDOne,
			UserIDTwo:  m.UserIDTwo,
		}
		if m.Content != nil {
			if err := json.Unmarshal([]byte(*m.Content), &message.Content); err != nil {
				return nil, err
			}
		}

		result.Messages = append(result.Messages, message)
	}

	return result, nil
}

// ListGroupUsers retrieves a group's users with optional state, limit, and cursor parameters.
func (c *Client) ListGroupUsers(session *Session, groupId string, state *int, limit *int, cursor *string) (*GroupUserList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.ApiClient.ListGroupUsers(session.Token, groupId, state, limit, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &GroupUserList{
		GroupUsers: []GroupUser{},
		Cursor:     apiResponse.Cursor,
	}

	if apiResponse.GroupUsers == nil {
		return result, nil
	}

	for _, gu := range apiResponse.GroupUsers {
		groupUser := GroupUser{
			User: &User{
				AvatarURL:    gu.User.AvatarURL,
				CreateTime:   timeToStringPointer(*gu.User.CreateTime, time.RFC3339),
				DisplayName:  gu.User.DisplayName,
				EdgeCount:    gu.User.EdgeCount,
				FacebookID:   gu.User.FacebookID,
				GameCenterID: gu.User.GameCenterID,
				GoogleID:     gu.User.GoogleID,
				ID:           gu.User.ID,
				LangTag:      gu.User.LangTag,
				Location:     gu.User.Location,
				Online:       gu.User.Online,
				SteamID:      gu.User.SteamID,
				Timezone:     gu.User.Timezone,
				UpdateTime:   timeToStringPointer(*gu.User.UpdateTime, time.RFC3339),
				Username:     gu.User.Username,
				Metadata:     nil,
			},
			State: gu.State,
		}

		if gu.User.Metadata != nil {
			if err := json.Unmarshal([]byte(*gu.User.Metadata), &groupUser.User.Metadata); err != nil {
				return nil, err
			}
		}

		result.GroupUsers = append(result.GroupUsers, groupUser)
	}

	return result, nil
}

// ListUserGroups lists a user's groups.
func (c *Client) ListUserGroups(session *Session, userId string, state *int, limit *int, cursor *string) (*UserGroupList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.ApiClient.ListUserGroups(session.Token, userId, state, limit, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &UserGroupList{
		UserGroups: []UserGroup{},
		Cursor:     apiResponse.Cursor,
	}

	if apiResponse.UserGroups == nil {
		return result, nil
	}

	for _, ug := range *apiResponse.UserGroups {
		userGroup := UserGroup{
			Group: &Group{
				AvatarURL:   ug.Group.AvatarURL,
				CreateTime:  timeToStringPointer(*ug.Group.CreateTime, time.RFC3339),
				CreatorID:   ug.Group.CreatorID,
				Description: ug.Group.Description,
				EdgeCount:   ug.Group.EdgeCount,
				ID:          ug.Group.ID,
				LangTag:     ug.Group.LangTag,
				MaxCount:    ug.Group.MaxCount,
				Metadata:    nil,
				Name:        ug.Group.Name,
				Open:        ug.Group.Open,
				UpdateTime:  timeToStringPointer(*ug.Group.UpdateTime, time.RFC3339),
			},
			State: ug.State,
		}

		if ug.Group.Metadata != nil {
			if err := json.Unmarshal([]byte(*ug.Group.Metadata), &userGroup.Group.Metadata); err != nil {
				return nil, err
			}
		}

		result.UserGroups = append(result.UserGroups, userGroup)
	}

	return result, nil
}

// ListGroups retrieves a list of groups based on the given filters.
func (c *Client) ListGroups(session *Session, name *string, cursor *string, limit *int) (*GroupList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.ApiClient.ListGroups(session.Token, name, cursor, limit, nil, nil, nil, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &GroupList{
		Groups: []Group{},
		Cursor: apiResponse.Cursor,
	}

	if apiResponse.Groups == nil {
		return result, nil
	}

	for _, ug := range apiResponse.Groups {
		group := Group{
			AvatarURL:   ug.AvatarURL,
			CreateTime:  timeToStringPointer(*ug.CreateTime, time.RFC3339),
			CreatorID:   ug.CreatorID,
			Description: ug.Description,
			EdgeCount:   nil,
			ID:          ug.ID,
			LangTag:     ug.LangTag,
			MaxCount:    ug.MaxCount,
			Metadata:    nil,
			Name:        ug.Name,
			Open:        ug.Open,
			UpdateTime:  timeToStringPointer(*ug.UpdateTime, time.RFC3339),
		}

		// Convert optional fields
		if ug.EdgeCount != nil {
			group.EdgeCount = ug.EdgeCount
		}
		if ug.Metadata != nil {
			if err := json.Unmarshal([]byte(*ug.Metadata), &group.Metadata); err != nil {
				return nil, err
			}
		}

		result.Groups = append(result.Groups, group)
	}

	return result, nil
}

// LinkApple adds an Apple ID to the social profiles on the current user's account.
func (c *Client) LinkApple(session *Session, request *ApiAccountApple) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkApple(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkCustom adds a custom ID to the social profiles on the current user's account.
func (c *Client) LinkCustom(session *Session, request *ApiAccountCustom) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkCustom(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkDevice adds a device ID to the social profiles on the current user's account.
func (c *Client) LinkDevice(session *Session, request *ApiAccountDevice) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkDevice(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkEmail adds an email and password to the social profiles on the current user's account.
func (c *Client) LinkEmail(session *Session, request *ApiAccountEmail) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkEmail(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkFacebook adds a Facebook ID to the social profiles on the current user's account.
func (c *Client) LinkFacebook(session *Session, request *ApiAccountFacebook) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkFacebook(session.Token, *request, nil, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkFacebookInstant adds Facebook Instant to the social profiles on the current user's account.
func (c *Client) LinkFacebookInstant(session *Session, request *ApiAccountFacebookInstantGame) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkFacebookInstantGame(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkGoogle adds a Google account to the social profiles on the current user's account.
func (c *Client) LinkGoogle(session *Session, request *ApiAccountGoogle) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkGoogle(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkGameCenter adds GameCenter to the social profiles on the current user's account.
func (c *Client) LinkGameCenter(session *Session, request *ApiAccountGameCenter) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkGameCenter(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// LinkSteam adds Steam to the social profiles on the current user's account.
func (c *Client) LinkSteam(session *Session, request *ApiLinkSteamRequest) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return false, err
		}
	}

	response, err := c.ApiClient.LinkSteam(session.Token, *request, make(map[string]string))
	if err != nil {
		return false, err
	}

	return response != nil, nil
}

// ListFriends lists all friends for the current user.
func (c *Client) ListFriends(session *Session, state *int, limit *int, cursor *string) (*Friends, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListFriends(session.Token, limit, state, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &Friends{
		Friends: []Friend{},
		Cursor:  response.Cursor,
	}

	if response.Friends == nil {
		return result, nil
	}

	for _, f := range response.Friends {
		friend := Friend{
			User: &User{
				AvatarURL:             f.User.AvatarURL,
				CreateTime:            timeToStringPointer(*f.User.CreateTime, time.RFC3339),
				DisplayName:           f.User.DisplayName,
				EdgeCount:             f.User.EdgeCount,
				FacebookID:            f.User.FacebookID,
				GameCenterID:          f.User.GameCenterID,
				GoogleID:              f.User.GoogleID,
				ID:                    f.User.ID,
				LangTag:               f.User.LangTag,
				Location:              f.User.Location,
				Online:                f.User.Online,
				SteamID:               f.User.SteamID,
				Timezone:              f.User.Timezone,
				UpdateTime:            timeToStringPointer(*f.User.UpdateTime, time.RFC3339),
				Username:              f.User.Username,
				Metadata:              nil,
				FacebookInstantGameID: f.User.FacebookInstantGameID,
			},
			State: f.State,
		}

		if f.User.Metadata != nil {
			if err := json.Unmarshal([]byte(*f.User.Metadata), &friend.User.Metadata); err != nil {
				return nil, err
			}
		}

		result.Friends = append(result.Friends, friend)
	}

	return result, nil
}

// ListFriendsOfFriends lists the friends of friends for the current user.
func (c *Client) ListFriendsOfFriends(session *Session, limit *int, cursor *string) (*FriendsOfFriends, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListFriendsOfFriends(session.Token, limit, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &FriendsOfFriends{
		FriendsOfFriends: []FriendOfFriend{},
		Cursor:           response.Cursor,
	}

	if response.FriendsOfFriends == nil {
		return result, nil
	}

	for _, f := range response.FriendsOfFriends {
		friendOfFriend := FriendOfFriend{
			Referrer: f.Referrer,
			User: &User{
				AvatarURL:             f.User.AvatarURL,
				CreateTime:            timeToStringPointer(*f.User.CreateTime, time.RFC3339),
				DisplayName:           f.User.DisplayName,
				EdgeCount:             f.User.EdgeCount,
				FacebookID:            f.User.FacebookID,
				GameCenterID:          f.User.GameCenterID,
				GoogleID:              f.User.GoogleID,
				ID:                    f.User.ID,
				LangTag:               f.User.LangTag,
				Location:              f.User.Location,
				Online:                f.User.Online,
				SteamID:               f.User.SteamID,
				Timezone:              f.User.Timezone,
				UpdateTime:            timeToStringPointer(*f.User.UpdateTime, time.RFC3339),
				Username:              f.User.Username,
				Metadata:              nil,
				FacebookInstantGameID: f.User.FacebookInstantGameID,
			},
		}

		if f.User.Metadata != nil {
			if err := json.Unmarshal([]byte(*f.User.Metadata), &friendOfFriend.User.Metadata); err != nil {
				return nil, err
			}
		}

		result.FriendsOfFriends = append(result.FriendsOfFriends, friendOfFriend)
	}

	return result, nil
}

// ListLeaderboardRecords lists the leaderboard records with optional ownerIds, pagination, and expiry filters.
func (c *Client) ListLeaderboardRecords(session *Session, leaderboardId string, ownerIds []string, limit *int, cursor *string, expiry *string) (*LeaderboardRecordList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListLeaderboardRecords(session.Token, leaderboardId, ownerIds, limit, cursor, expiry, make(map[string]string))
	if err != nil {
		return nil, err
	}

	list := &LeaderboardRecordList{
		NextCursor:   response.NextCursor,
		PrevCursor:   response.PrevCursor,
		RankCount:    stringPointerToIntPointer(response.RankCount),
		OwnerRecords: []LeaderboardRecord{},
		Records:      []LeaderboardRecord{},
	}

	if response.OwnerRecords != nil {
		for _, o := range response.OwnerRecords {
			metadata := map[string]interface{}{}
			if o.Metadata != nil {
				if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err != nil {
					return nil, err
				}
			}

			list.OwnerRecords = append(list.OwnerRecords, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata:      metadata,
				NumScore:      o.NumScore,
				OwnerID:       o.OwnerID,
				Rank:          stringPointerToIntPointer(o.Rank),
				Score:         stringPointerToIntPointer(o.Score),
				SubScore:      stringPointerToIntPointer(o.Subscore),
				UpdateTime:    timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:      o.Username,
				MaxNumScore:   o.MaxNumScore,
			})
		}
	}

	if response.Records != nil {
		for _, o := range response.Records {
			metadata := map[string]interface{}{}
			if o.Metadata != nil {
				if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err != nil {
					return nil, err
				}
			}
			list.Records = append(list.Records, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata:      metadata,
				NumScore:      o.NumScore,
				OwnerID:       o.OwnerID,
				Rank:          stringPointerToIntPointer(o.Rank),
				Score:         stringPointerToIntPointer(o.Score),
				SubScore:      stringPointerToIntPointer(o.Subscore),
				UpdateTime:    timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:      o.Username,
				MaxNumScore:   o.MaxNumScore,
			})
		}
	}

	return list, nil
}

func (c *Client) ListLeaderboardRecordsAroundOwner(session *Session, leaderboardId string, ownerId string, limit *int, expiry *string, cursor *string) (*LeaderboardRecordList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListLeaderboardRecordsAroundOwner(session.Token, leaderboardId, ownerId, limit, expiry, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	list := &LeaderboardRecordList{
		NextCursor:   response.NextCursor,
		PrevCursor:   response.PrevCursor,
		RankCount:    stringPointerToIntPointer(response.RankCount),
		OwnerRecords: []LeaderboardRecord{},
		Records:      []LeaderboardRecord{},
	}

	if response.OwnerRecords != nil {
		for _, o := range response.OwnerRecords {
			metadata := map[string]interface{}{}
			if o.Metadata != nil {
				if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err != nil {
					return nil, err
				}
			}
			list.OwnerRecords = append(list.OwnerRecords, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata:      metadata,
				NumScore:      o.NumScore,
				OwnerID:       o.OwnerID,
				Rank:          stringPointerToIntPointer(o.Rank),
				Score:         stringPointerToIntPointer(o.Score),
				SubScore:      stringPointerToIntPointer(o.Subscore),
				UpdateTime:    timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:      o.Username,
				MaxNumScore:   o.MaxNumScore,
			})
		}
	}

	if response.Records != nil {
		for _, o := range response.Records {
			metadata := map[string]interface{}{}
			if o.Metadata != nil {
				if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err != nil {
					return nil, err
				}
			}
			list.Records = append(list.Records, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata:      metadata,
				NumScore:      o.NumScore,
				OwnerID:       o.OwnerID,
				Rank:          stringPointerToIntPointer(o.Rank),
				Score:         stringPointerToIntPointer(o.Score),
				SubScore:      stringPointerToIntPointer(o.Subscore),
				UpdateTime:    timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:      o.Username,
				MaxNumScore:   o.MaxNumScore,
			})
		}
	}

	return list, nil
}

// ListMatches fetches a list of running matches.
func (c *Client) ListMatches(session *Session, limit *int, authoritative *bool, label *string, minSize *int, maxSize *int, query *string) (*ApiMatchList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListMatches(session.Token, limit, authoritative, label, minSize, maxSize, query, make(map[string]string))
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListNotifications fetches a list of notifications.
func (c *Client) ListNotifications(session *Session, limit *int, cacheableCursor *string) (*NotificationList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListNotifications(session.Token, limit, cacheableCursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &NotificationList{
		CacheableCursor: response.CacheableCursor,
		Notifications:   []Notification{},
	}

	if response.Notifications == nil {
		return result, nil
	}

	for _, n := range response.Notifications {
		var content map[string]interface{}
		if n.Content != nil {
			if err := json.Unmarshal([]byte(*n.Content), &content); err != nil {
				return nil, err
			}
		}

		result.Notifications = append(result.Notifications, Notification{
			Code:       n.Code,
			CreateTime: timeToStringPointer(*n.CreateTime, time.RFC3339),
			ID:         n.ID,
			Persistent: n.Persistent,
			SenderID:   n.SenderID,
			Subject:    n.Subject,
			Content:    content,
		})
	}

	return result, nil
}

// ListStorageObjects retrieves a list of storage objects.
func (c *Client) ListStorageObjects(session *Session, collection string, userID *string, limit *int, cursor *string) (*StorageObjectList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListStorageObjects(session.Token, collection, userID, limit, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &StorageObjectList{
		Objects: []StorageObject{},
		Cursor:  response.Cursor,
	}

	if response.Objects == nil {
		return result, nil
	}

	for _, o := range response.Objects {
		var value interface{}
		if o.Value != nil {
			if err := json.Unmarshal([]byte(*o.Value), &value); err != nil {
				return nil, err
			}
		}

		result.Objects = append(result.Objects, StorageObject{
			Collection: o.Collection,
			Key:        o.Key,
			PermissionRead: func() *int {
				if o.PermissionRead != nil {
					return o.PermissionRead
				} else {
					defaultValue := 0
					return &defaultValue
				}
			}(),
			PermissionWrite: func() *int {
				if o.PermissionRead != nil {
					return o.PermissionWrite
				} else {
					defaultValue := 0
					return &defaultValue
				}
			}(),
			Value: func() map[string]interface{} {
				if v, ok := value.(map[string]interface{}); ok {
					return v
				}
				return nil
			}(),
			Version:    o.Version,
			UserID:     o.UserID,
			CreateTime: timeToStringPointer(*o.CreateTime, time.RFC3339),
			UpdateTime: timeToStringPointer(*o.UpdateTime, time.RFC3339),
		})
	}

	return result, nil
}

// ListTournaments retrieves a list of current or upcoming tournaments.
func (c *Client) ListTournaments(session *Session, categoryStart *int, categoryEnd *int, startTime *int64, endTime *int64, limit *int, cursor *string) (*TournamentList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.ApiClient.ListTournaments(session.Token, categoryStart, categoryEnd, startTime, endTime, limit, cursor, make(map[string]string))
	if err != nil {
		return nil, err
	}

	result := &TournamentList{
		Cursor:      response.Cursor,
		Tournaments: []Tournament{},
	}

	if response.Tournaments == nil {
		return result, nil
	}

	for _, o := range response.Tournaments {
		var metadata map[string]interface{}
		if o.Metadata != nil {
			if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err != nil {
				return nil, err
			}
		}

		result.Tournaments = append(result.Tournaments, Tournament{
			ID:          o.ID,
			Title:       o.Title,
			Description: o.Description,
			Duration: func() *int {
				if o.Duration != nil {
					return o.Duration
				}
				defaultValue := 0
				return &defaultValue
			}(),
			Category: func() *int {
				if o.Category != nil {
					return o.Category
				}
				defaultValue := 0
				return &defaultValue
			}(),
			SortOrder: func() *int {
				if o.SortOrder != nil {
					return o.SortOrder
				}
				defaultValue := 0
				return &defaultValue
			}(),
			Size: func() *int {
				if o.Size != nil {
					return o.Size
				}
				defaultValue := 0
				return &defaultValue
			}(),
			MaxSize: func() *int {
				if o.MaxSize != nil {
					return o.MaxSize
				}
				defaultValue := 0
				return &defaultValue
			}(),
			MaxNumScore: func() *int {
				if o.MaxNumScore != nil {
					return o.MaxNumScore
				}
				defaultValue := 0
				return &defaultValue
			}(),
			CanEnter: o.CanEnter,
			EndActive: func() *int {
				if o.EndActive != nil {
					val := int(*o.EndActive)
					return &val
				}
				defaultValue := 0
				return &defaultValue
			}(),
			NextReset: func() *int {
				if o.NextReset != nil {
					val := int(*o.NextReset)
					return &val
				}
				defaultValue := 0
				return &defaultValue
			}(),
			Metadata:      metadata,
			CreateTime:    timeToStringPointer(*o.CreateTime, time.RFC3339),
			StartTime:     timeToStringPointer(*o.StartTime, time.RFC3339),
			EndTime:       timeToStringPointer(*o.EndTime, time.RFC3339),
			StartActive:   int64PointerToIntPointer(o.StartActive),
			Authoritative: o.Authoritative,
		})
	}

	return result, nil
}

// ListSubscriptions lists user subscriptions.
func (c *Client) ListSubscriptions(session *Session, cursor *string, limit *int) (*SubscriptionList, error) {
	if c.AutoRefreshSession && session.IsExpired(time.Now().Unix()+c.ExpiredTimespanMs/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
			return nil, err
		}
	}

	apiSubscriptionList, err := c.ApiClient.ListSubscriptions(
		session.Token, ApiListSubscriptionsRequest{
			Cursor: cursor,
			Limit:  limit,
		},
		make(map[string]string),
	)
	if err != nil {
		return nil, err
	}

	subscriptionList := &SubscriptionList{
		Cursor:     apiSubscriptionList.Cursor,
		PrevCursor: apiSubscriptionList.PrevCursor,
		ValidatedSubscriptions: func(subs []ApiValidatedSubscription) []ValidatedSubscription {
			validatedSubs := make([]ValidatedSubscription, len(subs))
			for i, sub := range subs {
				validatedSubs[i] = ValidatedSubscription{
					Active:                sub.Active,
					CreateTime:            sub.CreateTime,
					Environment:           intPointerToStringPointer((*int)(sub.Environment)),
					ExpiryTime:            sub.ExpiryTime,
					OriginalTransactionID: sub.OriginalTransactionID,
					ProductID:             sub.ProductID,
					ProviderNotification:  sub.ProviderNotification,
					ProviderResponse:      sub.ProviderResponse,
					PurchaseTime:          sub.PurchaseTime,
					RefundTime:            sub.RefundTime,
					Store:                 intPointerToStringPointer((*int)(sub.Store)),
					UpdateTime:            sub.UpdateTime,
					UserID:                sub.UserID,
				}
			}
			return validatedSubs
		}(apiSubscriptionList.ValidatedSubscriptions),
	}

	return subscriptionList, nil
}

// ListTournamentRecords lists tournament records from a given tournament.
func (c *Client) ListTournamentRecords(
	session *Session,
	tournamentId string,
	ownerIds []string,
	limit *int,
	cursor *string,
	expiry *string,
) (*TournamentRecordList, error) {
	// Refresh the session if auto-refresh is enabled and the session is expired.
	if c.AutoRefreshSession && session.IsExpired(time.Now().Unix()+c.ExpiredTimespanMs/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
			return nil, err
		}
	}

	// Call the API to list tournament records.
	apiTournamentRecordList, err := c.ApiClient.ListTournamentRecords(
		session.Token,
		tournamentId,
		ownerIds,
		limit,
		cursor,
		expiry,
		make(map[string]string),
	)
	if err != nil {
		return nil, err
	}

	// Prepare the response object.
	list := &TournamentRecordList{
		NextCursor:   apiTournamentRecordList.NextCursor,
		PrevCursor:   apiTournamentRecordList.PrevCursor,
		OwnerRecords: []LeaderboardRecord{},
		Records:      []LeaderboardRecord{},
	}

	// Process owner records.
	if apiTournamentRecordList.OwnerRecords != nil {
		for _, o := range apiTournamentRecordList.OwnerRecords {
			list.OwnerRecords = append(list.OwnerRecords, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata: func() map[string]interface{} {
					if o.Metadata == nil {
						return nil
					}
					var metadata map[string]interface{}
					if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err == nil {
						return metadata
					}
					return nil
				}(),
				NumScore:    o.NumScore,
				OwnerID:     o.OwnerID,
				Rank:        stringPointerToIntPointer(o.Rank),
				Score:       stringPointerToIntPointer(o.Score),
				SubScore:    stringPointerToIntPointer(o.Subscore),
				UpdateTime:  timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:    o.Username,
				MaxNumScore: o.MaxNumScore,
			})
		}
	}

	// Process records.
	if apiTournamentRecordList.Records != nil {
		for _, r := range apiTournamentRecordList.Records {
			list.Records = append(list.Records, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*r.ExpiryTime, time.RFC3339),
				LeaderboardID: r.LeaderboardID,
				Metadata: func() map[string]interface{} {
					if r.Metadata == nil {
						return nil
					}
					var metadata map[string]interface{}
					if err := json.Unmarshal([]byte(*r.Metadata), &metadata); err == nil {
						return metadata
					}
					return nil
				}(),
				NumScore:    r.NumScore,
				OwnerID:     r.OwnerID,
				Rank:        stringPointerToIntPointer(r.Rank),
				Score:       stringPointerToIntPointer(r.Score),
				SubScore:    stringPointerToIntPointer(r.Subscore),
				UpdateTime:  timeToStringPointer(*r.UpdateTime, time.RFC3339),
				Username:    r.Username,
				MaxNumScore: r.MaxNumScore,
			})
		}
	}

	return list, nil
}

// ListTournamentRecordsAroundOwner lists tournament records around a specific owner.
func (c *Client) ListTournamentRecordsAroundOwner(
	session *Session,
	tournamentId string,
	ownerId string,
	limit *int,
	expiry *string,
	cursor *string,
) (*TournamentRecordList, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
			return nil, err
		}
	}

	// Call the API to get tournament records around owner.
	apiTournamentRecordList, err := c.ApiClient.ListTournamentRecordsAroundOwner(
		session.Token,
		tournamentId,
		ownerId,
		limit,
		expiry,
		cursor,
		make(map[string]string),
	)
	if err != nil {
		return nil, err
	}

	// Prepare the response object.
	list := &TournamentRecordList{
		NextCursor:   apiTournamentRecordList.NextCursor,
		PrevCursor:   apiTournamentRecordList.PrevCursor,
		OwnerRecords: []LeaderboardRecord{},
		Records:      []LeaderboardRecord{},
	}

	// Process owner records.
	if apiTournamentRecordList.OwnerRecords != nil {
		for _, o := range apiTournamentRecordList.OwnerRecords {
			list.OwnerRecords = append(list.OwnerRecords, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*o.ExpiryTime, time.RFC3339),
				LeaderboardID: o.LeaderboardID,
				Metadata: func() map[string]interface{} {
					if o.Metadata == nil {
						return nil
					}
					var metadata map[string]interface{}
					if err := json.Unmarshal([]byte(*o.Metadata), &metadata); err == nil {
						return metadata
					}
					return nil
				}(),
				NumScore:    o.NumScore,
				OwnerID:     o.OwnerID,
				Rank:        stringPointerToIntPointer(o.Rank),
				Score:       stringPointerToIntPointer(o.Score),
				SubScore:    stringPointerToIntPointer(o.Subscore),
				UpdateTime:  timeToStringPointer(*o.UpdateTime, time.RFC3339),
				Username:    o.Username,
				MaxNumScore: o.MaxNumScore,
			})
		}
	}

	// Process records.
	if apiTournamentRecordList.Records != nil {
		for _, r := range apiTournamentRecordList.Records {
			list.Records = append(list.Records, LeaderboardRecord{
				ExpiryTime:    timeToStringPointer(*r.ExpiryTime, time.RFC3339),
				LeaderboardID: r.LeaderboardID,
				Metadata: func() map[string]interface{} {
					if r.Metadata == nil {
						return nil
					}
					var metadata map[string]interface{}
					if err := json.Unmarshal([]byte(*r.Metadata), &metadata); err == nil {
						return metadata
					}
					return nil
				}(),
				NumScore:    r.NumScore,
				OwnerID:     r.OwnerID,
				Rank:        stringPointerToIntPointer(r.Rank),
				Score:       stringPointerToIntPointer(r.Score),
				SubScore:    stringPointerToIntPointer(r.Subscore),
				UpdateTime:  timeToStringPointer(*r.UpdateTime, time.RFC3339),
				Username:    r.Username,
				MaxNumScore: r.MaxNumScore,
			})
		}
	}

	return list, nil
}

// PromoteGroupUsers promotes the users in a group to the next role up.
func (c *Client) PromoteGroupUsers(session *Session, groupId string, ids []string) (bool, error) {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
			return false, err
		}
	}

	success, err := c.ApiClient.PromoteGroupUsers(session.Token, groupId, ids, make(map[string]string))
	if err != nil {
		return false, err
	}
	return success.(bool), nil
}

// ValidatePurchaseApple validates an Apple IAP receipt.
func (c *Client) ValidatePurchaseApple(bearerToken string, body *ApiValidatePurchaseAppleRequest) (*ApiValidatePurchaseResponse, error) {
	if body == nil {
		return nil, fmt.Errorf("'body' is a required parameter but is nil")
	}

	urlPath := "/v2/iap/purchase/apple"
	queryParams := url.Values{}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize body to JSON: %w", err)
	}

	fullUrl := c.ApiClient.buildFullUrl(c.ApiClient.BasePath, urlPath, queryParams)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	client := &http.Client{
		Timeout: time.Duration(c.ApiClient.TimeoutMs) * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var apiResponse ApiValidatePurchaseResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response JSON: %w", err)
		}
		return &apiResponse, nil
	}

	// Handle non-2xx responses
	respBody, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(respBody))
}

// SessionRefresh refreshes a user's session using a refresh token retrieved from a previous authentication request.
func (c *Client) SessionRefresh(session *Session, vars map[string]string) (*Session, error) {
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
		Token: &session.RefreshToken,
		Vars:  vars,
	}, make(map[string]string))

	if err != nil {
		return nil, err
	}

	session.Update(*apiSession.Token, *apiSession.RefreshToken)
	return session, nil
}
