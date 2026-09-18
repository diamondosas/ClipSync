package events

import (
	"sync"
	"testing"
	"time"
)

func TestEventBus_PublishAndSubscribe(t *testing.T) {
	bus := NewEventBus()

	var receivedClips []string
	var mu sync.Mutex

	sub := bus.SubscribeClips(func(evt ClipReceivedEvent) {
		mu.Lock()
		defer mu.Unlock()
		receivedClips = append(receivedClips, evt.Content)
	})
	defer sub.Unsubscribe()

	bus.PublishClipReceived(ClipReceivedEvent{
		Content: "Hello from peer device",
		FromIP:  "192.168.1.50",
	})

	// Give a brief moment for async delivery if channels are used
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(receivedClips) != 1 || receivedClips[0] != "Hello from peer device" {
		t.Fatalf("Expected to receive clip 'Hello from peer device', got: %v", receivedClips)
	}
}

func TestEventBus_DeviceEvents(t *testing.T) {
	bus := NewEventBus()

	var discovered []string
	var lost []string
	var mu sync.Mutex

	sub1 := bus.SubscribeDevices(func(evt DeviceEvent) {
		mu.Lock()
		defer mu.Unlock()
		if evt.Type == DeviceConnected {
			discovered = append(discovered, evt.Device.Name)
		} else if evt.Type == DeviceDisconnected {
			lost = append(lost, evt.Device.Name)
		}
	})
	defer sub1.Unsubscribe()

	bus.PublishDeviceEvent(DeviceEvent{
		Type: DeviceConnected,
		Device: DeviceInfo{
			Name: "MacBook Pro",
			IP:   "192.168.1.100",
		},
	})

	bus.PublishDeviceEvent(DeviceEvent{
		Type: DeviceDisconnected,
		Device: DeviceInfo{
			Name: "MacBook Pro",
			IP:   "192.168.1.100",
		},
	})

	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(discovered) != 1 || discovered[0] != "MacBook Pro" {
		t.Errorf("Expected discovered MacBook Pro, got: %v", discovered)
	}
	if len(lost) != 1 || lost[0] != "MacBook Pro" {
		t.Errorf("Expected lost MacBook Pro, got: %v", lost)
	}
}

func TestEventBus_ToastEvents(t *testing.T) {
	bus := NewEventBus()

	var toastMsg string
	var mu sync.Mutex

	sub := bus.SubscribeToast(func(msg string) {
		mu.Lock()
		defer mu.Unlock()
		toastMsg = msg
	})
	defer sub.Unsubscribe()

	bus.PublishToast("Synced new clip from Linux PC")

	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if toastMsg != "Synced new clip from Linux PC" {
		t.Fatalf("Expected toast msg 'Synced new clip from Linux PC', got: %q", toastMsg)
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := NewEventBus()

	var count int
	var mu sync.Mutex

	sub := bus.SubscribeToast(func(msg string) {
		mu.Lock()
		defer mu.Unlock()
		count++
	})

	bus.PublishToast("First")
	time.Sleep(10 * time.Millisecond)

	sub.Unsubscribe()

	bus.PublishToast("Second")
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("Expected exactly 1 event after unsubscribe, got: %d", count)
	}
}
