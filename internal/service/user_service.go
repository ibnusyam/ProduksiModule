package service

import (
	"Steril-App/internal/repository"
	"Steril-App/model"
	"errors"
	"fmt"
)

type UserService struct {
	UserRepo   *repository.UserRepository   // Repo utama (Ganti nama field biar jelas)
	FingerRepo *repository.FingerRepository // <--- Repo "yang lain" ditambahkan di sini
}

func NewUserService(userRepo *repository.UserRepository, fingerRepo *repository.FingerRepository) *UserService {
	return &UserService{
		UserRepo:   userRepo,
		FingerRepo: fingerRepo,
	}
}

var ErrUserAlreadyExists = errors.New("user sudah terdaftar")

func (s *UserService) CreateUser(data *model.CreateUserRequest) error {
	fmt.Println(data)
	isExist, err := s.UserRepo.IsUserExist(data)
	if err != nil {
		return fmt.Errorf("gagal menjalankan method is user")
	}
	fmt.Println(isExist)
	if isExist != "" {
		return ErrUserAlreadyExists
	}
	err = s.UserRepo.CreateUser(data)
	if err != nil {
		// fmt.Println(err)
		return fmt.Errorf("gagal menjalankan method create user")
	}

	const RequiredFingerSlots = 3

	for i := 0; i < RequiredFingerSlots; i++ {
		fingerId, err := s.FingerRepo.FindEmptyFingerSlot()
		if err != nil {
			return fmt.Errorf("gagal melakukan pencarian slot kosong :%w", err)
		}
		err = s.FingerRepo.AddFingerData(data.NIK, fingerId)
		if err != nil {
			return fmt.Errorf("gagal menambahkan data id pada user : %w", err)
		}
	}

	return nil
}

func (s *UserService) DeleteUser(id string) error {
	err := s.UserRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("gagal Menjalankan service delete user :%w", err)
	}

	return nil
}

func (s *UserService) GetAllUser() ([]model.UserResponse, error) {
	rows, err := s.UserRepo.GetAllUser()
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan service :%w", err)
	}

	return rows, nil
}
