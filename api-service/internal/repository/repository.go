package repository

import (
	"api-service/internal/database"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{
		db: database.DB,
	}
}

func (r *Repository[T]) GetAll(filters map[string]interface{}, orderBy string, limit, offset int, preloads []string) ([]T, int, error) {
	var items []T
	var count int64

	db := r.db.Model(new(T))

	for key, value := range filters {
		if str, ok := value.(string); ok && str != "" {
			if str == "null" {
				db = db.Where(fmt.Sprintf("\"%s\" IS NULL", key))
			} else if strings.HasSuffix(key, "_like") {
				actualKey := strings.TrimSuffix(key, "_like")
				db = db.Where(fmt.Sprintf("\"%s\" LIKE ?", actualKey), "%"+str)
			} else if strings.Contains(str, ",") {
				db = db.Where(fmt.Sprintf("\"%s\" IN ?", key), strings.Split(str, ","))
			} else {
				db = db.Where(fmt.Sprintf("\"%s\" = ?", key), value)
			}
		} else if value != nil {
			db = db.Where(fmt.Sprintf("\"%s\" = ?", key), value)
		}
	}

	for _, preload := range preloads {
		if preload != "" {
			db = db.Preload(preload)
		}
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if orderBy != "" {
		db = db.Order(orderBy)
	}

	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err := db.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, int(count), nil
}

func (r *Repository[T]) GetByID(id interface{}, preloads []string) (*T, error) {
	var item T
	db := r.db.Model(new(T))

	for _, preload := range preloads {
		if preload != "" {
			db = db.Preload(preload)
		}
	}

	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository[T]) Create(item *T) error {
	return r.db.Create(item).Error
}

func (r *Repository[T]) Update(id interface{}, item *T) error {
	return r.db.Model(new(T)).Where("id = ?", id).Updates(item).Error
}

func (r *Repository[T]) Patch(id interface{}, updates map[string]interface{}) error {
	return r.db.Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository[T]) Delete(id interface{}) error {
	return r.db.Delete(new(T), "id = ?", id).Error
}

func (r *Repository[T]) GetDistinct(field string, filters map[string]interface{}) ([]string, error) {
	var results []string
	db := r.db.Model(new(T))

	for key, value := range filters {
		if value != "" {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	err := db.Distinct(field).Order(field + " asc").Pluck(field, &results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
