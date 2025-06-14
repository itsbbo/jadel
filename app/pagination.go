package app

import (
	"net/url"

	"github.com/itsbbo/jadel/gonertia"
)

type PaginationMode string

const (
	PaginationPrev PaginationMode = "prev"
	PaginationNext PaginationMode = "next"

	PaginationDefaultLimit = 10
)

func ToPaginationMode(mode string) PaginationMode {
	switch mode {
	case string(PaginationPrev):
		return PaginationPrev
	case string(PaginationNext):
		return PaginationNext
	default:
		return PaginationNext
	}
}

type HasID interface {
	GetID() string
}

func ToPaginationProps[T HasID](path string, limit int, items []T) gonertia.Props {
	if len(items) == 0 {
		return gonertia.Props{
			"items":       items,
			"perPage":     limit,
			"prevPageURL": "",
			"nextPageURL": "",
		}
	}

	if len(items) < limit {
		q := url.Values{}
		q.Set("mode", string(PaginationPrev))
		q.Set("id", items[0].GetID())

		prevURL := url.URL{Path: path, RawQuery: q.Encode()}

		return gonertia.Props{
			"items":       items,
			"perPage":     limit,
			"prevPageURL": prevURL.String(),
			"nextPageURL": "",
		}
	}

	prevQuery := url.Values{}
	prevQuery.Set("mode", string(PaginationPrev))
	prevQuery.Set("id", items[0].GetID())
	prevURL := url.URL{
		Path:     path,
		RawQuery: prevQuery.Encode(),
	}

	nextQuery := url.Values{}
	nextQuery.Set("mode", string(PaginationNext))
	nextQuery.Set("id", items[len(items)-1].GetID())
	nextURL := url.URL{
		Path:     path,
		RawQuery: nextQuery.Encode(),
	}

	return gonertia.Props{
		"items":       items,
		"perPage":     limit,
		"prevPageURL": prevURL.String(),
		"nextPageURL": nextURL.String(),
	}
}
