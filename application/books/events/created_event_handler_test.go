package events_test

import (
	"context"
	"fmt"
	"library/application"
	"library/application/books/events"
	"library/domain"
	"library/random"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBookCreatedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()

	log.SetOutput(writer)

	event := events.NewBookCreatedEvent(
		domain.NewBookID(),
		random.String(20),
		random.String(20),
		random.String(120),
		domain.NewAuthorID())
	eventHandler := &events.BookCreatedEventHandler{}
	expected := fmt.Sprintf("Book created [%v, %s, %s, %s, %v]", event.BookID, event.Title, event.ISBN, event.Description, event.AuthorID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
