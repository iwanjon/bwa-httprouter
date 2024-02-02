package user

import (
	"bwahttprouter/helper"
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(ctx context.Context, input RegisterUser) (User, error)
	LoginUser(ctx context.Context, input LoginInput) (User, error)
	CheckEmailAvailable(ctx context.Context, input CheckEmailInput) (bool, error)
	SaveAvatar(ctx context.Context, id int, filelocation string) (User, error)
	GetUserById(ctx context.Context, id int) (User, error)
}

type service struct {
	repo     Repository
	Validate *validator.Validate
}

func NewService(r Repository, v *validator.Validate) *service {
	return &service{repo: r, Validate: v}
}

func (s *service) RegisterUser(ctx context.Context, input RegisterUser) (User, error) {
	err := s.Validate.Struct(input)
	helper.PanicIfError(err, "error create user validation")
	var user User

	user.Email = input.Email
	user.Name = input.Name
	// user.PasswordHash = input.Password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	fmt.Println(passwordHash)
	user.PasswordHash = string(passwordHash)
	user.Occupation = input.Occupation
	fmt.Println(user, "this is user")
	u, err := s.repo.Save(ctx, user)
	helper.PanicIfError(err, " error register user service ")

	return u, nil
}

func (s *service) LoginUser(ctx context.Context, input LoginInput) (User, error) {
	err := s.Validate.Struct(input)
	helper.PanicIfError(err, "error create user validation")
	user, err := s.repo.FindByEmail(ctx, input.Email)

	helper.PanicIfError(err, "error in login serivce")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	helper.PanicIfError(err, " error in login dervice compare hash")

	return user, nil

}

func (s *service) CheckEmailAvailable(ctx context.Context, input CheckEmailInput) (bool, error) {
	err := s.Validate.Struct(input)
	helper.PanicIfError(err, "error validate input email cek")
	directory := http.Dir("./")
	fileServer := http.FileServer(directory)
	fmt.Println(fileServer, "makanan")

	user, err := s.repo.FindByEmail(ctx, input.Email)
	helper.PanicIfError(err, "error finding email")

	if user.ID == 0 {
		return true, nil
	}

	return false, nil

}

func (s *service) SaveAvatar(ctx context.Context, id int, filelocation string) (User, error) {
	user, err := s.repo.FindById(ctx, id)
	helper.PanicIfError(err, "error in srervice save avatar")

	user.AvatarFileName = filelocation
	fmt.Println(user, "user before update")
	updatedUser, err := s.repo.Update(ctx, user)
	helper.PanicIfError(err, "error in srervice save avatar")

	return updatedUser, nil

}

func (s *service) GetUserById(ctx context.Context, id int) (User, error) {
	user, err := s.repo.FindById(ctx, id)
	helper.PanicIfError(err, "error in finding user y id")

	return user, nil

}

// func (s *service) GetUserById(id int) (User, error) {
// 	user, err := s.repo.FindById(id)
// 	if err != nil {
// 		return user, err
// 	}
// 	if user.ID == 0 {
// 		return user, errors.New("no user found")
// 	}
// 	return user, nil

// }

// func (s *service) SaveAvatar(id int, filelocation string) (User, error) {

// 	user, err := s.repo.FindById(id)
// 	if err != nil {
// 		return user, err
// 	}
// 	user.AvatarFileName = filelocation
// 	updateeduser, err := s.repo.Update(user)
// 	if err != nil {
// 		return updateeduser, err
// 	}
// 	return updateeduser, nil

// }
