package models

import (
	"time"

	"github.com/yonraz/gochat_notifications/constants"
)

type WsMessage struct {
	ID      	string 					`json:"id"`			
	Content 	string					`json:"content"`
	Sender 		string					`json:"sender"`		
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.Notification	`json:"type"`
	CreatedAt time.Time  				`json:"createdAt"`
	UpdatedAt time.Time  				`json:"updatedAt"`
}

