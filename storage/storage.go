package storage

type Store struct {
	db *DBStore
}

func NewStore() (*Store, error) {
	dbStore, err := NewDB()
	if err != nil {
		return nil, err
	}

	return &Store{
		db: dbStore,
	}, nil
}

func (s *Store) Init() error {
	if err := s.db.Init(); err != nil {
		return err
	}
	return nil
}
