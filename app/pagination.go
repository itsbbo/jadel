package app

import (
	"github.com/itsbbo/jadel/gonertia"
	"github.com/oklog/ulid/v2"
)

const (
	PaginationDefaultLimit = 10
)

type PaginationRequest struct {
	Limit  int
	PrevID ulid.ULID
	NextID ulid.ULID
	UserID ulid.ULID
}

type HasID interface {
	GetID() string
}

func ToPaginationProps[T HasID](param PaginationRequest, items []T) gonertia.Props {
	limit := param.Limit
	if limit <= 0 {
		limit = PaginationDefaultLimit
	}

	if len(items) == 0 || (!param.PrevID.IsZero()) && (!param.NextID.IsZero()) {
		return gonertia.Props{
			"items":  []T{},
			"prevId": "",
			"nextId": "",
		}
	}

	// no prev and next = first page
	if param.PrevID.IsZero() && param.NextID.IsZero() {
		if len(items) > limit {
			items = items[:limit]

			return gonertia.Props{
				"items":  items,
				"prevId": "",
				"nextId": items[len(items)-1].GetID(),
			}
		}

		return gonertia.Props{
			"items":  items,
			"prevId": "",
			"nextId": "",
		}
	}

	// next mode
	if param.PrevID.IsZero() {
		if len(items) > limit {
			items = items[:limit]

			return gonertia.Props{
				"items":  items,
				"prevId": items[0].GetID(),
				"nextId": items[len(items)-1].GetID(),
			}
		}

		return gonertia.Props{
			"items":  items,
			"prevId": items[0].GetID(),
			"nextId": "",
		}
	}

	// prev mode
	if len(items) > limit {
		items = items[1:]

		return gonertia.Props{
			"items":  items,
			"prevId": items[0].GetID(),
			"nextId": items[len(items)-1].GetID(),
		}
	}

	return gonertia.Props{
		"items":  items,
		"prevId": "",
		"nextId": items[len(items)-1].GetID(),
	}
}
