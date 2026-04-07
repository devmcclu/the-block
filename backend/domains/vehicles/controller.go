package vehicles

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/devmcclu/the-block/backend/database"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"gorm.io/gorm"
)

// maxAuctionDurationHours is the fixed duration for all auctions (30 days).
const maxAuctionDurationHours = 720

// minBidIncrement is the minimum amount a new bid must exceed the current bid by.
const minBidIncrement = 100

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

	bidsGroup := fuego.Group(s, "/bids")
	fuego.Get(bidsGroup, "/", rs.getAllBids)

	fuego.Get(vehiclesGroup, "/{id}", rs.getVehicle)
	fuego.Put(vehiclesGroup, "/{id}", rs.putVehicle)
	fuego.Post(vehiclesGroup, "/{id}/buy", rs.buyNowVehicle)
	fuego.Delete(vehiclesGroup, "/{id}", rs.deleteVehicle)
}

// mapServiceErr translates domain/DB errors into fuego HTTP errors.
func mapServiceErr(err error, context string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fuego.NotFoundError{Detail: fmt.Sprintf("%s: not found", context)}
	}
	if errors.Is(err, ErrBidTooLow) || errors.Is(err, ErrAuctionEnded) || errors.Is(err, ErrNoBuyNow) {
		return fuego.ConflictError{Detail: err.Error()}
	}
	return fuego.HTTPError{
		Err:    fmt.Errorf("%s: %w", context, err),
		Detail: "an internal error occurred",
	}
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
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fuego.BadRequestError{Detail: "condition_min must be a finite number"}
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
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fuego.BadRequestError{Detail: "condition_max must be a finite number"}
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

	vehicles, err := rs.VehiclesService.GetAllVehicles(filters)
	if err != nil {
		return nil, mapServiceErr(err, "list vehicles")
	}
	return vehicles, nil
}

func (rs VehiclesResources) getConfig(c fuego.ContextNoBody) (database.AuctionConfig, error) {
	return database.AuctionConfig{
		MaxAuctionDurationHours: maxAuctionDurationHours,
		MinBidIncrement:         minBidIncrement,
	}, nil
}

func (rs VehiclesResources) getFilterOptions(c fuego.ContextNoBody) (database.VehicleFilterOptions, error) {
	opts, err := rs.VehiclesService.GetFilterOptions()
	if err != nil {
		return database.VehicleFilterOptions{}, mapServiceErr(err, "get filter options")
	}
	return opts, nil
}

func (rs VehiclesResources) postVehicle(c fuego.ContextWithBody[database.VehicleCreate]) (database.Vehicle, error) {
	body, err := c.Body()
	if err != nil {
		return database.Vehicle{}, err
	}

	vehicle, err := rs.VehiclesService.CreateVehicle(body)
	if err != nil {
		return database.Vehicle{}, mapServiceErr(err, "create vehicle")
	}
	return vehicle, nil
}

func (rs VehiclesResources) getVehicle(c fuego.ContextNoBody) (database.Vehicle, error) {
	id := c.PathParam("id")

	vehicle, err := rs.VehiclesService.GetVehicle(id)
	if err != nil {
		return database.Vehicle{}, mapServiceErr(err, fmt.Sprintf("get vehicle %s", id))
	}
	return vehicle, nil
}

func (rs VehiclesResources) putVehicle(c fuego.ContextWithBody[database.VehicleUpdate]) (database.Vehicle, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return database.Vehicle{}, err
	}

	vehicle, err := rs.VehiclesService.UpdateVehicle(id, body)
	if err != nil {
		return database.Vehicle{}, mapServiceErr(err, fmt.Sprintf("update vehicle %s", id))
	}
	return vehicle, nil
}

func (rs VehiclesResources) deleteVehicle(c fuego.ContextNoBody) (any, error) {
	result, err := rs.VehiclesService.DeleteVehicle(c.PathParam("id"))
	if err != nil {
		return nil, mapServiceErr(err, fmt.Sprintf("delete vehicle %s", c.PathParam("id")))
	}
	return result, nil
}

func (rs VehiclesResources) buyNowVehicle(c fuego.ContextNoBody) (database.Vehicle, error) {
	id := c.PathParam("id")

	vehicle, err := rs.VehiclesService.BuyNow(id)
	if err != nil {
		return database.Vehicle{}, mapServiceErr(err, fmt.Sprintf("buy vehicle %s", id))
	}
	return vehicle, nil
}

func (rs VehiclesResources) getAllBids(c fuego.ContextNoBody) ([]database.Bid, error) {
	bids, err := rs.VehiclesService.GetAllBids()
	if err != nil {
		return nil, mapServiceErr(err, "list bids")
	}
	return bids, nil
}
