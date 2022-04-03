package city

import (
	"example.com/with_gin/pkg/pagination"
	"net/http"

	"example.com/with_gin/internal/city"
	"github.com/gin-gonic/gin"
)

type CityController struct {
	cityService *city.CityService
}

func NewCityController(service *city.CityService) *CityController {
	return &CityController{
		cityService: service,
	}
}

// GetAllCities godoc
// @Summary Gets all cities with paginated result
// @Tags City
// @Accept  json
// @Produce  json
// @Param page query int false "Page Index"
// @Param pageSize query int false "Page Size"
// @Success 200 {object} pagination.Pages
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /city [get]
func (c *CityController) GetAllCities(g *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(g)
	items, count := c.cityService.GetAll(pageIndex, pageSize)
	paginatedResult := pagination.NewFromGinRequest(g, count)
	paginatedResult.Items = items

	g.JSON(http.StatusCreated, paginatedResult)
}

// CreateCity godoc
// @Summary Creates a new city
// @Tags City
// @Accept  json
// @Produce  json
// @Param createRequest body CreateCityRequest true "City informations"
// @Success 200 {object} CityResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /city [post]
func (c *CityController) CreateCity(g *gin.Context) {
	var req CreateCityRequest
	if err := g.ShouldBind(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Check your request body.",
		})
		g.Abort()
		return
	}
	city := city.NewCity(req.Name, req.Code, req.CountryCode)
	err := c.cityService.Create(city)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error_message": err.Error(),
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusCreated, CityResponse{
		Name:        city.Name,
		Code:        city.Code,
		CountryCode: city.CountryCode,
	})
}

// city?name=sdfsdf
func (c *CityController) GetQueryString(g *gin.Context) {
	name := g.Query("name")
	code := g.Query("code")
	countryCode, isOk := g.GetQuery("countryCode")
	if isOk {
		g.JSON(http.StatusOK, gin.H{
			"name":         name,
			"code":         code,
			"country_code": countryCode,
		})
	} else {
		g.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Country code cannot be null or empty.",
		})
	}
}
