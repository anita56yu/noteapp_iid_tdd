package api

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager manages WebSocket connections.
type ConnectionManager struct {
	connections map[string][]*websocket.Conn
	mutex       sync.Mutex
}

// NewConnectionManager creates a new ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string][]*websocket.Conn),
	}
}

// Add adds a new WebSocket connection for a given note ID.
func (cm *ConnectionManager) Add(noteID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.connections[noteID] = append(cm.connections[noteID], conn)
}

// Remove removes a WebSocket connection.
func (cm *ConnectionManager) Remove(noteID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connections := cm.connections[noteID]
	for i, c := range connections {
		if c == conn {
			cm.connections[noteID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	if len(cm.connections[noteID]) == 0 {
		delete(cm.connections, noteID)
	}
}

// Broadcast sends a message to all clients connected to a specific note.
func (cm *ConnectionManager) Broadcast(noteID string, message []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, conn := range cm.connections[noteID] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// Handle error, e.g., remove the connection if it's closed.
		}
	}
}

// Close closes all connections.
func (cm *ConnectionManager) Close() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for noteID, conns := range cm.connections {
		for _, conn := range conns {
			conn.Close()
		}
		delete(cm.connections, noteID)
	}
}

// CloseById closes all connections for a specific note ID.
func (cm *ConnectionManager) CloseById(noteID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, conn := range cm.connections[noteID] {
		conn.Close()
	}
	delete(cm.connections, noteID)
}
