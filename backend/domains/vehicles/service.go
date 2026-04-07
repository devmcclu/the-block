package vehicles

import (
	"errors"
	"fmt"
	"time"

	"github.com/devmcclu/the-block/backend/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrBidTooLow is returned when a bid does not exceed the current highest bid.
var ErrBidTooLow = errors.New("bid not higher than current bid")

// ErrAuctionEnded is returned when a bid is placed after the auction has ended.
var ErrAuctionEnded = errors.New("auction has ended")

type VehiclesService interface {
	GetVehicle(id string) (database.Vehicle, error)
	GetAllVehicles(filters database.VehicleFilter) ([]database.Vehicle, error)
	GetFilterOptions() (database.VehicleFilterOptions, error)
	CreateVehicle(database.VehicleCreate) (database.Vehicle, error)
	UpdateVehicle(id string, input database.VehicleUpdate) (database.Vehicle, error)
	DeleteVehicle(id string) (any, error)
}

type RealVehiclesService struct {
	DB                      *gorm.DB
	MaxAuctionDurationHours int
	MinBidIncrement         int
}

func (s RealVehiclesService) GetVehicle(id string) (database.Vehicle, error) {
	var vehicle database.Vehicle
	err := s.DB.Preload("DamageNotes").Preload("Images").
		Where("external_id = ?", id).First(&vehicle).Error
	return vehicle, err
}

func (s RealVehiclesService) GetAllVehicles(filters database.VehicleFilter) ([]database.Vehicle, error) {
	var vehicles []database.Vehicle
	q := s.DB.Preload("DamageNotes").Preload("Images")

	if filters.YearMin != nil {
		q = q.Where("year >= ?", *filters.YearMin)
	}
	if filters.YearMax != nil {
		q = q.Where("year <= ?", *filters.YearMax)
	}
	if len(filters.Makes) > 0 {
		q = q.Where("make IN ?", filters.Makes)
	}
	if len(filters.Models) > 0 {
		q = q.Where("model IN ?", filters.Models)
	}
	if len(filters.BodyStyles) > 0 {
		q = q.Where("body_style IN ?", filters.BodyStyles)
	}
	if len(filters.ExteriorColors) > 0 {
		q = q.Where("exterior_color IN ?", filters.ExteriorColors)
	}
	if len(filters.InteriorColors) > 0 {
		q = q.Where("interior_color IN ?", filters.InteriorColors)
	}
	if len(filters.Transmissions) > 0 {
		q = q.Where("transmission IN ?", filters.Transmissions)
	}
	if len(filters.Drivetrains) > 0 {
		q = q.Where("drivetrain IN ?", filters.Drivetrains)
	}
	if len(filters.FuelTypes) > 0 {
		q = q.Where("fuel_type IN ?", filters.FuelTypes)
	}
	if len(filters.TitleStatuses) > 0 {
		q = q.Where("title_status IN ?", filters.TitleStatuses)
	}
	if filters.OdometerMin != nil {
		q = q.Where("odometer_km >= ?", *filters.OdometerMin)
	}
	if filters.OdometerMax != nil {
		q = q.Where("odometer_km <= ?", *filters.OdometerMax)
	}
	if filters.ConditionMin != nil {
		q = q.Where("condition_grade >= ?", *filters.ConditionMin)
	}
	if filters.ConditionMax != nil {
		q = q.Where("condition_grade <= ?", *filters.ConditionMax)
	}

	switch filters.Sort {
	case "price_asc":
		q = q.Order("current_bid ASC")
	case "price_desc":
		q = q.Order("current_bid DESC")
	case "year_desc":
		q = q.Order("year DESC")
	case "year_asc":
		q = q.Order("year ASC")
	case "bids_desc":
		q = q.Order("bid_count DESC")
	case "bids_asc":
		q = q.Order("bid_count ASC")
	// All auctions share the same maxAuctionDurationHours, so earlier start
	// time implies earlier end time. Sorting by auction_start is equivalent
	// to sorting by auction end.
	case "ending_soon":
		q = q.Order("auction_start ASC")
	case "ending_last":
		q = q.Order("auction_start DESC")
	}

	err := q.Find(&vehicles).Error
	return vehicles, err
}

func (s RealVehiclesService) GetFilterOptions() (database.VehicleFilterOptions, error) {
	var opts database.VehicleFilterOptions

	if err := s.DB.Model(&database.Vehicle{}).Distinct("make").Order("make").Pluck("make", &opts.Makes).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("model").Order("model").Pluck("model", &opts.Models).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("body_style").Order("body_style").Pluck("body_style", &opts.BodyStyles).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("exterior_color").Order("exterior_color").Pluck("exterior_color", &opts.ExteriorColors).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("interior_color").Order("interior_color").Pluck("interior_color", &opts.InteriorColors).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("transmission").Order("transmission").Pluck("transmission", &opts.Transmissions).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("drivetrain").Order("drivetrain").Pluck("drivetrain", &opts.Drivetrains).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("fuel_type").Order("fuel_type").Pluck("fuel_type", &opts.FuelTypes).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Distinct("title_status").Order("title_status").Pluck("title_status", &opts.TitleStatuses).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MIN(year)").Scan(&opts.YearMin).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MAX(year)").Scan(&opts.YearMax).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MIN(odometer_km)").Scan(&opts.OdometerMin).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MAX(odometer_km)").Scan(&opts.OdometerMax).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MIN(condition_grade)").Scan(&opts.ConditionMin).Error; err != nil {
		return opts, err
	}
	if err := s.DB.Model(&database.Vehicle{}).Select("MAX(condition_grade)").Scan(&opts.ConditionMax).Error; err != nil {
		return opts, err
	}

	return opts, nil
}

