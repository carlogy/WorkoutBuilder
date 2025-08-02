package services

import (
	"database/sql"
	"time"
)

func NoneNullToNullString(s string) sql.NullString {
	nullString := sql.NullString{String: s, Valid: true}
	return nullString
}

func NullStringToString(s sql.NullString) *string {

	if s.Valid {
		return &s.String
	}
	return nil
}

func NullTimeToTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
