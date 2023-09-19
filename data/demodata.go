package data

import (
	"time"

	"gorm.io/gorm"
)

func dataDown(tx *gorm.DB) {
	must(tx.Exec("DELETE FROM `doctors`").Error)
	must(tx.Exec("DELETE FROM `doctor_routine`").Error)
	must(tx.Exec("DELETE FROM `doctor_schedule`").Error)
	must(tx.Exec("DELETE FROM `doctor_recurring`").Error)
	must(tx.Exec("DELETE FROM `occupied_slots`").Error)
}

func dataUp(tx *gorm.DB) {
	nowDate := time.Now()
	y, m, d := nowDate.Date()
	now := time.Date(y, m, d, 0, 0, 0, 0, nowDate.Location())
	genSchedule := func(from, to, days int, offset int) []DoctorSchedule {
		routine := make([]DoctorSchedule, 0)
		for i := 0; i < days; i++ {
			workDay := DoctorSchedule{
				From: from,
				To:   to,
				DoctorRoutine: []DoctorRoutine{
					{
						Date: now.AddDate(0, 0, i+offset+1).UnixMilli(),
					},
				},
			}
			routine = append(routine, workDay)
		}
		return routine
	}

	doctors := []Doctor{
		{
			Name:     "Toby Cannon",
			Category: "Dantist",
			SlotSize: 20,
			Price:    45,
			ImageURL: "",
			Gap:      20,
			DoctorSchedule: append(
				genSchedule(8*60, 16*60, 7, -1),
				DoctorSchedule{
					// every week day from 9:00-17:00
					From: 9 * 60,
					To:   17 * 60,
					// sun, sat - holidays
					DoctorRecurringRoutine: []DoctorRecurringRoutine{
						{
							WeekDay: 0,
						},
						{
							WeekDay: 1,
						},
					},
				},
			),
		},
		{
			Name:     "Larissa Tillery",
			Category: "Allergist",
			SlotSize: 45,
			Price:    120,
			ImageURL: "",
			Gap:      5,
			DoctorSchedule: append(
				[]DoctorSchedule{
					{
						// mon, wed, fri 9:00-14:00
						From: 9 * 60,
						To:   14 * 60,
						DoctorRecurringRoutine: []DoctorRecurringRoutine{
							{
								WeekDay: 1,
							},
							{
								WeekDay: 3,
							},
							{
								WeekDay: 5,
							},
						},
					},
					{
						// tue, thu 15:00-20:00
						From: 15 * 60,
						To:   20 * 60,
						DoctorRecurringRoutine: []DoctorRecurringRoutine{
							{
								WeekDay: 2,
							},
							{
								WeekDay: 4,
							},
						},
					},
				},
				genSchedule(8*60, 14*60, 5, 1)...,
			),
		},
		{
			Name:     "Raymond Castro",
			Category: "Oculist",
			SlotSize: 25,
			Price:    35,
			ImageURL: "",
			Gap:      0,
			DoctorSchedule: []DoctorSchedule{
				{
					// mon, wen, wed fri 8:00-17:00
					From: 8 * 60,
					To:   17 * 60,
					DoctorRecurringRoutine: []DoctorRecurringRoutine{
						{
							WeekDay: 1,
						},
						{
							WeekDay: 2,
						},
						{
							WeekDay: 3,
						},
					},
				},
			},
		},
		{
			Name:     "Dexter Porter",
			Category: "Oculist",
			SlotSize: 25,
			Price:    40,
			ImageURL: "",
			Gap:      5,
			DoctorSchedule: []DoctorSchedule{
				{
					// every week day from 9:00-17:00
					From: 9 * 60,
					To:   17 * 60,
					// sun, sat - holidays
					DoctorRecurringRoutine: []DoctorRecurringRoutine{
						{
							WeekDay: 0,
						},
						{
							WeekDay: 1,
						},
					},
				},
			},
		},
		{
			Name:     "Alan Boren",
			Category: "Dantist",
			SlotSize: 60,
			Price:    175,
			ImageURL: "",
			Gap:      10,
			DoctorSchedule: append(
				[]DoctorSchedule{
					{
						// every week day from 9:00-17:00
						From: 9 * 60,
						To:   17 * 60,
						// sun, sat - holidays
						DoctorRecurringRoutine: []DoctorRecurringRoutine{
							{
								WeekDay: 0,
							},
							{
								WeekDay: 1,
							},
						},
					},
				},
				genSchedule(14*60, 21*60, 7, 0)...,
			),
		},
	}

	err := tx.Create(doctors).Error
	if err != nil {
		panic(err)
	}
}
