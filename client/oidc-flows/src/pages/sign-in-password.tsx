import { LoadingButton } from '@mui/lab';
import {
  Box,
  FormControl,
  Link,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api, ApiError, isApiError } from '../api';
import {
  LoginModelsLoginPasswordErrorResponse,
  LoginModelsLoginPasswordRequest,
} from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { AuthSocialButtons } from '../components/auth/AuthSocialButtons/AuthSocialButtons';
import { RoutePaths } from '../constants/routes';
import { useManifest } from '../contexts/ManifestContext/ManifestContext';
import { useNotification } from '../contexts/NotificationContext/NotificationContext';
import { PageProps } from '../types';
import { withPreventDefault } from '../utils/links';

type ErrorReason = 'wrong_password' | 'unknown';

const ErrorReasons: Record<ErrorReason, string> = {
  wrong_password: 'The password you entered is incorrect.',
  unknown: 'Something went wrong. Please try again.',
};

export const SignInPasswordPage: React.FC<PageProps<{ email: string }>> = ({
  pageProps,
  onNavigate,
}) => {
  const manifest = useManifest();
  const { showNotification } = useNotification();

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsLoginPasswordRequest>({
    defaultValues: {
      email: pageProps?.email,
      password: '',
    },
  });

  const { mutateAsync, isLoading } = useMutation(api.loginPasswordCreate, {
    onSuccess: ({ data }) => {
      if (data.directive === 'displayVerifyCodePage') {
        return onNavigate(RoutePaths.VerifyCode, {
          email: data.email,
          code: data.directiveEmailCodeChallenge?.code,
        });
      }
    },
  });

  async function onSubmit(values: LoginModelsLoginPasswordRequest) {
    try {
      await mutateAsync(values);
    } catch (error) {
      if (isApiError(error)) {
        const apiError =
          error as ApiError<LoginModelsLoginPasswordErrorResponse>;

        const reason: ErrorReason =
          (apiError.response?.data.reason as ErrorReason) || 'unknown';

        return showNotification(ErrorReasons[reason], 'error');
      }

      showNotification(ErrorReasons.unknown, 'error');
    }
  }

  return (
    <AuthLayout title="Sign in">
      <Box component="form" onSubmit={handleSubmit(onSubmit)}>
        <FormControl>
          <TextField
            {...register('email', { required: 'You must enter your email.' })}
            error={getFieldState('email').invalid}
            helperText={errors.email?.message}
            label="Email address"
            placeholder="Enter your email"
            InputProps={{
              readOnly: true,
            }}
          />
        </FormControl>
        <FormControl>
          <TextField
            {...register('password', {
              required: 'You must enter your password.',
            })}
            error={getFieldState('password').invalid}
            helperText={errors.password?.message}
            label="Password"
            type="password"
            placeholder="Enter your password"
          />
        </FormControl>
        <Link
          href="#"
          onClick={withPreventDefault(() =>
            onNavigate(RoutePaths.ForgotPassword),
          )}
        >
          Forgot Password?
        </Link>
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography>Sign in with socials</Typography>
              <AuthSocialButtons />
              {manifest.data?.passkey_enabled && (
                <>
                  <Typography textAlign="center">or</Typography>
                  <Link
                    href="#"
                    onClick={withPreventDefault(() =>
                      onNavigate(RoutePaths.SignInPassKeys, pageProps),
                    )}
                  >
                    Sign in with pass keys
                  </Link>
                </>
              )}
            </Stack>
            <Stack
              direction="row"
              spacing={2}
              sx={{ marginLeft: 'auto', alignItems: 'center' }}
            >
              <Link
                href="#"
                onClick={withPreventDefault(() =>
                  onNavigate(RoutePaths.SignUp),
                )}
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
