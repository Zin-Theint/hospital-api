package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zin-Theint/hospital-api/internal/model"
)

// -------- search filter --------

type PatientSearchFilter struct {
	// IDs
	NationalID *string
	PassportID *string

	// English names
	FirstNameEN  *string
	MiddleNameEN *string
	LastNameEN   *string

	// Thai names
	FirstNameTH  *string
	MiddleNameTH *string
	LastNameTH   *string

	DateOfBirth *string // yyyy-mm-dd

	PhoneNumber *string
	Email       *string

	HospitalID int // injected from JWT, never nil
}

// -------- repository --------

type PatientRepo struct {
	DB *pgxpool.Pool
}

func (r PatientRepo) Search(ctx context.Context, f PatientSearchFilter) ([]model.Patient, error) {
	var (
		conds []string
		args  []any
		i     = 1
	)

	add := func(cond string, val *string) {
		if val != nil && *val != "" {
			conds = append(conds, cond)
			args = append(args, *val)
			i++
		}
	}

	add("national_id=$"+strconv.Itoa(i), f.NationalID)
	add("passport_id=$"+strconv.Itoa(i), f.PassportID)

	add("first_name_en  ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.FirstNameEN)
	add("middle_name_en ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.MiddleNameEN)
	add("last_name_en   ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.LastNameEN)

	add("first_name_th  ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.FirstNameTH)
	add("middle_name_th ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.MiddleNameTH)
	add("last_name_th   ILIKE '%'||$"+strconv.Itoa(i)+"||'%'", f.LastNameTH)

	add("phone_number=$"+strconv.Itoa(i), f.PhoneNumber)
	add("email=$"+strconv.Itoa(i), f.Email)

	if f.DateOfBirth != nil && *f.DateOfBirth != "" {
		conds = append(conds, "date_of_birth=$"+strconv.Itoa(i))
		args = append(args, *f.DateOfBirth)
		i++
	}

	// Every search is always scoped to the caller’s hospital
	args = append(args, f.HospitalID)
	where := "WHERE hospital_id=$" + strconv.Itoa(i)
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ") + " AND hospital_id=$" + strconv.Itoa(i)
	}

	rows, err := r.DB.Query(ctx, `
		SELECT id,
		       first_name_th,  middle_name_th,  last_name_th,
		       first_name_en,  middle_name_en,  last_name_en,
		       date_of_birth, patient_hn,
		       national_id,   passport_id,
		       phone_number,  email, gender,
		       hospital_id
		FROM   patients `+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Patient
	for rows.Next() {
		var p model.Patient
		err = rows.Scan(
			&p.ID,
			&p.FirstNameTH, &p.MiddleNameTH, &p.LastNameTH,
			&p.FirstNameEN, &p.MiddleNameEN, &p.LastNameEN,
			&p.DateOfBirth, &p.PatientHN,
			&p.NationalID, &p.PassportID,
			&p.PhoneNumber, &p.Email, &p.Gender,
			&p.HospitalID,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}
