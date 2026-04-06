package cars

import (
	"github.com/go-fuego/fuego"
)

type CarsResources struct {
	// TODO add resources
	CarsService CarsService
}

func (rs CarsResources) Routes(s *fuego.Server) {
	carsGroup := fuego.Group(s, "/cars")

	fuego.Get(carsGroup, "/", rs.getAllCars)
	fuego.Post(carsGroup, "/", rs.postCars)

	fuego.Get(carsGroup, "/{id}", rs.getCars)
	fuego.Put(carsGroup, "/{id}", rs.putCars)
	fuego.Delete(carsGroup, "/{id}", rs.deleteCars)
}

func (rs CarsResources) getAllCars(c fuego.ContextNoBody) ([]Cars, error) {
	return rs.CarsService.GetAllCars()
}

func (rs CarsResources) postCars(c fuego.ContextWithBody[CarsCreate]) (Cars, error) {
	body, err := c.Body()
	if err != nil {
		return Cars{}, err
	}

	return rs.CarsService.CreateCars(body)
}

func (rs CarsResources) getCars(c fuego.ContextNoBody) (Cars, error) {
	id := c.PathParam("id")

	return rs.CarsService.GetCars(id)
}

func (rs CarsResources) putCars(c fuego.ContextWithBody[CarsUpdate]) (Cars, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return Cars{}, err
	}

	return rs.CarsService.UpdateCars(id, body)
}

func (rs CarsResources) deleteCars(c fuego.ContextNoBody) (any, error) {
	return rs.CarsService.DeleteCars(c.PathParam("id"))
}
