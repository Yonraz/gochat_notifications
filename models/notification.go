package models

import "github.com/yonraz/gochat_notifications/constants"

type Notification struct {
	Sender  string
	Content string
	Type    constants.Notification
}