package main

import (
	"log"

	"github.com/devmcclu/the-block/backend/database"
	"github.com/devmcclu/the-block/backend/domains/vehicles"
	"github.com/go-fuego/fuego"
)

func main() {
	db, err := database.InitDB("vehicles.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	if err := database.SeedIfEmpty(db, "vehicles.json"); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	s := fuego.NewServer()

	vehicleResources := vehicles.VehiclesResources{
		VehiclesService: vehicles.RealVehiclesService{
			DB:                      db,
			MaxAuctionDurationHours: 720,
			MinBidIncrement:         100,
		},
	}
	vehicleResources.Routes(s)

	if err := s.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
