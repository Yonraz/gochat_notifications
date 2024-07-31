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

