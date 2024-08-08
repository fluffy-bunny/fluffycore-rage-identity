import React, { useState } from 'react';

import { RoutePaths } from './constants/routes';
import { ForgotPasswordPage } from './pages/forgot-password';
import { ResetPasswordPage } from './pages/reset-password';
import { SignInPage } from './pages/sign-in';
import { SignInPasskeysPage } from './pages/sign-in-passkeys';
import { SignInPasswordPage } from './pages/sign-in-password';
import { SignUpPage } from './pages/sign-up';
import { VerifyCodePage } from './pages/verify-code';
import { AppRoute, PageProps } from './types';

const pages: Record<
  AppRoute | 'default',
  React.FunctionComponent<PageProps>
> = {
  default: SignInPage,
  [RoutePaths.SignIn]: SignInPage,
  [RoutePaths.SignInPassword]: SignInPasswordPage,
  [RoutePaths.SignInPassKeys]: SignInPasskeysPage,
  [RoutePaths.SignUp]: SignUpPage,
  [RoutePaths.VerifyCode]: VerifyCodePage,
  [RoutePaths.ForgotPassword]: ForgotPasswordPage,
  [RoutePaths.ResetPassword]: ResetPasswordPage,
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
    <PageComponent
      pageProps={currentPageState.pageProps}
      onNavigate={(route, pageProps) =>
        setCurrentPageState({ route, pageProps })
      }
    />
  );
}
