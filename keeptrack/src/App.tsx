import React from "react";
import { Route, Routes } from "react-router-dom";
import { RoutePaths } from "./constants/routes";
import { AuthLayout } from "./components/auth/AuthLayout/AuthLayout";

const SignInPage = React.lazy(() =>
  import("./pages/sign-in").then((module) => ({ default: module.Page }))
);

const SignUpPage = React.lazy(() =>
  import("./pages/sign-up").then((module) => ({ default: module.Page }))
);

function App() {
  return (
    <Routes>
      <Route path={RoutePaths.SignIn} element={<AuthLayout />}>
        <Route index element={<SignInPage />} />
      </Route>
      <Route path={RoutePaths.SignUp} element={<AuthLayout />}>
        <Route index element={<SignUpPage />} />
      </Route>
    </Routes>
  );
}

export default App;
