package service

import (
	"log"

	"github.com/adetiamarhadi/golang-rest-api/dto"
	"github.com/adetiamarhadi/golang-rest-api/entity"
	"github.com/adetiamarhadi/golang-rest-api/repository"
	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userId string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}

	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("failed to map %v:", err)
	}

	updatedUser := service.userRepository.UpdateUser(userToUpdate)

	return updatedUser
}

func (service *userService) Profile(userId string) entity.User {
	return service.userRepository.ProfileUser(userId)
}
