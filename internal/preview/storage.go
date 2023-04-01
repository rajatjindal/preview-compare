package preview

import (
	"encoding/json"

	"github.com/fermyon/spin/sdk/go/key_value"
)

type Store struct {
	kvstore key_value.Store
}

func NewStore() (*Store, error) {
	store, err := key_value.Open("default")
	if err != nil {
		return nil, err
	}

	return &Store{
		kvstore: store,
	}, nil
}

func (s *Store) Close() {
	key_value.Close(s.kvstore)
}

func (s *Store) set(preview *PreviewRequest) error {
	raw, err := json.Marshal(preview)
	if err != nil {
		return err
	}

	return key_value.Set(s.kvstore, preview.Id, raw)
}

func (s *Store) get(id string) (*PreviewRequest, error) {
	raw, err := key_value.Get(s.kvstore, id)
	if err != nil {
		return nil, err
	}

	preview := &PreviewRequest{}
	err = json.Unmarshal(raw, &preview)
	if err != nil {
		return nil, err
	}

	return preview, nil
}
