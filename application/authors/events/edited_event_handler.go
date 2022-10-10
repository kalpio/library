package events

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type AuthorEditedEventHandler struct {
}

func (*AuthorEditedEventHandler) Handle(_ context.Context, event *AuthorEditedEvent) error {
	log.Infof("Author edited [%s, %s, %s, %s]",
		event.AuthorID,
		event.FirstName,
		event.MiddleName,
		event.LastName)

	return nil
}
