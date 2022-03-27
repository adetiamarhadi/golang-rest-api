package service

import (
	"fmt"
	"log"

	"github.com/adetiamarhadi/golang-rest-api/dto"
	"github.com/adetiamarhadi/golang-rest-api/entity"
	"github.com/adetiamarhadi/golang-rest-api/repository"
	"github.com/mashingan/smapping"
)

// BookService is ...
type BookService interface {
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookUpdateDTO) entity.Book
	Delete(b entity.Book)
	All() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

// NewBookService is a service that have a responsibility for ...
func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

// Insert ...
func (service *bookService) Insert(b dto.BookCreateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("failed map %v: ", err)
	}

	res := service.bookRepository.InsertBook(book)
	return res
}

// Update ...
func (service *bookService) Update(b dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("failed map %v: ", err)
	}

	res := service.bookRepository.UpdateBook(book)
	return res
}

// Delete ...
func (service *bookService) Delete(b entity.Book) {
	service.bookRepository.DeleteBook(b)
}

// All ...
func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}

// FindById ...
func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

// IsAllowedToEdit ...
func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	book := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", book.UserID)
	return userID == id
}