package city

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{
		db: db,
	}
}

func (r *CityRepository) GetAll(pageIndex, pageSize int) ([]City, int) {
	var cities []City
	var count int64

	r.db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&cities).Count(&count)

	return cities, int(count)
}

func (r *CityRepository) GetByCountryCode(countryCode string) []City {
	var cities []City
	r.db.Where("CountryCode = ?", countryCode).Order("Id desc,name").Find(&cities)

	return cities
}

func (r *CityRepository) GetByCountryCodeOrCityCode(code string) []City {
	var cities []City
	r.db.Where("CountryCode = ?", code).Or("Code = ?", code).Find(&cities)
	return cities
}
func (r *CityRepository) GetByCityCode(code string) []City {
	var cities []City
	r.db.Where("Code = ?", code).Find(&cities)
	return cities
}
func (r *CityRepository) GetByName(name string) []City {
	var cities []City
	r.db.Where("Name LIKE ?", "%"+name+"%").Find(&cities)

	return cities
}

func (r *CityRepository) GetByNameWithRawSQL(name string) []City {
	var cities []City
	r.db.Raw("SELECT * FROM City WHERE Name LIKE ?", "%"+name+"%").Scan(&cities)

	return cities
}

func (r *CityRepository) GetById(id int) City {
	var city City
	result := r.db.First(&city, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("City not found with id : %d", id)
		return City{}
	}
	return city
}

// Eğer pointer verilmezse reflect.Value.Set using unaddressable value hatası alınır
func (r *CityRepository) Create(c *City) error {
	result := r.db.Create(c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Update(c City) error {
	result := r.db.Save(c)
	//r.db.Model(&c).Update("name", "deneme")

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Delete(c City) error {
	result := r.db.Delete(c)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *CityRepository) DeleteById(id int) error {
	result := r.db.Delete(&City{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Migration() {
	r.db.AutoMigrate(&City{})
	//https://gorm.io/docs/migration.html#content-inner
}

func (r *CityRepository) InsertSampleData() {
	cities := []City{
		{Name: "Adana", Code: "01", CountryCode: "TR"},
		{Name: "Adıyaman", Code: "02", CountryCode: "TR"},
		{Name: "Ankara", Code: "06", CountryCode: "TR"},
		{Name: "İstanbul", Code: "34", CountryCode: "TR"},
		{Name: "İzmir", Code: "35", CountryCode: "TR"},
	}

	for _, c := range cities {
		r.db.Where(City{Code: c.Code}).Attrs(City{Code: c.Code, Name: c.Name}).FirstOrCreate(&c)
	}
}
