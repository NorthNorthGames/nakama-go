package nakama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ApiOperator int

const (
	NoOverride ApiOperator = iota
	Best
	Set
	Increment
	Decrement
)

type ApiStoreEnvironment int

const (
	Unknown ApiStoreEnvironment = iota
	Sandbox
	Production
)

type ApiStoreProvider int

const (
	AppleAppStore ApiStoreProvider = iota
	GooglePlayStore
	HuaweiAppGallery
	FacebookInstantStore
)

type FriendsOfFriendsListFriendOfFriend struct {
	Referrer string   // The user who referred its friend.
	User     *ApiUser // User.
}

type GroupUserListGroupUser struct {
	State int      // Their relationship to the group.
	User  *ApiUser // User.
}

type UserGroupListUserGroup struct {
	Group *ApiGroup // Group.
	State int       // The user's relationship to the group.
}

type WriteLeaderboardRecordRequestLeaderboardRecordWrite struct {
	Metadata string      // Optional record metadata.
	Operator ApiOperator // Operator override.
	Score    string      // The score value to submit.
	Subscore string      // An optional secondary value.
}

type WriteTournamentRecordRequestTournamentRecordWrite struct {
	Metadata string      // A JSON object of additional properties (optional).
	Operator ApiOperator // Operator override.
	Score    string      // The score value to submit.
	Subscore string      // An optional secondary value.
}

type ApiAccount struct {
	CustomID    string             // The custom id in the user's account.
	Devices     []ApiAccountDevice // The devices which belong to the user's account.
	DisableTime string             // The time when the user's account was disabled/banned.
	Email       string             // The email address of the user.
	User        *ApiUser           // The user object.
	VerifyTime  string             // The time when the user's email was verified.
	Wallet      string             // The user's wallet data.
}

