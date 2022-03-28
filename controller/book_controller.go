package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/adetiamarhadi/golang-rest-api/dto"
	"github.com/adetiamarhadi/golang-rest-api/entity"
	"github.com/adetiamarhadi/golang-rest-api/helper"
	"github.com/adetiamarhadi/golang-rest-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// BookController ...
type BookController interface {
	All(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService service.JWTService
}

// NewBookController ...
func NewBookController(bs service.BookService, jwt service.JWTService) BookController {
	return &bookController{
		bookService: bs,
		jwtService: jwt,
	}
}

func (controller *bookController) getUserIDByToken(token string) string {
	aToken, err := controller.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}

func (controller *bookController) All(ctx *gin.Context) {
	books := controller.bookService.All()
	res := helper.BuildResponse(true, "OK", books)
	ctx.JSON(http.StatusOK, res)
}

func (controller *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("no param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	book := controller.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("data not found", "no data with given id", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", book)
		ctx.JSON(http.StatusOK, res)
	}
}

func (controller *bookController) Insert(ctx *gin.Context) {
	var bookCreateDto dto.BookCreateDTO
	errDto := ctx.ShouldBind(&bookCreateDto)
	if errDto != nil {
		res := helper.BuildErrorResponse("failed to process request", errDto.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := controller.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			bookCreateDto.UserID = convertedUserID
		}
		result := controller.bookService.Insert(bookCreateDto)
		res := helper.BuildResponse(true, "OK", result)
		ctx.JSON(http.StatusCreated, res)
	}
}

func (controller *bookController) Update(ctx *gin.Context) {
	var bookUpdateDto dto.BookUpdateDTO
	errDto := ctx.ShouldBind(&bookUpdateDto)
	if errDto != nil {
		res := helper.BuildErrorResponse("failed to process request", errDto.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	userID := controller.getUserIDByToken(authHeader)
	if (controller.bookService.IsAllowedToEdit(userID, bookUpdateDto.ID)) {
		id, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			bookUpdateDto.ID = id
		}
		result := controller.bookService.Update(bookUpdateDto)
		res := helper.BuildResponse(true, "OK", result)
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("you dont have permission", "you are not the owner", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
	}
}

func (controller *bookController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("no param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book
	book.ID = id

	authHeader := ctx.GetHeader("Authorization")
	userID := controller.getUserIDByToken(authHeader)
	if (controller.bookService.IsAllowedToEdit(userID, book.ID)) {
		controller.bookService.Delete(book)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("you dont have permission", "you are not the owner", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, res)
	}
}