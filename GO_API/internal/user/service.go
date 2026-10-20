package user

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("user not found")
	ErrName     = errors.New("name is required (max 120 characters)")
	ErrEmail    = errors.New("email is invalid")
	ErrCPF      = errors.New("cpf must contain exactly 11 digits")
	ErrBirthday = errors.New("birthday must use YYYY-MM-DD")
)

type Repository interface {
	List(context.Context) ([]User, error)
	Get(context.Context, int64) (User, error)
	Create(context.Context, Input) (User, error)
	Update(context.Context, int64, Input) (User, error)
	Delete(context.Context, int64) (bool, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) List(ctx context.Context) ([]User, error) {
	return service.repository.List(ctx)
}

func (service *Service) Get(ctx context.Context, id int64) (User, error) {
	user, err := service.repository.Get(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (service *Service) Create(ctx context.Context, input Input) (User, error) {
	cleanInput, err := validate(input)
	if err != nil {
		return User{}, err
	}
	return service.repository.Create(ctx, cleanInput)
}

func (service *Service) Update(ctx context.Context, id int64, input Input) (User, error) {
	cleanInput, err := validate(input)
	if err != nil {
		return User{}, err
	}
	return service.repository.Update(ctx, id, cleanInput)
}

func (service *Service) Delete(ctx context.Context, id int64) error {
	deleted, err := service.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func validate(input Input) (Input, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.CPF = strings.TrimSpace(input.CPF)

	if input.Name == "" || len(input.Name) > 120 {
		return Input{}, ErrName
	}
	if len(input.Email) > 254 || !emailPattern.MatchString(input.Email) {
		return Input{}, ErrEmail
	}
	if !cpfPattern.MatchString(input.CPF) {
		return Input{}, ErrCPF
	}
	if _, err := time.Parse("2006-01-02", input.Birthday); err != nil {
		return Input{}, ErrBirthday
	}
	return input, nil
}
