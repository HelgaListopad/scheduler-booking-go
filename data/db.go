package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Path         string
	ResetOnStart bool
}

type DAO struct {
	db             *gorm.DB
	Doctors        *doctorsDAO
	DoctorsRoutine *doctorsRoutineDAO
	OccupiedSlots  *occupiedSlotsDAO
}

func NewDAO(config DBConfig) *DAO {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Doctor{})
	db.AutoMigrate(&DoctorRoutine{})
	db.AutoMigrate(&OccupiedSlot{})

	dao := DAO{db: db}
	dao.Doctors = newDoctorsDAO(db)
	dao.DoctorsRoutine = newDoctorsRoutineDAO(db)
	dao.OccupiedSlots = newOccupiedSlotsDAO(db)

	if config.ResetOnStart {
		db.Transaction(func(tx *gorm.DB) error {
			dataDown(tx)
			dataUp(tx)
			return nil
		})
	}

	return &dao
}

func (d *DAO) GetDB() *gorm.DB {
	return d.db
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
