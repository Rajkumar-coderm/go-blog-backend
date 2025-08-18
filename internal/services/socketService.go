package services

import (
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type OnlineUser struct {
	UID      string `json:"uid"` // Changed from ID to UID to match Flutter
	Name     string `json:"name"`
	Email    string `json:"email"`
	SocketID string `json:"socket_id"` // Added to track socket connection
}

// Global user list - using socketID as key instead of email
var onlineUsers = make(map[string]OnlineUser)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func InitSocketServer() {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("New connection:", s.ID())
		return nil
	})

	// Fixed register event - expecting uid, name, email from Flutter
	server.OnEvent("/", "register", func(s socketio.Conn, data map[string]string) {
		uid := data["uid"] // Changed from "id" to "uid"
		name := data["name"]
		email := data["email"]

		// Validate required fields
		if uid == "" || name == "" || email == "" {
			fmt.Println("Invalid registration data:", data)
			s.Emit("error", "Missing required fields")
			return
		}

		// Store user with socket ID as key
		onlineUsers[s.ID()] = OnlineUser{
			UID:      uid,
			Name:     name,
			Email:    email,
			SocketID: s.ID(),
		}

		fmt.Printf("User %s (%s) registered with socket %s\n", name, uid, s.ID())
		fmt.Println("Current online users:", len(onlineUsers))

		// Send online users to the newly registered user
		s.Emit("online_users_response", getOnlineUserList())

		// Broadcast updated user list to all connected clients
		broadcastOnlineUsers(server)
	})

	// Fixed get_online_users event - supporting both with and without acknowledgment
	server.OnEvent("/", "get_online_users", func(s socketio.Conn) {
		userList := getOnlineUserList()
		fmt.Printf("Sending online users to %s: %d users\n", s.ID(), len(userList))

		// Send response back to requesting client
		s.Emit("online_users_response", userList)
	})

	// Alternative: Handle get_online_users with acknowledgment callback
	server.OnEvent("/", "get_online_users_ack", func(s socketio.Conn, ack func([]OnlineUser)) {
		userList := getOnlineUserList()
		fmt.Printf("Sending online users via ack to %s: %d users\n", s.ID(), len(userList))
		ack(userList)
	})

	server.OnEvent("/", "send_message", func(s socketio.Conn, data map[string]string) {
		from := data["from"]
		to := data["to"]
		message := data["message"]

		// Find recipient by UID instead of email
		var recipientSocketID string
		for socketID, user := range onlineUsers {
			if user.UID == to {
				recipientSocketID = socketID
				break
			}
		}

		if recipientSocketID != "" {
			fmt.Printf("Sending message from %s to %s (socket: %s)\n", from, to, recipientSocketID)
			server.BroadcastToRoom("/", recipientSocketID, "receive_message", map[string]string{
				"from":    from,
				"message": message,
				"to":      to,
			})
		} else {
			fmt.Printf("User %s not found online\n", to)
			s.Emit("error", "Recipient not online")
		}
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Printf("Socket %s disconnected: %s\n", s.ID(), reason)

		// Find and remove user by socket ID
		if user, exists := onlineUsers[s.ID()]; exists {
			fmt.Printf("User %s (%s) went offline\n", user.Name, user.UID)
			delete(onlineUsers, s.ID())

			// Broadcast updated user list
			broadcastOnlineUsers(server)
		}
	})

	// Add error handling
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Printf("Socket error for %s: %v\n", s.ID(), e)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	fmt.Println("Socket server serving at localhost:3000")
	fmt.Print(http.ListenAndServe(":3000", nil))
}

func getOnlineUserList() []OnlineUser {
	users := []OnlineUser{}
	for _, user := range onlineUsers {
		users = append(users, user)
	}
	fmt.Printf("Getting online user list: %d users\n", len(users))
	return users
}

func broadcastOnlineUsers(server *socketio.Server) {
	userList := getOnlineUserList()
	fmt.Printf("Broadcasting online users to all clients: %d users\n", len(userList))
	server.BroadcastToNamespace("/", "online_users_response", userList)
}
