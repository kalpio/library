package events_test

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/application/books/events"
	"library/domain"
	"testing"
)

func TestBookDeletedEventHandler_Handle(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()
	log.SetOutput(writer)

	event := &events.BookDeletedEvent{BookID: domain.NewBookID()}
	eventHandler := &events.BookDeletedEventHandler{}
	expected := fmt.Sprintf("Book deleted: [%s]", event.BookID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
