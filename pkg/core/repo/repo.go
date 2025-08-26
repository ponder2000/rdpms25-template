package repository

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Deletable interface {
	DeleteWithTx(ctx context.Context, tx boil.ContextExecutor, ids ...int) error
	Delete(ctx context.Context, ids ...int) error
}

type Creatable[T any] interface {
	InsertWithTx(ctx context.Context, tx boil.ContextExecutor, newObj T, cols boil.Columns) (T, error)
	Insert(ctx context.Context, newObj T, cols boil.Columns) (T, error)
}

type Updatable[T any] interface {
	UpdateWithTx(ctx context.Context, tx boil.ContextExecutor, newObj T, cols boil.Columns) (T, error)
	Update(ctx context.Context, newObj T, cols boil.Columns) (T, error)
}

type Readable[T any] interface {
	FetchOne(ctx context.Context, q []qm.QueryMod) (T, error)
	FetchAll(ctx context.Context, q []qm.QueryMod) ([]T, error)

	FetchOneWithTx(ctx context.Context, tx boil.ContextExecutor, q []qm.QueryMod) (T, error)
	FetchAllWithTx(ctx context.Context, tx boil.ContextExecutor, q []qm.QueryMod) ([]T, error)

	Count(ctx context.Context, q []qm.QueryMod) (int64, error)
}

type CRUD[T any] interface {
	Creatable[T]
	Readable[T]
	Updatable[T]
	Deletable
}

type ReadableView[T any] interface {
	FetchOne(ctx context.Context, q []qm.QueryMod) (T, error)
	FetchAll(ctx context.Context, q []qm.QueryMod) ([]T, error)
	Count(ctx context.Context, q []qm.QueryMod) (int64, error)
}
