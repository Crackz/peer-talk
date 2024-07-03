package member

// import (
// 	"fmt"
// 	"peer-talk/internal/hub"
// 	"peer-talk/internal/room"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// type Handler struct {
// 	hub *hub.Hub
// }

// func NewHandler(hub *hub.Hub) *Handler {
// 	return &Handler{
// 		hub: hub,
// 	}
// }

// type CreateRoomRequest struct {
// 	ID string `json:"id"`
// }

// func (h *Handler) createRoomIfNotExists(roomId string) *room.Room {

// 	existedRoom, isRoomExists := h.hub.Rooms[roomId]
// 	if isRoomExists {
// 		return existedRoom
// 	}

// 	room := room.NewRoom(roomId)

// 	h.hub.Rooms[roomId] = room
// 	return room
// }

// func (h *Handler) CreateRoom(c *gin.Context) {
// 	var req CreateRoomRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.createRoomIfNotExists(req.ID)

// 	c.JSON(http.StatusOK, req)

// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func (h *Handler) JoinRoom(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	roomId := c.Param("roomId")
// 	clientId := c.Query("userId")
// 	username := c.Query("username")

// 	h.createRoomIfNotExists(roomId)

// 	client := &Client{
// 		Conn:     conn,
// 		Message:  make(chan *Message),
// 		ID:       clientId,
// 		RoomID:   roomId,
// 		Username: username,
// 	}

// 	h.hub.Register <- client
// 	h.hub.Broadcast <- NewWelcomeMessage(roomId, username)

// 	go client.writeMessage()

// 	client.readMessage(h.hub)

// }

// type RoomResponse struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// func (h *Handler) GetRooms(c *gin.Context) {
// 	var rooms = make([]RoomResponse, 0)

// 	for _, room := range h.hub.Rooms {
// 		rooms = append(rooms, RoomResponse{
// 			ID: room.ID,
// 		})
// 	}

// 	c.JSON(http.StatusOK, rooms)
// }

// type ClientResponse struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// }

// func (h *Handler) GetClients(c *gin.Context) {
// 	roomId := c.Param("roomId")

// 	room, isRoomExist := h.hub.Rooms[roomId]
// 	if !isRoomExist {
// 		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("couldn't find room %v", roomId)})
// 		return
// 	}
// 	var clients = make([]ClientResponse, 0)
// 	for _, client := range room.Clients {
// 		clients = append(clients, ClientResponse{
// 			ID:       client.ID,
// 			Username: client.Username,
// 		})
// 	}

// 	c.JSON(http.StatusOK, clients)
// }
