package dto

import (
	"fmt"

	"github.com/ponder2000/rdpms25-template/pkg/util/generic"
)

type PageResponse[T any] struct {
	TotalRecords int64 `json:"total_records"`
	Offset       int   `json:"offset"`
	Limit        int   `json:"limit"`
	Items        []T   `json:"items"`
}

func TransformPageResponse[U any, V any](resp *PageResponse[U], transform func(U) V) *PageResponse[V] {
	return &PageResponse[V]{
		Limit:        resp.Limit,
		Offset:       resp.Offset,
		TotalRecords: resp.TotalRecords,
		Items:        generic.Mapper(resp.Items, transform),
	}
}

type PageRequest struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	PageColumn string `json:"page_column"`
	Order      string `json:"order"`
}

func (r *PageRequest) UpdateDefaultIfNotExist(pageCol, order string) {
	if r.PageColumn == "" {
		r.PageColumn = pageCol
	}
	if r.Order == "" {
		r.Order = order
	}
}

func (r *PageRequest) OrderByClause() string {
	return fmt.Sprintf("%s %s", r.PageColumn, r.Order)
}
