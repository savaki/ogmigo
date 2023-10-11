package ogmigo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func timeout(delay time.Duration) http.HandlerFunc {
	var upgrader = websocket.Upgrader{} // use default options
	return func(w http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			select {
			case <-req.Context().Done():
			case <-time.After(delay):
			}

			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func TestClient_query(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	defer listener.Close()

	go func() {
		_ = http.Serve(listener, timeout(time.Minute))
	}()

	parts := strings.Split(listener.Addr().String(), ":")
	port := parts[len(parts)-1]

	client := New(WithEndpoint(fmt.Sprintf("ws://127.0.0.1:%v", port)))

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err = client.SubmitTx(ctx, "fffefdfc")
	if ok := errors.Is(err, context.DeadlineExceeded); !ok {
		t.Fatalf("expected context.Canceled; got %v", err)
	}
}
