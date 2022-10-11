package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/domain"
	"library/random"
	"testing"
)

func TestBookCreatedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()

	log.SetOutput(writer)

	event := NewBookCreatedEvent(
		domain.BookID(uuid.New().String()),
		random.String(20),
		random.String(20),
		random.String(120),
		domain.AuthorID(uuid.New().String()))
	eventHandler := &BookCreatedEventHandler{}
	expected := fmt.Sprintf("Book created [%v, %s, %s, %s, %v]", event.BookID, event.Title, event.ISBN, event.Description, event.AuthorID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
