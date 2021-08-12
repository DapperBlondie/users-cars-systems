package models

type StatusIdentifier struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// Users holding users data
type Users struct {
	ID           int     `json:"-"`
	CompleteName string  `json:"complete_name"`
	Sex          bool    `json:"sex"`
	BirthDay     string  `json:"birth_day"`
	Password     string  `json:"password"`
	UsersCars    []*Cars `json:"-"`
}

// Cars holding cars data
type Cars struct {
	ID          int    `json:"-"`
	NumberPlate string `json:"number_plate"`
	Color       string `json:"color"`
	VIN         string `json:"vin"`
	OwnerID     int    `json:"owner_id"`
}
