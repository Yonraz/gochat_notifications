# Notifications Service

This service is responsible for pushing notifications to the client. It listens for events emitted by other services and notifies the client of them by employing a persistent websocket connection.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
- [API Endpoints](#api-endpoints)
  - [Get Users](#get-users)
<!-- - [Event Consumers](#event-consumers)
  - [User Registered Handler](#user-registered-handler)
  - [User Logged-in Handler](#user-logged-in-handler)
  - [User Signed-out Handler](#user-signed-out-handler) -->
- [License](#license)

## Features

- Persistent, real time two-way communication through websockets.
- Consumes RabbitMQ events.
- Notifies client of user online/offline/message-sent/message-read events.
- Maintains a Redis cache that holds client information.

## Technologies

- Go
- Gin Framework
- Gorilla websocket
- GORM
- RabbitMQ
- Redis

# NotificationsController API Documentation

## Overview

The `NotificationsController` provides endpoints for managing notifications in Redis. Notifications are identified by their IDs and are stored in Redis hashes, keyed by username.

## Endpoints

### Get Notifications for User

**Endpoint**: `GET /notifications/:username`

Retrieves all notifications for a given user.

#### Request Parameters

- **username**: The username of the user whose notifications are being requested.

#### Responses

- **200 OK**: Returns a list of notifications.
- **400 Bad Request**: If the `username` parameter is not provided.
- **500 Internal Server Error**: If there's an error accessing Redis.

#### Response Body

```json
{
  "notifs": [
    {
      "id": "notification_id_1",
      "content": "Notification content",
      "sender": "sender_username",
      "receiver": "receiver_username",
      "status": "notification_status",
      "type": "notification_type",
      "read": false,
      "sent": true
    },
    ...
  ]
}
```
#### Example Request:

```bash
GET /notifications/johndoe
```
#### Response:

```json
{
  "notifs": [
    {
      "id": "12345",
      "content": "You have a new message.",
      "sender": "alice",
      "receiver": "johndoe",
      "status": "delivered",
      "type": "message",
      "read": false,
      "sent": true
    }
  ]
}
```
### Set Notifications for User
__Endpoint:__ `POST /notifications`

Sets or updates notifications for a user.

#### Request Body
```json
{
  "notifs": [
    {
      "id": "notification_id_1",
      "content": "Notification content",
      "sender": "sender_username",
      "receiver": "receiver_username",
      "status": "notification_status",
      "type": "notification_type",
      "read": false,
      "sent": true
    },
    ...
  ],
  "username": "username"
}
```
#### Responses
- **201 Created**: If notifications are successfully inserted.
- **400 Bad Request**: If the request body is invalid or missing fields.
- **500 Internal Server Error**: If there's an error setting notifications in Redis.

#### Example Request:

`POST /notifications`
```json
{
  "notifs": [
    {
      "id": "67890",
      "content": "Your profile has been updated.",
      "sender": "system",
      "receiver": "johndoe",
      "status": "read",
      "type": "system",
      "read": true,
      "sent": true
    }
  ],
  "username": "johndoe"
}
```
#### Response:

```json
{
  "message": "inserted notifications into redis"
}
```

#### Delete Notifications
**Endpoint**: `DELETE /notifications/:username`

Deletes specific notifications for a user.

#### Request Parameters
**username**: `The username of the user whose notifications are being deleted.`
**Request Body**
```json
{
  "notifs": [
    {
      "id": "notification_id_1",
      "content": "Notification content",
      "sender": "sender_username",
      "receiver": "receiver_username",
      "status": "notification_status",
      "type": "notification_type",
      "read": false,
      "sent": true
    },
    ...
  ]
}
```
Responses
- **200 OK**: If notifications are successfully deleted.
- **400 Bad Request**: If the username parameter is not provided or no notifications are provided.
- **500 Internal Server Error**: If there's an error deleting notifications in Redis.

#### Example Request:

```json
DELETE /notifications/johndoe
{
  "notifs": [
    {
      "id": "12345"
    }
  ]
}
```
#### Response:

```json
{
  "message": "successfully deleted notifications"
}
```
### Models
##### WsMessage
Represents a notification message.
Fields
```go
ID: string  //The unique identifier of the notification.
Content: string  //The content of the notification.
Sender: string  //The username of the sender.
Receiver: string  //The username of the receiver.
Status: string  //The status of the notification (e.g., "delivered").
Type: string  //The type of the notification (e.g., "message").
Read: bool  //Whether the notification has been read.
Sent: bool  //Whether the notification has been sent.
```
This documentation provides detailed information about how to use the NotificationsController API, including request formats, responses, and examples for each endpoint.