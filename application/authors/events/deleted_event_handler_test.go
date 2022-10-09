package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthorDeletedEventHandler_Handle_LogValidMessage(t *testing.T) {
	ass := assert.New(t)
	writer := &logWriter{}
	log.SetOutput(writer)

	event := &AuthorDeletedEvent{AuthorID: uuid.New()}
	eventHandler := &AuthorDeletedEventHandler{}
	expected := fmt.Sprintf("Author deleted [%s]", event.AuthorID)
	err := eventHandler.Handle(context.Background(), event)

	ass.NoError(err)
	ass.Contains(writer.message, expected)
}
