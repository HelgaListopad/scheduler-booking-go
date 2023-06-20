package data

import (
	"scheduler-booking/common"
	"time"

	"gorm.io/gorm"
)

type doctorsRoutineDAO struct {
	db *gorm.DB
}

func newDoctorsRoutineDAO(db *gorm.DB) *doctorsRoutineDAO {
	return &doctorsRoutineDAO{db}
}

func (d *doctorsRoutineDAO) GetOne(id int) (DoctorRoutine, error) {
	data := DoctorRoutine{}
	err := d.db.Find(&data, id).Error
	return data, err
}

func (d *doctorsRoutineDAO) GetAll(notExpired bool) ([]DoctorRoutine, error) {
	data := make([]DoctorRoutine, 0)
	var err error
	if notExpired {
		err = d.db.Find(&data, "start_date > ?", time.Now().UnixMilli()).Error
	} else {
		err = d.db.Find(&data).Error
	}
	return data, err
}

func (d *doctorsRoutineDAO) GetRoutineByTime(doctorId int, date time.Time) ([]DoctorRoutine, error) {
	data := make([]DoctorRoutine, 0)
	err := d.db.Find(&data, "start_date > ? AND end_date < ?", date, date).Error
	return data, err
}

func (d *doctorsRoutineDAO) Add(doctorId int, start, end *common.JDate) (int, error) {
	data := DoctorRoutine{
		DoctorID:  doctorId,
		StartDate: start,
		EndDate:   end,
	}
	err := d.db.Save(&data).Error
	return data.ID, err
}

func (d *doctorsRoutineDAO) Update(id, doctorId int, start, end *common.JDate) error {
	err := d.db.Where("id = ?", id).Updates(&DoctorRoutine{
		DoctorID:  doctorId,
		StartDate: start,
		EndDate:   end,
	}).Error

	return err
}

func (d *doctorsRoutineDAO) Delete(id int) error {
	return d.db.Delete(&DoctorRoutine{}, id).Error
}
