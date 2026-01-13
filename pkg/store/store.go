package store

type Store struct {
	VersionSeq uint64
	Data       map[string][]*Entry
}

func NewStore() *Store {
	return &Store{
		Data:       make(map[string][]*Entry, 0),
		VersionSeq: 0,
	}
}
