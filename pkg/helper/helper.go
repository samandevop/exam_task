package helper

import (
	"database/sql"
	"strconv"
	"strings"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" {
			if v == "" {
				var nullStr sql.NullString
				// oldsize := len(namedQuery)
				namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))

				// if oldsize != len(namedQuery) {
				args = append(args, nullStr)
				i++
				// }
			} else {
				// oldsize := len(namedQuery)
				namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))

				// if oldsize != len(namedQuery) {
				args = append(args, v)
				i++
				// }
			}
		}
	}

	return namedQuery, args
}
