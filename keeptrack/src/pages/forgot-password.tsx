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
import { LoadingButton } from "@mui/lab";
import { forgotPassword } from "../api";
import { AuthLayout } from "../components/auth/AuthLayout/AuthLayout";
import { RoutePaths } from "../constants/routes";
import { getCSRF } from "../utils/cookies";

export const ForgotPasswordPage = ({
  onNavigate,
}: {
  onNavigate(route: string): void;
}) => {
  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<{ email: string }>();
  const { mutateAsync, isLoading } = useMutation(
    forgotPassword.forgotPasswordCreate
  );

  return (
    <AuthLayout>
      <Typography variant="h4" component="h1" gutterBottom>
        Forgot password
      </Typography>
      <Box
        component="form"
        onSubmit={handleSubmit((values) => {
          // @ts-ignore
          mutateAsync({
            ...values,
            // @ts-ignore
            csrf: getCSRF(),
          });
        })}
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
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack
              spacing={2}
              direction="row"
              sx={{ marginLeft: "auto", alignItems: "center" }}
            >
              <Link
                component="button"
                onClick={() => onNavigate(RoutePaths.SignIn)}
              >
                Cancel
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
