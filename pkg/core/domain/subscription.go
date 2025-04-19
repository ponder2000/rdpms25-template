package domain

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type OperationType string

const (
	OPERATION_Insert OperationType = "Insert"
	OPERATION_Update OperationType = "Update"
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

	ownedChannel := make(chan *OperationEvent[T])
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

	for id, ownedChannel := range s.subscriptions {
		go func(id string, channel chan<- *OperationEvent[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
			defer cancel()

			select {
			case <-ctx.Done():
				slog.Error("unable to send data! unsubscribing...", "err", ctx.Err().Error(), "id", id)
				_ = s.UnSubscribe(id)
			case ownedChannel <- operation:
				// slog.Debug("sent data", "id", id)
			}
		}(id, ownedChannel)
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
