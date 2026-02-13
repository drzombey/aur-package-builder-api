package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/drzombey/aur-package-builder-api/cmd/api/config"
	repo "github.com/drzombey/aur-package-builder-api/cmd/api/db/repo"
	"github.com/drzombey/aur-package-builder-api/pkg/aur"
	"github.com/drzombey/aur-package-builder-api/pkg/docker"
	"github.com/drzombey/aur-package-builder-api/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageRepository is a mock implementation of PackageRepository
type MockPackageRepository struct {
	mock.Mock
}

func (m *MockPackageRepository) GetAurPackageInfoByName(name string) (*aur.Package, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aur.Package), args.Error(1)
}

func (m *MockPackageRepository) GetPackageFromAur(name string) ([]aur.Package, error) {
	args := m.Called(name)
	return args.Get(0).([]aur.Package), args.Error(1)
}

func (m *MockPackageRepository) GetAlreadyBuildPackageByAurIdAndVersion(id int64, version string) (*aur.Package, error) {
	args := m.Called(id, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aur.Package), args.Error(1)
}

func (m *MockPackageRepository) GetAlreadyBuildPackages() (*[]aur.Package, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]aur.Package), args.Error(1)
}

func (m *MockPackageRepository) AddAurPackage(model aur.Package) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockPackageRepository) DeleteAurPackage(model aur.Package) error {
	args := m.Called(model)
	return args.Error(0)
}

// MockStorageProvider is a mock implementation of storage.Provider
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) AddFile(ctx context.Context, filename string, stream io.Reader) error {
	args := m.Called(ctx, filename, stream)
	return args.Error(0)
}

func (m *MockStorageProvider) GetFile(ctx context.Context, filename string) (*os.File, error) {
	args := m.Called(ctx, filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockStorageProvider) DeleteFile(ctx context.Context, filename string) error {
	args := m.Called(ctx, filename)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createTestHandler() (*Handler, *MockPackageRepository, *MockStorageProvider) {
	mockRepo := &MockPackageRepository{}
	mockStorage := &MockStorageProvider{}
	
	cfg := &config.AppConfig{
		Docker: docker.DockerConfig{
			ContainerImage: "test-image",
			Username:       "testuser",
			Password:       "testpass",
			Auth:           true,
		},
		PackagePath: "/test/path",
	}
	
	var packageRepo repo.PackageRepository = mockRepo
	var storageProvider storage.Provider = mockStorage
	handler := NewHandler(cfg, storageProvider, packageRepo)
	return handler, mockRepo, mockStorage
}

func TestHandleGetAlreadyBuildPackages_Success(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	// Setup expected packages
	expectedPackages := []aur.Package{
		{ID: 1, Name: "test-package", Version: "1.0.0"},
		{ID: 2, Name: "another-package", Version: "2.0.0"},
	}
	
	mockRepo.On("GetAlreadyBuildPackages").Return(&expectedPackages, nil)
	
	// Setup gin router
	router := setupTestRouter()
	router.GET("/packages", handler.HandleGetAlreadyBuildPackages)
	
	// Create request
	req, _ := http.NewRequest("GET", "/packages", nil)
	w := httptest.NewRecorder()
	
	// Perform request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []aur.Package
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPackages, response)
	
	mockRepo.AssertExpectations(t)
}

func TestHandleGetAlreadyBuildPackages_Error(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	mockRepo.On("GetAlreadyBuildPackages").Return((*[]aur.Package)(nil), fmt.Errorf("database error"))
	
	router := setupTestRouter()
	router.GET("/packages", handler.HandleGetAlreadyBuildPackages)
	
	req, _ := http.NewRequest("GET", "/packages", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandleGetAurPackageByName_Success(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	expectedPackages := []aur.Package{
		{ID: 1, Name: "test-package", Version: "1.0.0"},
	}
	
	mockRepo.On("GetPackageFromAur", "test-package").Return(expectedPackages, nil)
	
	router := setupTestRouter()
	router.GET("/aur-package", handler.HandleGetAurPackageByName)
	
	req, _ := http.NewRequest("GET", "/aur-package?packageName=test-package", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []aur.Package
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPackages, response)
	
	mockRepo.AssertExpectations(t)
}

func TestHandleGetAurPackageByName_MissingParameter(t *testing.T) {
	handler, _, _ := createTestHandler()
	
	router := setupTestRouter()
	router.GET("/aur-package", handler.HandleGetAurPackageByName)
	
	req, _ := http.NewRequest("GET", "/aur-package", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGetAurPackageByName_Error(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	mockRepo.On("GetPackageFromAur", "test-package").Return([]aur.Package{}, fmt.Errorf("AUR API error"))
	
	router := setupTestRouter()
	router.GET("/aur-package", handler.HandleGetAurPackageByName)
	
	req, _ := http.NewRequest("GET", "/aur-package?packageName=test-package", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandleBuildPackage_Success(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	newPackage := aur.Package{
		ID:          1,
		Name:        "test-package",
		PackageBase: "test-package",
		Version:     "1.0.0",
	}
	
	// Mock that package doesn't exist yet
	mockRepo.On("GetAlreadyBuildPackageByAurIdAndVersion", int64(1), "1.0.0").Return((*aur.Package)(nil), nil)
	mockRepo.On("AddAurPackage", mock.MatchedBy(func(pkg aur.Package) bool {
		return pkg.ID == 1 && pkg.Version == "1.0.0"
	})).Return(nil)
	
	router := setupTestRouter()
	router.POST("/build", handler.HandleBuildPackage)
	
	jsonBody, _ := json.Marshal(newPackage)
	req, _ := http.NewRequest("POST", "/build", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusCreated), response["status"])
	assert.Contains(t, response["msg"], "Package currently creating")
	assert.Contains(t, response, "processId")
	
	mockRepo.AssertExpectations(t)
}

func TestHandleBuildPackage_PackageAlreadyExists(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	newPackage := aur.Package{
		ID:      1,
		Name:    "test-package",
		Version: "1.0.0",
	}
	
	existingPackage := &aur.Package{
		ID:      1,
		Name:    "test-package",
		Version: "1.0.0",
	}
	
	mockRepo.On("GetAlreadyBuildPackageByAurIdAndVersion", int64(1), "1.0.0").Return(existingPackage, nil)
	
	router := setupTestRouter()
	router.POST("/build", handler.HandleBuildPackage)
	
	jsonBody, _ := json.Marshal(newPackage)
	req, _ := http.NewRequest("POST", "/build", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["msg"], "Package already exist")
	
	mockRepo.AssertExpectations(t)
}

func TestHandleBuildPackage_InvalidJSON(t *testing.T) {
	handler, _, _ := createTestHandler()
	
	router := setupTestRouter()
	router.POST("/build", handler.HandleBuildPackage)
	
	req, _ := http.NewRequest("POST", "/build", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleBuildPackage_DatabaseError(t *testing.T) {
	handler, mockRepo, _ := createTestHandler()
	
	newPackage := aur.Package{
		ID:      1,
		Name:    "test-package",
		Version: "1.0.0",
	}
	
	mockRepo.On("GetAlreadyBuildPackageByAurIdAndVersion", int64(1), "1.0.0").Return((*aur.Package)(nil), fmt.Errorf("database error"))
	
	router := setupTestRouter()
	router.POST("/build", handler.HandleBuildPackage)
	
	jsonBody, _ := json.Marshal(newPackage)
	req, _ := http.NewRequest("POST", "/build", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}