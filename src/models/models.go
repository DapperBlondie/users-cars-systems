package models

import "time"

// Users holding users data
type Users struct {
	ID           int       `json:"id"`
	CompleteName string    `json:"complete_name"`
	Sex          bool      `json:"sex"`
	BirthDay     time.Time `json:"birth_day"`
	Password     string    `json:"password"`
	UsersCars    []*Cars   `json:"-"`
}

// Cars holding cars data
type Cars struct {
	ID          int    `json:"id"`
	NumberPlate string `json:"number_plate"`
	Color       string `json:"color"`
	VIN         string `json:"vin"`
	OwnerID     int    `json:"owner_id"`
}
