package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zin-Theint/hospital-api/internal/repository"
	"github.com/Zin-Theint/hospital-api/internal/service"
)

type PatientHandler struct {
	Service service.PatientService
}

func (h PatientHandler) Search(c *gin.Context) {
	hospitalID := c.GetInt("hospitalID")

	filter := repository.PatientSearchFilter{
		HospitalID: hospitalID,

		// IDs
		NationalID: ptr(c.Query("national_id")),
		PassportID: ptr(c.Query("passport_id")),

		// English names
		FirstNameEN:  ptr(c.Query("first_name_en")),
		MiddleNameEN: ptr(c.Query("middle_name_en")),
		LastNameEN:   ptr(c.Query("last_name_en")),

		// Thai names
		FirstNameTH:  ptr(c.Query("first_name_th")),
		MiddleNameTH: ptr(c.Query("middle_name_th")),
		LastNameTH:   ptr(c.Query("last_name_th")),

		// Others
		DateOfBirth: ptr(c.Query("date_of_birth")),
		PhoneNumber: ptr(c.Query("phone_number")),
		Email:       ptr(c.Query("email")),
	}

	patients, err := h.Service.Search(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patients)
}

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
