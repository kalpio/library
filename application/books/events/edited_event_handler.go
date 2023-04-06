package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type BookEditedEventHandler struct {
}

func (*BookEditedEventHandler) Handle(_ context.Context, event *BookEditedEvent) error {
	log.Infof("Book edited [%v, %s, %s, %s, %v]", event.BookID, event.Title, event.ISBN, event.Description, event.AuthorID)

	return nil
}
