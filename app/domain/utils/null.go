package utils

import (
	"database/sql"
	"time"
)

func NullStringToAny(ns sql.NullString) any {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func NullInt64ToAny(ni sql.NullInt64) any {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func NullTimeToAny(nt sql.NullTime) any {
	if nt.Valid {
		return nt.Time
	}
	return nil
}

func ToNullInt64(val *int64) sql.NullInt64 {
	if val != nil && *val != 0 {
		return sql.NullInt64{Int64: *val, Valid: true}
	}
	return sql.NullInt64{Valid: false}
}

func ToNullTime(val *time.Time) sql.NullTime {
	if val != nil {
		return sql.NullTime{Time: *val, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func ToNullString(val string) sql.NullString {
	if val != "" {
		return sql.NullString{String: val, Valid: true}
	}
	return sql.NullString{Valid: false}
}
