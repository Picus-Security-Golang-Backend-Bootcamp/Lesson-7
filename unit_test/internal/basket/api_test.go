package basket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(usecase *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.Default()
	e.RemoveExtraSlash = true // varsayılan olarak false bu değer true setlenmediğinde test sırasında sürekli 404 durum kodu döner.
	mockRepo := &mockRepository{}
	RegisterHandlers(e, mockRepo)
	loadData(mockRepo)

	usecase.Run("Get Operations", func(t *testing.T) {
		tests := []struct {
			name       string
			args       string
			wantStatus int
		}{
			{name: "Get Basket_1", args: "ID_1", wantStatus: http.StatusOK},
			{name: "Get Basket_2", args: "", wantStatus: http.StatusNotFound},
			{name: "Get Basket_3", args: "INVALID_BASKET_ID", wantStatus: http.StatusNotFound},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", fmt.Sprintf("basket/%s", tt.args), nil)
				e.ServeHTTP(w, req)
				res := w.Result()
				defer res.Body.Close()
				assert.Equal(t, tt.wantStatus, w.Code)
				assert.NotNil(t, w.Body.String())

				t.Logf("Response:%v", w)

			})
		}

	})
}
