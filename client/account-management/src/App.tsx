import React, { useState } from 'react';

import { RoutePaths } from './constants/routes';
import { ManifestProvider } from './contexts/ManifestContext/ManifestContext';
import { UserProvider } from './contexts/UserContext/UserContext';
import { HomePage } from './pages';
import { UserProfilePasskeysManagementPage } from './pages/profile/passkeys-management';
import { UserProfilePersonalInformationPage } from './pages/profile/personal-information';
import { UserProfileSecuritySettingsPage } from './pages/profile/security-settings';
import { ResetPasswordPage } from './pages/reset-password';
import { VerifyCodePage } from './pages/verify-code';
import { AppRoute, PageProps } from './types';

const pages: Record<
  AppRoute | 'default',
  React.FunctionComponent<PageProps>
> = {
  default: HomePage,
  [RoutePaths.Root]: HomePage,
  [RoutePaths.VerifyCode]: VerifyCodePage,
  [RoutePaths.ResetPassword]: ResetPasswordPage,

  [RoutePaths.ProfilePersonalInformation]: UserProfilePersonalInformationPage,
  [RoutePaths.ProfileSecuritySettings]: UserProfileSecuritySettingsPage,
  [RoutePaths.ProfilePasskeysManagement]: UserProfilePasskeysManagementPage,
};

export function App() {
  const [currentPageState, setCurrentPageState] = useState<{
    route: string;
    pageProps: any;
  }>({
    route: 'default',
    pageProps: undefined,
  });

  const PageComponent = pages[currentPageState.route];

  if (!PageComponent) {
    return <div>Page not found</div>;
  }

  return (
    <ManifestProvider>
      <UserProvider>
        <PageComponent
          currentPage={currentPageState.route}
          pageProps={currentPageState.pageProps}
          onNavigate={(route, pageProps) =>
            setCurrentPageState({ route, pageProps })
          }
        />
      </UserProvider>
    </ManifestProvider>
  );
}
