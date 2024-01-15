package services

import (
	"golang.org/x/crypto/bcrypt"
	"vision/daos"
	"vision/repositories"
	"vision/types"
)

type UserService interface {
	GetUsers() []*daos.User
	GetUserByID(userID *types.ID) (*daos.User, error)
	GetUserByPhone(phone string) (*daos.User, error)
	GetUserByEmail(email string) (*daos.User, error)
	CreateNewUser(data *daos.User) (*daos.User, error)
	DeleteUser(userID *types.ID) error
}

type userService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(
	userRepository repositories.UserRepositoryInterface,
) UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

func (s *userService) GetUsers() []*daos.User {
	//TODO implement me
	panic("implement me")
}

func (s *userService) GetUserByID(userID *types.ID) (*daos.User, error) {
	return s.UserRepository.FindByID(userID)
}

func (s *userService) GetUserByPhone(phone string) (*daos.User, error) {
	return s.UserRepository.FindByPhoneOrEmail(&daos.User{Phone: phone})

}

func (s *userService) GetUserByEmail(email string) (*daos.User, error) {
	return s.UserRepository.FindByPhoneOrEmail(&daos.User{Email: email})
}

func (s *userService) CreateNewUser(user *daos.User) (*daos.User, error) {
	user.Initialize()
	passBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passBytes)
	/*isExists, err := u.UserRepository.FindByPhoneOrEmail(user)
	if isExists {
		return user, fmt.Errorf("already Exists")
	}*/
	err = s.UserRepository.Save(user)
	return user, err
}

func (s *userService) DeleteUser(userID *types.ID) error {
	user := &daos.User{Base: daos.Base{ID: userID}}
	return s.UserRepository.Delete(user)
}
