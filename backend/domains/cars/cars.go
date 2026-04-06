package cars

type Cars struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CarsCreate struct {
	Name string `json:"name"`
}

type CarsUpdate struct {
	Name string `json:"name"`
}

type CarsService interface {
	GetCars(id string) (Cars, error)
	CreateCars(CarsCreate) (Cars, error)
	GetAllCars() ([]Cars, error)
	UpdateCars(id string, input CarsUpdate) (Cars, error)
	DeleteCars(id string) (any, error)
}

type RealCarsService struct {
	CarsService
}

func (s RealCarsService) GetCars(id string) (Cars, error) {
	return Cars{}, nil
}

func (s RealCarsService) CreateCars(input CarsCreate) (Cars, error) {
	return Cars{}, nil
}

func (s RealCarsService) GetAllCars() ([]Cars, error) {
	return []Cars{}, nil
}

func (s RealCarsService) UpdateCars(id string, input CarsUpdate) (Cars, error) {
	return Cars{}, nil
}

func (s RealCarsService) DeleteCars(id string) (any, error) {
	return nil, nil
}
