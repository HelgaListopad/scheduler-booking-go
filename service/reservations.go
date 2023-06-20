package service

import (
	"scheduler-booking/data"
	"fmt"
	"time"
)

type reservationsService struct {
	dao *data.DAO
}

type ReservationForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Details string `json:"details"`
}

type Reservation struct {
	DoctorID int             `json:"doctor"`
	Date     int64           `json:"date"`
	Form     ReservationForm `json:"data"`
}

func (s *reservationsService) GetAll() ([]data.OccupiedSlot, error) {
	records, err := s.dao.OccupiedSlots.GetWithQuery(data.Query{})
	return records, err
}

func (s *reservationsService) Add(r Reservation) (int, error) {
	// check if reservation time is available and has not expired
	err := s.checkIfReservationIsAvailable(r.DoctorID, r.Date)
	if err != nil {
		return 0, err
	}
	// check if doctor's schedule has exactly requested time
	err = s.checkIfReservationIsCorrect(r.DoctorID, r.Date)
	if err != nil {
		return 0, err
	}

	id, err := s.dao.OccupiedSlots.Add(
		r.DoctorID,
		r.Date,
		r.Form.Name,
		r.Form.Email,
		r.Form.Details,
	)

	return id, err
}

func (s *reservationsService) Delete(id int) error {
	slot, err := s.dao.OccupiedSlots.GetOne(id)
	if err != nil {
		return err
	}
	if slot.Date < time.Now().UnixMilli() {
		return fmt.Errorf("cannot delete reservation that time has expired")
	}
	return s.dao.OccupiedSlots.Delete(id)
}

func (s *reservationsService) checkIfReservationIsAvailable(doctorId int, date int64) error {
	slots, err := s.dao.OccupiedSlots.GetWithQuery(data.Query{
		DoctorID:  doctorId,
		EqualDate: date,
	})
	if err != nil {
		return err
	}
	if len(slots) > 0 && slots[0].ID != 0 {
		return fmt.Errorf("this time is already booked")
	}
	if date < time.Now().UnixMilli() {
		return fmt.Errorf("booking time has expired")
	}

	return err
}

func (s *reservationsService) checkIfReservationIsCorrect(doctorId int, date int64) error {
	doctor, err := s.dao.Doctors.GetOne(doctorId)
	if err != nil {
		return err
	}
	if doctor.ID == 0 {
		return fmt.Errorf("doctor with id %d not found", doctorId)
	}

	w, err := s.dao.DoctorsRoutine.GetRoutineByTime(doctor.ID, time.UnixMilli(date))
	if err != nil {
		return err
	}
	l := len(w)
	if l == 0 {
		return fmt.Errorf("booking time not valid: doctor's schedule undefined")
	}
	if len(w) > 1 {
		return fmt.Errorf("doctor's schedule is ambiguously")
	}

	schedule := w[0]
	checkDate := schedule.StartDate.Date().UnixMilli()
	endDate := schedule.EndDate.Date().UnixMilli()
	step := int64((doctor.SlotSize + doctor.Gap) * 60 * 1000)
	exists := false
	for checkDate <= endDate {
		if checkDate == date {
			exists = true
			break
		}
		checkDate += step
	}
	if !exists {
		return fmt.Errorf("reservation time doesn't exist in schedule of the selected doctor")
	}

	return nil
}
