package data

import (
	"fmt"

	"gorm.io/gorm"
)

type occupiedSlotsDAO struct {
	db *gorm.DB
}

type Query struct {
	DoctorID  int
	EqualDate int64
	MinDate   int64
	MaxDate   int64
}

func newOccupiedSlotsDAO(db *gorm.DB) *occupiedSlotsDAO {
	return &occupiedSlotsDAO{db}
}

func (d *occupiedSlotsDAO) GetOne(id int) (OccupiedSlot, error) {
	slot := OccupiedSlot{}
	err := d.db.Find(&slot, id).Error
	return slot, err
}

func (d *occupiedSlotsDAO) GetWithQuery(query Query) ([]OccupiedSlot, error) {
	args, err := query.sql()
	if err != nil {
		return nil, err
	}

	slots := make([]OccupiedSlot, 0)
	err = d.db.Find(&slots, args...).Error
	return slots, err
}

func (d *occupiedSlotsDAO) Add(doctor int, date int64, name, email, details string) (int, error) {
	record := OccupiedSlot{
		DoctorID:      doctor,
		Date:          date,
		ClientName:    name,
		ClientEmail:   email,
		ClientDetails: details,
	}
	err := d.db.Save(&record).Error
	return record.ID, err
}

func (d *occupiedSlotsDAO) Delete(id int) error {
	return d.db.Delete(&OccupiedSlot{}, id).Error
}

func (q Query) sql() ([]any, error) {
	if q.EqualDate != 0 && (q.MinDate != 0 && q.MaxDate != 0) {
		return nil, fmt.Errorf("query not valid: EqualDate or MinDate/MaxDate must be defined, not both")
	}

	sql := ""
	args := make([]any, 0)

	if q.DoctorID != 0 {
		sql += "doctor_id = ?"
		args = append(args, q.DoctorID)
	}
	if q.EqualDate != 0 {
		sql += " date = ?"
		args = append(args, q.EqualDate)
	}
	if q.MaxDate != 0 {
		sql += " date >= ?"
		args = append(args, q.MinDate)
	}
	if q.EqualDate != 0 {
		sql += " date <= ?"
		args = append(args, q.MaxDate)
	}

	if sql != "" {
		args = append([]any{sql}, args)
	}

	return args, nil
}
