package events_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/application/authors/events"
	"testing"
)

func TestAuthorDeletedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := application.NewTestLogWriter()
	log.SetOutput(writer)

	event := &events.AuthorDeletedEvent{AuthorID: uuid.New()}
	eventHandler := &events.AuthorDeletedEventHandler{}
	expected := fmt.Sprintf("Author deleted [%s]", event.AuthorID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.GetMessage(), expected)
}
