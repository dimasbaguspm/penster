-- name: UpsertRateCurrency :one
INSERT INTO rate_currencies (from_currency, to_currency, rate, rate_date)
VALUES (@from_currency, @to_currency, @rate, @rate_date)
ON CONFLICT (from_currency, to_currency, rate_date)
DO UPDATE SET rate = EXCLUDED.rate
RETURNING id, from_currency, to_currency, rate, rate_date, created_at;

-- name: GetRateCurrency :one
SELECT id, from_currency, to_currency, rate, rate_date, created_at
FROM rate_currencies
WHERE from_currency = @from_currency
  AND to_currency = @to_currency
  AND rate_date = @rate_date;

-- name: ListRateCurrencies :many
SELECT id, from_currency, to_currency, rate, rate_date, created_at
FROM rate_currencies
WHERE
    ($1::text IS NULL OR from_currency = $1::text)
    AND ($2::text IS NULL OR to_currency = $2::text)
ORDER BY rate_date DESC, from_currency, to_currency
LIMIT NULLIF($3, 0)
OFFSET $4;

-- name: CountRateCurrencies :one
SELECT count(*) as total
FROM rate_currencies
WHERE
    (@from_currency::text IS NULL OR from_currency = @from_currency)
    AND (@to_currency::text IS NULL OR to_currency = @to_currency);

-- name: PruneOldRates :exec
DELETE FROM rate_currencies
WHERE rate_date < @older_than;
