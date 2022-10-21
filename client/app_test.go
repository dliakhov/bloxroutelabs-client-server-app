package client

import (
	"testing"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
)

func Test_createRandomCommand(t *testing.T) {
	type args struct {
		commandType models.CommandType
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Command
		wantErr bool
	}{
		{
			name: "should return command add item",
			args: args{commandType: models.CommandType_AddItem},
			want: &models.Command{
				Type:        models.CommandType_AddItem,
				ItemID:      1,
				ItemPayload: "A",
			},
			wantErr: false,
		},
		{
			name: "should return command remove item",
			args: args{commandType: models.CommandType_RemoveItem},
			want: &models.Command{
				Type:   models.CommandType_RemoveItem,
				ItemID: 1,
			},
			wantErr: false,
		},
		{
			name: "should return command get all items",
			args: args{commandType: models.CommandType_GetAllItems},
			want: &models.Command{
				Type: models.CommandType_GetAllItems,
			},
			wantErr: false,
		},
		{
			name: "should return command get item",
			args: args{commandType: models.CommandType_GetItem},
			want: &models.Command{
				Type:   models.CommandType_GetItem,
				ItemID: 1,
			},
			wantErr: false,
		},
		{
			name: "should return error when command type is invalid",
			args: args{commandType: 5},
			want: &models.Command{
				Type: 5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createRandomCommand(tt.args.commandType)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("createRandomCommand() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if got.Type != tt.want.Type {
				t.Errorf("type is not correct")
				return
			}

			if tt.want.ItemPayload != "" && got.ItemPayload == "" {
				t.Errorf("got item payload shouldn't be empty")
				return
			}

			if tt.want.ItemPayload == "" && got.ItemPayload != "" {
				t.Errorf("got item payload should be empty")
				return
			}

			if tt.want.ItemPayload != "" && len(got.ItemPayload) > 1 {
				t.Errorf("got item payload should have only one character")
				return
			}

			if tt.want.ItemID != 0 && got.ItemID == 0 {
				t.Errorf("got item id should not be 0")
				return
			}

			if tt.want.ItemID == 0 && got.ItemID != 0 {
				t.Errorf("got item id should be 0")
				return
			}
		})
	}
}
