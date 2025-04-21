//nolint:gochecknoglobals // test data
package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	repository_mocks "github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/mocks"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service"
	"go.uber.org/mock/gomock"
)

const mockBorrowerID = "80f10962-b346-436f-9813-3e20f1ffc7e8"

var mockBorrower = func() *dbmodels.Borrower {
	return &dbmodels.Borrower{Name: "New Borrower"}
}
var mockBorrowerWithID = func() *dbmodels.Borrower {
	return &dbmodels.Borrower{ID: mockBorrowerID, Name: "New Borrower"}
}

func TestNewBorrowers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := repository_mocks.NewMockBorrowers(ctrl)

	type args struct {
		borrowersRepository repository.Borrowers
	}
	tests := []struct {
		name string
		args args
		want service.Borrowers
	}{
		{
			name: "happy",
			args: args{
				borrowersRepository: repositoryMock,
			},
			want: service.NewBorrowers(repositoryMock),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.NewBorrowers(tt.args.borrowersRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBorrowers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_borrowers_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		borrowersRepository func() repository.Borrowers
	}
	type args struct {
		newBorrower *dbmodels.Borrower
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dbmodels.Borrower
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockBorrower()).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				newBorrower: mockBorrower(),
			},
			want:    mockBorrowerWithID(),
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockBorrower()).Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				newBorrower: mockBorrower(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := service.NewBorrowers(tt.fields.borrowersRepository())
			got, err := b.Create(t.Context(), tt.args.newBorrower)
			if (err != nil) != tt.wantErr {
				t.Errorf("borrowers.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("borrowers.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_borrowers_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		borrowersRepository func() repository.Borrowers
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dbmodels.Borrower
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				id: mockBorrowerID,
			},
			want:    mockBorrowerWithID(),
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				id: mockBorrowerID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := service.NewBorrowers(tt.fields.borrowersRepository())
			got, err := b.GetByID(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("borrowers.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("borrowers.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_borrowers_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		borrowersRepository func() repository.Borrowers
	}
	tests := []struct {
		name    string
		fields  fields
		want    []dbmodels.Borrower
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetAll(gomock.Any()).Return([]dbmodels.Borrower{*mockBorrowerWithID()}, nil)
					return mock
				},
			},
			want:    []dbmodels.Borrower{*mockBorrowerWithID()},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				borrowersRepository: func() repository.Borrowers {
					mock := repository_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("error"))
					return mock
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := service.NewBorrowers(tt.fields.borrowersRepository())
			got, err := b.GetAll(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("borrowers.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("borrowers.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
