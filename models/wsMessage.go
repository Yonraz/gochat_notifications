package models

import (
	"github.com/yonraz/gochat_notifications/constants"
	"gorm.io/gorm"
)

type WsMessage struct {
	ID      	string 					`json:"id"`			
	Content 	string					`json:"content"`
	Sender 		string					`json:"sender"`		
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.Notification	`json:"type"`
	gorm.Model
}

