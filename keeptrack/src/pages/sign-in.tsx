import { useForm } from "react-hook-form";
import {
  Box,
  FormControl,
  Link,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation } from "react-query";
import { api } from "../api";
import { LoadingButton } from "@mui/lab";
import { LoginModelsLoginPhaseOneRequest } from "../api/Api";
import { AuthLayout } from "../components/auth/AuthLayout/AuthLayout";
import { AuthSocialButtons } from "../components/auth/AuthSocialButtons/AuthSocialButtons";
import { RoutePaths } from "../constants/routes";

export const SignInPage = ({
  onNavigate,
}: {
  onNavigate(route: string): void;
}) => {
  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsLoginPhaseOneRequest>();
  const { mutateAsync, isLoading } = useMutation(
    async (values: LoginModelsLoginPhaseOneRequest) => {
      const { data } = await api.loginPhaseOneCreate(values);

      return api.startExternalLoginCreate({
        // @ts-ignore
        slug: data.directiveStartExternalLogin.slug,
        directive: "login",
      });
    },
    {
      onSuccess: (data) => {
        if (data.data.redirectUri) {
          window.location.href = data.data.redirectUri;
        }
      },
    }
  );

  return (
    <AuthLayout>
      <Typography variant="h4" component="h1" gutterBottom>
        Sign in
      </Typography>
      <Box
        component="form"
        onSubmit={handleSubmit((values) => mutateAsync(values))}
      >
        <FormControl>
          <TextField
            {...register("email", { required: "You must enter your email." })}
            error={getFieldState("email").invalid}
            helperText={errors.email?.message}
            label="Email address"
            placeholder="Enter your email"
          />
        </FormControl>
        <Link
          component="button"
          onClick={() => onNavigate(RoutePaths.ForgotPassword)}
        >
          Forgot Password?
        </Link>
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography>Sign in with socials</Typography>
              <AuthSocialButtons />
            </Stack>
            <Stack
              direction="row"
              spacing={2}
              sx={{ marginLeft: "auto", alignItems: "center" }}
            >
              <Link
                component="button"
                onClick={() => onNavigate(RoutePaths.SignUp)}
              >
                Sign Up
              </Link>
              <LoadingButton
                loading={isLoading}
                type="submit"
                variant="contained"
              >
                Next
              </LoadingButton>
            </Stack>
          </Stack>
        </FormControl>
      </Box>
    </AuthLayout>
  );
};
