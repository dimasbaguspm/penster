// AppLayout
export { default as AppLayout } from './app-layout/app-layout.vue';
export type { AppLayoutProps, AppLayoutEmits } from './app-layout/types';

// TopBar
export { default as TopBar } from './top-bar/top-bar.vue';
export type { TopBarProps } from './top-bar/types';

// SideBar
export { default as SideBar } from './side-bar/side-bar.vue';
export type { SideBarProps, SideBarEmits } from './side-bar/types';

// BottomBar
export { default as BottomBar } from './bottom-bar/bottom-bar.vue';
export type { BottomBarProps } from './bottom-bar/types';

// Config
export { defaultNavItems } from './config';
export type { NavItem } from './config';