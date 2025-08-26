package service

import (
	"context"

	"github.com/ponder2000/rdpms25-template/pkg/core/domain"
	"github.com/ponder2000/rdpms25-template/pkg/core/dto"
	"github.com/ponder2000/rdpms25-template/pkg/util/sqlhelper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Deletable interface {
	Delete(ctx context.Context, ids ...int) error
}

type Creatable[T any] interface {
	Save(ctx context.Context, newObj T) (T, error)
}

type Updatable[T any] interface {
	Edit(ctx context.Context, newObj T, cols boil.Columns) (T, error)
}

type Readable[T any] interface {
	GetOne(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (T, error)
	GetAll(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) ([]T, error)
}

type PaginatedReadable[T any] interface {
	GetPaginated(ctx context.Context, pageRequest *dto.PageRequest, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (*dto.PageResponse[T], error)
}

type CRUD[T any] interface {
	Creatable[T]
	Updatable[T]
	Readable[T]
	Deletable
}

type PaginatedCRUD[T any] interface {
	CRUD[T]
	PaginatedReadable[T]
}

type ReadableView[T any] interface {
	Readable[T]
	PaginatedReadable[T]
}

type Subscription[T any] interface {
	Subscribe(id string) (<-chan *domain.OperationEvent[T], error)
	UnSubscribe(id string) error
}
