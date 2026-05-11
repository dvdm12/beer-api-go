package repository

import (
	"context"
	"errors"
	"testing"

	"dataanalysis/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockCollection struct {
	findResult    []models.Beer
	findErr       error
	findOneResult *models.Beer
	findOneErr    error
}

func (m *mockCollection) Find(_ context.Context, _ interface{}, _ ...*options.FindOptions) (*mongo.Cursor, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return mongo.NewCursorFromDocuments(nil, nil, nil)
}

func (m *mockCollection) FindOne(_ context.Context, _ interface{}, _ ...*options.FindOneOptions) *mongo.SingleResult {
	if m.findOneErr != nil {
		return mongo.NewSingleResultFromDocument(nil, m.findOneErr, nil)
	}
	if m.findOneResult != nil {
		return mongo.NewSingleResultFromDocument(m.findOneResult, nil, nil)
	}
	return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}

func TestNewAnalysisRepository_NilLogger(t *testing.T) {
	repo := NewAnalysisRepository(&mockCollection{}, nil)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.logger)
}

func TestNewAnalysisRepository_WithLogger(t *testing.T) {
	logger := &noopLogger{}
	repo := NewAnalysisRepository(&mockCollection{}, logger)
	assert.NotNil(t, repo)
	assert.Equal(t, logger, repo.logger)
}

func TestBeerFilter_Validate_Valid(t *testing.T) {
	min, max := 4.0, 8.0
	from, to := 1900, 2000
	filter := BeerFilter{MinAlcohol: &min, MaxAlcohol: &max, FromYear: &from, ToYear: &to}
	assert.NoError(t, filter.Validate())
}

func TestBeerFilter_Validate_EmptyFilter(t *testing.T) {
	assert.NoError(t, BeerFilter{}.Validate())
}

