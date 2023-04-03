package preview

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Change struct {
	Link     string `json:"link"`
	LinkText string `json:"linkText"`
}
type PreviewRequest struct {
	Id       string   `json:"id"`
	ThisBase string   `json:"this"`
	ThatBase string   `json:"that"`
	Changes  []Change `json:"changes"`
}

func (s *Store) CreatePreview(ctx context.Context, req *PreviewRequest) (*PreviewRequest, error) {
	req.Id = fmt.Sprintf("preq-%s", uuid.New().String())
	return req, s.set(req)
}

func (s *Store) GetPreview(ctx context.Context, id string) (*PreviewRequest, error) {
	return s.get(id)
}
