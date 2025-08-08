package services

import (
	"database/sql"
	"time"
)

func NoneNullToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{String: "", Valid: false}
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

func NullInttoInt(i sql.NullInt64) *int {
	if i.Valid {
		num := int(i.Int64)
		return &num
	}
	return nil
}

func NoneNullIntToNullInt(i *int) sql.NullInt64 {

	if i != nil {
		return sql.NullInt64{
			Int64: int64(*i),
			Valid: true,
		}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}
