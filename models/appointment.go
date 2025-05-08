package models

type Appointment struct {
	AppointmentID int    `json:"appointment_id"`
	UserID        int    `json:"user_id"`
	PropertyID    int    `json:"property_id"`
	Time          string `json:"time"`
	Date          string `json:"date"`
	Mobile        string `json:"mobile"`
	Address       string `json:"address"`
}
