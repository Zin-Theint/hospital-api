package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-github-name/hospital-api/internal/model"
)

type PatientSearchFilter struct {
	NationalID  *string
	PassportID  *string
	FirstName   *string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string // yyyy‑mm‑dd
	PhoneNumber *string
	Email       *string
	HospitalID  int
}

type PatientRepo struct {
	DB *pgxpool.Pool
}

func (r PatientRepo) Search(ctx context.Context, f PatientSearchFilter) ([]model.Patient, error) {
	var (
		conds []string
		args  []any
		i     = 1
	)

	// helper to add conditions
	add := func(cond string, val *string) {
		if val != nil && *val != "" {
			conds = append(conds, cond)
			args = append(args, *val)
			i++
		}
	}

	add("national_id=$"+itoa(i), f.NationalID)
	add("passport_id=$"+itoa(i), f.PassportID)
	add("first_name_en ILIKE '%'||$"+itoa(i)+"||'%'", f.FirstName)
	add("middle_name_en ILIKE '%'||$"+itoa(i)+"||'%'", f.MiddleName)
	add("last_name_en  ILIKE '%'||$"+itoa(i)+"||'%'", f.LastName)
	add("phone_number=$"+itoa(i), f.PhoneNumber)
	add("email=$"+itoa(i), f.Email)
	if f.DateOfBirth != nil && *f.DateOfBirth != "" {
		conds = append(conds, "date_of_birth=$"+itoa(i))
		args = append(args, *f.DateOfBirth)
		i++
	}

	args = append(args, f.HospitalID)
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ") + " AND hospital_id=$" + itoa(i)
	} else {
		where = "WHERE hospital_id=$" + itoa(i)
	}

	rows, err := r.DB.Query(ctx,
		`SELECT id,first_name_th,middle_name_th,last_name_th,
		        first_name_en,middle_name_en,last_name_en,
		        date_of_birth,patient_hn,national_id,passport_id,
		        phone_number,email,gender,hospital_id
		 FROM patients `+where, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Patient
	for rows.Next() {
		var p model.Patient
		err = rows.Scan(
			&p.ID, &p.FirstNameTH, &p.MiddleNameTH, &p.LastNameTH,
			&p.FirstNameEN, &p.MiddleNameEN, &p.LastNameEN,
			&p.DateOfBirth, &p.PatientHN, &p.NationalID, &p.PassportID,
			&p.PhoneNumber, &p.Email, &p.Gender, &p.HospitalID,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func itoa(i int) string { return strconv.Itoa(i) }
