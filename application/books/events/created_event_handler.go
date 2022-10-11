package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type BookCreatedEventHandler struct {
}

func (*BookCreatedEventHandler) Handle(_ context.Context, event *BookCreatedEvent) error {
	log.Infof("Book created [%v, %s, %s, %s, %v]", event.BookID, event.Title, event.ISBN, event.Description, event.AuthorID)

	return nil
}
