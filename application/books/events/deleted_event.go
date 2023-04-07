package events

import "library/domain"

type BookDeletedEvent struct {
	BookID domain.BookID
}
