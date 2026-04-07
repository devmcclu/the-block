package vehicles

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/devmcclu/the-block/backend/database"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// maxAuctionDurationHours is the fixed duration for all auctions (30 days).
const maxAuctionDurationHours = 720

var validSorts = map[string]bool{
	"":            true,
	"price_asc":   true,
	"price_desc":  true,
	"year_asc":    true,
	"year_desc":   true,
	"bids_asc":    true,
	"bids_desc":   true,
	"ending_soon": true,
	"ending_last": true,
}

type VehiclesResources struct {
	VehiclesService VehiclesService
}

func (rs VehiclesResources) Routes(s *fuego.Server) {
	vehiclesGroup := fuego.Group(s, "/vehicles")

	fuego.Get(vehiclesGroup, "/", rs.getAllVehicles,
		option.QueryInt("year_min", "Minimum year"),
		option.QueryInt("year_max", "Maximum year"),
		option.QueryArray("make", "Filter by make", reflect.String),
		option.QueryArray("model", "Filter by model", reflect.String),
		option.QueryArray("body_style", "Filter by body style", reflect.String),
		option.QueryArray("exterior_color", "Filter by exterior color", reflect.String),
		option.QueryArray("interior_color", "Filter by interior color", reflect.String),
		option.QueryArray("transmission", "Filter by transmission", reflect.String),
		option.QueryArray("drivetrain", "Filter by drivetrain", reflect.String),
		option.QueryArray("fuel_type", "Filter by fuel type", reflect.String),
		option.QueryArray("title_status", "Filter by title status", reflect.String),
		option.QueryInt("odometer_min", "Minimum odometer (km)"),
		option.QueryInt("odometer_max", "Maximum odometer (km)"),
		option.Query("condition_min", "Minimum condition grade"),
		option.Query("condition_max", "Maximum condition grade"),
		option.Query("sort", "Sort order: price_asc, price_desc, year_asc, year_desc, bids_asc, bids_desc, ending_soon, ending_last"),
	)
	fuego.Post(vehiclesGroup, "/", rs.postVehicle)

	fuego.Get(vehiclesGroup, "/config", rs.getConfig)
	fuego.Get(vehiclesGroup, "/filters", rs.getFilterOptions)

	fuego.Get(vehiclesGroup, "/{id}", rs.getVehicle)
	fuego.Put(vehiclesGroup, "/{id}", rs.putVehicle)
	fuego.Delete(vehiclesGroup, "/{id}", rs.deleteVehicle)
}

func (rs VehiclesResources) getAllVehicles(c fuego.ContextNoBody) ([]database.Vehicle, error) {
	filters := database.VehicleFilter{
		Makes:          c.QueryParamArr("make"),
		Models:         c.QueryParamArr("model"),
		BodyStyles:     c.QueryParamArr("body_style"),
		ExteriorColors: c.QueryParamArr("exterior_color"),
		InteriorColors: c.QueryParamArr("interior_color"),
		Transmissions:  c.QueryParamArr("transmission"),
		Drivetrains:    c.QueryParamArr("drivetrain"),
		FuelTypes:      c.QueryParamArr("fuel_type"),
		TitleStatuses:  c.QueryParamArr("title_status"),
	}

	if c.QueryParam("year_min") != "" {
		if v, err := c.QueryParamIntErr("year_min"); err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid year_min: %s", err)}
		} else {
			if v < 1900 || v > 2100 {
				return nil, fuego.BadRequestError{Detail: "year_min must be between 1900 and 2100"}
			}
			filters.YearMin = &v
		}
	}
	if c.QueryParam("year_max") != "" {
		if v, err := c.QueryParamIntErr("year_max"); err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid year_max: %s", err)}
		} else {
			if v < 1900 || v > 2100 {
				return nil, fuego.BadRequestError{Detail: "year_max must be between 1900 and 2100"}
			}
			filters.YearMax = &v
		}
	}
	if filters.YearMin != nil && filters.YearMax != nil && *filters.YearMin > *filters.YearMax {
		return nil, fuego.BadRequestError{Detail: "year_min must not exceed year_max"}
	}

	if c.QueryParam("odometer_min") != "" {
		if v, err := c.QueryParamIntErr("odometer_min"); err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid odometer_min: %s", err)}
		} else {
			if v < 0 {
				return nil, fuego.BadRequestError{Detail: "odometer_min must not be negative"}
			}
			filters.OdometerMin = &v
		}
	}
	if c.QueryParam("odometer_max") != "" {
		if v, err := c.QueryParamIntErr("odometer_max"); err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid odometer_max: %s", err)}
		} else {
			if v < 0 {
				return nil, fuego.BadRequestError{Detail: "odometer_max must not be negative"}
			}
			filters.OdometerMax = &v
		}
	}
	if filters.OdometerMin != nil && filters.OdometerMax != nil && *filters.OdometerMin > *filters.OdometerMax {
		return nil, fuego.BadRequestError{Detail: "odometer_min must not exceed odometer_max"}
	}

	if s := c.QueryParam("condition_min"); s != "" {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid condition_min: %s", s)}
		}
		if v < 0 || v > 5 {
			return nil, fuego.BadRequestError{Detail: "condition_min must be between 0 and 5"}
		}
		filters.ConditionMin = &v
	}
	if s := c.QueryParam("condition_max"); s != "" {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid condition_max: %s", s)}
		}
		if v < 0 || v > 5 {
			return nil, fuego.BadRequestError{Detail: "condition_max must be between 0 and 5"}
		}
		filters.ConditionMax = &v
	}
	if filters.ConditionMin != nil && filters.ConditionMax != nil && *filters.ConditionMin > *filters.ConditionMax {
		return nil, fuego.BadRequestError{Detail: "condition_min must not exceed condition_max"}
	}

	filters.Sort = c.QueryParam("sort")
	if !validSorts[filters.Sort] {
		return nil, fuego.BadRequestError{Detail: fmt.Sprintf("invalid sort value: %s", filters.Sort)}
	}

	return rs.VehiclesService.GetAllVehicles(filters)
}

func (rs VehiclesResources) getConfig(c fuego.ContextNoBody) (database.AuctionConfig, error) {
	return database.AuctionConfig{MaxAuctionDurationHours: maxAuctionDurationHours}, nil
}

func (rs VehiclesResources) getFilterOptions(c fuego.ContextNoBody) (database.VehicleFilterOptions, error) {
	return rs.VehiclesService.GetFilterOptions()
}

func (rs VehiclesResources) postVehicle(c fuego.ContextWithBody[database.VehicleCreate]) (database.Vehicle, error) {
	body, err := c.Body()
	if err != nil {
		return database.Vehicle{}, err
	}

	return rs.VehiclesService.CreateVehicle(body)
}

func (rs VehiclesResources) getVehicle(c fuego.ContextNoBody) (database.Vehicle, error) {
	id := c.PathParam("id")

	return rs.VehiclesService.GetVehicle(id)
}

func (rs VehiclesResources) putVehicle(c fuego.ContextWithBody[database.VehicleUpdate]) (database.Vehicle, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return database.Vehicle{}, err
	}

	return rs.VehiclesService.UpdateVehicle(id, body)
}

func (rs VehiclesResources) deleteVehicle(c fuego.ContextNoBody) (any, error) {
	return rs.VehiclesService.DeleteVehicle(c.PathParam("id"))
}
