package constants

type Queues string
type RoutingKey string
type Exchange string
type MessageType string
type RedisSet string
type Notification string

const (
	UserEventsExchange    Exchange = "UserEventsExchange"
	MessageEventsExchange Exchange = "MessageEventsExchange"
)

const (
	UserLoggedInKey     RoutingKey = "user.logged.in"
	UserSignedoutKey    RoutingKey = "user.signed.out"
	MessageSentKey      RoutingKey = "message.sent"
	MessageDeliveredKey RoutingKey = "message.delivered"
	MessageReadKey      RoutingKey = "message.read"
)

const (
	UserLoginQueue        Queues = "UserLoginQueue"
	UserSignoutQueue      Queues = "UserSignoutQueue"
	MessageSentQueue      Queues = "MessageSentQueue"
	MessageDeliveredQueue Queues = "MessageDeliveredQueue"
	MessageReadQueue      Queues = "MessageReadQueue"
)

const (
	UserOnline      Notification = "user.online"
	UserOffline     Notification = "user.offline"
	UserSentMessage Notification = "user.sent.message"
)

const (
	NotificationClients RedisSet = "notifications:clients"
)