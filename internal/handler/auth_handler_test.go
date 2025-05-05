package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Zin-Theint/hospital-api/internal/testutil"
)

func TestStaffRegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	db, err := testutil.NewTestDB(ctx)
	require.NoError(t, err)
	defer db.Terminate(ctx)

	r := testutil.NewTestRouter(db.Pool)

	type loginResp struct {
		Token string `json:"token"`
	}

	// 1) register
	regBody := gin.H{"username": "alice", "password": "pass123", "hospital": 1}
	w := perform(r, "POST", "/staff/create", regBody)
	require.Equal(t, http.StatusCreated, w.Code)

	// duplicate username should fail
	w = perform(r, "POST", "/staff/create", regBody)
	require.Equal(t, http.StatusInternalServerError, w.Code)

	// 2)login happy path
	w = perform(r, "POST", "/staff/login", gin.H{"username": "alice", "password": "pass123"})
	require.Equal(t, http.StatusOK, w.Code)

	var lr loginResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &lr))
	require.NotEmpty(t, lr.Token)

	// 3) login wrong password
	w = perform(r, "POST", "/staff/login", gin.H{"username": "alice", "password": "wrong"})
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func perform(r http.Handler, method, path string, body gin.H) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
