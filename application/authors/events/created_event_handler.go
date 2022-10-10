package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type AuthorCreatedEventHandler struct {
}

func (*AuthorCreatedEventHandler) Handle(_ context.Context, event *AuthorCreatedEvent) error {
	log.Infof("Author created [%s, %s, %s, %s]",
		event.AuthorID,
		event.FirstName,
		event.MiddleName,
		event.LastName)

	return nil
}
