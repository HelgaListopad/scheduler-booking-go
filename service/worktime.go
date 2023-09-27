package service

import (
	"fmt"
	"scheduler-booking/common"
	"scheduler-booking/data"
	"time"
)

type worktimeService struct {
	dao *data.DAO
}

type WorktimeO struct {
	DoctorID int   `json:"doctor_id"`
	From     int   `json:"from"`
	To       int   `json:"to"`
	Days     []int `json:"days"`
	Dates    []int `json:"dates"`
}

type Worktime struct {
	DoctorID  int           `json:"doctor_id"`
	StartDate *common.JDate `json:"start_date"`
	EndDate   *common.JDate `json:"end_date"`
}

type DoctorRoutineStr struct {
	ID        int    `json:"id"`
	DoctorID  int    `json:"doctor_id"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

// returns records for the Scheduler Doctors View
func (s *worktimeService) GetRoutine() ([]DoctorRoutineStr, error) {
	schedule, err := s.dao.DoctorsSchedule.GetAllSchedule()
	out := make([]DoctorRoutineStr, 0)
	for _, sch := range schedule {
		if len(sch.DoctorRoutine) == 0 {
			continue
		}

		loc := time.Now().Location()

		routine := sch.DoctorRoutine[0]

		strFormat := "2006-01-02 15:04:05"

		y, m, d := time.UnixMilli(routine.Date).Date()
		fh := sch.From / 60
		fm := sch.From % 60
		th := sch.To / 60
		tm := sch.To % 60

		r := DoctorRoutineStr{
			ID:        sch.ID,
			DoctorID:  sch.DoctorID,
			StartDate: time.Date(y, m, d, fh, fm, 0, 0, loc).Format(strFormat),
			EndDate:   time.Date(y, m, d, th, tm, 0, 0, loc).Format(strFormat),
		}

		out = append(out, r)
	}

	return out, err
}

// adds doctor's schedule for the specific day
func (s *worktimeService) Add(data Worktime) (int, error) {
	if err := data.validate(); err != nil {
		return 0, err
	}

	date := data.StartDate.UnixMilli()
	from := data.StartDate.Hour()*60 + data.StartDate.Minute()
	to := data.EndDate.Hour()*60 + data.EndDate.Minute()

	id, err := s.dao.DoctorsSchedule.AddRoutineOnDate(data.DoctorID, from, to, date)

	return id, err
}

// updates doctor's schedule for the specifc day
func (s *worktimeService) UpdateDateSchedule(scheduleId int, data Worktime) error {
	schedule, err := s.dao.DoctorsSchedule.GetOne(scheduleId)
	if err != nil {
		return err
	}
	if schedule.ID == 0 {
		return fmt.Errorf("schedule with id %d not found", scheduleId)
	}
	if err := data.validate(); err != nil {
		return err
	}

	date := data.StartDate.UnixMilli()
	from := data.StartDate.Hour()*60 + data.StartDate.Minute()
	to := data.EndDate.Hour()*60 + data.EndDate.Minute()

	err = s.dao.DoctorsSchedule.UpdateDateSchedule(scheduleId, data.DoctorID, from, to, date)

	return err
}

// delets doctor's schedule for the specific day
func (s *worktimeService) Delete(id int) error {
	return s.dao.DoctorsSchedule.Delete(id)
}

func (w Worktime) validate() error {
	if w.StartDate.UnixMilli() < time.Now().UnixMilli() {
		return fmt.Errorf("cannot set work time in the past")
	}
	if w.StartDate.UnixMilli() >= w.EndDate.UnixMilli() {
		return fmt.Errorf("invalid time interval")
	}
	return nil
}
