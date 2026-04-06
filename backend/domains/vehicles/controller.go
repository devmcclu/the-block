package vehicles

import (
	"github.com/devmcclu/the-block/backend/database"
	"github.com/go-fuego/fuego"
)

type VehiclesResources struct {
	VehiclesService VehiclesService
}

func (rs VehiclesResources) Routes(s *fuego.Server) {
	vehiclesGroup := fuego.Group(s, "/vehicles")

	fuego.Get(vehiclesGroup, "/", rs.getAllVehicles)
	fuego.Post(vehiclesGroup, "/", rs.postVehicle)

	fuego.Get(vehiclesGroup, "/{id}", rs.getVehicle)
	fuego.Put(vehiclesGroup, "/{id}", rs.putVehicle)
	fuego.Delete(vehiclesGroup, "/{id}", rs.deleteVehicle)
}

func (rs VehiclesResources) getAllVehicles(c fuego.ContextNoBody) ([]database.Vehicle, error) {
	return rs.VehiclesService.GetAllVehicles()
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
