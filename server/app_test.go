package server

import (
	"errors"
	"testing"

	"github.com/dliakhov/bloxroutelabs/client-server-app/models"
	service "github.com/dliakhov/bloxroutelabs/client-server-app/server/service"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func TestApp_ProcessMessage(t *testing.T) {
	command := models.Command{
		Type: models.CommandType_GetAllItems,
	}

	commandBodyBytes, err := proto.Marshal(&command)
	if err != nil {
		t.Fail()
		return
	}

	type fields struct {
		config      Configurations
		itemService func(ctrl *gomock.Controller) service.ItemService
	}
	type args struct {
		d amqp.Delivery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should process command successfully",
			fields: fields{
				config: Configurations{},
				itemService: func(ctrl *gomock.Controller) service.ItemService {
					itemService := service.NewMockItemService(ctrl)
					itemService.EXPECT().ProcessItemCommand(gomock.Any(), commandBodyBytes).Return(nil)
					return itemService
				},
			},
			args: args{
				d: amqp.Delivery{
					Body: commandBodyBytes,
					Headers: map[string]interface{}{
						"X-Trace-ID": "trace_id",
					},
				},
			},
		},
		{
			name: "should process command successfully without trace id",
			fields: fields{
				config: Configurations{},
				itemService: func(ctrl *gomock.Controller) service.ItemService {
					itemService := service.NewMockItemService(ctrl)
					itemService.EXPECT().ProcessItemCommand(gomock.Any(), commandBodyBytes).Return(nil)
					return itemService
				},
			},
			args: args{
				d: amqp.Delivery{
					Body: commandBodyBytes,
				},
			},
		},
		{
			name: "should not process command when body is not correct",
			fields: fields{
				config: Configurations{},
				itemService: func(ctrl *gomock.Controller) service.ItemService {
					itemService := service.NewMockItemService(ctrl)
					itemService.EXPECT().ProcessItemCommand(gomock.Any(), []byte{1, 2, 3}).Return(nil)
					return itemService
				},
			},
			args: args{
				d: amqp.Delivery{
					Body: []byte{1, 2, 3},
					Headers: map[string]interface{}{
						"X-Trace-ID": "trace_id",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when processing command returns error",
			fields: fields{
				config: Configurations{},
				itemService: func(ctrl *gomock.Controller) service.ItemService {
					itemService := service.NewMockItemService(ctrl)
					itemService.EXPECT().ProcessItemCommand(gomock.Any(), commandBodyBytes).Return(errors.New("cannot process command"))
					return itemService
				},
			},
			args: args{
				d: amqp.Delivery{
					Body: []byte{1, 2, 3},
					Headers: map[string]interface{}{
						"X-Trace-ID": "trace_id",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewApp(tt.fields.config, nil)
			err := a.ProcessMessage(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error occured: %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
