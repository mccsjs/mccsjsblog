package service

import (
	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/repository"
)

type ComicService struct {
	repo *repository.ComicRepository
}

func NewComicService(repo *repository.ComicRepository) *ComicService {
	return &ComicService{repo: repo}
}

func (s *ComicService) ListForWeb() ([]dto.ComicResponse, error) {
	comics, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return toComicResponses(comics), nil
}

func (s *ComicService) List(req *dto.ListComicRequest) ([]dto.ComicResponse, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	comics, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}
	return toComicResponses(comics), total, nil
}

func (s *ComicService) GetByID(id uint) (*dto.ComicResponse, error) {
	comic, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := toComicResponse(comic)
	return &resp, nil
}

func (s *ComicService) Create(req *dto.CreateComicRequest) (*dto.ComicResponse, error) {
	comic := &model.Comic{
		Name:  req.Name,
		URL:   req.URL,
		Cover: req.Cover,
		Sort:  req.Sort,
	}
	if err := s.repo.Create(comic); err != nil {
		return nil, err
	}
	resp := toComicResponse(comic)
	return &resp, nil
}

func (s *ComicService) Update(id uint, req *dto.UpdateComicRequest) (*dto.ComicResponse, error) {
	comic, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		comic.Name = req.Name
	}
	if req.URL != "" {
		comic.URL = req.URL
	}
	if req.Cover != "" {
		comic.Cover = req.Cover
	}
	if req.Sort != nil {
		comic.Sort = *req.Sort
	}
	if err := s.repo.Update(comic); err != nil {
		return nil, err
	}
	resp := toComicResponse(comic)
	return &resp, nil
}

func (s *ComicService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func toComicResponse(c *model.Comic) dto.ComicResponse {
	return dto.ComicResponse{
		ID:        c.ID,
		Name:      c.Name,
		URL:       c.URL,
		Cover:     c.Cover,
		Sort:      c.Sort,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toComicResponses(comics []model.Comic) []dto.ComicResponse {
	result := make([]dto.ComicResponse, len(comics))
	for i, c := range comics {
		result[i] = toComicResponse(&c)
	}
	return result
}
