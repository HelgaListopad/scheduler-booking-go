package data

import (
	"scheduler-booking/common"
)

type Doctor struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float32 `json:"price"`
	Gap      int     `json:"gap"`
	SlotSize int     `json:"slot_size"`
	ImageURL string  `json:"image_url"`

	DoctorsRoutine []DoctorRoutine `json:"-"`
	OccupiedSlots  []OccupiedSlot  `json:"-"`
}

type DoctorRoutine struct {
	ID        int           `json:"id"`
	DoctorID  int           `json:"doctor_id"`
	StartDate *common.JDate `json:"start_date"`
	EndDate   *common.JDate `json:"end_date"`

	Doctor Doctor `json:"-"`
}

type OccupiedSlot struct {
	ID            int    `json:"id"`
	DoctorID      int    `json:"doctor_id"`
	Date          int64  `json:"date"`
	ClientName    string `json:"client_name"`
	ClientEmail   string `json:"client_email"`
	ClientDetails string `json:"client_details"`

	Doctor Doctor `json:"-"`
}

func (DoctorRoutine) TableName() string {
	return "doctors_routine"
}

func (OccupiedSlot) TableName() string {
	return "occupied_slots"
}
