package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/random"
	"testing"
)

func TestAuthorEditedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()

	log.SetOutput(writer)

	event := &AuthorEditedEvent{
		AuthorID:   uuid.UUID{},
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
	}
	eventHandler := &AuthorEditedEventHandler{}
	expected := fmt.Sprintf("Author edited [%s, %s, %s, %s]",
		event.AuthorID,
		event.FirstName,
		event.MiddleName,
		event.LastName)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
