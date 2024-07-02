import React, { useState } from 'react';

import { RoutePaths } from './constants/routes';
import { ForgotPasswordPage } from './pages/forgot-password';
import { ResetPasswordPage } from './pages/reset-password';
import { SignInPage } from './pages/sign-in';
import { SignInPasswordPage } from './pages/sign-in-password';
import { SignUpPage } from './pages/sign-up';
import { VerifyCodePage } from './pages/verify-code';
import { AppRoute, PageProps } from './types';

const pages: Record<AppRoute, React.FunctionComponent<PageProps>> = {
  [RoutePaths.SignIn]: SignInPage,
  [RoutePaths.SignIn]: SignInPage,
  [RoutePaths.SignInPassword]: SignInPasswordPage,
  [RoutePaths.SignUp]: SignUpPage,
  [RoutePaths.VerifyCode]: VerifyCodePage,
  [RoutePaths.ForgotPassword]: ForgotPasswordPage,
  [RoutePaths.ResetPassword]: ResetPasswordPage,
};

export function App() {
  const [currentPageState, setCurrentPageState] = useState<{
    route: AppRoute;
    pageProps: any;
  }>({
    route: RoutePaths.SignIn,
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
