package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/adetiamarhadi/golang-rest-api/dto"
	"github.com/adetiamarhadi/golang-rest-api/helper"
	"github.com/adetiamarhadi/golang-rest-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(ctx *gin.Context) {
	var userUpdateDto dto.UserUpdateDTO

	errDto := ctx.ShouldBind(&userUpdateDto)
	if errDto != nil {
		res := helper.BuildErrorResponse("failed to process request", errDto.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}

	userUpdateDto.ID = id

	u := c.userService.Update(userUpdateDto)
	res := helper.BuildResponse(true, "OK!", u)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	user := c.userService.Profile(fmt.Sprintf("%v", claims["user_id"]))
	res := helper.BuildResponse(true, "OK", user)
	ctx.JSON(http.StatusOK, res)
}
