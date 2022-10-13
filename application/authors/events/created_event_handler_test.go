package events_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/application/authors/events"
	"library/random"
	"testing"
)

func TestAuthorCreatedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()

	log.SetOutput(writer)

	event := &events.AuthorCreatedEvent{
		AuthorID:   uuid.UUID{},
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
	}
	eventHandler := &events.AuthorCreatedEventHandler{}
	expected := fmt.Sprintf("Author created [%s, %s, %s, %s]",
		event.AuthorID,
		event.FirstName,
		event.MiddleName,
		event.LastName)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
