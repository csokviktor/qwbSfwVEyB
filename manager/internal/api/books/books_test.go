package books

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/csokviktor/lib_manager/internal/api"
	"github.com/csokviktor/lib_manager/internal/repository"
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"github.com/csokviktor/lib_manager/internal/service"
	"github.com/csokviktor/lib_manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mocks.NewMockBooks(ctrl)

	type args struct {
		booksService service.Books
	}
	tests := []struct {
		name string
		args args
		want Routes
	}{
		{
			name: "happy",
			args: args{
				booksService: serviceMock,
			},
			want: &routes{
				booksService: serviceMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoutes(tt.args.booksService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routes_CreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockID := "ab7c7712-d10c-412a-87fd-7340b8361e31"

	type args struct {
		body []byte
	}
	type fields struct {
		booksService func() service.Books
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
				booksService: func() service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), &dbmodels.Book{
							Title:    "John Doe",
							AuthorID: mockID,
						}).
						Return(&dbmodels.Book{
							ID:       "mockid",
							Title:    "John Doe",
							AuthorID: mockID,
						}, nil)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"title":    "John Doe",
						"authorID": mockID,
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(
					APIBook{
						ID:       "mockid",
						Title:    "John Doe",
						AuthorID: mockID,
					},
				)
				return jsonBody
			},
			wantCode: 201,
		},
		{
			name: "parse_error",
			fields: fields{
				booksService: func() service.Books {
					mock := mocks.NewMockBooks(ctrl)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"unknown": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "'BookCreateValidator.Title' Error",
				})
				return jsonBody
			},
			wantCode: 400,
		},
		{
			name: "create_error",
			fields: fields{
				booksService: func() service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"title":    "John Doe",
						"authorID": mockID,
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
			handler := NewRoutes(tt.fields.booksService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(tt.args.body))
			handler.CreateBook(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode != 400 {
				assert.Equal(t, w.Body.Bytes(), tt.want())
			}
		})
	}
}

func Test_routes_GetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		booksService func() service.Books
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
				booksService: func() service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						GetAll(gomock.Any()).
						Return([]dbmodels.Book{{
							ID:    "mockid",
							Title: "John Doe",
						}}, nil)
					return mock
				},
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(
					[]APIBook{{
						ID:    "mockid",
						Title: "John Doe",
					}},
				)
				return jsonBody
			},
			wantCode: 200,
		},
		{
			name: "internal_error",
			fields: fields{
				booksService: func() service.Books {
					mock := mocks.NewMockBooks(ctrl)
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
			handler := NewRoutes(tt.fields.booksService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/books", nil)
			handler.GetBooks(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want())
		})
	}
}

func Test_routes_BorrowBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookID := "ab7c7712-d10c-412a-87fd-7340b8361e31"
	mockBorrowerID := "ab7c7712-d10c-412a-87fd-7340b8361e32"

	type fields struct {
		booksService func(borrowerID, bookID string) service.Books
	}
	type args struct {
		bookID     string
		borrowerID string
		body       []byte
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
				booksService: func(borrowerID, bookID string) service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						Borrow(gomock.Any(), borrowerID, bookID).
						Return(nil)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"borrowerID": mockBorrowerID,
					})
					return jsonBody
				}(),
			},
			want: func(id string) []byte {
				jsonBody, _ := json.Marshal(nil)
				return jsonBody
			},
			wantCode: 201,
		},
		{
			name: "invalid_book_uuid",
			fields: fields{
				booksService: func(borrowerID, bookID string) service.Books {
					mock := mocks.NewMockBooks(ctrl)
					return mock
				},
			},
			args: args{
				bookID:     "invalid_uuid",
				borrowerID: mockBorrowerID,
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"borrowerID": mockBorrowerID,
					})
					return jsonBody
				}(),
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
			name: "invalid_borrower_uuid",
			fields: fields{
				booksService: func(borrowerID, bookID string) service.Books {
					mock := mocks.NewMockBooks(ctrl)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: "invalid_uuid",
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"borrowerID": "invalid_uuid",
					})
					return jsonBody
				}(),
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
			name: "not_found_error",
			fields: fields{
				booksService: func(borrowerID, bookID string) service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						Borrow(gomock.Any(), borrowerID, bookID).
						Return(repository.NotFoundError{})
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"borrowerID": mockBorrowerID,
					})
					return jsonBody
				}(),
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
				booksService: func(borrowerID, bookID string) service.Books {
					mock := mocks.NewMockBooks(ctrl)
					mock.EXPECT().
						Borrow(gomock.Any(), borrowerID, bookID).
						Return(errors.New("error"))
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"borrowerID": mockBorrowerID,
					})
					return jsonBody
				}(),
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
			handler := NewRoutes(tt.fields.booksService(tt.args.borrowerID, tt.args.bookID))

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{{Key: "id", Value: tt.args.bookID}}
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/books"+tt.args.bookID, bytes.NewBuffer(tt.args.body))
			handler.BorrowBook(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode != 400 {
				assert.Equal(t, w.Body.Bytes(), tt.want(tt.args.bookID))
			}
		})
	}
}
