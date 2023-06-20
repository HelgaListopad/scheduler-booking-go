package service

import (
	"scheduler-booking/common"
	"scheduler-booking/data"
	"fmt"
	"time"
)

type worktimeService struct {
	dao *data.DAO
}

type Worktime struct {
	DoctorID  int           `json:"doctor_id"`
	StartDate *common.JDate `json:"start_date"`
	EndDate   *common.JDate `json:"end_date"`
}

func (s *worktimeService) GetAll() ([]data.DoctorRoutine, error) {
	records, err := s.dao.DoctorsRoutine.GetAll(false)
	return records, err
}

func (s *worktimeService) Add(data Worktime) (int, error) {
	if err := data.validate(); err != nil {
		return 0, err
	}
	id, err := s.dao.DoctorsRoutine.Add(data.DoctorID, data.StartDate, data.EndDate)
	return id, err
}

func (s *worktimeService) Update(id int, data Worktime) error {
	routine, err := s.dao.DoctorsRoutine.GetOne(id)
	if err != nil {
		return err
	}
	if routine.ID == 0 {
		return fmt.Errorf("record not found")
	}
	if err := data.validate(); err != nil {
		return err
	}

	err = s.dao.DoctorsRoutine.Update(id, data.DoctorID, data.StartDate, data.EndDate)
	return err
}

func (s *worktimeService) Delete(id int) error {
	w, err := s.dao.DoctorsRoutine.GetOne(id)
	if err != nil {
		return err
	}
	if w.StartDate.Date().UnixMilli() < time.Now().UnixMilli() {
		return fmt.Errorf("cannot delete schedule from the past")
	}

	data, err := s.dao.OccupiedSlots.GetWithQuery(data.Query{
		DoctorID: w.DoctorID,
		MinDate:  w.StartDate.Date().UnixMilli(),
		MaxDate:  w.EndDate.Date().UnixMilli(),
	})
	if err != nil {
		return err
	}
	if len(data) > 0 {
		return fmt.Errorf("cannot delete schedule until it has reservations")
	}

	return s.dao.DoctorsRoutine.Delete(id)
}

func (w Worktime) validate() error {
	if w.StartDate.Date().UnixMilli() < time.Now().UnixMilli() {
		return fmt.Errorf("cannot set work time in the past")
	}
	if w.StartDate.Date().UnixMilli() >= w.EndDate.Date().UnixMilli() {
		return fmt.Errorf("invalid time interval")
	}
	return nil
}
