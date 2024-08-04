package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yonraz/gochat_notifications/ws"
)

type NotificationsController struct {
	RedisClient *redis.Client
}

func NewNotificationsController(client *redis.Client) *NotificationsController {
	return &NotificationsController{
		RedisClient: client,
	}
}

func (c *NotificationsController) GetNotificationsForUser(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		log.Printf("username param not provided!\n")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username param not provided",
		})
		return
	}

	notifKeys, err := c.RedisClient.HGetAll(context.Background(), "notifications:"+username).Result()
	if err != nil {
		log.Printf("error accessing redis!\n")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error accessing redis",
		})
		return
	}

	var notifs []*ws.Message
	for _, notifJson := range notifKeys {
		var notif ws.Message
		err := json.Unmarshal([]byte(notifJson), &notif)
		if err != nil {
			log.Printf("error unmarshaling notifs!\n")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error unmarshaling notifications",
			})
			return
		}
		notifs = append(notifs, &notif)
	}

	// Optionally clear notifications after fetching them
	_, err = c.RedisClient.Del(context.Background(), "notifications:"+username).Result()
	if err != nil {
		log.Printf("error clearing notifs!\n")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error clearing notifs",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notifs": notifs,
	})
}

func (c *NotificationsController) SetNotificationsForUser(ctx *gin.Context) {
	var body struct {
		Notifs   []*ws.Message `json:"notifs"`
		Username string              `json:"username"`
	}

	if ctx.BindJSON(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	for _, notif := range body.Notifs {
		notifJson, err := json.Marshal(notif)
		if err != nil {
			log.Printf("error marshaling notifs!\n")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error marshaling notifications",
			})
			return
		}

		_, err = c.RedisClient.HSet(context.Background(), "notifications:"+body.Username, notif.ID, notifJson).Result()
		if err != nil {
			log.Printf("error setting notifs!\n")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error setting notifications",
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "inserted notifications into redis",
	})
}

func (c *NotificationsController) DeleteNotifications(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		log.Printf("username param not provided!\n")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username param not provided",
		})
		return
	}

	var body struct {
		Notifs []*ws.Message `json:"notifs"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		log.Printf("error binding JSON: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	if len(body.Notifs) == 0 {
		log.Printf("no notifications provided for deletion\n")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "no notifications provided",
		})
		return
	}

	for _, notif := range body.Notifs {
		_, err := c.RedisClient.HDel(context.Background(), "notifications:"+username, notif.ID).Result()
		if err != nil {
			log.Printf("error deleting notification ID %s: %v\n", notif.ID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete notification",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted notifications",
	})
}
