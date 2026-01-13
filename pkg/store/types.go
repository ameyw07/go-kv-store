package store

import "time"

type Entry struct {
	Value          string
	Version        uint64
	ExpirationDate *time.Time
	HasExpiry      bool
}

func GetNewEntry(val string, version uint64, expirationDate *time.Time, hasExpiry bool) *Entry {
	return &Entry{
		Value:          val,
		Version:        version,
		ExpirationDate: expirationDate,
		HasExpiry:      hasExpiry,
	}
}
