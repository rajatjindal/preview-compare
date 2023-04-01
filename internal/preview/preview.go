package preview

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PreviewRequest struct {
	Id       string `json:"id"`
	ThisBase string `json:"this"`
	ThatBase string `json:"that"`
}

func (s *Store) CreatePreview(ctx context.Context, req *PreviewRequest) (*PreviewRequest, error) {
	req.Id = fmt.Sprintf("preq:%s", uuid.New().String())
	return req, s.set(req)
}
