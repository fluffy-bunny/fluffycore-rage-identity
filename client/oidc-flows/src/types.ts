import { RoutePaths } from './constants/routes';

export type AppRoute = (typeof RoutePaths)[keyof typeof RoutePaths];

export interface PageProps<T = any> {
  pageProps?: T | undefined;
  onNavigate<T>(route: AppRoute, pageProps?: T): void;
}

export enum AppType {
  Auth = 'auth',
  Profile = 'profile',
}
