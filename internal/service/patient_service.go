package service

import (
	"context"

	"github.com/Zin-Theint/hospital-api/internal/model"
	"github.com/Zin-Theint/hospital-api/internal/repository"
)

type PatientService struct {
	Repo repository.PatientRepo
}

func (p PatientService) Search(ctx context.Context, f repository.PatientSearchFilter) ([]model.Patient, error) {
	return p.Repo.Search(ctx, f)
}
