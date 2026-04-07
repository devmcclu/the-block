package database

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	ExternalID        string         `json:"external_id" gorm:"uniqueIndex;not null"`
	VIN               string         `json:"vin"`
	Year              int            `json:"year"`
	Make              string         `json:"make"`
	VehicleModel      string         `json:"model" gorm:"column:model"`
	Trim              string         `json:"trim"`
	BodyStyle         string         `json:"body_style"`
	ExteriorColor     string         `json:"exterior_color"`
	InteriorColor     string         `json:"interior_color"`
	Engine            string         `json:"engine"`
	Transmission      string         `json:"transmission"`
	Drivetrain        string         `json:"drivetrain"`
	OdometerKM        int            `json:"odometer_km"`
	FuelType          string         `json:"fuel_type"`
	ConditionGrade    float64        `json:"condition_grade"`
	ConditionReport   string         `json:"condition_report"`
	TitleStatus       string         `json:"title_status"`
	Province          string         `json:"province"`
	City              string         `json:"city"`
	AuctionStart      string         `json:"auction_start"`
	StartingBid       int            `json:"starting_bid"`
	ReservePrice      *int           `json:"reserve_price"`
	BuyNowPrice       *int           `json:"buy_now_price"`
	SellingDealership string         `json:"selling_dealership"`
	Lot               string         `json:"lot"`
	CurrentBid        int            `json:"current_bid"`
	BidCount          int            `json:"bid_count"`
	DamageNotes       []DamageNote   `json:"damage_notes" gorm:"foreignKey:VehicleID"`
	Images            []VehicleImage `json:"images" gorm:"foreignKey:VehicleID"`
}

type DamageNote struct {
	gorm.Model
	VehicleID uint   `json:"vehicle_id"`
	Note      string `json:"note"`
}

type VehicleImage struct {
	gorm.Model
	VehicleID uint   `json:"vehicle_id"`
	URL       string `json:"url"`
}

type VehicleCreate struct {
	VIN               string   `json:"vin"`
	Year              int      `json:"year"`
	Make              string   `json:"make"`
	Model             string   `json:"model"`
	Trim              string   `json:"trim"`
	BodyStyle         string   `json:"body_style"`
	ExteriorColor     string   `json:"exterior_color"`
	InteriorColor     string   `json:"interior_color"`
	Engine            string   `json:"engine"`
	Transmission      string   `json:"transmission"`
	Drivetrain        string   `json:"drivetrain"`
	OdometerKM        int      `json:"odometer_km"`
	FuelType          string   `json:"fuel_type"`
	ConditionGrade    float64  `json:"condition_grade"`
	ConditionReport   string   `json:"condition_report"`
	TitleStatus       string   `json:"title_status"`
	Province          string   `json:"province"`
	City              string   `json:"city"`
	AuctionStart      string   `json:"auction_start"`
	StartingBid       int      `json:"starting_bid"`
	ReservePrice      *int     `json:"reserve_price"`
	BuyNowPrice       *int     `json:"buy_now_price"`
	SellingDealership string   `json:"selling_dealership"`
	Lot               string   `json:"lot"`
	DamageNotes       []string `json:"damage_notes"`
	Images            []string `json:"images"`
}

type VehicleUpdate struct {
	CurrentBid *int `json:"current_bid"`
	BidCount   *int `json:"bid_count"`
}

type VehicleFilter struct {
	YearMin         *int
	YearMax         *int
	Makes           []string
	Models          []string
	BodyStyles      []string
	ExteriorColors  []string
	InteriorColors  []string
	Transmissions   []string
	Drivetrains     []string
	FuelTypes       []string
	TitleStatuses   []string
	OdometerMin     *int
	OdometerMax     *int
	ConditionMin    *float64
	ConditionMax    *float64
	Sort            string
}

type VehicleFilterOptions struct {
	YearMin        int      `json:"year_min"`
	YearMax        int      `json:"year_max"`
	Makes          []string `json:"makes"`
	Models         []string `json:"models"`
	BodyStyles     []string `json:"body_styles"`
	ExteriorColors []string `json:"exterior_colors"`
	InteriorColors []string `json:"interior_colors"`
	Transmissions  []string `json:"transmissions"`
	Drivetrains    []string `json:"drivetrains"`
	FuelTypes      []string `json:"fuel_types"`
	TitleStatuses  []string `json:"title_statuses"`
	OdometerMin    int      `json:"odometer_min"`
	OdometerMax    int      `json:"odometer_max"`
	ConditionMin   float64  `json:"condition_min"`
	ConditionMax   float64  `json:"condition_max"`
}
