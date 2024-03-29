package events

import (
	"context"
	"reflect"
	"sync"

	"github.com/mehdihadeli/go-mediatr"
)

type notifications struct {
	m      sync.RWMutex
	events map[reflect.Type][]any
}

func newNotifications() *notifications {
	return &notifications{
		m:      sync.RWMutex{},
		events: map[reflect.Type][]any{},
	}
}

func (n *notifications) add(notification any) {
	n.m.Lock()
	defer n.m.Unlock()

	notificationType := reflect.TypeOf(notification)
	if _, ok := n.events[notificationType]; ok {
		n.events[notificationType] = append(n.events[notificationType], notification)
		return
	}

	n.events[notificationType] = []any{notification}
}

func (n *notifications) clear() {
	n.m.Lock()
	defer n.m.Unlock()

	for t := range n.events {
		delete(n.events, t)
	}
}

var notificationsObj = newNotifications()

func GetEvents[TNotification any]() []TNotification {
	notificationsObj.m.RLock()
	defer notificationsObj.m.RUnlock()

	notificationType := getType[TNotification]()
	if events, ok := notificationsObj.events[notificationType]; ok {
		var result []TNotification
		for _, event := range events {
			eventValue := event.(TNotification)
			result = append(result, eventValue)
		}
		return result
	}

	return nil
}

func getType[T interface{}]() reflect.Type {
	var obj T
	return reflect.TypeOf(&obj).Elem()
}

func Clear() {
	notificationsObj.clear()
}

func Publish[TNotification any](ctx context.Context, notification TNotification) {
	// TODO(kalpio): don't swallow the error
	if err := mediatr.Publish(ctx, notification); err == nil {
		notificationsObj.add(notification)
	}
}
