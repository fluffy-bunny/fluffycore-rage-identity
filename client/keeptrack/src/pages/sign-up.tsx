import { LoadingButton } from '@mui/lab';
import {
  Box,
  Button,
  FormControl,
  Link,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api, isApiError } from '../api';
import { LoginModelsSignupRequest } from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { AuthSocialButtons } from '../components/auth/AuthSocialButtons/AuthSocialButtons';
import { RoutePaths } from '../constants/routes';
import { useNotification } from '../contexts/NotificationContext/NotificationContext';
import { useExternalLogin } from '../hooks/useExternalLogin';
import { PageProps } from '../types';
import { withPreventDefault } from '../utils/links';

export const SignUpPage: React.FC<PageProps> = ({ onNavigate }) => {
  const { showNotification } = useNotification();
  const [executeExternalLogin] = useExternalLogin();

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsSignupRequest>();

  const { mutateAsync, isLoading } = useMutation(
    async (values: LoginModelsSignupRequest) => {
      const response = await api.signupCreate(values, {
        withCredentials: true,
        withXSRFToken: true,
      });

      if (response.data.directive === 'displayVerifyCodePage') {
        return onNavigate(RoutePaths.VerifyCode, {
          email: response.data.email,
          code: response.data.directiveEmailCodeChallenge?.code,
        });
      }

      if (response.data.directive === 'startExternalLogin') {
        return executeExternalLogin(
          response.data.directiveStartExternalLogin?.slug,
        );
      }
    },
  );

  async function onSubmit(values: LoginModelsSignupRequest) {
    try {
      await mutateAsync(values);
    } catch (error) {
      if (isApiError(error)) {
        const responseData = (error as any).response?.data;
        showNotification(
          responseData?.message || 'Something went wrong. Please try again.',
          'error',
        );
      }
    }
  }

  return (
    <AuthLayout title="Sign up">
      <Box component="form" onSubmit={handleSubmit(onSubmit)}>
        <FormControl>
          <TextField
            {...register('email', { required: 'You must enter your email.' })}
            error={getFieldState('email').invalid}
            helperText={errors.email?.message}
            label="Email address"
            placeholder="Enter your email"
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
            </Stack>
            <Stack direction="row" spacing={1} sx={{ marginLeft: 'auto' }}>
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
