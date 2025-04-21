//nolint:gochecknoglobals // test data
package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/csokviktor/lib_manager/internal/repository"
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	repository_mocks "github.com/csokviktor/lib_manager/internal/repository/mocks"
	"github.com/csokviktor/lib_manager/internal/service"
	service_mocks "github.com/csokviktor/lib_manager/internal/service/mocks"
	"github.com/csokviktor/lib_manager/internal/util"
	"go.uber.org/mock/gomock"
)

const mockBookID = "80f10962-b346-436f-9813-3e20f1ffc7e7"

var mockBook = func() *dbmodels.Book {
	return &dbmodels.Book{Title: "New Book", AuthorID: mockAuthorID}
}
var mockBookWithID = func() *dbmodels.Book {
	return &dbmodels.Book{ID: mockAuthorID, Title: "New Book", AuthorID: mockAuthorID}
}
var mockBookWithBorrowerID = func() *dbmodels.Book {
	return &dbmodels.Book{
		ID:         mockAuthorID,
		Title:      "New Book",
		BorrowerID: util.AsPointer(mockBorrowerID),
		AuthorID:   mockAuthorID,
	}
}

func TestNewBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := repository_mocks.NewMockBooks(ctrl)
	authorsServiceMock := service_mocks.NewMockAuthors(ctrl)
	borrowersServiceMock := service_mocks.NewMockBorrowers(ctrl)

	type args struct {
		bookRepository   repository.Books
		authorsService   service.Authors
		borrowersService service.Borrowers
	}
	tests := []struct {
		name string
		args args
		want service.Books
	}{
		{
			name: "happy",
			args: args{
				bookRepository:   repositoryMock,
				authorsService:   authorsServiceMock,
				borrowersService: borrowersServiceMock,
			},
			want: service.NewBooks(repositoryMock, authorsServiceMock, borrowersServiceMock),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.NewBooks(
				tt.args.bookRepository,
				tt.args.authorsService,
				tt.args.borrowersService,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_books_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		booksRepository func() repository.Books
		authorsService  func() service.Authors
	}
	type args struct {
		newBook *dbmodels.Book
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dbmodels.Book
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockBook()).Return(mockBookWithID(), nil)
					return mock
				},
				authorsService: func() service.Authors {
					mock := service_mocks.NewMockAuthors(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBook().AuthorID).Return(nil, nil)
					return mock
				},
			},
			args: args{
				newBook: mockBook(),
			},
			want:    mockBookWithID(),
			wantErr: false,
		},

		{
			name: "error",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().Create(gomock.Any(), mockBook()).Return(nil, errors.New("error"))
					return mock
				},
				authorsService: func() service.Authors {
					mock := service_mocks.NewMockAuthors(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBook().AuthorID).Return(nil, nil)
					return mock
				},
			},
			args: args{
				newBook: mockBook(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := service.NewBooks(tt.fields.booksRepository(), tt.fields.authorsService(), nil)
			got, err := b.Create(t.Context(), tt.args.newBook)
			if (err != nil) != tt.wantErr {
				t.Errorf("books.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("books.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_books_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		booksRepository func() repository.Books
	}
	tests := []struct {
		name    string
		fields  fields
		want    []dbmodels.Book
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().GetAll(gomock.Any()).Return([]dbmodels.Book{*mockBookWithID()}, nil)
					return mock
				},
			},
			want:    []dbmodels.Book{*mockBookWithID()},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
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
			b := service.NewBooks(tt.fields.booksRepository(), nil, nil)
			got, err := b.GetAll(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("books.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("books.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_books_Borrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		booksRepository  func() repository.Books
		borrowersService func() service.Borrowers
	}
	type args struct {
		borrowerID string
		bookID     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBookID).Return(mockBookWithID(), nil)
					mock.EXPECT().Update(gomock.Any(), mockBookWithBorrowerID()).Return(nil, nil)
					return mock
				},
				borrowersService: func() service.Borrowers {
					mock := service_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
			},
			wantErr: false,
		},
		{
			name: "service_error",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					return mock
				},
				borrowersService: func() service.Borrowers {
					mock := service_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
			},
			wantErr: true,
		},
		{
			name: "repository_eror",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBookID).Return(nil, errors.New("error"))
					return mock
				},
				borrowersService: func() service.Borrowers {
					mock := service_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
			},
			wantErr: true,
		},
		{
			name: "book_already_borrowed",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBookID).Return(mockBookWithBorrowerID(), nil)
					return mock
				},
				borrowersService: func() service.Borrowers {
					mock := service_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
			},
			wantErr: true,
		},
		{
			name: "update_error",
			fields: fields{
				booksRepository: func() repository.Books {
					mock := repository_mocks.NewMockBooks(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBookID).Return(mockBookWithID(), nil)
					mock.EXPECT().Update(gomock.Any(), mockBookWithBorrowerID()).Return(nil, errors.New("error"))
					return mock
				},
				borrowersService: func() service.Borrowers {
					mock := service_mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().GetByID(gomock.Any(), mockBorrowerID).Return(mockBorrowerWithID(), nil)
					return mock
				},
			},
			args: args{
				bookID:     mockBookID,
				borrowerID: mockBorrowerID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := service.NewBooks(tt.fields.booksRepository(), nil, tt.fields.borrowersService())
			if err := b.Borrow(t.Context(), tt.args.borrowerID, tt.args.bookID); (err != nil) != tt.wantErr {
				t.Errorf("books.Borrow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
