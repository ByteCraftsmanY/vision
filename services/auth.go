package services

import (
	"math/rand"
	"strconv"
	"time"
	"vision/daos"
	"vision/dtos"
	"vision/repositories"
	"vision/types"
)

type AuthService interface {
	GetAuth(data *daos.Auth) (*daos.Auth, error)
	GetAuthByPhone(phone string) (*daos.Auth, error)
	CreateNewAuth(data *daos.Auth) (*daos.Auth, error)
	CreateNewAuthWithTTL(form *dtos.OTPCreateRequest, duration time.Duration) (*daos.Auth, error)
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

func (s *authService) GetAuthByPhone(phone string) (*daos.Auth, error) {
	return s.AuthRepository.Find(&daos.Auth{Phone: phone})
}

func (s *authService) GetAuth(data *daos.Auth) (*daos.Auth, error) {
	return s.AuthRepository.Find(data)
}

func (s *authService) CreateNewAuth(data *daos.Auth) (*daos.Auth, error) {
	//data.Initialize()
	err := s.AuthRepository.Save(data)
	return data, err
}

func (s *authService) CreateNewAuthWithTTL(form *dtos.OTPCreateRequest, duration time.Duration) (*daos.Auth, error) {
	//data.Initialize()
	currentTime := time.Now()
	randomNumber := rand.Intn(9000) + 1000
	data := &daos.Auth{
		Phone:     form.Phone,
		Code:      strconv.Itoa(randomNumber),
		CreatedAt: &currentTime,
	}
	err := s.AuthRepository.SaveWithTTL(data, duration)
	return data, err
}

func (s *authService) DeleteAuth(id *types.ID) error {
	data := daos.Auth{}
	return s.AuthRepository.Delete(&data)
}
