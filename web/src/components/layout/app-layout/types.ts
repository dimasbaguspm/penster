import type { TopBarProps } from '../top-bar/types';
import type { SideBarProps } from '../side-bar/types';
import type { BottomBarProps } from '../bottom-bar/types';

export interface AppLayoutProps {
  topBar?: Partial<TopBarProps>;
  sideBar?: Partial<SideBarProps> & { enabled?: boolean };
  bottomBar?: Partial<BottomBarProps> & { enabled?: boolean };
  contentMaxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '7xl' | 'full';
  footerEnabled?: boolean;
}

export interface AppLayoutEmits {
  (e: 'mobile-menu-toggle'): void;
}