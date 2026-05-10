package controllers

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "dataanalysis/internal/errors"
	"dataanalysis/internal/models"
	"dataanalysis/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	strongestResult    *models.Beer
	weakestResult      *models.Beer
	oldestResult       *models.Beer
	newestResult       *models.Beer
	statsResult        *models.GeneralStats
	statsByBrandResult []models.BrandStats
	err                error
}

func (m *mockService) GetStrongest(_ context.Context) (*models.Beer, error) {
	return m.strongestResult, m.err
}
func (m *mockService) GetWeakest(_ context.Context) (*models.Beer, error) {
	return m.weakestResult, m.err
}
func (m *mockService) GetOldest(_ context.Context) (*models.Beer, error) {
	return m.oldestResult, m.err
}
func (m *mockService) GetNewest(_ context.Context) (*models.Beer, error) {
	return m.newestResult, m.err
}
func (m *mockService) GetStats(_ context.Context) (*models.GeneralStats, error) {
	return m.statsResult, m.err
}
func (m *mockService) GetStatsByBrand(_ context.Context) ([]models.BrandStats, error) {
	return m.statsByBrandResult, m.err
}

func init() {
	gin.SetMode(gin.TestMode)
}

type testSetup struct {
	router *gin.Engine
}

func newTestSetup(svc services.AnalysisServiceInterface) *testSetup {
	router := gin.New()
	ctrl := NewAnalysisController(svc)

	router.GET("/beers/rankings/strongest", ctrl.GetStrongest)
	router.GET("/beers/rankings/weakest", ctrl.GetWeakest)
	router.GET("/beers/rankings/oldest", ctrl.GetOldest)
	router.GET("/beers/rankings/newest", ctrl.GetNewest)
	router.GET("/beers/stats", ctrl.GetStats)
	router.GET("/beers/stats/brand", ctrl.GetStatsByBrand)

	return &testSetup{router: router}
}

func (ts *testSetup) do(method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	ts.router.ServeHTTP(w, req)
	return w
}

func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return body
}

func noResultsErr() error {
	return fmt.Errorf("%w", services.ErrNoResults)
}

func internalErr() error {
	return apperrors.Internal(stderrors.New("db failure"))
}

func TestGetStrongest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Chimay", Alcohol: 9.0}
	ts := newTestSetup(&mockService{strongestResult: beer})

	w := ts.do("GET", "/beers/rankings/strongest")

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Chimay", body["name"])
}

func TestGetStrongest_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/rankings/strongest")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
}

func TestGetStrongest_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/rankings/strongest")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "ERR_INTERNAL", body["code"])
}

func TestGetWeakest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Hoegaarden", Alcohol: 4.9}
	ts := newTestSetup(&mockService{weakestResult: beer})

	w := ts.do("GET", "/beers/rankings/weakest")

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Hoegaarden", body["name"])
}

func TestGetWeakest_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/rankings/weakest")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
}

func TestGetWeakest_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/rankings/weakest")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOldest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Chimay", Year: 1862}
	ts := newTestSetup(&mockService{oldestResult: beer})

	w := ts.do("GET", "/beers/rankings/oldest")

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Chimay", body["name"])
}

func TestGetOldest_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/rankings/oldest")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
}

func TestGetOldest_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/rankings/oldest")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetNewest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Hoegaarden", Year: 1966}
	ts := newTestSetup(&mockService{newestResult: beer})

	w := ts.do("GET", "/beers/rankings/newest")

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Hoegaarden", body["name"])
}

func TestGetNewest_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/rankings/newest")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
}

func TestGetNewest_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/rankings/newest")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetStats_Success(t *testing.T) {
	stats := &models.GeneralStats{Total: 10, TotalBrands: 3, AvgAlcohol: 6.5}
	ts := newTestSetup(&mockService{statsResult: stats})

	w := ts.do("GET", "/beers/stats")

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, float64(10), body["total"])
	assert.Equal(t, float64(3), body["total_brands"])
}

func TestGetStats_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/stats")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
	assert.Equal(t, "no beers found in database", body["message"])
}

func TestGetStats_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/stats")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "ERR_INTERNAL", body["code"])
}

func TestGetStatsByBrand_Success(t *testing.T) {
	stats := []models.BrandStats{
		{Brand: "Chimay", Count: 3, AvgAlcohol: 8.0},
		{Brand: "Duvel", Count: 1, AvgAlcohol: 8.5},
	}
	ts := newTestSetup(&mockService{statsByBrandResult: stats})

	w := ts.do("GET", "/beers/stats/brand")

	assert.Equal(t, http.StatusOK, w.Code)
	var body []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Len(t, body, 2)
}

func TestGetStatsByBrand_EmptyCollection(t *testing.T) {
	ts := newTestSetup(&mockService{err: noResultsErr()})

	w := ts.do("GET", "/beers/stats/brand")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMPTY_COLLECTION", body["code"])
}

func TestGetStatsByBrand_InternalError(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/stats/brand")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "ERR_INTERNAL", body["code"])
}

func TestWriteError_UnknownError(t *testing.T) {
	ts := newTestSetup(&mockService{err: stderrors.New("unexpected")})

	w := ts.do("GET", "/beers/stats")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "INTERNAL_ERROR", body["code"])
	assert.Equal(t, "internal server error", body["message"])
}

func TestWriteError_AppErrorPreferredOverSentinel(t *testing.T) {
	ts := newTestSetup(&mockService{err: internalErr()})

	w := ts.do("GET", "/beers/stats")

	body := parseBody(t, w)
	assert.Equal(t, "ERR_INTERNAL", body["code"])
}
