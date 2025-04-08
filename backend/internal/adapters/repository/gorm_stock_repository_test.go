package repository_test

import (
	"regexp"
	"stock-information/internal/adapters/repository"
	"stock-information/internal/domain"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetRecommendations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creando mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error abriendo GORM DB: %v", err)
	}

	repo := &repository.GormCockroachRepo{DB: gormDB}

	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"recommendations\"$").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery("^SELECT \\* FROM \"recommendations\" ORDER BY ticker asc LIMIT \\$1$").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}).
			AddRow("AAPL", "Apple", "Buy", "Neutral", "Buy", 100.0, 150.0).
			AddRow("GOOG", "Google", "Hold", "Neutral", "Hold", 2000.0, 2500.0))

	recs, totalCount, err := repo.GetRecommendations(1, 10, "asc")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(recs))
	assert.Equal(t, 2, totalCount)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas: %v", err)
	}
}

func TestGetAllRecommendations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM DB: %v", err)
	}

	repo := &repository.GormCockroachRepo{DB: gormDB}

	mock.ExpectQuery(`SELECT \* FROM "recommendations"`).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}).
			AddRow("AAPL", "Apple", "Buy", "Neutral", "Buy", 100.0, 150.0).
			AddRow("GOOG", "Google", "Hold", "Neutral", "Hold", 2000.0, 2500.0))

	recs, err := repo.GetAllRecommendations()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(recs))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestSaveRecommendations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creando mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM DB: %v", err)
	}

	repo := &repository.GormCockroachRepo{DB: gormDB}

	recommendations := []domain.Recommendation{
		{Ticker: "AAPL", Brokerage: "Goldman Sachs", Company: "Apple", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 150.0, TargetTo: 180.0},
		{Ticker: "GOOG", Brokerage: "Morgan Stanley", Company: "Google", Action: "Hold", RatingFrom: "Buy", RatingTo: "Hold", TargetFrom: 2500.0, TargetTo: 2700.0},
	}

	for _, rec := range recommendations {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "recommendations" ("ticker","brokerage","company","action","rating_from","rating_to","target_from","target_to") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
			WithArgs(rec.Ticker, rec.Brokerage, rec.Company, rec.Action, rec.RatingFrom, rec.RatingTo, rec.TargetFrom, rec.TargetTo).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	for _, rec := range recommendations {
		err = repo.DB.Create(&rec).Error
		assert.NoError(t, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas: %v", err)
	}
}

func TestGetRecommendationByTicker(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creando mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM DB: %v", err)
	}

	repo := &repository.GormCockroachRepo{DB: gormDB}

	ticker := "AAPL"
	expectedRecommendation := domain.Recommendation{
		Ticker:     "AAPL",
		Brokerage:  "Goldman Sachs",
		Company:    "Apple Inc.",
		Action:     "Buy",
		RatingFrom: "Neutral",
		RatingTo:   "Buy",
		TargetFrom: 150.0,
		TargetTo:   180.0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recommendations" WHERE ticker = $1 ORDER BY "recommendations"."ticker" LIMIT $2`)).
		WithArgs(ticker, 1).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "brokerage", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}).
			AddRow(expectedRecommendation.Ticker, expectedRecommendation.Brokerage, expectedRecommendation.Company, expectedRecommendation.Action, expectedRecommendation.RatingFrom, expectedRecommendation.RatingTo, expectedRecommendation.TargetFrom, expectedRecommendation.TargetTo))

	rec, err := repo.GetRecommendationByTicker(ticker)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendation, rec)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas: %v", err)
	}

	notFoundTicker := "MSFT"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recommendations" WHERE ticker = $1 ORDER BY "recommendations"."ticker" LIMIT $2`)).
		WithArgs(notFoundTicker, 1).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "brokerage", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}))

	_, err = repo.GetRecommendationByTicker(notFoundTicker)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas (not found): %v", err)
	}
}

