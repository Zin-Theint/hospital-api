package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Zin-Theint/hospital-api/internal/model"
	"github.com/Zin-Theint/hospital-api/internal/repository"
)

type StaffService struct {
	Repo repository.StaffRepo
}

func (s StaffService) Create(ctx context.Context, username, password string, hospitalID int) (*model.Staff, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	staff := &model.Staff{
		Username:     username,
		PasswordHash: string(hash),
		HospitalID:   hospitalID,
	}
	if err := s.Repo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func (s StaffService) Authenticate(ctx context.Context, username, password string) (*model.Staff, error) {
	staff, err := s.Repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return staff, nil
}
