package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/stretchr/testify/assert"
)

func Test_repoImpl_AddItem(t *testing.T) {
	type args struct {
		manipulateWithRepo func(r Repo)
	}
	tests := []struct {
		name      string
		args      args
		wantItems []models.Item
	}{
		{
			name: "should add one item",
			args: args{manipulateWithRepo: func(r Repo) {
				r.AddItem(models.Item{
					ID:      1,
					Payload: "A",
				})
			}},
			wantItems: []models.Item{{
				ID:      1,
				Payload: "A",
			}},
		},
		{
			name: "should add two items",
			args: args{manipulateWithRepo: func(r Repo) {
				r.AddItem(models.Item{
					ID:      1,
					Payload: "A",
				})
				r.AddItem(models.Item{
					ID:      2,
					Payload: "B",
				})
			}},
			wantItems: []models.Item{{
				ID:      1,
				Payload: "A",
			}, {
				ID:      2,
				Payload: "B",
			}},
		},
		{
			name: "should add 5 items in parallel",
			args: args{manipulateWithRepo: func(r Repo) {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func(i int) {
						r.AddItem(models.Item{
							ID:      int64(i),
							Payload: fmt.Sprintf("%c", 'A'+i),
						})
						wg.Done()
					}(i)
				}

				wg.Wait()
			}},
			wantItems: []models.Item{{
				ID:      0,
				Payload: "A",
			}, {
				ID:      1,
				Payload: "B",
			}, {
				ID:      2,
				Payload: "C",
			}, {
				ID:      3,
				Payload: "D",
			}, {
				ID:      4,
				Payload: "E",
			}},
		},
		{
			name: "should add 5 items in parallel and remove 2",
			args: args{manipulateWithRepo: func(r Repo) {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func(i int) {
						r.AddItem(models.Item{
							ID:      int64(i),
							Payload: fmt.Sprintf("%c", 'A'+i),
						})
						wg.Done()
					}(i)
				}

				wg.Wait()

				r.RemoveItem(0)
				r.RemoveItem(1)
			}},
			wantItems: []models.Item{{
				ID:      2,
				Payload: "C",
			}, {
				ID:      3,
				Payload: "D",
			}, {
				ID:      4,
				Payload: "E",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New().(*repoImpl)

			tt.args.manipulateWithRepo(r)

			gotItems := getAllItems(r)
			assert.Subset(t, tt.wantItems, gotItems)
		})
	}
}

func getAllItems(r *repoImpl) []models.Item {
	var items []models.Item
	r.storage.All(func(key, value any) bool {
		items = append(items, models.Item{
			ID:      key.(int64),
			Payload: value.(string),
		})
		return true
	})
	return items
}

func Test_repoImpl_GetItem(t *testing.T) {
	r := New().(*repoImpl)
	r.storage.Put(int64(1), "A")
	r.storage.Put(int64(2), "B")
	r.storage.Put(int64(4), 1234)

	type args struct {
		itemID int64
	}
	tests := []struct {
		name    string
		args    args
		want    models.Item
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should get item correctly",
			args: args{itemID: 1},
			want: models.Item{
				ID:      1,
				Payload: "A",
			},
			wantErr: assert.NoError,
		},
		{
			name: "should get second item correctly",
			args: args{itemID: 2},
			want: models.Item{
				ID:      2,
				Payload: "B",
			},
			wantErr: assert.NoError,
		},
		{
			name:    "should get error when item not exist",
			args:    args{itemID: 3},
			want:    models.Item{},
			wantErr: assert.Error,
		},
		{
			name:    "should get error when item has wrong type",
			args:    args{itemID: 4},
			want:    models.Item{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetItem(tt.args.itemID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetItem(%v)", tt.args.itemID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetItem(%v)", tt.args.itemID)
		})
	}
}

func Test_repoImpl_GetAllItems(t *testing.T) {
	tests := []struct {
		name     string
		initRepo func(r *repoImpl)
		want     []models.Item
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "should return all 2 items",
			initRepo: func(r *repoImpl) {
				r.storage.Put(int64(1), "A")
				r.storage.Put(int64(2), "B")
			},
			want: []models.Item{
				{
					ID:      1,
					Payload: "A",
				},
				{
					ID:      2,
					Payload: "B",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return all 4 items",
			initRepo: func(r *repoImpl) {
				r.storage.Put(int64(1), "A")
				r.storage.Put(int64(2), "B")
				r.storage.Put(int64(3), "B")
				r.storage.Put(int64(4), "B")
			},
			want: []models.Item{
				{
					ID:      1,
					Payload: "A",
				},
				{
					ID:      2,
					Payload: "B",
				},
				{
					ID:      3,
					Payload: "B",
				},
				{
					ID:      4,
					Payload: "B",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New().(*repoImpl)
			tt.initRepo(r)

			got, err := r.GetAllItems()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllItems()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllItems()")
		})
	}
}
