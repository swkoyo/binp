package storage

type Store struct {
	db    *DBStore
	cache *CacheStore
}

func NewStore() (*Store, error) {
	dbStore, err := NewDB()
	if err != nil {
		return nil, err
	}

	cacheStore := NewCache()

	return &Store{
		db:    dbStore,
		cache: cacheStore,
	}, nil
}

func (s *Store) Init() error {
	if err := s.db.Init(); err != nil {
		return err
	}
	return nil
}
