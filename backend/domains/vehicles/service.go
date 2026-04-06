package vehicles

import (
	"github.com/devmcclu/the-block/backend/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehiclesService interface {
	GetVehicle(id string) (database.Vehicle, error)
	GetAllVehicles() ([]database.Vehicle, error)
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

func (s RealVehiclesService) GetAllVehicles() ([]database.Vehicle, error) {
	var vehicles []database.Vehicle
	err := s.DB.Preload("DamageNotes").Preload("Images").Find(&vehicles).Error
	return vehicles, err
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
