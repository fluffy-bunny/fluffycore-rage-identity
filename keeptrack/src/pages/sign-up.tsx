import { useForm } from "react-hook-form";
import {
  Box,
  Button,
  FormControl,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation } from "react-query";
import { api, externalIdp } from "../api";
import { LoadingButton } from "@mui/lab";
import { LoginModelsSignupRequest } from "../api/Api";
import { AuthLayout } from "../components/auth/AuthLayout/AuthLayout";
import { AuthSocialButtons } from "../components/auth/AuthSocialButtons/AuthSocialButtons";
import { getCSRF } from "../utils/cookies";
import { RoutePaths } from "../constants/routes";

export const SignUpPage = ({
  onNavigate,
}: {
  onNavigate(route: string): void;
}) => {
  const { register, handleSubmit } = useForm<LoginModelsSignupRequest>();
  const { mutateAsync, isLoading } = useMutation(
    async (values: LoginModelsSignupRequest) => {
      const response = await api.signupCreate(values, {
        withCredentials: true,
        withXSRFToken: true,
      });

      await externalIdp.externalIdpCreate(
        {
          ...response.data.directiveFormPost?.formParams,
          // @ts-ignore
          csrf: getCSRF(),
        },
        {
          withCredentials: true,
          withXSRFToken: true,
        }
      );
    }
  );

  return (
    <AuthLayout>
      <Typography variant="h4" component="h1" gutterBottom>
        Sign up
      </Typography>
      <Box
        component="form"
        onSubmit={handleSubmit((values) => mutateAsync(values))}
      >
        <FormControl>
          <TextField
            {...register("email")}
            label="Email address"
            placeholder="Enter your email"
          />
        </FormControl>
        <FormControl>
          <TextField
            {...register("password")}
            label="Password"
            placeholder="Enter your password"
          />
        </FormControl>
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography>Sign in with socials</Typography>
              <AuthSocialButtons />
            </Stack>
            <Stack direction="row" spacing={1} sx={{ marginLeft: "auto" }}>
              <Button onClick={() => onNavigate(RoutePaths.SignIn)}>
                Sign In
              </Button>
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
