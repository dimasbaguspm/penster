export interface SideBarProps {
  open?: boolean;
  navItems?: Array<{
    label: string;
    to: string;
    icon?: string;
    badge?: string | number;
  }>;
  sections?: Array<{
    title?: string;
    items: SideBarProps['navItems'];
  }>;
}

export interface SideBarEmits {
  (e: 'update:open', value: boolean): void;
}