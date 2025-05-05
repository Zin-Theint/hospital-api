package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/Zin-Theint/hospital-api/internal/testutil"
)

func TestPatientSearch(t *testing.T) {
	ctx := context.Background()
	db, err := testutil.NewTestDB(ctx)
	require.NoError(t, err)
	defer db.Terminate(ctx)

	r := testutil.NewTestRouter(db.Pool)

	// seed one staff (id=1) & one patient
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO staff (username, password_hash, hospital_id)
		VALUES ('bob',
		        '$2y$10$hWx2L4XDT0yRCsNdxhthR.FsLYE0kAQN0eT82J36M.Fn03TaqDtXu', -- bcrypt("xxx")
		        1);

		INSERT INTO patients (
		    first_name_en, last_name_en,
		    hospital_id,   national_id,
		    date_of_birth
		) VALUES (
		    'John', 'Doe',
		    1,       '999',
		    '2000-01-01'   
		);
	`)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": 1,
		"hid": 1,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	})
	tok, _ := token.SignedString([]byte(testutil.JwtKey))

	// 1) happy search by national_id
	w := performAuth(r, "/patient/search?national_id=999", tok)
	require.Equal(t, http.StatusOK, w.Code)

	var got []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	require.Len(t, got, 1)

	// 2) search with wrong hospital (should return empty)
	badTok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": 2,
		"hid": 2,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}).SignedString([]byte(testutil.JwtKey))

	w = performAuth(r, "/patient/search?national_id=999", badTok)
	require.Equal(t, http.StatusOK, w.Code)

	_ = json.Unmarshal(w.Body.Bytes(), &got)
	require.Len(t, got, 0)
}

func performAuth(r http.Handler, path, tok string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
