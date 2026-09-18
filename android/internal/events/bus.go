package events

import (
	"sync"
	"sync/atomic"
)

type DeviceEventType int

const (
	DeviceConnected DeviceEventType = iota
	DeviceDisconnected
)

type DeviceInfo struct {
	Name string
	IP   string
}

type DeviceEvent struct {
	Type   DeviceEventType
	Device DeviceInfo
}

type ClipReceivedEvent struct {
	Content string
	FromIP  string
}

type Subscription interface {
	Unsubscribe()
}

type subHandle struct {
	id     uint64
	bus    *EventBus
	cancel func()
}

func (s *subHandle) Unsubscribe() {
	if s.cancel != nil {
		s.cancel()
	}
}

type EventBus struct {
	mu           sync.RWMutex
	nextID       uint64
	clipSubs     map[uint64]func(ClipReceivedEvent)
	deviceSubs   map[uint64]func(DeviceEvent)
	toastSubs    map[uint64]func(string)
}

func NewEventBus() *EventBus {
	return &EventBus{
		clipSubs:   make(map[uint64]func(ClipReceivedEvent)),
		deviceSubs: make(map[uint64]func(DeviceEvent)),
		toastSubs:  make(map[uint64]func(string)),
	}
}

func (b *EventBus) nextSubID() uint64 {
	return atomic.AddUint64(&b.nextID, 1)
}

func (b *EventBus) SubscribeClips(fn func(ClipReceivedEvent)) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.nextSubID()
	b.clipSubs[id] = fn
	return &subHandle{
		id:  id,
		bus: b,
		cancel: func() {
			b.mu.Lock()
			delete(b.clipSubs, id)
			b.mu.Unlock()
		},
	}
}

func (b *EventBus) PublishClipReceived(evt ClipReceivedEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, fn := range b.clipSubs {
		go fn(evt)
	}
}

func (b *EventBus) SubscribeDevices(fn func(DeviceEvent)) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.nextSubID()
	b.deviceSubs[id] = fn
	return &subHandle{
		id:  id,
		bus: b,
		cancel: func() {
			b.mu.Lock()
			delete(b.deviceSubs, id)
			b.mu.Unlock()
		},
	}
}

func (b *EventBus) PublishDeviceEvent(evt DeviceEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, fn := range b.deviceSubs {
		go fn(evt)
	}
}

func (b *EventBus) SubscribeToast(fn func(string)) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.nextSubID()
	b.toastSubs[id] = fn
	return &subHandle{
		id:  id,
		bus: b,
		cancel: func() {
			b.mu.Lock()
			delete(b.toastSubs, id)
			b.mu.Unlock()
		},
	}
}

func (b *EventBus) PublishToast(msg string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, fn := range b.toastSubs {
		go fn(msg)
	}
}
