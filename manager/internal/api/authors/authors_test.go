package authors

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/api"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/api/books"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mocks.NewMockAuthors(ctrl)

	type args struct {
		authorsService service.Authors
	}
	tests := []struct {
		name string
		args args
		want Routes
	}{
		{
			name: "happy",
			args: args{
				authorsService: serviceMock,
			},
			want: &routes{
				authorsService: serviceMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoutes(tt.args.authorsService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routes_CreateAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		body []byte
	}
	type fields struct {
		authorsService func() service.Authors
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		want     func() []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				authorsService: func() service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), &dbmodels.Author{Name: "John Doe"}).
						Return(&dbmodels.Author{ID: "mockid", Name: "John Doe"}, nil)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(
					APIAuthor{
						ID:    "mockid",
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				)
				return jsonBody
			},
			wantCode: 201,
		},
		{
			name: "parse_error",
			fields: fields{
				authorsService: func() service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name_unknown": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "Key: 'AuthorCreateValidator.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				})
				return jsonBody
			},
			wantCode: 400,
		},
		{
			name: "create_error",
			fields: fields{
				authorsService: func() service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), &dbmodels.Author{Name: "John Doe"}).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.authorsService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/authors", bytes.NewBuffer(tt.args.body))
			handler.CreateAuthor(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want())
		})
	}
}

func Test_routes_GetAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockID := "ab7c7712-d10c-412a-87fd-7340b8361e31"

	type fields struct {
		authorsService func(id string) service.Authors
	}
	type args struct {
		authorID string
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		want     func(id string) []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				authorsService: func(id string) service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(&dbmodels.Author{ID: id, Name: "John Doe"}, nil)
					return mock
				},
			},
			args: args{
				authorID: mockID,
			},
			want: func(id string) []byte {
				jsonBody, _ := json.Marshal(
					APIAuthor{
						ID:    id,
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				)
				return jsonBody
			},
			wantCode: 200,
		},
		{
			name: "invalid_uuid",
			fields: fields{
				authorsService: func(_ string) service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					return mock
				},
			},
			args: args{
				authorID: "invalid_uuid",
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "invalid UUID length: 12",
				})
				return jsonBody
			},
			wantCode: 400,
		},
		{
			name: "not_found",
			fields: fields{
				authorsService: func(id string) service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(nil, repository.NotFoundError{})
					return mock
				},
			},
			args: args{
				authorID: mockID,
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "not found",
				})
				return jsonBody
			},
			wantCode: 404,
		},
		{
			name: "internal_error",
			fields: fields{
				authorsService: func(id string) service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				authorID: mockID,
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.authorsService(tt.args.authorID))

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{{Key: "id", Value: tt.args.authorID}}
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/authors"+tt.args.authorID, nil)
			handler.GetAuthor(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want(tt.args.authorID))
		})
	}
}

func Test_routes_GetAuthors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		authorsService func() service.Authors
	}
	tests := []struct {
		name     string
		fields   fields
		want     func() []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				authorsService: func() service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						GetAll(gomock.Any()).
						Return(
							[]dbmodels.Author{{ID: "mockid", Name: "John Doe"}},
							nil,
						)
					return mock
				},
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal([]APIAuthor{
					{
						ID:    "mockid",
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				},
				)
				return jsonBody
			},
			wantCode: 200,
		},
		{
			name: "internal_error",
			fields: fields{
				authorsService: func() service.Authors {
					mock := mocks.NewMockAuthors(ctrl)
					mock.EXPECT().
						GetAll(gomock.Any()).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.authorsService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/authors", nil)
			handler.GetAuthors(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want())
		})
	}
}