type ApiAccountApple struct {
	Token string            // The ID token received from Apple to validate.
	Vars  map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountCustom struct {
	ID   string            // A custom identifier.
	Vars map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountDevice struct {
	ID   string            // A device identifier. Should be obtained by platform-specific device API.
	Vars map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountEmail struct {
	Email    string            // A valid RFC-5322 email address.
	Password string            // A password for the user account. Ignored with unlink operations.
	Vars     map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountFacebook struct {
	Token string            // The OAuth token from Facebook to access their profile API.
	Vars  map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountFacebookInstantGame struct {
	SignedPlayerInfo string            // The signed player info from Facebook.
	Vars             map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountGameCenter struct {
	BundleID     string            // Bundle ID (generated by GameCenter).
	PlayerID     string            // Player ID (generated by GameCenter).
	PublicKeyURL string            // The URL for the public encryption key.
	Salt         string            // A random "NSString" used to compute the hash and keep it randomized.
	Signature    string            // The verification signature data generated.
	Timestamp    string            // Time since UNIX epoch when the signature was created.
	Vars         map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountGoogle struct {
	Token string            // The OAuth token received from Google to access their profile API.
	Vars  map[string]string // Extra information that will be bundled in the session token.
}

type ApiAccountSteam struct {
	Token string            // The account token received from Steam to access their profile API.
	Vars  map[string]string // Extra information that will be bundled in the session token.
}

type ApiChannelMessage struct {
	ChannelID  string // The channel this message belongs to.
	Code       int    // The code representing a message type or category.
	Content    string // The content payload.
	CreateTime string // The time when the message was created.
	GroupID    string // The ID of the group, or empty if not a group channel.
	MessageID  string // The unique ID of this message.
	Persistent bool   // True if the message was persisted to history; false otherwise.
	RoomName   string // The name of the chat room, or empty if it was not a chat room.
	SenderID   string // Message sender, usually a user ID.
	UpdateTime string // The time when the message was last updated.
	UserIDOne  string // The ID of the first DM user, or empty if it was not a DM chat.
	UserIDTwo  string // The ID of the second DM user, or empty if it was not a DM chat.
	Username   string // The username of the message sender, if any.
}

type ApiChannelMessageList struct {
	CacheableCursor string              // Cacheable cursor to list newer messages. Durable and designed to be stored, unlike next/prev cursors.
	Messages        []ApiChannelMessage // A list of messages.
	NextCursor      string              // The cursor to send when retrieving the next page, if any.
	PrevCursor      string              // The cursor to send when retrieving the previous page, if any.
}

type ApiCreateGroupRequest struct {
	AvatarURL   string // A URL for an avatar image.
	Description string // A description for the group.
	LangTag     string // The language expected to be a tag which follows the BCP-47 spec.
	MaxCount    int    // Maximum number of group members.
	Name        string // A unique name for the group.
	Open        bool   // Mark a group as open or not, where only admins can accept members.
}

type ApiDeleteStorageObjectId struct {
	Collection string // The collection which stores the object.
	Key        string // The key of the object within the collection.
	Version    string // The version hash of the object.
}

type ApiDeleteStorageObjectsRequest struct {
	ObjectIds []ApiDeleteStorageObjectId // Batch of storage objects.
}

type ApiEvent struct {
	External   bool              // True if the event came directly from a client call, false otherwise.
	Name       string            // An event name, type, category, or identifier.
	Properties map[string]string // Arbitrary event property values.
	Timestamp  string            // The time when the event was triggered.
}

type ApiFriend struct {
	State      int      // The friend status. One of "Friend.State".
	UpdateTime string   // Time of the latest relationship update.
	User       *ApiUser // The user object.
}

type ApiFriendList struct {
	Cursor  string      // Cursor for the next page of results, if any.
	Friends []ApiFriend // The friend objects.
}

type ApiFriendsOfFriendsList struct {
	Cursor           string                               // Cursor for the next page of results, if any.
	FriendsOfFriends []FriendsOfFriendsListFriendOfFriend // User friends of friends.
}

type ApiGroup struct {
	AvatarURL   string // A URL for an avatar image.
	CreateTime  string // The UNIX time when the group was created.
	CreatorID   string // The ID of the user who created the group.
	Description string // A description for the group.
	EdgeCount   int    // The current count of all members in the group.
	ID          string // The ID of the group.
	LangTag     string // The language expected to be a tag which follows the BCP-47 spec.
	MaxCount    int    // The maximum number of members allowed.
	Metadata    string // Additional information, stored as a JSON object.
	Name        string // The unique name of the group.
	Open        bool   // Anyone can join open groups; otherwise, only admins can accept members.
	UpdateTime  string // The UNIX time when the group was last updated.
}

type ApiGroupList struct {
	Cursor string     // A cursor used to get the next page.
	Groups []ApiGroup // One or more groups.
}

type ApiGroupUserList struct {
	Cursor     string                   // Cursor for the next page of results, if any.
	GroupUsers []GroupUserListGroupUser // User-role pairs for a group.
}

type ApiLeaderboardRecord struct {
	CreateTime    string // The UNIX time when the leaderboard record was created.
	ExpiryTime    string // The UNIX time when the leaderboard record expires.
	LeaderboardID string // The ID of the leaderboard this score belongs to.
	MaxNumScore   int    // The maximum number of score updates allowed by the owner.
	Metadata      string // Metadata.
	NumScore      int    // The number of submissions to this score record.
	OwnerID       string // The ID of the score owner, usually a user or group.
	Rank          string // The rank of this record.
	Score         string // The score value.
	Subscore      string // An optional subscore value.
	UpdateTime    string // The UNIX time when the leaderboard record was updated.
	Username      string // The username of the score owner, if the owner is a user.
}

type ApiLeaderboardRecordList struct {
	NextCursor   string                 // The cursor to send when retrieving the next page, if any.
	OwnerRecords []ApiLeaderboardRecord // A batched set of leaderboard records belonging to specified owners.
	PrevCursor   string                 // The cursor to send when retrieving the previous page, if any.
	RankCount    string                 // The total number of ranks available.
	Records      []ApiLeaderboardRecord // A list of leaderboard records.
}

type ApiLinkSteamRequest struct {
	Account *ApiAccountSteam // The Steam account details.
	Sync    bool             // Import Steam friends for the user.
}

type ApiListSubscriptionsRequest struct {
	Cursor string // Cursor for paginated subscriptions, if any.
	Limit  int    // Maximum number of subscriptions to retrieve (optional).
}

type ApiMatch struct {
	Authoritative bool   // True if it's a server-managed authoritative match, false otherwise.
	HandlerName   string // Handler name for the match, if any.
	Label         string // Match label, if any.
	MatchID       string // The ID of the match, used to join a match.
	Size          int    // Current number of users in the match.
	TickRate      int    // Tick rate of the match, if any.
}

type ApiMatchList struct {
	Matches []ApiMatch // A number of matches corresponding to a list operation.
}

type ApiNotification struct {
	Code       int    // Category code for this notification.
	Content    string // Content of the notification in JSON.
	CreateTime string // The time when the notification was created.
	ID         string // ID of the notification.
	Persistent bool   // True if this notification was persisted to the database.
	SenderID   string // ID of the sender, if a user; otherwise empty.
	Subject    string // Subject of the notification.
}

type ApiNotificationList struct {
	CacheableCursor string            // Cursor to paginate notifications.
	Notifications   []ApiNotification // Collection of notifications.
}

type ApiReadStorageObjectId struct {
	Collection string // The collection which stores the object.
	Key        string // The key of the object within the collection.
	UserID     string // The user owner of the object.
}

type ApiReadStorageObjectsRequest struct {
	ObjectIds []ApiReadStorageObjectId // Batch of storage objects.
}

type ApiRpc struct {
	HTTPKey string // The authentication key used when executed as a non-client HTTP request.
	ID      string // The identifier of the function.
	Payload string // The payload of the function which must be a JSON object.
}

type ApiSession struct {
	Created      bool   // True if the corresponding account was just created, false otherwise.
	RefreshToken string // Refresh token that can be used for session token renewal.
	Token        string // Authentication credentials.
}

type ApiSessionLogoutRequest struct {
	RefreshToken string // Refresh token to invalidate.
	Token        string // Session token to log out.
}

type ApiSessionRefreshRequest struct {
	Token string            // Refresh token.
	Vars  map[string]string // Extra information that will be bundled in the session token.
}

type ApiStorageObject struct {
	Collection      string // The collection which stores the object.
	CreateTime      string // The UNIX time when the object was created.
	Key             string // The key of the object within the collection.
	PermissionRead  int    // The read access permissions for the object.
	PermissionWrite int    // The write access permissions for the object.
	UpdateTime      string // The UNIX time when the object was last updated.
	UserID          string // The user owner of the object.
	Value           string // The value of the object.
	Version         string // The version hash of the object.
}

type ApiStorageObjectAck struct {
	Collection string // The collection which stores the object.
	CreateTime string // The UNIX time when the object was created.
	Key        string // The key of the object within the collection.
	UpdateTime string // The UNIX time when the object was last updated.
	UserID     string // The owner of the object.
	Version    string // The version hash of the object.
}

type ApiStorageObjectAcks struct {
	Acks []ApiStorageObjectAck // Batch of storage write acknowledgements.
}

type ApiStorageObjectList struct {
	Cursor  string             // The cursor for the next page of results, if any.
	Objects []ApiStorageObject // The list of storage objects.
}

type ApiStorageObjects struct {
	Objects []ApiStorageObject // The batch of storage objects.
}

type ApiSubscriptionList struct {
	Cursor                 string                     // The cursor to send when retrieving the next page, if any.
	PrevCursor             string                     // The cursor to send when retrieving the previous page, if any.
	ValidatedSubscriptions []ApiValidatedSubscription // Stored validated subscriptions.
}

type ApiTournament struct {
	Authoritative bool        // Whether the leaderboard was created authoritatively or not.
	CanEnter      bool        // True if the tournament is active and can enter. A computed value.
	Category      int         // The category of the tournament. e.g. "vip" could be category 1.
	CreateTime    string      // The UNIX time when the tournament was created.
	Description   string      // The description of the tournament. May be blank.
	Duration      int         // Duration of the tournament in seconds.
	EndActive     int         // The UNIX time when the tournament stops being active until the next reset.
	EndTime       string      // The UNIX time when the tournament will be stopped.
	ID            string      // The ID of the tournament.
	MaxNumScore   int         // The maximum score updates allowed per player for the current tournament.
	MaxSize       int         // The maximum number of players for the tournament.
	Metadata      string      // Additional information stored as a JSON object.
	NextReset     int         // The UNIX time when the tournament is next playable.
	Operator      ApiOperator // Operator.
	PrevReset     int         // The UNIX time when the tournament was last reset.
	Size          int         // The current number of players in the tournament.
	SortOrder     int         // ASC (0) or DESC (1) sort mode of scores in the tournament.
	StartActive   int         // The UNIX time when the tournament starts being active.
	StartTime     string      // The UNIX time when the tournament will start.
	Title         string      // The title of the tournament.
}

type ApiTournamentList struct {
	Cursor      string          // A pagination cursor (optional).
	Tournaments []ApiTournament // The list of tournaments returned.
}

type ApiTournamentRecordList struct {
	NextCursor   string                 // The cursor to send when retrieving the next page.
	OwnerRecords []ApiLeaderboardRecord // A batched set of tournament records belonging to specified owners.
	PrevCursor   string                 // The cursor to send when retrieving the previous page.
	RankCount    string                 // The total number of ranks available.
	Records      []ApiLeaderboardRecord // A list of tournament records.
}

type ApiUpdateAccountRequest struct {
	AvatarURL   string // A URL for an avatar image.
	DisplayName string // The display name of the user.
	LangTag     string // The language expected to be a tag that follows the BCP-47 spec.
	Location    string // The location set by the user.
	Timezone    string // The timezone set by the user.
	Username    string // The username of the user's account.
}

type ApiUpdateGroupRequest struct {
	AvatarURL   string // Avatar URL.
	Description string // Description string.
	GroupID     string // The ID of the group to update.
	LangTag     string // Language tag.
	Name        string // Name.
	Open        bool   // True if anyone can join, false otherwise only admins can approve members.
}

type ApiUser struct {
	AppleID               string // The Apple Sign-In ID in the user's account.
	AvatarURL             string // A URL for an avatar image.
	CreateTime            string // The UNIX time when the user was created.
	DisplayName           string // The display name of the user.
	EdgeCount             int    // Number of related edges to this user.
	FacebookID            string // The Facebook ID in the user's account.
	FacebookInstantGameID string // The Facebook Instant Game ID in the user's account.
	GameCenterID          string // The Apple Game Center ID in the user's account.
	GoogleID              string // The Google ID in the user's account.
	ID                    string // The ID of the user's account.
	LangTag               string // The language expected to be a tag that follows the BCP-47 spec.
	Location              string // The location set by the user.
	Metadata              string // Additional information stored as a JSON object.
	Online                bool   // Indicates whether the user is currently online.
	SteamID               string // The Steam ID in the user's account.
	Timezone              string // The timezone set by the user.
	UpdateTime            string // The UNIX time when the user was last updated.
	Username              string // The username of the user's account.
}

type ApiUserGroupList struct {
	Cursor     string                   // Cursor for the next page of results, if any.
	UserGroups []UserGroupListUserGroup // Group-role pairs for a user.
}

type ApiUsers struct {
	Users []ApiUser // The list of user objects.
}

type ApiValidatePurchaseAppleRequest struct {
	Persist bool   // True to persist the validated purchase.
	Receipt string // Base64 encoded Apple receipt data payload.
}

type ApiValidatePurchaseFacebookInstantRequest struct {
	Persist       bool   // True to persist the validated purchase.
	SignedRequest string // Base64 encoded Facebook Instant signedRequest receipt data payload.
}

type ApiValidatePurchaseGoogleRequest struct {
	Persist  bool   // True to persist the validated purchase.
	Purchase string // JSON encoded Google purchase payload.
}

type ApiValidatePurchaseHuaweiRequest struct {
	Persist   bool   // True to persist the validated purchase.
	Purchase  string // JSON encoded Huawei InAppPurchaseData.
	Signature string // InAppPurchaseData signature.
}

type ApiValidatePurchaseResponse struct {
	ValidatedPurchases []ApiValidatedPurchase // Newly seen validated purchases.
}

type ApiValidateSubscriptionAppleRequest struct {
	Persist bool   // True to persist the subscription.
	Receipt string // Base64 encoded Apple receipt data payload.
}

type ApiValidateSubscriptionGoogleRequest struct {
	Persist bool   // True to persist the subscription.
	Receipt string // JSON encoded Google purchase payload.
}

type ApiValidateSubscriptionResponse struct {
	ValidatedSubscription *ApiValidatedSubscription // The validated subscription.
}

type ApiValidatedPurchase struct {
	CreateTime       string              // The time when the receipt validation was stored in the DB.
	Environment      ApiStoreEnvironment // Whether the purchase was done in production or sandbox environment.
	ProductID        string              // Purchase product ID.
	ProviderResponse string              // Raw provider validation response.
	PurchaseTime     string              // The time when the purchase was made.
	RefundTime       string              // The time when the purchase was refunded, if applicable.
	SeenBefore       bool                // True if the purchase was already validated.
	Store            ApiStoreProvider    // The store where the purchase was made.
	TransactionID    string              // The transaction ID for the purchase.
	UpdateTime       string              // The time when the receipt validation was last updated in the DB.
	UserID           string              // The ID of the user who made the purchase.
}

type ApiValidatedSubscription struct {
	Active                bool                // Whether the subscription is currently active or not.
	CreateTime            string              // The time when the receipt validation was stored in the DB.
	Environment           ApiStoreEnvironment // The environment in which the subscription was made.
	ExpiryTime            string              // The expiration time of the subscription.
	OriginalTransactionID string              // The original transaction ID for the subscription.
	ProductID             string              // Subscription product ID.
	ProviderNotification  string              // The raw provider notification body.
	ProviderResponse      string              // The raw provider validation response body.
	PurchaseTime          string              // The time when the subscription was purchased.
	RefundTime            string              // The time when the subscription was refunded, if applicable.
	Store                 ApiStoreProvider    // The store where the subscription was made.
	UpdateTime            string              // The time when the receipt validation was updated in the DB.
	UserID                string              // The user ID for the subscription.
}

type ApiWriteStorageObject struct {
	Collection      string // The collection to store the object.
	Key             string // The key for the object within the collection.
	PermissionRead  int    // The read access permissions for the object.
	PermissionWrite int    // The write access permissions for the object.
	Value           string // The value of the object.
	Version         string // The version hash of the object to check.
}

type ApiWriteStorageObjectsRequest struct {
	Objects []ApiWriteStorageObject // The objects to store on the server.
}

type NakamaApi struct {
	ServerKey string
	BasePath  string
	TimeoutMs int
}

// Healthcheck is a healthcheck function that load balancers can use to check the service.
func (api *NakamaApi) Healthcheck(bearerToken string, options map[string]string) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/healthcheck"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// DeleteAccount deletes the current user's account.
func (api *NakamaApi) DeleteAccount(bearerToken string, options map[string]string) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// GetAccount fetches the current user's account.
func (api *NakamaApi) GetAccount(bearerToken string, options map[string]string) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}
