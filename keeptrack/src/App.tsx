import React, { useState } from 'react';

import { MainLayout } from './components/MainLayout/MainLayout';
import { ProfileLayout } from './components/profile/ProfileLayout/ProfileLayout';
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
import { AppRoute, AppType, PageProps } from './types';

const pages: Record<
  AppType,
  Record<AppRoute | 'default', React.FunctionComponent<PageProps>>
> = {
  [AppType.Auth]: {
    default: SignInPage,
    [RoutePaths.SignIn]: SignInPage,
    [RoutePaths.SignInPassword]: SignInPasswordPage,
    [RoutePaths.SignUp]: SignUpPage,
    [RoutePaths.VerifyCode]: VerifyCodePage,
    [RoutePaths.ForgotPassword]: ForgotPasswordPage,
    [RoutePaths.ResetPassword]: ResetPasswordPage,
  },
  [AppType.Profile]: {
    default: UserProfilePersonalInformationPage,
    [RoutePaths.ProfilePersonalInformation]: UserProfilePersonalInformationPage,
    [RoutePaths.ProfileSecuritySettings]: UserProfileSecuritySettingsPage,
    [RoutePaths.ProfilePasskeysManagement]: UserProfilePasskeysManagementPage,
  },
};

export function App({ app }: { app: AppType | null }) {
  const currentApp = app === null ? AppType.Auth : app;
  const [currentPageState, setCurrentPageState] = useState<{
    route: string;
    pageProps: any;
  }>({
    route: 'default',
    pageProps: undefined,
  });

  const PageComponent = pages[currentApp][currentPageState.route];

  const components = {
    [AppType.Auth]: (
      <PageComponent
        pageProps={currentPageState.pageProps}
        onNavigate={(route, pageProps) =>
          setCurrentPageState({ route, pageProps })
        }
      />
    ),
    [AppType.Profile]: (
      <UserProvider>
        <MainLayout>
          <ProfileLayout
            currentPage={currentPageState.route}
            onNavigate={(route) =>
              setCurrentPageState({ route, pageProps: undefined })
            }
          >
            <PageComponent
              pageProps={currentPageState.pageProps}
              onNavigate={(route, pageProps) =>
                setCurrentPageState({ route, pageProps })
              }
            />
          </ProfileLayout>
        </MainLayout>
      </UserProvider>
    ),
  };

  return components[currentApp] || <div>Page not found</div>;
}
