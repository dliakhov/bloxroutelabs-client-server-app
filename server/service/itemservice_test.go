package service

import (
	"context"
	"testing"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	"github.com/dliakhov/bloxroutelabs/client-server-app/server/repository"
	"github.com/golang/mock/gomock"
)

func Test_itemServiceImpl_ProcessItemCommand(t *testing.T) {
	type fields struct {
		repo func(c *gomock.Controller) repository.Repo
	}
	type args struct {
		ctx     context.Context
		command *models.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should process add item command",
			fields: fields{repo: func(ctrl *gomock.Controller) repository.Repo {
				repo := repository.NewMockRepo(ctrl)
				repo.EXPECT().AddItem(models.Item{
					ID:      1,
					Payload: "A",
				}).Return(nil)

				return repo
			}},
			args: args{
				ctx: context.Background(),
				command: &models.Command{
					Type:        models.CommandType_AddItem,
					ItemID:      1,
					ItemPayload: "A",
				},
			},
		},
		{
			name: "should process remove item command",
			fields: fields{repo: func(ctrl *gomock.Controller) repository.Repo {
				repo := repository.NewMockRepo(ctrl)
				repo.EXPECT().RemoveItem(int64(1)).Return(nil)

				return repo
			}},
			args: args{
				ctx: context.Background(),
				command: &models.Command{
					Type:   models.CommandType_RemoveItem,
					ItemID: 1,
				},
			},
		},
		{
			name: "should process get item command",
			fields: fields{repo: func(ctrl *gomock.Controller) repository.Repo {
				repo := repository.NewMockRepo(ctrl)
				repo.EXPECT().GetItem(int64(1)).Return(models.Item{
					ID:      1,
					Payload: "A",
				}, nil)

				return repo
			}},
			args: args{
				ctx: context.Background(),
				command: &models.Command{
					Type:   models.CommandType_GetItem,
					ItemID: 1,
				},
			},
		},
		{
			name: "should process get all items command",
			fields: fields{repo: func(ctrl *gomock.Controller) repository.Repo {
				repo := repository.NewMockRepo(ctrl)
				repo.EXPECT().GetAllItems().Return([]models.Item{{
					ID:      1,
					Payload: "A",
				}}, nil)

				return repo
			}},
			args: args{
				ctx: context.Background(),
				command: &models.Command{
					Type: models.CommandType_GetAllItems,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := New(tt.fields.repo(ctrl))
			if err := i.ProcessItemCommand(tt.args.ctx, tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("ProcessItemCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
