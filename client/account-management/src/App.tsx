import React, { useState } from 'react';

import { RoutePaths } from './constants/routes';
import { UserProvider } from './contexts/UserContext/UserContext';
import { ForgotPasswordPage } from './pages/forgot-password';
import { UserProfilePasskeysManagementPage } from './pages/profile/passkeys-management';
import { UserProfilePersonalInformationPage } from './pages/profile/personal-information';
import { UserProfileSecuritySettingsPage } from './pages/profile/security-settings';
import { ResetPasswordPage } from './pages/reset-password';
import { SignInPage } from './pages/sign-in';
import { SignInPasswordPage } from './pages/sign-in-password';
import { SignUpPage } from './pages/sign-up';
import { VerifyCodePage } from './pages/verify-code';
import { AppRoute, PageProps } from './types';

const pages: Record<
  AppRoute | 'default',
  React.FunctionComponent<PageProps>
> = {
  default: UserProfilePersonalInformationPage,
  [RoutePaths.SignIn]: SignInPage,
  [RoutePaths.SignInPassword]: SignInPasswordPage,
  [RoutePaths.SignUp]: SignUpPage,
  [RoutePaths.VerifyCode]: VerifyCodePage,
  [RoutePaths.ForgotPassword]: ForgotPasswordPage,
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
    <UserProvider>
      <PageComponent
        pageProps={currentPageState.pageProps}
        onNavigate={(route, pageProps) =>
          setCurrentPageState({ route, pageProps })
        }
      />
    </UserProvider>
  );
}
