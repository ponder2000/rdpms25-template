package sqlhelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ponder2000/rdpms25-template/pkg/util/generic"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ComplexFilterAppend(filters []qm.QueryMod, newFilter string, val any) []qm.QueryMod {
	return append(filters, qm.Where(newFilter, val))
}

func ParentChildRelationFilterAppend(filters []qm.QueryMod, rootIds []int, queryColName, descendantViewName, queryType string) []qm.QueryMod {
	rootArrayParam := IntArrayToPostgresArray(rootIds)

	var q qm.QueryMod
	switch queryType {
	case "direct":
		q = qm.Where(fmt.Sprintf(`%s in (select descendant_id from %s where ancestor_id = ANY(%s) and depth <= 1)`, queryColName, descendantViewName, rootArrayParam))
	case "all":
		q = qm.Where(fmt.Sprintf(`%s in (select descendant_id from %s where ancestor_id = ANY(%s) )`, queryColName, descendantViewName, rootArrayParam))
	default:
		q = qm.Where(fmt.Sprintf(`%s = ANY(%s)`, queryColName, rootArrayParam))
	}
	return append(filters, q)
}

func InnerQueryParentChildRelation(rootIds []int, descendantViewName, queryType string) string {
	rootArrayParam := IntArrayToPostgresArray(rootIds)

	var q string
	switch queryType {
	case "direct":
		q = fmt.Sprintf(`(select descendant_id from %s where ancestor_id = ANY(%s) and depth <= 1)`, descendantViewName, rootArrayParam)
	case "all":
		q = fmt.Sprintf(`(select descendant_id from %s where ancestor_id = ANY(%s) )`, descendantViewName, rootArrayParam)
	default:
		q = fmt.Sprintf(`(%s)`, strings.Join(generic.Mapper(rootIds, func(i int) string { return strconv.Itoa(i) }), ","))
	}
	return q
}

func ComplexInStringFilterAppend(filters []qm.QueryMod, newFilter string, rawVal string) []qm.QueryMod {
	if len(rawVal) <= 0 {
		return filters
	}

	arrayVal := strings.Split(rawVal, ",")
	for i := range arrayVal {
		arrayVal[i] = fmt.Sprintf("'%s'", arrayVal[i])
	}

	filters = append(filters, qm.Where(fmt.Sprintf("%s in (%s)", newFilter, strings.Join(arrayVal, ","))))
	return filters
}
