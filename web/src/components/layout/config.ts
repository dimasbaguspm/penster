export interface NavItem {
  label: string;
  to: string;
  icon: string;
  badge?: string | number;
}

export const defaultNavItems: NavItem[] = [
  { label: "Dashboard", to: "/", icon: "layout-dashboard" },
  { label: "Accounts", to: "/accounts", icon: "wallet" },
  { label: "Transactions", to: "/transactions", icon: "arrow-left-right" },
  { label: "Drafts", to: "/drafts", icon: "file-text", badge: 0 },
  { label: "Reports", to: "/reports", icon: "bar-chart-2" },
];