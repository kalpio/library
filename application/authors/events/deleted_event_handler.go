package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type AuthorDeletedEventHandler struct {
}

func (*AuthorDeletedEventHandler) Handle(_ context.Context, event *AuthorDeletedEvent) error {
	log.Infof("Author deleted [%s]", event.AuthorID)

	return nil
}
