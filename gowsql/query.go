package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/thoas/go-funk"
)

func BuildMultipleNamedValues(keys string, values []interface{}) (query string, args []interface{}, err error) {
	if !strings.ContainsRune(keys, ':') {
		return "", nil, errors.New("Did not find named parameter in keys")
	}
	queryParts := []string{}
	for _, a := range values {
		q, a, err := sqlx.Named(keys, a)
		if err != nil {
			return "", nil, err
		}
		queryParts = append(queryParts, q)
		args = append(args, a...)
	}
	return strings.Join(queryParts, ","), args, err
}

// PrefixColumns adds a table "prefix" before each column found in
// comma-separated "fields"
func PrefixColumns(fields, prefix string) string {
	if len(fields) <= 0 {
		return ""
	}

	tokens := strings.Split(fields, ",")
	prefixedTokens := funk.Map(tokens, func(token string) string {
		return fmt.Sprintf("%s.%s", prefix, token)
	}).([]string)

	return strings.Join(prefixedTokens, ",")
}

func BuildNamedColumns(columns string) string {
	arrColumns := strings.Split(columns, ",")

	if len(arrColumns) == 1 && arrColumns[0] == "" {
		return ""
	}

	var s string
	for _, col := range arrColumns {
		s += ":" + col + ","
	}
	s = strings.TrimSuffix(s, ",")
	return s
}

type OnDuplicateMode int

var OnDuplicateOverride OnDuplicateMode = 0
var OnDuplicateIncrement OnDuplicateMode = 1

func BuildOnDuplicateClause(columns string, mode OnDuplicateMode) string {
	arrColumns := strings.Split(columns, ",")

	if len(arrColumns) == 1 && arrColumns[0] == "" {
		return ""
	}

	var s string
	for _, col := range arrColumns {
		switch mode {
		case OnDuplicateIncrement:
			s += fmt.Sprintf("\t\t%s = %s + VALUES(%s),\n", col, col, col)
			break
		case OnDuplicateOverride:
			fallthrough
		default:
			s += fmt.Sprintf("\t\t%s = VALUES(%s),\n", col, col)
		}
	}
	s = strings.TrimSuffix(s, ",\n")

	return s
}

func BuildOnDuplicateClauseIncrement(columns string, srcTable string) string {
	arrColumns := strings.Split(columns, ",")

	if len(arrColumns) == 1 && arrColumns[0] == "" {
		return ""
	}

	if srcTable != "" {
		srcTable += "."
	}

	var s string
	for _, col := range arrColumns {
		s += fmt.Sprintf("\t\t%s = %s + VALUES(%s),\n", col, srcTable+col, col)
	}
	s = strings.TrimSuffix(s, ",\n")

	return s
}
