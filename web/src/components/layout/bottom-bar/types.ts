export interface BottomBarProps {
  items?: Array<{
    label: string;
    to: string;
    icon: string;
    active?: boolean;
  }>;
}