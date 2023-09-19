package service

import (
	"scheduler-booking/data"
	"time"
)

type unitsService struct {
	dao *data.DAO
}

type Unit struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Subtitle string  `json:"subtitle"`
	Details  string  `json:"details"`
	Preview  string  `json:"preview"`
	Price    float32 `json:"price"`

	Slots          []Schedule `json:"slots"`
	AvailableSlots []int64    `json:"availableSlots,omitempty"`
	UsedSlots      []int64    `json:"usedSlots,omitempty"`
}

type Schedule struct {
	From  int     `json:"from"`
	To    int     `json:"to"`
	Size  int     `json:"size"`
	Gap   int     `json:"gap"`
	Days  []int   `json:"days,omitempty"`
	Dates []int64 `json:"dates,omitempty"`
}

func (s *unitsService) GetAll() ([]Unit, error) {
	doctors, err := s.dao.Doctors.GetAll(true)
	if err != nil {
		return nil, err
	}

	units := make([]Unit, len(doctors))

	for i := range doctors {
		d := &doctors[i]

		usedSlots := make([]int64, len(d.OccupiedSlots))
		for j := range d.OccupiedSlots {
			usedSlots[j] = d.OccupiedSlots[j].Date
		}

		schedule := make([]Schedule, 0)
		for j := range d.DoctorSchedule {
			sch := &d.DoctorSchedule[j]

			days := make([]int, len(sch.DoctorRecurringRoutine))
			for i := range sch.DoctorRecurringRoutine {
				days[i] = sch.DoctorRecurringRoutine[i].WeekDay
			}

			dates := make([]int64, 0)
			for i := range sch.DoctorRoutine {
				if sch.DoctorRoutine[i].Date > time.Now().UnixMilli() {
					dates = append(dates, sch.DoctorRoutine[i].Date)
				}
			}
			if len(dates) == 0 && len(days) == 0 {
				// skip this rule as it is already expired
				continue
			}

			schedule = append(schedule, Schedule{
				From:  sch.From,
				To:    sch.To,
				Size:  d.SlotSize,
				Gap:   d.Gap,
				Days:  days,
				Dates: dates,
			})
		}

		units[i] = Unit{
			ID:        d.ID,
			Title:     d.Name,
			Subtitle:  "",
			Details:   "",
			Category:  d.Category,
			Price:     d.Price,
			Preview:   d.ImageURL,
			UsedSlots: usedSlots,
			Slots:     schedule,
		}
	}

	return units, nil
}
