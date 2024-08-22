# WebSocket Chat in Go

WebSocket chat application implemented in Go, consisting of two services, Service A and Service B, that communicate with each other via WebSockets.

## Features

- Real-time messaging using WebSockets
- Console-based message input
- Automatic reconnection on connection loss

## Prerequisites

- Go 1.15 or later
- `github.com/gorilla/websocket` package

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/hojamuhammet/websocket-chat.git
   ```

2. **Install dependencies:**
   ```sh
   go get github.com/gorilla/websocket
   ```

## Usage

### Running the Services

1. **Start Service A:**
   ```sh
   cd service_a
   go run service_a.go
   ```

2. **Start Service B:**
   ```sh
   cd service_b
   go run service_b.go
   ```

### Sending Messages

Enter a message in the console where each service is running to send a message to the other service.
