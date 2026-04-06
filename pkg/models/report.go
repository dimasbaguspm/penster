package models

import "time"

type ReportSummary struct {
	TotalBalance   int64     `json:"total_balance"`
	TotalExpenses  int64     `json:"total_expenses"`
	TotalIncome    int64     `json:"total_income"`
	TotalTransfers int64     `json:"total_transfers"`
	BaseCurrency   string    `json:"base_currency"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
}

type CategoryBreakdown struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Type         string `json:"type"`
	Total        int64  `json:"total"`
}

type AccountBreakdown struct {
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Type        string `json:"type"`
	Total       int64  `json:"total"`
}

type TrendDataPoint struct {
	Date  string `json:"date"`
	Type  string `json:"type"`
	Total int64  `json:"total"`
}

type ReportByAccount struct {
	Accounts    []AccountBreakdown `json:"accounts"`
	PeriodStart time.Time          `json:"period_start"`
	PeriodEnd   time.Time          `json:"period_end"`
}

type ReportByCategory struct {
	Categories  []CategoryBreakdown `json:"categories"`
	PeriodStart time.Time           `json:"period_start"`
	PeriodEnd   time.Time           `json:"period_end"`
}

type ReportTrends struct {
	DataPoints  []TrendDataPoint `json:"data_points"`
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
}
