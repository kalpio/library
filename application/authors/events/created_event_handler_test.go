package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/random"
	"testing"
)

func TestAuthorCreatedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := &logWriter{}

	log.SetOutput(writer)

	event := &AuthorCreatedEvent{
		AuthorID:   uuid.UUID{},
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
	}
	eventHandler := &AuthorCreatedEventHandler{}
	expected := fmt.Sprintf("Author created [%s, %s, %s, %s]",
		event.AuthorID,
		event.FirstName,
		event.MiddleName,
		event.LastName)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.message, expected)
}
