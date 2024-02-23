package common

var (
	// PurchaseTopic is the subscribed topic for new purchase
	PurchaseTopic = "purchase"
	// PurchaseResultTopic is the topic to which we publish new purchase result
	PurchaseResultTopic = "purchase.result"

	// ReplyTopic is saga step reply topic
	ReplyTopic = "reply"
)
