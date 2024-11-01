package models

import (
	"database/sql"
	"strings"
)

func newNullString(s string) sql.NullString {
	if len(strings.TrimSpace(s)) == 0 {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
