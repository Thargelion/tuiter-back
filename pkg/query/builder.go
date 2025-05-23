package query

import (
	"fmt"
	"strings"
)

type Builder string

func (q Builder) OrderBy(order string) Builder {
	return Builder(fmt.Sprintf("%s ORDER BY %s", q, order))
}

func (q Builder) Paginated(limit int, offset int) Builder {
	return Builder(fmt.Sprintf("%s LIMIT %d OFFSET %d;", q, limit, offset))
}

func (q Builder) Where(condition string) Builder {
	return Builder(fmt.Sprintf("%s WHERE %s", q, condition))
}

func (q Builder) WithParams(params Params) Builder {
	var conditions []string
	for k, v := range params {
		conditions = append(conditions, fmt.Sprintf("WHERE %s = %s", k, v))
	}

	return Builder(fmt.Sprintf("%s %s", strings.Join(conditions, " AND "), q))
}

func (q Builder) String() string {
	return string(q)
}
