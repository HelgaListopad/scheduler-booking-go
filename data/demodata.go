package data

import (
	"encoding/json"
	"os"

	"gorm.io/gorm"
)

func dataDown(tx *gorm.DB) {
	must(tx.Exec("DELETE FROM `doctors`").Error)
	must(tx.Exec("DELETE FROM `doctors_routine`").Error)
	must(tx.Exec("DELETE FROM `occupied_slots`").Error)
}

func dataUp(tx *gorm.DB) {
	doctors := make([]Doctor, 0)
	must(parseJSON(&doctors, "./demodata/doctors.json"))
	must(tx.Create(&doctors).Error)

	doctors_routine := make([]DoctorRoutine, 0)
	must(parseJSON(&doctors_routine, "./demodata/doctors_routine.json"))
	must(tx.Create(&doctors_routine).Error)

	occupied_slots := make([]OccupiedSlot, 0)
	must(parseJSON(&occupied_slots, "./demodata/occupied_slots.json"))
	must(tx.Create(&occupied_slots).Error)
}

func parseJSON(dest any, path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &dest)
}
