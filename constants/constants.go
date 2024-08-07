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
	UserLoginQueue        Queues = "NOTIFICATIONS_SRV_UserLoginQueue"
	UserSignoutQueue      Queues = "NOTIFICATIONS_SRV_UserSignoutQueue"
	MessageSentQueue      Queues = "NOTIFICATIONS_SRV_MessageSentQueue"
	MessageDeliveredQueue Queues = "NOTIFICATIONS_SRV_MessageDeliveredQueue"
	MessageReadQueue      Queues = "NOTIFICATIONS_SRV_MessageReadQueue"
)

const (
	UserOnline      Notification = "user.online"
	UserOffline     Notification = "user.offline"
	UserSentMessage Notification = "user.sent.message"
)

const (
	NotificationClients RedisSet = "notifications:clients"
)