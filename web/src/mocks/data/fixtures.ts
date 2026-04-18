import type {
  HandlerHealthResponse,
  ModelsAccount,
  ModelsCategory,
  ModelsDraft,
  ModelsReportByAccount,
  ModelsReportByCategory,
  ModelsReportSummary,
  ModelsReportTrends,
  ModelsTransaction,
} from "@/api/types";
import { ModelsAccountType, ModelsCategoryType, ModelsTransactionType } from "@/api/types";

export const mockAccounts: ModelsAccount[] = [
  {
    id: "acc-001",
    name: "Checking Account",
    type: ModelsAccountType.AccountTypeExpense,
    balance: 5230.5,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2026-04-10T12:30:00Z",
  },
  {
    id: "acc-002",
    name: "Savings Account",
    type: ModelsAccountType.AccountTypeExpense,
    balance: 15000.0,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2026-04-10T12:30:00Z",
  },
  {
    id: "acc-003",
    name: "Salary Income",
    type: ModelsAccountType.AccountTypeIncome,
    balance: 8500.0,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2026-04-01T09:00:00Z",
  },
];

export const mockCategories: ModelsCategory[] = [
  {
    id: "cat-001",
    name: "Groceries",
    type: ModelsCategoryType.CategoryTypeExpense,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2025-01-15T08:00:00Z",
  },
  {
    id: "cat-002",
    name: "Rent",
    type: ModelsCategoryType.CategoryTypeExpense,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2025-01-15T08:00:00Z",
  },
  {
    id: "cat-003",
    name: "Salary",
    type: ModelsCategoryType.CategoryTypeIncome,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2025-01-15T08:00:00Z",
  },
  {
    id: "cat-004",
    name: "Transfer",
    type: ModelsCategoryType.CategoryTypeTransfer,
    created_at: "2025-01-15T08:00:00Z",
    updated_at: "2025-01-15T08:00:00Z",
  },
];

export const mockTransactions: ModelsTransaction[] = [
  {
    id: "txn-001",
    account_id: "acc-001",
    category_id: "cat-001",
    amount: 150.75,
    currency: "USD",
    title: "Weekly Groceries",
    transaction_type: ModelsTransactionType.TransactionTypeExpense,
    created_at: "2026-04-05T10:00:00Z",
    updated_at: "2026-04-05T10:00:00Z",
  },
  {
    id: "txn-002",
    account_id: "acc-003",
    category_id: "cat-003",
    amount: 5000.0,
    currency: "USD",
    title: "Monthly Salary",
    transaction_type: ModelsTransactionType.TransactionTypeIncome,
    created_at: "2026-04-01T09:00:00Z",
    updated_at: "2026-04-01T09:00:00Z",
  },
  {
    id: "txn-003",
    account_id: "acc-001",
    category_id: "cat-002",
    amount: 1200.0,
    currency: "USD",
    title: "Monthly Rent",
    transaction_type: ModelsTransactionType.TransactionTypeExpense,
    created_at: "2026-04-02T08:00:00Z",
    updated_at: "2026-04-02T08:00:00Z",
  },
];

export const mockDrafts: ModelsDraft[] = [
  {
    id: "dft-001",
    account_id: "acc-001",
    category_id: "cat-001",
    amount: 89.99,
    currency: "USD",
    title: "Amazon Purchase",
    transaction_type: "expense",
    source: "manual",
    status: "pending",
    created_at: "2026-04-15T14:00:00Z",
    updated_at: "2026-04-15T14:00:00Z",
  },
  {
    id: "dft-002",
    account_id: "acc-001",
    category_id: "cat-004",
    amount: 500.0,
    currency: "USD",
    title: "Savings Transfer",
    transaction_type: "transfer",
    transfer_account_id: "acc-002",
    source: "manual",
    status: "pending",
    created_at: "2026-04-15T15:00:00Z",
    updated_at: "2026-04-15T15:00:00Z",
  },
];

export const mockReportSummary: ModelsReportSummary = {
  total_balance: 23730.5,
  total_income: 8500.0,
  total_expenses: 5230.5,
  total_transfers: 15000.0,
  base_currency: "USD",
  period_start: "2026-03-01",
  period_end: "2026-04-18",
};

export const mockReportByCategory: ModelsReportByCategory = {
  period_start: "2026-03-01",
  period_end: "2026-04-18",
  categories: [
    { category_id: "cat-001", category_name: "Groceries", total: 650.0, type: "expense" },
    { category_id: "cat-002", category_name: "Rent", total: 2400.0, type: "expense" },
    { category_id: "cat-003", category_name: "Salary", total: 8500.0, type: "income" },
  ],
};

export const mockReportByAccount: ModelsReportByAccount = {
  period_start: "2026-03-01",
  period_end: "2026-04-18",
  accounts: [
    { account_id: "acc-001", account_name: "Checking Account", total: 1350.75, type: "expense" },
    { account_id: "acc-003", account_name: "Salary Income", total: 8500.0, type: "income" },
  ],
};

export const mockReportTrends: ModelsReportTrends = {
  period_start: "2026-03-01",
  period_end: "2026-04-18",
  data_points: [
    { date: "2026-03-01", total: 5000.0, type: "income" },
    { date: "2026-03-15", total: -350.0, type: "expense" },
    { date: "2026-04-01", total: 8500.0, type: "income" },
    { date: "2026-04-10", total: -150.75, type: "expense" },
  ],
};

export const mockHealth: HandlerHealthResponse = {
  status: "ok",
  timestamp: new Date().toISOString(),
  version: "1.0.0",
};
