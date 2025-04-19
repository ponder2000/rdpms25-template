package sqlhelper

import (
	"fmt"

	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ponder2000/rdpms25-template/pkg/util/generic"
)

func IntArrayToPostgresArray(intSlice []int) string {
	idsStrSlice := generic.Mapper(intSlice, func(elm int) string { return strconv.Itoa(elm) })
	return fmt.Sprintf(`ARRAY[%s]`, strings.Join(idsStrSlice, ","))
}

func StringArrayFormat(commaSeperaetdVal string) string {
	vals := strings.Split(commaSeperaetdVal, ",")
	for i := range vals {
		vals[i] = fmt.Sprintf(`'%s'`, vals[i])
	}
	return fmt.Sprintf(`(%s)`, strings.Join(vals, ","))
}

func IntArrayFormat(commaSeperaetdVal string) string {
	return fmt.Sprintf(`(%s)`, commaSeperaetdVal)
}

func JsonContainQuery(columnName, commaSeperaetdVal string) string {
	queries := make([]string, 0)

	for _, kv := range strings.Split(commaSeperaetdVal, ",") {
		tmp := strings.Split(kv, ":")
		if len(tmp) == 2 {
			key, value := tmp[0], tmp[1]
			value = "%" + value + "%"
			queries = append(queries, fmt.Sprintf(`"%s"->>'%s' ilike '%s'`, columnName, key, value))
		} else {
			slog.Warn("Invalid json contain filter")
		}
	}
	return strings.Join(queries, " and ")
}

func SelectAlias(colName, alias string) string {
	return fmt.Sprintf(`%s as %s`, colName, alias)
}

func DateTimeString(t time.Time) string {
	return fmt.Sprintf(`%d-%d-%d %d:%d:%d`, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
