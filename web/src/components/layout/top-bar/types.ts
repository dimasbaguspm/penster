export interface TopBarProps {
  logo?: { text: string; href?: string };
  navItems?: Array<{
    label: string;
    to: string;
    icon?: string;
  }>;
  showMobileMenu?: boolean;
  sticky?: boolean;
}