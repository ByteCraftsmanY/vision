package services

import (
	"golang.org/x/crypto/bcrypt"
	"vision/entities"
	"vision/repositories"
	"vision/types"
)

type UserService interface {
	GetUsers() []*entities.User
	GetUserByID(userID *types.ID) (*entities.User, error)
	GetUserByPhone(phone string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	CreateNewUser(data *entities.User) (*entities.User, error)
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

func (s *userService) GetUsers() []*entities.User {
	//TODO implement me
	panic("implement me")
}

func (s *userService) GetUserByID(userID *types.ID) (*entities.User, error) {
	return s.UserRepository.FindByID(userID)
}

func (s *userService) GetUserByPhone(phone string) (*entities.User, error) {
	return s.UserRepository.FindByPhoneOrEmail(&entities.User{Phone: phone})

}

func (s *userService) GetUserByEmail(email string) (*entities.User, error) {
	return s.UserRepository.FindByPhoneOrEmail(&entities.User{Email: email})
}

func (s *userService) CreateNewUser(user *entities.User) (*entities.User, error) {
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
	user := &entities.User{Base: entities.Base{ID: userID}}
	return s.UserRepository.Delete(user)
}
