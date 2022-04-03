package city

import (
	"example.com/with_gin/pkg/database_handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	db         *gorm.DB
	repository *CityRepository
	city       *City
}

func (s *Suite) SetupSuite() {
	conString := "root:Password123!@tcp(127.0.0.1:3306)/location_test?parseTime=True&loc=Local"
	s.db = database_handler.NewMySQLDB(conString)
	s.repository = NewCityRepository(s.db)
	// Migrate Table
	for _, val := range getModels() {
		s.db.AutoMigrate(val)
	}
}

func (s *Suite) TestCityRepository_CreateCity() {
	tests := []struct {
		tag  string
		city *City
	}{
		{"Adana", NewCity("Adana", "01", "TR")},
		{"Ankara", NewCity("Ankara", "02", "TR")},
		{"Paris", NewCity("Paris", "123", "FR")},
		{"Chicago", NewCity("Chicago", "1234", "US")},
	}
	for _, test := range tests {
		err := s.repository.Create(test.city)
		assert.Equal(s.T(), nil, err, "Error should be nil")
	}

}
func (s *Suite) TestCityRepository_GetByCountryCode() {

	c := s.repository.GetByCountryCode("TR")

	assert.Equal(s.T(), 2, len(c), "Return values length should be 2")
}

// Run After All Test Done
func (t *Suite) TearDownSuite() {
	sqlDB, _ := t.db.DB()
	defer sqlDB.Close()

	// Drop Table
	for _, val := range getModels() {
		t.db.Migrator().DropTable(val)
	}
}

func getModels() []interface{} {
	return []interface{}{&City{}}
}
