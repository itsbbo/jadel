package app

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

type MockItem struct {
	ID string
}

func (m MockItem) GetID() string {
	return m.ID
}

func createMockItems(count int) []MockItem {
	items := make([]MockItem, count)
	for i := range count {
		items[i] = MockItem{ID: ulid.Make().String()}
	}
	return items
}

func TestToPaginationProps(t *testing.T) {
	tests := []struct {
		name           string
		param          PaginationRequest
		itemCount      int
		expectedLen    int
		hasPrevID      bool
		hasNextID      bool
		prevIDIndex    int  // -1 means no prevId
		nextIDIndex    int  // -1 means no nextId
	}{
		{
			name: "Empty items",
			param: PaginationRequest{
				Limit: 10,
			},
			itemCount:   0,
			expectedLen: 0,
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "Both PrevID and NextID are set (invalid state)",
			param: PaginationRequest{
				Limit:  10,
				PrevID: ulid.Make(),
				NextID: ulid.Make(),
			},
			itemCount:   5,
			expectedLen: 0,
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "First page with more items than limit",
			param: PaginationRequest{
				Limit: 10,
			},
			itemCount:   15,
			expectedLen: 10,
			hasPrevID:   false,
			hasNextID:   true,
			prevIDIndex: -1,
			nextIDIndex: 9,
		},
		{
			name: "First page with fewer items than limit",
			param: PaginationRequest{
				Limit: 10,
			},
			itemCount:   5,
			expectedLen: 5,
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "Next page with more items than limit",
			param: PaginationRequest{
				Limit:  10,
				NextID: ulid.Make(),
			},
			itemCount:   15,
			expectedLen: 10,
			hasPrevID:   true,
			hasNextID:   true,
			prevIDIndex: 0,
			nextIDIndex: 9,
		},
		{
			name: "Next page with fewer items than limit",
			param: PaginationRequest{
				Limit:  10,
				NextID: ulid.Make(),
			},
			itemCount:   5,
			expectedLen: 5,
			hasPrevID:   true,
			hasNextID:   false,
			prevIDIndex: 0,
			nextIDIndex: -1,
		},
		{
			name: "Prev page with more items than limit",
			param: PaginationRequest{
				Limit:  10,
				PrevID: ulid.Make(),
			},
			itemCount:   15,
			expectedLen: 14, // 15 - 1 (skip first item in prev mode)
			hasPrevID:   true,
			hasNextID:   true,
			prevIDIndex: 1,
			nextIDIndex: 14,
		},
		{
			name: "Prev page with fewer items than limit",
			param: PaginationRequest{
				Limit:  10,
				PrevID: ulid.Make(),
			},
			itemCount:   5,
			expectedLen: 5,
			hasPrevID:   false,
			hasNextID:   true,
			prevIDIndex: -1,
			nextIDIndex: 4,
		},
		{
			name: "Edge case: exactly at limit",
			param: PaginationRequest{
				Limit: 10,
			},
			itemCount:   10,
			expectedLen: 10,
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "Edge case: single item",
			param: PaginationRequest{
				Limit: 10,
			},
			itemCount:   1,
			expectedLen: 1,
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "Edge case: limit is zero",
			param: PaginationRequest{
				Limit: 0,
			},
			itemCount:   5,
			expectedLen: 5, // Default limit is used
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
		{
			name: "Edge case: limit is negative",
			param: PaginationRequest{
				Limit: -1,
			},
			itemCount:   5,
			expectedLen: 5, // Default limit is used
			hasPrevID:   false,
			hasNextID:   false,
			prevIDIndex: -1,
			nextIDIndex: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items := createMockItems(tc.itemCount)
			props := ToPaginationProps(tc.param, items)

			assert.Len(t, props["items"], tc.expectedLen)

			// Check prevId
			if tc.hasPrevID {
				assert.Equal(t, items[tc.prevIDIndex].ID, props["prevId"])
			} else {
				assert.Equal(t, "", props["prevId"])
			}

			// Check nextId
			if tc.hasNextID {
				assert.Equal(t, items[tc.nextIDIndex].ID, props["nextId"])
			} else {
				assert.Equal(t, "", props["nextId"])
			}
		})
	}
}
