package events

import "context"

// Беглая попытка в event manager
//type EventCallback func(ctx context.Context, data any) error не прокатило в mock

type EventManager struct {
	Subscribers map[string][]func(ctx context.Context, data any) error
}

func NewEventManager() *EventManager {
	return &EventManager{
		Subscribers: make(map[string][]func(ctx context.Context, data any) error),
	}
}

func (em *EventManager) Subscribe(event string, fn func(ctx context.Context, data any) error) {
	em.Subscribers[event] = append(em.Subscribers[event], fn)
}

func (em *EventManager) Trigger(ctx context.Context, event string, data any) error {
	if subscribers, ok := em.Subscribers[event]; ok {
		for _, fn := range subscribers {
			err := fn(ctx, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
