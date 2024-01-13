package services

import (
	"math/rand"
	"strconv"
	"time"
	"vision/dtos"
	"vision/entities"
	"vision/repositories"
	"vision/types"
)

type AuthService interface {
	GetAuth(data *entities.Auth) (*entities.Auth, error)
	GetAuthByPhone(phone string) (*entities.Auth, error)
	CreateNewAuth(data *entities.Auth) (*entities.Auth, error)
	CreateNewAuthWithTTL(form *dtos.OTPCreateRequest, duration time.Duration) (*entities.Auth, error)
	DeleteAuth(userID *types.ID) error
}

func NewAuthService(
	authRepository repositories.AuthRepositoryInterface,
) AuthService {
	return &authService{
		AuthRepository: authRepository,
	}
}

type authService struct {
	AuthRepository repositories.AuthRepositoryInterface
}

func (s *authService) GetAuthByPhone(phone string) (*entities.Auth, error) {
	return s.AuthRepository.Find(&entities.Auth{Phone: phone})
}

func (s *authService) GetAuth(data *entities.Auth) (*entities.Auth, error) {
	return s.AuthRepository.Find(data)
}

func (s *authService) CreateNewAuth(data *entities.Auth) (*entities.Auth, error) {
	//data.Initialize()
	err := s.AuthRepository.Save(data)
	return data, err
}

func (s *authService) CreateNewAuthWithTTL(form *dtos.OTPCreateRequest, duration time.Duration) (*entities.Auth, error) {
	//data.Initialize()
	currentTime := time.Now()
	randomNumber := rand.Intn(9000) + 1000
	data := &entities.Auth{
		Phone:     form.Phone,
		Code:      strconv.Itoa(randomNumber),
		CreatedAt: &currentTime,
	}
	err := s.AuthRepository.SaveWithTTL(data, duration)
	return data, err
}

func (s *authService) DeleteAuth(id *types.ID) error {
	data := entities.Auth{}
	return s.AuthRepository.Delete(&data)
}
