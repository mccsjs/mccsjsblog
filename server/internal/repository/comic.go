package repository

import (
	"blog/internal/dto"
	"blog/internal/model"

	"gorm.io/gorm"
)

type ComicRepository struct {
	db *gorm.DB
}

func NewComicRepository(db *gorm.DB) *ComicRepository {
	return &ComicRepository{db: db}
}

func (r *ComicRepository) List(req *dto.ListComicRequest) ([]model.Comic, int64, error) {
	var comics []model.Comic
	var total int64

	query := r.db.Model(&model.Comic{})
	query.Count(&total)

	if err := query.Order("sort ASC, id DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&comics).Error; err != nil {
		return nil, 0, err
	}
	return comics, total, nil
}

func (r *ComicRepository) ListAll() ([]model.Comic, error) {
	var comics []model.Comic
	err := r.db.Order("sort ASC, id DESC").Find(&comics).Error
	return comics, err
}

func (r *ComicRepository) GetByID(id uint) (*model.Comic, error) {
	var comic model.Comic
	err := r.db.First(&comic, id).Error
	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func (r *ComicRepository) Create(comic *model.Comic) error {
	return r.db.Create(comic).Error
}

func (r *ComicRepository) Update(comic *model.Comic) error {
	return r.db.Save(comic).Error
}

func (r *ComicRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comic{}, id).Error
}
