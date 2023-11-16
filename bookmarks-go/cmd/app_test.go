package cmd_test

import (
	"encoding/json"
	"github.com/sivaprasadreddy/bookmarks/cmd"
	"github.com/sivaprasadreddy/bookmarks/internal/config"
	"github.com/sivaprasadreddy/bookmarks/internal/domain"
	"github.com/sivaprasadreddy/bookmarks/testsupport"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	PgContainer *testsupport.PostgresContainer
	cfg         config.AppConfig
	app         *cmd.App
	router      http.Handler
}

func (suite *ControllerTestSuite) SetupSuite() {
	suite.PgContainer = testsupport.InitPostgresContainer()
	cfg, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	suite.cfg = cfg

	suite.app = cmd.NewApp(suite.cfg)
	suite.router = suite.app.Router
}

func (suite *ControllerTestSuite) TearDownSuite() {
	suite.PgContainer.CloseFn()
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (suite *ControllerTestSuite) TestGetAllBookmarks() {
	t := suite.T()
	req, _ := http.NewRequest("GET", "/api/bookmarks", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	actualResponseJSON := w.Body.String()
	assert.NotEqual(t, "[]", actualResponseJSON)
}

func (suite *ControllerTestSuite) TestGetBookmarkByID() {
	t := suite.T()
	req, _ := http.NewRequest(http.MethodGet, "/api/bookmarks/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Bookmark
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, "How To Remove Docker Containers, Images, Volumes, and Networks", response.Title)
	assert.Equal(t, "https://linuxize.com/post/how-to-remove-docker-images-containers-volumes-and-networks/", response.Url)
	assert.NotNil(t, response.CreatedAt)
}

func (suite *ControllerTestSuite) TestCreateBookmark() {
	t := suite.T()
	reqBody := strings.NewReader(`
		{
			"title": "Test Post title",
			"url":     "https://example.com"
		}
	`)

	req, _ := http.NewRequest(http.MethodPost, "/api/bookmarks", reqBody)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.Bookmark
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, "Test Post title", response.Title)
	assert.Equal(t, "https://example.com", response.Url)
	assert.NotNil(t, response.CreatedAt)
}

func (suite *ControllerTestSuite) TestUpdateBookmark() {
	t := suite.T()
	reqBody := strings.NewReader(`
		{
			"title": "Test Updated title",
			"url":   "https://example2.com"
		}
	`)

	req, _ := http.NewRequest(http.MethodPut, "/api/bookmarks/1", reqBody)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func (suite *ControllerTestSuite) TestDeleteBookmark() {
	t := suite.T()

	req, _ := http.NewRequest(http.MethodDelete, "/api/bookmarks/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
