import React, { useState } from "react";
import { RoutePaths } from "./constants/routes";
import { SignInPage } from "./pages/sign-in";
import { SignUpPage } from "./pages/sign-up";

const pages = {
  [RoutePaths.SignIn]: SignInPage,
  [RoutePaths.SignUp]: SignUpPage,
};

export function App() {
  const [currentPage, setCurrentPage] = useState(RoutePaths.SignIn);

  const PageComponent = pages[currentPage];

  return <PageComponent onNavigate={setCurrentPage} />;
}
