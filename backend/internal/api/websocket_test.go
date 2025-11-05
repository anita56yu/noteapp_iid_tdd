package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestConnectionManager(t *testing.T) {
	// Arrange
	cm := NewConnectionManager()
	// defer cm.Close()
	noteID := "test-note"

	// Create a test server to handle WebSocket upgrades
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		// When the test is over, the connection will be closed.
		// We need to make sure the connection is removed from the manager.
		conn.SetCloseHandler(func(code int, text string) error {
			cm.Remove(noteID, conn)
			return nil
		})
		cm.Add(noteID, conn)
	}))
	defer server.Close()

	// Create two WebSocket clients
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer conn2.Close()

	// Give the server a moment to add the connections
	time.Sleep(100 * time.Millisecond)

	// Act & Assert: Broadcast to both connections
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, msg, err := conn1.ReadMessage()
		if err != nil {
			// Don't fail the test if the connection is closed before the message is read
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				t.Errorf("conn1 failed to read message: %v", err)
			}
			return
		}
		if string(msg) != "hello" {
			t.Errorf("conn1 expected 'hello', got '%s'", string(msg))
		}
	}()
	go func() {
		defer wg.Done()
		_, msg, err := conn2.ReadMessage()
		if err != nil {
			t.Errorf("conn2 failed to read message: %v", err)
			return
		}
		if string(msg) != "hello" {
			t.Errorf("conn2 expected 'hello', got '%s'", string(msg))
		}
	}()

	cm.Broadcast(noteID, []byte("hello"))
	wg.Wait()

	// Act & Assert: Remove one connection and broadcast again
	conn1.Close()
	// Give the server a moment to process the close handler
	time.Sleep(100 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, msg, err := conn2.ReadMessage()
		if err != nil {
			t.Errorf("conn2 failed to read message after conn1 removed: %v", err)
			return
		}
		if string(msg) != "world" {
			t.Errorf("conn2 expected 'world', got '%s'", string(msg))
		}
	}()

	cm.Broadcast(noteID, []byte("world"))
	wg.Wait()
}
