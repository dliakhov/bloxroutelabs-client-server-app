package repository

import (
	"errors"
	"sync"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

//go:generate mockgen -package=repository -source=repo.go -destination=repo_mock.go
type Repo interface {
	AddItem(item models.Item) error
	RemoveItem(itemID int64) error
	GetItem(itemID int64) (models.Item, error)
	GetAllItems() ([]models.Item, error)
}

type repoImpl struct {
	storage *linkedhashmap.Map
	rwMx    sync.RWMutex
}

func New() Repo {
	return &repoImpl{
		storage: linkedhashmap.New(),
	}
}

func (r *repoImpl) AddItem(item models.Item) error {
	r.rwMx.Lock()
	defer r.rwMx.Unlock()

	r.storage.Put(item.ID, item.Payload)
	return nil
}

func (r *repoImpl) RemoveItem(itemID int64) error {
	r.rwMx.Lock()
	defer r.rwMx.Unlock()

	r.storage.Remove(itemID)
	return nil
}

func (r *repoImpl) GetItem(itemID int64) (models.Item, error) {
	r.rwMx.RLock()
	defer r.rwMx.RUnlock()

	value, ok := r.storage.Get(itemID)
	if !ok {
		return models.Item{}, nil
	}

	payload, ok := value.(string)
	if !ok {
		return models.Item{}, errors.New("item has not correct type")
	}
	return models.Item{
		ID:      itemID,
		Payload: payload,
	}, nil
}

func (r *repoImpl) GetAllItems() ([]models.Item, error) {
	r.rwMx.RLock()
	defer r.rwMx.RUnlock()

	var items []models.Item
	r.storage.All(func(key, value any) bool {
		items = append(items, models.Item{
			ID:      key.(int64),
			Payload: value.(string),
		})
		return true
	})

	return items, nil
}
