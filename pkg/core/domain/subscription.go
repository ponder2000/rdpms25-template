package domain

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type OperationType string

const (
	OPERATION_Insert OperationType = "Insert"
	OPERATION_Update OperationType = "Update"
	OPERATION_Upsert OperationType = "Upsert"
	OPERATION_Delete OperationType = "Delete"
	OPERATION_Event  OperationType = "Event"
)

type OperationEvent[T any] struct {
	Type OperationType `json:"operation"`
	Obj  T             `json:"item"`
}

type SubscriptionHandler[T any] struct {
	sync.RWMutex

	timeout       time.Duration
	subscriptions map[string]chan<- *OperationEvent[T]
}

func NewSubscriptionHandlerWithTimeout[T any](timeout time.Duration) *SubscriptionHandler[T] {
	return &SubscriptionHandler[T]{timeout: timeout, subscriptions: make(map[string]chan<- *OperationEvent[T])}
}

func NewSubscriptionHandler[T any]() *SubscriptionHandler[T] {
	return &SubscriptionHandler[T]{timeout: time.Second * 10, subscriptions: make(map[string]chan<- *OperationEvent[T])}
}

func (s *SubscriptionHandler[T]) Subscribe(id string) (<-chan *OperationEvent[T], error) {
	s.Lock()
	defer s.Unlock()

	_, alreadyExists := s.subscriptions[id]
	if alreadyExists {
		return nil, fmt.Errorf("subscription already exists with id %s", id)
	}

	ownedChannel := make(chan *OperationEvent[T], 100)
	s.subscriptions[id] = ownedChannel
	// applog.InfoF("subscribing id=%s", id)
	return ownedChannel, nil
}

func (s *SubscriptionHandler[T]) UnSubscribe(id string) error {
	s.Lock()
	defer s.Unlock()

	ownedChannel, exists := s.subscriptions[id]
	if !exists {
		return fmt.Errorf("no subscriptions found with id %s", id)
	}
	close(ownedChannel)
	delete(s.subscriptions, id)
	// applog.InfoF("unsubscribing id=%s", id)
	return nil
}

func (s *SubscriptionHandler[T]) notify(operation *OperationEvent[T]) {
	s.RLock()
	defer s.RUnlock()

	for subId, ownedChannel := range s.subscriptions {
		select {
		case ownedChannel <- operation:
		default:
			slog.Error("discarding subscription notification msg cause subscriber unable to handle backpressure", "sub_id", subId, "msg", operation)
		}
	}
}

func (s *SubscriptionHandler[T]) NotifyInsert(obj T) {
	s.notify(&OperationEvent[T]{Type: OPERATION_Insert, Obj: obj})
}

func (s *SubscriptionHandler[T]) NotifyUpdate(obj T) {
	s.notify(&OperationEvent[T]{Type: OPERATION_Update, Obj: obj})
}

func (s *SubscriptionHandler[T]) NotifyDelete(obj T) {
	s.notify(&OperationEvent[T]{Type: OPERATION_Delete, Obj: obj})
}

func (s *SubscriptionHandler[T]) NotifyCustom(t OperationType, obj T) {
	s.notify(&OperationEvent[T]{Type: t, Obj: obj})
}
