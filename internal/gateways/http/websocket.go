package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

	go func() {
		defer conn.Close(websocket.StatusNormalClosure, "Connection is closed")
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
					return
				}

				err = wsjson.Write(c, conn, event)
				if err != nil {
					return
				}
			}
		}
	}()

	return nil
}

func (h *WebSocketHandler) Shutdown() error {
	h.closeChan <- struct{}{}

	return nil
}
