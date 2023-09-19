package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type doctorsScheduleDAO struct {
	db *gorm.DB
}

func newDoctorsScheduleDAO(db *gorm.DB) *doctorsScheduleDAO {
	return &doctorsScheduleDAO{db}
}

func (d *doctorsScheduleDAO) GetOne(id int) (DoctorSchedule, error) {
	data := DoctorSchedule{}
	err := d.db.
		Preload("DoctorRoutine").
		Find(&data, id).Error
	return data, err
}

func (d *doctorsScheduleDAO) GetAllSchedule() ([]DoctorSchedule, error) {
	sch := make([]DoctorSchedule, 0)

	err := d.db.
		Preload("DoctorRoutine", "date > ?", time.Now().UnixMilli()).
		Find(&sch).Error

	return sch, err
}

func (d *doctorsScheduleDAO) AddRoutineOnDate(doctorId, from, to int, date int64) (int, error) {
	if date == 0 {
		return 0, errors.New("date argument not defined")
	}

	schedule := DoctorSchedule{
		From:     from,
		To:       to,
		DoctorID: doctorId,
		DoctorRoutine: []DoctorRoutine{
			{
				Date: date,
			},
		},
	}

	err := d.db.Save(&schedule).Error

	return schedule.ID, err
}

func (d *doctorsScheduleDAO) UpdateDateSchedule(id, docId, from, to int, date int64) (err error) {

	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	schedule, err := d.GetOne(id)
	if err != nil {
		return err
	}

	if schedule.DoctorID == docId {
		// delete routine in schedule becaus we dont know routine id
		err = tx.Delete(&DoctorRoutine{}, "schedule_id = ?", id).Error
		if err != nil {
			return err
		}
		// update routine by creating new one
		schedule.From = from
		schedule.To = to
		schedule.DoctorRoutine = []DoctorRoutine{
			{
				Date: date,
			},
		}
		err = tx.Save(&schedule).Error
	} else {
		// delete schedule at all for this old doctor and create schedule for docId
		err = tx.Delete(&DoctorSchedule{}, id).Error
		if err != nil {
			return err
		}

		schedule = DoctorSchedule{
			DoctorID: docId,
			From:     from,
			To:       to,
			DoctorRoutine: []DoctorRoutine{
				{
					Date: date,
				},
			},
		}
		err = tx.Create(&schedule).Error
	}

	return err
}

func (d *doctorsScheduleDAO) Delete(id int) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Delete(&DoctorRoutine{}, "where schedule_id = ?", id).Error
	if err != nil {
		return err
	}

	err = tx.Delete(&DoctorRecurringRoutine{}, "where schedule_id = ?", id).Error
	if err != nil {
		return err
	}

	err = tx.Delete(&DoctorSchedule{}, id).Error

	return err
}
