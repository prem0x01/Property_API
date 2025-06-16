package handlers

import (
	"Property_App/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, %v", err)
	}

	InitAppointmentHandler(db)
	return mock, func() {
		db.Close()
	}
}

func TestAddAppointment(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	appointment := models.Appointment{
		UserID:     1,
		PropertyID: 2,
		Time:       "12:30:00",
		Date:       "2025-06-05",
		Mobile:     "1234567890",
		Address:    "Test Address",
	}

	mock.ExpectPrepare("INSERT INTO appointment").
		ExpectQuery().
		WithArgs(appointment.UserID, appointment.PropertyID, appointment.Time, appointment.Date, appointment.Mobile, appointment.Address).
		WillReturnRows(sqlmock.NewRows([]string{"appointment_id"}).AddRow(1))

	body, _ := json.Marshal(appointment)
	req := httptest.NewRequest(http.MethodPost, "/appointment", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	AppointmentHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

}

func TestDeleteAppointment(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	mock.ExpectPrepare("DELETE FROM appointment").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodDelete, "/appointment/1", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/appointment/{id}", deleteAppointment)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}
