package http

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var activeWebsockets = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_websockets",
		Help: "Number of currently active WebSocket connections",
	},
)

type WebSocketHandler struct {
	useCases  UseCases
	closeChan chan struct{}
}

func NewWebSocketHandler(useCases UseCases) *WebSocketHandler {
	return &WebSocketHandler{
		useCases:  useCases,
		closeChan: make(chan struct{}, 1),
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context, id int64) error {
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	ctx := conn.CloseRead(c)
	ticker := time.NewTicker(1 * time.Second)

	activeWebsockets.Inc()

	go func() {
		defer func() {
			err := conn.Close(websocket.StatusNormalClosure, "Connection is closed")
			if err != nil {
				log.Printf("error while closing connection: %v", err)
			}
			activeWebsockets.Dec()
		}()
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-h.closeChan:
				return
			case <-ticker.C:
				event, err := h.useCases.Event.GetLastEventBySensorID(c, id)
				if err != nil {
					continue
				}

				err = wsjson.Write(c, conn, event)
				if err != nil {
					log.Printf("error while sending json: %v", err)
					return
				}
			}
		}
	}()

	return nil
}

func (h *WebSocketHandler) Shutdown() error {
	defer close(h.closeChan)

	h.closeChan <- struct{}{}
	return nil
}
