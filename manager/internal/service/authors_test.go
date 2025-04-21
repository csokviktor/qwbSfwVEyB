//nolint:gochecknoglobals // test data
package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/mocks"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service"
	"go.uber.org/mock/gomock"
)

const mockAuthorID = "80f10962-b346-436f-9813-3e20f1ffc7e6"

var mockAuthor = func() *dbmodels.Author {
	return &dbmodels.Author{Name: "New Author"}
}
var mockAuthorWithID = func() *dbmodels.Author {
	return &dbmodels.Author{ID: mockAuthorID, Name: "New Author"}
}

func TestNewAuthors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := mocks.NewMockAuthors(ctrl)

	type args struct {
		authorsRepository repository.Authors
	}
	tests := []struct {
		name string
		args args
		want service.Authors
	}{
		{
			name: "happy",
			args: args{
				authorsRepository: repositoryMock,
			},
			want: service.NewAuthors(repositoryMock),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.NewAuthors(tt.args.authorsRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authors_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		authorsRepository func() repository.Authors
	}
	type args struct {
		newAuthor *dbmodels.Author
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dbmodels.Author
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockAuthor()).Return(mockAuthorWithID(), nil)
					return mock
				},
			},
			args: args{
				newAuthor: mockAuthor(),
			},
			want:    mockAuthorWithID(),
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockAuthor()).Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				newAuthor: mockAuthor(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := service.NewAuthors(tt.fields.authorsRepository())
			got, err := a.Create(t.Context(), tt.args.newAuthor)
			if (err != nil) != tt.wantErr {
				t.Errorf("authors.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authors.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authors_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		authorsRepository func() repository.Authors
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dbmodels.Author
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockAuthorID).Return(mockAuthorWithID(), nil)
					return mock
				},
			},
			args: args{
				id: mockAuthorID,
			},
			want:    mockAuthorWithID(),
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockAuthorID).Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				id: mockAuthorID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := service.NewAuthors(tt.fields.authorsRepository())
			got, err := a.GetByID(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("authors.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authors.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authors_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		authorsRepository func() repository.Authors
	}
	tests := []struct {
		name    string
		fields  fields
		want    []dbmodels.Author
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().GetAll(gomock.Any()).Return([]dbmodels.Author{*mockAuthorWithID()}, nil)
					return mock
				},
			},
			want:    []dbmodels.Author{*mockAuthorWithID()},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				authorsRepository: func() repository.Authors {
					mock := mocks.NewMockAuthors(ctrl)
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
			a := service.NewAuthors(tt.fields.authorsRepository())
			got, err := a.GetAll(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("authors.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authors.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
