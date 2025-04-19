package sqlhelper

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ponder2000/rdpms25-template/pkg/util/parser"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Filter struct {
	ColName  string
	ColVal   interface{}
	Operator string
	IsArray  bool
	IsNum    bool
}

type FilterSlice []*Filter

func (f *Filter) whereValue() (string, error) {
	switch f.ColVal.(type) {
	case int:
		return fmt.Sprintf("%d", f.ColVal.(int)), nil
	case string:
		pref := "'"
		suffix := "'"
		if f.isLikeOperator() {
			pref = "'%"
			suffix = "%'"
		}
		return pref + f.ColVal.(string) + suffix, nil
	case bool:
		return fmt.Sprintf("%v", f.ColVal), nil
	default:
		return "", fmt.Errorf("invalid Type for filter")
	}
}

func (f *Filter) isLikeOperator() bool {
	return strings.Contains(f.Operator, "like")

}

func FilterQueryBuilder(table string, filters FilterSlice) []qm.QueryMod {
	queries := make([]qm.QueryMod, 0)
	for _, f := range filters {
		if val, e := f.whereValue(); e != nil {
			slog.Error("unable to filter", "filter", f)
		} else {
			if f.IsArray {
				val = val[1 : len(val)-1]
				arrayVal := strings.Split(val, ",")
				if !f.IsNum {
					for i := range arrayVal {
						arrayVal[i] = "'" + arrayVal[i] + "'"
					}
				}
				queries = append(queries, qm.Where(fmt.Sprintf(`"%s"."%s" %s (%s)`, table, f.ColName, f.Operator, strings.Join(arrayVal, ","))))
			} else {
				queries = append(queries, qm.Where(fmt.Sprintf(`"%s"."%s" %s %s`, table, f.ColName, f.Operator, val)))
			}

		}
	}
	slog.Debug("query formed", "table", table, "query", queries)
	return queries
}

// StringAppendToFilter
// helper append function for string queries
func StringAppendToFilter(filters FilterSlice, colName, operator, colValue string) FilterSlice {
	if len(colValue) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValue, Operator: operator})
	}
	return filters
}

// BoolAppendToFilter
// helper append function for bool queries
func BoolAppendToFilter(filters FilterSlice, colName, operator string, colValue string) FilterSlice {
	if v, e := parser.ExtractBool(colValue); e != nil {
		return filters
	} else {
		return append(filters, &Filter{ColName: colName, ColVal: v, Operator: operator})
	}
}

// IntAppendToFilter
// helper append function for string queries
func IntAppendToFilter(filters FilterSlice, colName, operator string, colValue string) FilterSlice {
	if num, e := parser.ExtractInt(colValue, true); e != nil {
		return filters
	} else {
		return append(filters, &Filter{ColName: colName, ColVal: num, Operator: operator})
	}
}

func InStringsAppendToFilter(filters FilterSlice, colName, colValues string) FilterSlice {
	if len(colValues) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValues, Operator: "in", IsArray: true})
	}
	return filters
}

func InIntAppendToFilter(filters FilterSlice, colName, colValues string) FilterSlice {
	if len(colValues) > 0 {
		filters = append(filters, &Filter{ColName: colName, ColVal: colValues, Operator: "in", IsArray: true, IsNum: true})
	}
	return filters
}
