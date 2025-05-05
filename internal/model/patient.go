package model

import (
	"database/sql"
)

type Patient struct {
	ID           int
	FirstNameTH  sql.NullString
	MiddleNameTH sql.NullString
	LastNameTH   sql.NullString
	FirstNameEN  sql.NullString
	MiddleNameEN sql.NullString
	LastNameEN   sql.NullString
	DateOfBirth  sql.NullTime
	PatientHN    sql.NullString
	NationalID   sql.NullString
	PassportID   sql.NullString
	PhoneNumber  sql.NullString
	Email        sql.NullString
	Gender       sql.NullString
	HospitalID   int
}
