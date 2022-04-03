package city

type CityService struct {
	r CityRepository
}

func NewCityService(r CityRepository) *CityService {
	return &CityService{
		r: r,
	}
}

func (c *CityService) Create(city *City) error {
	existCity := c.r.GetByName(city.Name)
	if existCity != nil && len(existCity) > 0 {
		return ErrCityExistWithName
	}

	existCity = c.r.GetByCityCode(city.Code)
	if existCity != nil && len(existCity) > 0 {
		return ErrCityExistWithCityCode
	}

	err := c.r.Create(city)
	if err != nil {
		return err
	}

	return nil
}

func (c *CityService) GetAll(pageIndex, pageSize int) ([]City, int) {
	items, count := c.r.GetAll(pageIndex, pageSize)

	return items, count
}
