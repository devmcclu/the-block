package database

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type jsonVehicle struct {
	ID                string   `json:"id"`
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
	CurrentBid        int      `json:"current_bid"`
	BidCount          int      `json:"bid_count"`
	DamageNotes       []string `json:"damage_notes"`
	Images            []string `json:"images"`
}

func SeedIfEmpty(db *gorm.DB, jsonPath string) error {
	var count int64
	if err := db.Model(&Vehicle{}).Count(&count).Error; err != nil {
		return fmt.Errorf("counting vehicles: %w", err)
	}
	if count > 0 {
		return nil
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("reading seed file: %w", err)
	}

	var raw []jsonVehicle
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshaling seed data: %w", err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, jv := range raw {
			vehicle := Vehicle{
				ExternalID:        jv.ID,
				VIN:               jv.VIN,
				Year:              jv.Year,
				Make:              jv.Make,
				VehicleModel:      jv.Model,
				Trim:              jv.Trim,
				BodyStyle:         jv.BodyStyle,
				ExteriorColor:     jv.ExteriorColor,
				InteriorColor:     jv.InteriorColor,
				Engine:            jv.Engine,
				Transmission:      jv.Transmission,
				Drivetrain:        jv.Drivetrain,
				OdometerKM:        jv.OdometerKM,
				FuelType:          jv.FuelType,
				ConditionGrade:    jv.ConditionGrade,
				ConditionReport:   jv.ConditionReport,
				TitleStatus:       jv.TitleStatus,
				Province:          jv.Province,
				City:              jv.City,
				AuctionStart:      jv.AuctionStart,
				StartingBid:       jv.StartingBid,
				ReservePrice:      jv.ReservePrice,
				BuyNowPrice:       jv.BuyNowPrice,
				SellingDealership: jv.SellingDealership,
				Lot:               jv.Lot,
				CurrentBid:        jv.CurrentBid,
				BidCount:          jv.BidCount,
			}

			for _, note := range jv.DamageNotes {
				vehicle.DamageNotes = append(vehicle.DamageNotes, DamageNote{Note: note})
			}
			for _, url := range jv.Images {
				vehicle.Images = append(vehicle.Images, VehicleImage{URL: url})
			}

			if err := tx.Create(&vehicle).Error; err != nil {
				return fmt.Errorf("inserting vehicle %s: %w", jv.ID, err)
			}
		}

		fmt.Printf("Seeded %d vehicles\n", len(raw))
		return nil
	})
}
