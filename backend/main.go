package main

import (
	"github.com/devmcclu/the-block/backend/domains/cars"
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	carResourses := cars.CarsResources{
		CarsService: cars.RealCarsService{},
	}

	carResourses.Routes(s)

	s.Run()
}
