package usecase

import (
	"noteapp/internal/usecase/noteuc"
	"sync"
	"testing"
)

func TestEventBus_Init(t *testing.T) {
	// Arrange
	eb := NewEventBus[noteuc.NoteEventPO]()

	// Assert
	if eb == nil {
		t.Error("Expected EventBus instance to be non-nil after initialization")
	}
}

func TestEventBus_SubscribeAndPublish(t *testing.T) {
	// Arrange
	eb := NewEventBus[noteuc.NoteEventPO]()
	var wg sync.WaitGroup
	wg.Add(2)
	eventReceived := false

	subscriber := func(event noteuc.NoteEventPO) {
		defer wg.Done()
		if event.EventType == "test-event" && event.NoteID == "note-1" {
			eventReceived = true
		}
	}

	// Act
	eb.Subscribe(subscriber)
	go func() {
		defer wg.Done()
		eb.Publish(noteuc.NoteEventPO{EventType: "test-event", NoteID: "note-1"})
	}()

	// Assert
	wg.Wait() // Wait for the subscriber to complete
	if !eventReceived {
		t.Error("Expected test event to be received with correct data")
	}
}

func TestEventBus_SubscribeAndPublish_MultipleSubscribers(t *testing.T) {
	// Arrange
	eb := NewEventBus[noteuc.NoteEventPO]()
	var wg sync.WaitGroup
	wg.Add(3)
	eventReceived := 0

	subscriber := func(event noteuc.NoteEventPO) {
		defer wg.Done()
		if event.EventType == "test-event" && event.NoteID == "note-1" {
			eventReceived += 1
		}
	}

	// Act
	eb.Subscribe(subscriber)
	eb.Subscribe(subscriber)
	go func() {
		defer wg.Done()
		eb.Publish(noteuc.NoteEventPO{EventType: "test-event", NoteID: "note-1"})
	}()

	// Assert
	wg.Wait() // Wait for the subscriber to complete
	if eventReceived != 2 {
		t.Error("Expected test event to be received with correct data")
	}
}
