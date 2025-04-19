package sqlhelper

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Loader struct {
	Relation string
}

type LoaderSlice []*Loader

func LoaderQueryBuilder(loaders LoaderSlice) []qm.QueryMod {
	queries := make([]qm.QueryMod, 0)
	for _, l := range loaders {
		queries = append(queries, qm.Load(l.Relation))
	}
	return queries
}

func LoaderAppend(loader LoaderSlice, relationName string, required string) LoaderSlice {
	if required == "true" {
		loader = append(loader, &Loader{Relation: relationName})
	}
	return loader
}