func TestGormCockroachRepo_GetRecommendationsByCompany_Simplified(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creando mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM DB: %v", err)
	}

	repo := &repository.GormCockroachRepo{DB: gormDB}

	company := "Apple"
	page := 1
	limit := 10
	sort := "asc"

	expectedRecommendations := []domain.Recommendation{
		{Ticker: "AAPL", Brokerage: "Goldman Sachs", Company: "Apple Inc.", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 150.0, TargetTo: 180.0},
		{Ticker: "APLE", Brokerage: "JP Morgan", Company: "Apple Pie Corp", Action: "Hold", RatingFrom: "Hold", RatingTo: "Hold", TargetFrom: 50.0, TargetTo: 60.0},
	}
	totalCount := 2

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recommendations" WHERE LOWER(company) LIKE $1`)).
		WithArgs("%" + strings.ToLower(company) + "%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalCount))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recommendations" WHERE LOWER(company) LIKE $1 ORDER BY ticker asc LIMIT $2`)).
		WithArgs("%"+strings.ToLower(company)+"%", limit).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "brokerage", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}).
			AddRow(expectedRecommendations[0].Ticker, expectedRecommendations[0].Brokerage, expectedRecommendations[0].Company, expectedRecommendations[0].Action, expectedRecommendations[0].RatingFrom, expectedRecommendations[0].RatingTo, expectedRecommendations[0].TargetFrom, expectedRecommendations[0].TargetTo).
			AddRow(expectedRecommendations[1].Ticker, expectedRecommendations[1].Brokerage, expectedRecommendations[1].Company, expectedRecommendations[1].Action, expectedRecommendations[1].RatingFrom, expectedRecommendations[1].RatingTo, expectedRecommendations[1].TargetFrom, expectedRecommendations[1].TargetTo))

	recs, count, err := repo.GetRecommendationsByCompany(company, page, limit, sort)
	assert.NoError(t, err)
	assert.Equal(t, totalCount, count)
	assert.Equal(t, expectedRecommendations, recs)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas: %v", err)
	}

	sort = "desc"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recommendations" WHERE LOWER(company) LIKE $1`)).
		WithArgs("%" + strings.ToLower(company) + "%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalCount))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recommendations" WHERE LOWER(company) LIKE $1 ORDER BY ticker desc LIMIT $2`)).
		WithArgs("%"+strings.ToLower(company)+"%", limit).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "brokerage", "company", "action", "rating_from", "rating_to", "target_from", "target_to"}).
			AddRow(expectedRecommendations[1].Ticker, expectedRecommendations[1].Brokerage, expectedRecommendations[1].Company, expectedRecommendations[1].Action, expectedRecommendations[1].RatingFrom, expectedRecommendations[1].RatingTo, expectedRecommendations[1].TargetFrom, expectedRecommendations[1].TargetTo).
			AddRow(expectedRecommendations[0].Ticker, expectedRecommendations[0].Brokerage, expectedRecommendations[0].Company, expectedRecommendations[0].Action, expectedRecommendations[0].RatingFrom, expectedRecommendations[0].RatingTo, expectedRecommendations[0].TargetFrom, expectedRecommendations[0].TargetTo))

	recs, count, err = repo.GetRecommendationsByCompany(company, page, limit, sort)
	assert.NoError(t, err)
	assert.Equal(t, totalCount, count)
	assert.Equal(t, []domain.Recommendation{expectedRecommendations[1], expectedRecommendations[0]}, recs)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas (desc): %v", err)
	}

	// Test with pagination (page 2)
	page = 2
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recommendations" WHERE LOWER(company) LIKE $1`)).
		WithArgs("%" + strings.ToLower(company) + "%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalCount))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recommendations" WHERE LOWER(company) LIKE $1 ORDER BY ticker asc LIMIT $2 OFFSET $3`)).
		WithArgs("%"+strings.ToLower(company)+"%", limit, (page-1)*limit).
		WillReturnRows(sqlmock.NewRows([]string{"ticker", "brokerage", "company", "action", "rating_from", "rating_to", "target_from", "target_to"})) // No rows for page 2 in this mock

	recs, count, err = repo.GetRecommendationsByCompany(company, page, limit, "asc")
	assert.NoError(t, err)
	assert.Equal(t, totalCount, count)
	assert.Empty(t, recs)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Hubo expectativas no cumplidas (page 2): %v", err)
	}
}
