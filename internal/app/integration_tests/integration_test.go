package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"userportal/internal/app/dto"
	"userportal/internal/app/handler"
	"userportal/internal/app/repository"
	"userportal/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDbInstance *sqlx.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestGetAllUsers(t *testing.T) {
	repo := repository.NewUserRepository(testDbInstance)
	service := service.NewUserService(repo)
	h := handler.NewUserHandler(service)

	engine := gin.Default()

	engine.GET("/test/users", h.GetAllUsers)

	request, err := http.NewRequest(http.MethodGet, "/test/users", nil)
	require.NoError(t, err)

	responseRecoder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecoder, request)

	var users []dto.User
	json.NewDecoder(responseRecoder.Body).Decode(&users)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecoder.Result().StatusCode)
	for _, user := range users {

		assert.Contains(t, user.FirstName, "radha geethika")

		//check response
	}

}
