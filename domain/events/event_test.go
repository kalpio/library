package events

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
)

type FirstDummyEvent struct {
	ID uint32
}

type FirstDummyEventHandler struct{}

func (*FirstDummyEventHandler) Handle(_ context.Context, _ *FirstDummyEvent) error {
	return nil
}

type SecondDummyEvent struct {
	ID uint32
}

type SecondDummyEventHandler struct{}

func (*SecondDummyEventHandler) Handle(_ context.Context, _ *SecondDummyEvent) error {
	return nil
}

func TestGetEventsShouldReturnNil(t *testing.T) {
	rand.Seed(time.Now().UnixMilli())
	Clear()
	ass := assert.New(t)

	events := GetEvents(&FirstDummyEvent{})

	ass.Nil(events)
}

func TestGetEventsShouldReturnFirstDummyEventOnce(t *testing.T) {
	rand.Seed(time.Now().UnixMilli())
	Clear()
	ass := assert.New(t)
	err := mediatr.RegisterNotificationHandler[*FirstDummyEvent](&FirstDummyEventHandler{})

	ass.NoError(err)

	id := rand.Uint32()
	Publish(context.Background(), &FirstDummyEvent{ID: id})

	events := GetEvents(&FirstDummyEvent{})

	ass.Len(events, 1)
	ass.Equal(id, events[0].ID)
}

func TestGetEventsShouldReturnFirstDummyEventFourTimes(t *testing.T) {
	rand.Seed(time.Now().UnixMilli())
	Clear()
	ass := assert.New(t)
	err := mediatr.RegisterNotificationHandler[*FirstDummyEvent](&FirstDummyEventHandler{})

	ass.NoError(err)

	ctx := context.Background()
	id0 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id0})

	id1 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id1})

	id2 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id2})

	id3 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id3})

	events := GetEvents(&FirstDummyEvent{})

	ass.Len(events, 4)
	ass.Equal(id0, events[0].ID)
	ass.Equal(id1, events[1].ID)
	ass.Equal(id2, events[2].ID)
	ass.Equal(id3, events[3].ID)
}

func TestGetEventsShouldReturnsCorrentEvents(t *testing.T) {
	rand.Seed(time.Now().UnixMilli())
	Clear()
	ass := assert.New(t)

	err := mediatr.RegisterNotificationHandler[*FirstDummyEvent](&FirstDummyEventHandler{})
	ass.NoError(err)

	err = mediatr.RegisterNotificationHandler[*SecondDummyEvent](&SecondDummyEventHandler{})
	ass.NoError(err)

	ctx := context.Background()
	id0 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id0})

	id1 := rand.Uint32()
	Publish(ctx, &FirstDummyEvent{ID: id1})

	id2 := rand.Uint32()
	Publish(ctx, &SecondDummyEvent{ID: id2})

	id3 := rand.Uint32()
	Publish(ctx, &SecondDummyEvent{ID: id3})

	firstDummyEvents := GetEvents(&FirstDummyEvent{})
	secondDummyEvents := GetEvents(&SecondDummyEvent{})

	ass.Len(firstDummyEvents, 2)
	ass.Len(secondDummyEvents, 2)

	ass.Equal(id0, firstDummyEvents[0].ID)
	ass.Equal(id1, firstDummyEvents[1].ID)
	ass.Equal(id2, secondDummyEvents[0].ID)
	ass.Equal(id3, secondDummyEvents[1].ID)
}
