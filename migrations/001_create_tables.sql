-- +goose Up
CREATE TABLE hospitals (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

INSERT INTO hospitals (name) VALUES ('Hospital A');

CREATE TABLE staff (
    id            SERIAL PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    hospital_id   INT NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE patients (
    id             SERIAL PRIMARY KEY,
    first_name_th  TEXT,
    middle_name_th TEXT,
    last_name_th   TEXT,
    first_name_en  TEXT,
    middle_name_en TEXT,
    last_name_en   TEXT,
    date_of_birth  DATE,
    patient_hn     TEXT,
    national_id    TEXT,
    passport_id    TEXT,
    phone_number   TEXT,
    email          TEXT,
    gender         CHAR(1),
    hospital_id    INT NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE
);

CREATE INDEX idx_patients_national_id  ON patients(national_id);
CREATE INDEX idx_patients_passport_id  ON patients(passport_id);

-- +goose Down
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS hospitals;
