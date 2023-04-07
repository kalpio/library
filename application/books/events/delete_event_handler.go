package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type BookDeletedEventHandler struct {
}

func (*BookDeletedEventHandler) Handle(_ context.Context, event *BookDeletedEvent) error {
	log.Infof("Book deleted: [%s]", event.BookID)

	return nil
}
