package pkg

type Entry struct {
	val       any
	version   uint64
	isExpired bool
}

func GetNewEntry(val any, version uint64, isExpired bool) *Entry {
	return &Entry{
		val:       val,
		version:   version,
		isExpired: isExpired,
	}
}
