package pkg

import (
	"time"
)

type Entry struct {
	val            any
	version        uint64
	expirationDate time.Time
}

func GetNewEntry(val any, version uint64, expirationDate time.Time) *Entry {
	return &Entry{
		val:            val,
		version:        version,
		expirationDate: expirationDate,
	}
}
