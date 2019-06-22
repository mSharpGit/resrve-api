package models

type floors struct {
	ID            int    `json:"id"`
	RestaurantID  int    `json:"restaurant_id"`
	FloorNumber   int    `json:"floor_number"`
	FloorName     string `json:"floor_name"`
	BackgroundPic string `json:"background_pic"`
	AddDate       int    `json:"add_date"`
}
