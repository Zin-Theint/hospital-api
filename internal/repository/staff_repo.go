package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-github-name/hospital-api/internal/model"
)

type StaffRepo struct {
	DB *pgxpool.Pool
}

func (r StaffRepo) Create(ctx context.Context, s *model.Staff) error {
	return r.DB.QueryRow(
		ctx,
		`INSERT INTO staff(username,password_hash,hospital_id)
		 VALUES($1,$2,$3) RETURNING id`,
		s.Username,
		s.PasswordHash,
		s.HospitalID,
	).Scan(&s.ID)
}

func (r StaffRepo) GetByUsername(ctx context.Context, u string) (*model.Staff, error) {
	var s model.Staff
	err := r.DB.QueryRow(
		ctx,
		`SELECT id,username,password_hash,hospital_id
		 FROM staff WHERE username=$1`, u,
	).Scan(&s.ID, &s.Username, &s.PasswordHash, &s.HospitalID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
