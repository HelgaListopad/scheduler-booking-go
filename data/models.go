package data

type Doctor struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Subtitle string  `json:"subtitle"`
	Details  string  `json:"details"`
	Category string  `json:"category"`
	Price    float32 `json:"price"`
	Gap      int     `json:"gap"`
	SlotSize int     `json:"slot_size"`
	ImageURL string  `json:"image_url"`

	DoctorSchedule []DoctorSchedule `json:"-"`
	OccupiedSlots  []OccupiedSlot   `json:"-"`
	Review         Review           `json:"-" gorm:"foreignkey:DoctorID"`
}

type Review struct {
	ID       int `json:"-"`
	Count    int `json:"count"`
	Star     int `json:"star"`
	DoctorID int `json:"-"`
}

type DoctorSchedule struct {
	ID       int
	DoctorID int
	From     int
	To       int

	DoctorRoutine          []DoctorRoutine          `gorm:"foreignkey:ScheduleID"`
	DoctorRecurringRoutine []DoctorRecurringRoutine `gorm:"foreignkey:ScheduleID"`
}

type DoctorRecurringRoutine struct {
	ID         int
	ScheduleID int
	WeekDay    int

	DoctorSchedule DoctorSchedule `gorm:"foreignkey:ScheduleID"`
}

type DoctorRoutine struct {
	ID         int
	ScheduleID int
	Date       int64

	DoctorSchedule DoctorSchedule `gorm:"foreignkey:ScheduleID"`
}

type OccupiedSlot struct {
	ID            int    `json:"id"`
	DoctorID      int    `json:"doctor_id"`
	Date          int64  `json:"date"`
	ClientName    string `json:"client_name"`
	ClientEmail   string `json:"client_email"`
	ClientDetails string `json:"client_details"`
}

func (DoctorRoutine) TableName() string {
	return "doctor_routine"
}

func (DoctorRecurringRoutine) TableName() string {
	return "doctor_recurring"
}

func (DoctorSchedule) TableName() string {
	return "doctor_schedule"
}

func (OccupiedSlot) TableName() string {
	return "occupied_slots"
}