func TestBeerFilter_Validate_MinExceedsMax(t *testing.T) {
	min, max := 10.0, 5.0
	err := BeerFilter{MinAlcohol: &min, MaxAlcohol: &max}.Validate()
	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func TestBeerFilter_Validate_FromYearExceedsToYear(t *testing.T) {
	from, to := 2020, 2000
	err := BeerFilter{FromYear: &from, ToYear: &to}.Validate()
	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func TestBeerFilter_Validate_NameAndNameLikeMutuallyExclusive(t *testing.T) {
	name, like := "Chimay", "chim"
	err := BeerFilter{Name: &name, NameLike: &like}.Validate()
	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func TestBeerFilter_Validate_MinOnlyValid(t *testing.T) {
	min := 4.0
	assert.NoError(t, BeerFilter{MinAlcohol: &min}.Validate())
}

func TestBeerFilter_Validate_MaxOnlyValid(t *testing.T) {
	max := 8.0
	assert.NoError(t, BeerFilter{MaxAlcohol: &max}.Validate())
}

func TestBeerFilter_Validate_EqualMinMax(t *testing.T) {
	v := 5.0
	assert.NoError(t, BeerFilter{MinAlcohol: &v, MaxAlcohol: &v}.Validate())
}

func TestBeerFilter_Validate_EqualYears(t *testing.T) {
	y := 2000
	assert.NoError(t, BeerFilter{FromYear: &y, ToYear: &y}.Validate())
}

func TestBuildQuery_EmptyFilter(t *testing.T) {
	q := buildQuery(BeerFilter{})
	assert.Empty(t, q)
}

func TestBuildQuery_Brand(t *testing.T) {
	brand := "Chimay"
	q := buildQuery(BeerFilter{Brand: &brand})
	assert.Equal(t, "Chimay", q["brand"])
}

func TestBuildQuery_Name(t *testing.T) {
	name := "Duvel"
	q := buildQuery(BeerFilter{Name: &name})
	assert.Equal(t, "Duvel", q["name"])
}

func TestBuildQuery_NameLike(t *testing.T) {
	like := "chim"
	q := buildQuery(BeerFilter{NameLike: &like})
	_, ok := q["name"]
	assert.True(t, ok)
}

func TestBuildQuery_AlcoholRange(t *testing.T) {
	min, max := 4.0, 8.0
	q := buildQuery(BeerFilter{MinAlcohol: &min, MaxAlcohol: &max})
	alcohol, ok := q["alcohol"].(bson.M)
	require.True(t, ok)
	assert.Equal(t, 4.0, alcohol["$gte"])
	assert.Equal(t, 8.0, alcohol["$lte"])
}

func TestBuildQuery_AlcoholMinOnly(t *testing.T) {
	min := 4.0
	q := buildQuery(BeerFilter{MinAlcohol: &min})
	alcohol, ok := q["alcohol"].(bson.M)
	require.True(t, ok)
	assert.Equal(t, 4.0, alcohol["$gte"])
	assert.Nil(t, alcohol["$lte"])
}

func TestBuildQuery_YearRange(t *testing.T) {
	from, to := 1900, 2000
	q := buildQuery(BeerFilter{FromYear: &from, ToYear: &to})
	year, ok := q["year"].(bson.M)
	require.True(t, ok)
	assert.Equal(t, 1900, year["$gte"])
	assert.Equal(t, 2000, year["$lte"])
}

func TestBuildQuery_NameLikeOverridesName(t *testing.T) {
	name, like := "Duvel", "duv"
	// BeerFilter prevents both, but buildQuery handles NameLike last
	q := buildQuery(BeerFilter{Name: &name, NameLike: &like})
	// NameLike overwrites Name in query
	_, isRegex := q["name"].(primitive.Regex)
	assert.True(t, isRegex)
}

func TestFind_InvalidFilter(t *testing.T) {
	min, max := 10.0, 5.0
	repo := NewAnalysisRepository(&mockCollection{}, nil)

	_, err := repo.Find(context.Background(), BeerFilter{MinAlcohol: &min, MaxAlcohol: &max})

	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func TestFind_CollectionError(t *testing.T) {
	repo := NewAnalysisRepository(&mockCollection{
		findErr: errors.New("connection refused"),
	}, nil)

	_, err := repo.Find(context.Background(), BeerFilter{})

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryNetwork, repoErr.Category)
}

func TestFind_EmptyResult(t *testing.T) {
	repo := NewAnalysisRepository(&mockCollection{}, nil)

	result, err := repo.Find(context.Background(), BeerFilter{})

	// cursor returns empty — no error, empty slice
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFindTop_InvalidSortField(t *testing.T) {
	repo := NewAnalysisRepository(&mockCollection{}, nil)

	_, err := repo.FindTop(context.Background(), SortField("invalid"), true)

	assert.ErrorIs(t, err, ErrInvalidFilter)
}

func TestFindTop_NotFound(t *testing.T) {
	err := MapMongoError(mongo.ErrNoDocuments, OpFindTop, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryNotFound, repoErr.Category)
	assert.Equal(t, OpFindTop, repoErr.Operation)
}

func TestFindTop_NetworkError(t *testing.T) {
	err := MapMongoError(errors.New("connection refused"), OpFindTop, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryNetwork, repoErr.Category)
}

func TestFindTop_Success(t *testing.T) {
	expected := &models.Beer{Name: "Chimay", Alcohol: 9.0}
	repo := NewAnalysisRepository(&mockCollection{findOneResult: expected}, nil)

	result, err := repo.FindTop(context.Background(), SortByAlcohol, true)

	require.NoError(t, err)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Alcohol, result.Alcohol)
}

func TestFindTop_AllValidSortFields(t *testing.T) {
	fields := []SortField{SortByAlcohol, SortByYear, SortByName}
	expected := &models.Beer{Name: "Duvel"}

	for _, field := range fields {
		repo := NewAnalysisRepository(&mockCollection{findOneResult: expected}, nil)
		result, err := repo.FindTop(context.Background(), field, true)
		require.NoError(t, err, "field: %s", field)
		assert.Equal(t, expected.Name, result.Name)
	}
}

func TestFindTop_DescVsAsc(t *testing.T) {
	beer := &models.Beer{Name: "Heineken", Alcohol: 5.0}
	repo := NewAnalysisRepository(&mockCollection{findOneResult: beer}, nil)

	// Both directions should succeed when collection has results
	resultDesc, errDesc := repo.FindTop(context.Background(), SortByAlcohol, true)
	resultAsc, errAsc := repo.FindTop(context.Background(), SortByAlcohol, false)

	require.NoError(t, errDesc)
	require.NoError(t, errAsc)
	assert.Equal(t, beer.Name, resultDesc.Name)
	assert.Equal(t, beer.Name, resultAsc.Name)
}

func TestRepoError_Error_WithOperation(t *testing.T) {
	err := newRepoError(CategoryNotFound, MsgNotFound, OpFind, errors.New("cause"))
	assert.Equal(t, "[NOT_FOUND] document not found (find): cause", err.Error())
}

func TestRepoError_Error_WithoutOperation(t *testing.T) {
	err := &RepoError{Category: CategoryUnknown, Message: MsgUnmapped, Cause: errors.New("cause")}
	assert.Equal(t, "[UNKNOWN] unmapped mongo error: cause", err.Error())
}

func TestRepoError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := newRepoError(CategoryTimeout, MsgTimeout, OpFindTop, cause)
	assert.Equal(t, cause, err.Unwrap())
}

func TestNoopLogger_Error(t *testing.T) {
	logger := &noopLogger{}
	assert.NotPanics(t, func() {
		logger.Error("test message", "key", "value")
	})
}

func TestMapMongoError_Nil(t *testing.T) {
	assert.Nil(t, MapMongoError(nil, OpFind, CollectionBeers, nil))
}

func TestMapMongoError_ErrNilCursor(t *testing.T) {
	err := MapMongoError(mongo.ErrNilCursor, OpFind, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryCursor, repoErr.Category)
	assert.Equal(t, MsgNilCursor, repoErr.Message)
}

func TestMapMongoError_ErrClientDisconnected(t *testing.T) {
	err := MapMongoError(mongo.ErrClientDisconnected, OpFind, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryNetwork, repoErr.Category)
	assert.Equal(t, MsgClientDisconn, repoErr.Message)
}

func TestMapMongoError_CommandError_Timeout(t *testing.T) {
	cmd := mongo.CommandError{Code: 50, Message: "MaxTimeMSExpired"}
	err := MapMongoError(cmd, OpFindTop, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryTimeout, repoErr.Category)
	assert.Equal(t, MsgTimeout, repoErr.Message)
}

func TestMapMongoError_CommandError_ServerInterrupt(t *testing.T) {
	for _, code := range []int32{11600, 11601} {
		cmd := mongo.CommandError{Code: code}
		err := MapMongoError(cmd, OpFind, CollectionBeers, nil)

		var repoErr *RepoError
		require.True(t, errors.As(err, &repoErr))
		assert.Equal(t, CategoryNetwork, repoErr.Category)
		assert.Equal(t, MsgServerInterrupt, repoErr.Message)
	}
}

func TestMapMongoError_CommandError_SocketError(t *testing.T) {
	for _, code := range []int32{9001, 211} {
		cmd := mongo.CommandError{Code: code}
		err := MapMongoError(cmd, OpFind, CollectionBeers, nil)

		var repoErr *RepoError
		require.True(t, errors.As(err, &repoErr))
		assert.Equal(t, CategoryNetwork, repoErr.Category)
		assert.Equal(t, MsgSocketError, repoErr.Message)
	}
}

func TestMapMongoError_CommandError_Unknown(t *testing.T) {
	cmd := mongo.CommandError{Code: 999, Message: "some error"}
	err := MapMongoError(cmd, OpFind, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryUnknown, repoErr.Category)
}

func TestMapByMessage_Timeout(t *testing.T) {
	for _, msg := range []string{"context deadline exceeded", "timed out", "timeout"} {
		err := MapMongoError(errors.New(msg), OpFind, CollectionBeers, nil)

		var repoErr *RepoError
		require.True(t, errors.As(err, &repoErr))
		assert.Equal(t, CategoryTimeout, repoErr.Category, "pattern: %s", msg)
	}
}

func TestMapByMessage_Decode(t *testing.T) {
	for _, msg := range []string{"cannot decode", "no decoder found"} {
		err := MapMongoError(errors.New(msg), OpFind, CollectionBeers, nil)

		var repoErr *RepoError
		require.True(t, errors.As(err, &repoErr))
		assert.Equal(t, CategoryDecode, repoErr.Category, "pattern: %s", msg)
	}
}

func TestMapByMessage_Unknown(t *testing.T) {
	err := MapMongoError(errors.New("some completely unknown error"), OpFind, CollectionBeers, nil)

	var repoErr *RepoError
	require.True(t, errors.As(err, &repoErr))
	assert.Equal(t, CategoryUnknown, repoErr.Category)
	assert.Equal(t, MsgUnmapped, repoErr.Message)
}

func TestMapMongoError_WithLogger(t *testing.T) {
	logged := false
	logger := &testLogger{onError: func() { logged = true }}

	MapMongoError(errors.New("some error"), OpFind, CollectionBeers, logger)

	assert.True(t, logged)
}

type testLogger struct {
	onError func()
}

func (l *testLogger) Error(_ string, _ ...any) {
	l.onError()
}
