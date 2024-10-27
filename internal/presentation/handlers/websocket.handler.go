package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
)

type WebSocketHandler struct {
	Upgrader    websocket.Upgrader
	MongoClient *mongo.Client
	Stop        context.CancelFunc
}

var totalConnections int
var maxTotalConnections = 15 // 設定最大總連線數
var totalConnectionsMutex sync.Mutex

// 連線池
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte, 100)
var clientsMutex sync.Mutex // 保護 clients map

var activeConnections = make(map[string]int)

var allowedOrigins = map[string]bool{
	"http://localhost":      true,
	"http://localhost:8080": true,
}

func NewWebSocketHandler(stop context.CancelFunc, mongoClient *mongo.Client) *WebSocketHandler {
	// 启动消息广播处理
	go handleBroadcast()

	return &WebSocketHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
		MongoClient: mongoClient,
		Stop:        stop,
	}
}

func (h *WebSocketHandler) HandleWebSocket(ctx *gin.Context) {
	conn, err := h.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.Error(ctx.Writer, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func() {
		// 減少計數器 & 修改連接池
		clientIP := ctx.Request.RemoteAddr
		activeConnections[clientIP]--

		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()

		err := conn.Close()
		if err != nil {
			return
		}
	}()

	// 增加連接池
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// 取得 MongoDB chatdb 中 messages collection
	messagesCollection := h.MongoClient.Database(viper.GetString("MONGO_DB_NAME")).Collection(viper.GetString("MONGO_COLLECTION_NAME"))
	cursor, err := messagesCollection.Find(context.TODO(), bson.M{})
	if err == nil {
		var messages []bson.M
		if err = cursor.All(context.TODO(), &messages); err == nil {
			for _, msg := range messages {
				content, ok := msg["content"].(string)
				if !ok {
					continue
				}
				err := conn.WriteMessage(websocket.TextMessage, []byte(content))
				if err != nil {
					log.Error().Err(err).Msg("Could not write history message")
					// graceful shutdown
					h.Stop()
					return
				}
			}
		}
	}

	// WebSocket 通信逻辑
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("Could not read message")
			break
		}

		// 保存消息到 MongoDB
		_, err = messagesCollection.InsertOne(context.TODO(), bson.M{"content": string(msg)})
		if err != nil {
			log.Error().Err(err).Msg("Could not save message")
			break
		}

		// 将消息发送到 broadcast 通道
		broadcast <- msg
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	log.Debug().Msgf("Origin: %s", origin)
	if !allowedOrigins[origin] {
		return false
	}

	// 每個 IP 限制最多 5 個連線
	clientIP := r.RemoteAddr
	if activeConnections[clientIP] >= 5 {
		return false
	}

	// 總連線數限制
	totalConnectionsMutex.Lock()
	defer totalConnectionsMutex.Unlock()

	if totalConnections >= maxTotalConnections {
		// 超過最大連線數
		log.Warn().Msgf("Connection limit reached, refusing connection from IP: %s", clientIP)
		return false
	}

	// 計數器們
	activeConnections[clientIP]++
	totalConnections++

	return true
}

//func handleBroadcast() {
//	for {
//		msg := <-broadcast
//		// 廣播給所有連線進來的客戶端
//		clientsMutex.Lock()
//		for client := range clients {
//			err := client.WriteMessage(websocket.TextMessage, msg)
//			if err != nil {
//				log.Error().Err(err).Msg("Could not write message to client")
//				err := client.Close()
//				if err != nil {
//					return
//				}
//				delete(clients, client)
//			}
//		}
//		clientsMutex.Unlock()
//	}
//}

func handleBroadcast() {
	for {
		msg := <-broadcast
		clientsMutex.Lock()
		clientsCopy := make(map[*websocket.Conn]bool)
		for client := range clients {
			clientsCopy[client] = true
		}
		clientsMutex.Unlock()

		// 廣播給所有連線進來的客戶端
		for client := range clientsCopy {
			go func(client *websocket.Conn) {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Error().Err(err).Msg("Could not write message to client")
					clientsMutex.Lock()
					client.Close()
					delete(clients, client)
					clientsMutex.Unlock()
				}
			}(client)
		}
	}
}
