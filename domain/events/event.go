package events

import (
	"context"
	"reflect"

	"github.com/mehdihadeli/go-mediatr"
)

type notifications struct {
	events map[reflect.Type][]any
}

func newNotifications() *notifications {
	return &notifications{
		events: map[reflect.Type][]any{},
	}
}

func (n *notifications) add(notification any) {
	notificationType := reflect.TypeOf(notification)
	if _, ok := n.events[notificationType]; ok {
		n.events[notificationType] = append(n.events[notificationType], notification)
		return
	}

	n.events[notificationType] = []any{notification}
}

func (n *notifications) clear() {
	for t := range n.events {
		delete(n.events, t)
	}
}

var notificationsObj = newNotifications()

func GetEvents[TNotification any](forType TNotification) []TNotification {
	notificationType := reflect.TypeOf(forType)
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

func Clear() {
	notificationsObj.clear()
}

func Publish[TNotification any](ctx context.Context, notification TNotification) {
	if err := mediatr.Publish(ctx, notification); err == nil {
		notificationsObj.add(notification)
	}
}
