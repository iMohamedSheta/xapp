package models

import "database/sql"

func nullStringToAny(ns sql.NullString) any {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func nullInt64ToAny(ni sql.NullInt64) any {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

func nullTimeToAny(nt sql.NullTime) any {
	if nt.Valid {
		return nt.Time
	}
	return nil
}