func (s RealVehiclesService) CreateVehicle(input database.VehicleCreate) (database.Vehicle, error) {
	vehicle := database.Vehicle{
		ExternalID:        uuid.New().String(),
		VIN:               input.VIN,
		Year:              input.Year,
		Make:              input.Make,
		VehicleModel:      input.Model,
		Trim:              input.Trim,
		BodyStyle:         input.BodyStyle,
		ExteriorColor:     input.ExteriorColor,
		InteriorColor:     input.InteriorColor,
		Engine:            input.Engine,
		Transmission:      input.Transmission,
		Drivetrain:        input.Drivetrain,
		OdometerKM:        input.OdometerKM,
		FuelType:          input.FuelType,
		ConditionGrade:    input.ConditionGrade,
		ConditionReport:   input.ConditionReport,
		TitleStatus:       input.TitleStatus,
		Province:          input.Province,
		City:              input.City,
		AuctionStart:      input.AuctionStart,
		StartingBid:       input.StartingBid,
		ReservePrice:      input.ReservePrice,
		BuyNowPrice:       input.BuyNowPrice,
		SellingDealership: input.SellingDealership,
		Lot:               input.Lot,
	}

	for _, note := range input.DamageNotes {
		vehicle.DamageNotes = append(vehicle.DamageNotes, database.DamageNote{Note: note})
	}
	for _, url := range input.Images {
		vehicle.Images = append(vehicle.Images, database.VehicleImage{URL: url})
	}

	err := s.DB.Create(&vehicle).Error
	return vehicle, err
}

func (s RealVehiclesService) UpdateVehicle(id string, input database.VehicleUpdate) (database.Vehicle, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var vehicle database.Vehicle
		if err := tx.Where("external_id = ?", id).First(&vehicle).Error; err != nil {
			return err
		}

		auctionStart, err := time.ParseInLocation("2006-01-02T15:04:05", vehicle.AuctionStart, time.UTC)
		if err != nil {
			return fmt.Errorf("invalid auction_start: %w", err)
		}
		auctionEnd := auctionStart.Add(time.Duration(s.MaxAuctionDurationHours) * time.Hour)
		if time.Now().UTC().After(auctionEnd) {
			return fmt.Errorf("%w: auction ended at %s", ErrAuctionEnded, auctionEnd.Format(time.RFC3339))
		}

		if input.BidAmount == nil || *input.BidAmount <= 0 {
			return fmt.Errorf("%w: bid amount must be greater than 0", ErrBidTooLow)
		}

		bidAmount := *input.BidAmount
		minBid := vehicle.StartingBid
		if vehicle.BidCount > 0 {
			minBid = vehicle.CurrentBid + s.MinBidIncrement
		}
		if bidAmount < minBid {
			return fmt.Errorf("%w: bid of %d is below the minimum of %d", ErrBidTooLow, bidAmount, minBid)
		}

		result := tx.Model(&vehicle).Where("current_bid = ?", vehicle.CurrentBid).Updates(map[string]any{
			"current_bid": bidAmount,
			"bid_count":   gorm.Expr("bid_count + 1"),
		})
		if result.RowsAffected == 0 && result.Error == nil {
			return fmt.Errorf("%w: another bid was placed concurrently, please retry", ErrBidTooLow)
		}
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return database.Vehicle{}, err
	}

	return s.GetVehicle(id)
}

func (s RealVehiclesService) DeleteVehicle(id string) (any, error) {
	var vehicle database.Vehicle
	if err := s.DB.Where("external_id = ?", id).First(&vehicle).Error; err != nil {
		return nil, err
	}
	return nil, s.DB.Select("DamageNotes", "Images").Delete(&vehicle).Error
}
