package jobs

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/scheduler"
	"github.com/dimasbaguspm/penster/pkg/models"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type RateCurrencyJob struct {
	cfg *config.Config
	svc *service.RateCurrencyService
}

var _ scheduler.Job = (*RateCurrencyJob)(nil)

func NewRateCurrencyJob(cfg *config.Config, svc *service.RateCurrencyService) *RateCurrencyJob {
	return &RateCurrencyJob{cfg: cfg, svc: svc}
}

func (j *RateCurrencyJob) Name() string {
	return "rate_currency"
}

func (j *RateCurrencyJob) Schedule() scheduler.Schedule {
	return scheduler.IntervalSchedule{
		Interval: j.cfg.RateCurrency.Interval,
	}
}

func (j *RateCurrencyJob) Run(ctx context.Context) error {
	log := observability.NewLogger(ctx, "scheduler", "jobs")
	ctx, span := observability.StartJobSpan(log.Context(), "rate_currency")
	defer span.End()

	log.Info("Running rate_currency job - fetching ECB rates")

	rates, err := fetchECBRates(ctx, j.cfg.RateCurrency.ECBURL)
	if err != nil {
		observability.RecordError(ctx, err)
		return fmt.Errorf("failed to fetch ECB rates: %w", err)
	}

	if len(rates) == 0 {
		return fmt.Errorf("no rates found in ECB data")
	}

	rateDate := time.Now().Truncate(24 * time.Hour)

	existing, err := j.svc.Get(ctx, "EUR", "USD", rateDate)
	if err != nil {
		log.Error("Failed to check existing rate", "error", err)
	}

	if existing != nil {
		log.Info("Rates already exist, skipping upsert", "date", rateDate.Format("2006-01-02"))
		return nil
	}

	log.Info("Upserting rates", "count", len(rates), "date", rateDate.Format("2006-01-02"))

	for currency, rate := range rates {
		req := &models.UpsertRateCurrencyRequest{
			FromCurrency: "EUR",
			ToCurrency:   currency,
			Rate:         rate,
			RateDate:     rateDate,
		}
		_, err := j.svc.Upsert(ctx, req)
		if err != nil {
			log.Error("Failed to upsert rate", "from", "EUR", "to", currency, "error", err)
		}
	}

	log.Info("Rate currency job completed")
	return nil
}

type ecbEnvelope struct {
	Cube ecbCube `xml:"Cube"`
}

type ecbCube struct {
	Time string   `xml:"time,attr"`
	Cube ecbCube2 `xml:"Cube"`
}

type ecbCube2 struct {
	Cubes []ecbRate `xml:"Cube"`
}

type ecbRate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

func fetchECBRates(ctx context.Context, url string) (map[string]float64, error) {
	ctx, span := observability.StartJobSpan(ctx, "fetch_ecb_rates")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ECB data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ECB returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return parseXML(data)
}

func parseXML(data []byte) (map[string]float64, error) {
	var envelope ecbEnvelope
	if err := xml.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("XML unmarshal failed: %w", err)
	}

	rates := make(map[string]float64)
	rates["EUR"] = 1.0

	for _, cube := range envelope.Cube.Cube.Cubes {
		if cube.Currency == "" {
			continue
		}
		if cube.Rate <= 0 {
			return nil, fmt.Errorf("invalid rate for %s: %f", cube.Currency, cube.Rate)
		}
		rates[cube.Currency] = cube.Rate
	}

	return rates, nil
}
