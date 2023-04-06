package events_test

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/application/books/events"
	"library/domain"
	"library/random"
	"testing"
)

func TestBookEditedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()

	log.SetOutput(writer)

	event := &events.BookEditedEvent{
		BookID:      domain.NewBookID(),
		Title:       random.String(120),
		ISBN:        random.String(13),
		Description: random.String(300),
		AuthorID:    domain.NewAuthorID(),
	}
	eventHandler := &events.BookEditedEventHandler{}
	expected := fmt.Sprintf("Book edited [%v, %s, %s, %s, %v]", event.BookID, event.Title, event.ISBN, event.Description, event.AuthorID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
