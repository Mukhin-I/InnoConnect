package transport

import (
	chatpb "innoconnect/pkg/pb/chat"
	meetingpb "innoconnect/pkg/pb/meeting"
	requestpb "innoconnect/pkg/pb/request"
	userpb "innoconnect/pkg/pb/user"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
    mu      sync.RWMutex
    clients map[int64]*websocket.Conn // userID -> websocket client
}

func NewHub() *Hub {
    return &Hub{
        clients: make(map[int64]*websocket.Conn),
    }
}

func (h *Hub) AddClient(userID int64, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.clients[userID] = conn
}

func (h *Hub) RemoveClient(userID int64) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if conn, ok := h.clients[userID]; ok {
        conn.Close()
        delete(h.clients, userID)
    }
}

func (h *Hub) GetClient(userID int64) (*websocket.Conn, bool) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    conn, ok := h.clients[userID]
    return conn, ok
}

func (h *Hub) Send(userID int64, msg any) error {
    conn, ok := h.GetClient(userID)
    if !ok {
        return nil // user is offline
    }

    return conn.WriteJSON(msg)
}

type Handler struct {
	meetingClient meetingpb.MeetingServiceClient
	requestClient requestpb.RequestServiceClient
	chatClient    chatpb.ChatServiceClient
	userClient    userpb.UserServiceClient
	wsHub *Hub
}

func NewHandler(
    meetingClient meetingpb.MeetingServiceClient,
    requestClient requestpb.RequestServiceClient,
    chatClient chatpb.ChatServiceClient,
    userClient userpb.UserServiceClient,
    wsHub *Hub,
) *Handler {
    return &Handler{
        meetingClient: meetingClient,
        requestClient: requestClient,
        chatClient:    chatClient,
        userClient:    userClient,
        wsHub:         wsHub,
    }
}