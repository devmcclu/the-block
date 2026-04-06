package vehicles

import (
	"github.com/devmcclu/the-block/backend/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehiclesService interface {
	GetVehicle(id string) (database.Vehicle, error)
	GetAllVehicles(filters database.VehicleFilter) ([]database.Vehicle, error)
	GetFilterOptions() (database.VehicleFilterOptions, error)
	CreateVehicle(database.VehicleCreate) (database.Vehicle, error)
	UpdateVehicle(id string, input database.VehicleUpdate) (database.Vehicle, error)
	DeleteVehicle(id string) (any, error)
}

type RealVehiclesService struct {
	DB *gorm.DB
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

	err := q.Find(&vehicles).Error
	return vehicles, err
}

func (s RealVehiclesService) GetFilterOptions() (database.VehicleFilterOptions, error) {
	var opts database.VehicleFilterOptions

	s.DB.Model(&database.Vehicle{}).Distinct("make").Order("make").Pluck("make", &opts.Makes)
	s.DB.Model(&database.Vehicle{}).Distinct("model").Order("model").Pluck("model", &opts.Models)
	s.DB.Model(&database.Vehicle{}).Distinct("body_style").Order("body_style").Pluck("body_style", &opts.BodyStyles)
	s.DB.Model(&database.Vehicle{}).Distinct("exterior_color").Order("exterior_color").Pluck("exterior_color", &opts.ExteriorColors)
	s.DB.Model(&database.Vehicle{}).Distinct("interior_color").Order("interior_color").Pluck("interior_color", &opts.InteriorColors)
	s.DB.Model(&database.Vehicle{}).Distinct("transmission").Order("transmission").Pluck("transmission", &opts.Transmissions)
	s.DB.Model(&database.Vehicle{}).Distinct("drivetrain").Order("drivetrain").Pluck("drivetrain", &opts.Drivetrains)
	s.DB.Model(&database.Vehicle{}).Distinct("fuel_type").Order("fuel_type").Pluck("fuel_type", &opts.FuelTypes)
	s.DB.Model(&database.Vehicle{}).Distinct("title_status").Order("title_status").Pluck("title_status", &opts.TitleStatuses)

	s.DB.Model(&database.Vehicle{}).Select("MIN(year)").Scan(&opts.YearMin)
	s.DB.Model(&database.Vehicle{}).Select("MAX(year)").Scan(&opts.YearMax)
	s.DB.Model(&database.Vehicle{}).Select("MIN(odometer_km)").Scan(&opts.OdometerMin)
	s.DB.Model(&database.Vehicle{}).Select("MAX(odometer_km)").Scan(&opts.OdometerMax)
	s.DB.Model(&database.Vehicle{}).Select("MIN(condition_grade)").Scan(&opts.ConditionMin)
	s.DB.Model(&database.Vehicle{}).Select("MAX(condition_grade)").Scan(&opts.ConditionMax)

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
	var vehicle database.Vehicle
	if err := s.DB.Where("external_id = ?", id).First(&vehicle).Error; err != nil {
		return database.Vehicle{}, err
	}

	updates := map[string]any{}
	if input.CurrentBid != nil {
		updates["current_bid"] = *input.CurrentBid
	}
	if input.BidCount != nil {
		updates["bid_count"] = *input.BidCount
	}

	if len(updates) > 0 {
		if err := s.DB.Model(&vehicle).Updates(updates).Error; err != nil {
			return database.Vehicle{}, err
		}
	}

	return s.GetVehicle(id)
}

func (s RealVehiclesService) DeleteVehicle(id string) (any, error) {
	var vehicle database.Vehicle
	if err := s.DB.Where("external_id = ?", id).First(&vehicle).Error; err != nil {
		return nil, err
	}
	return nil, s.DB.Delete(&vehicle).Error
}
