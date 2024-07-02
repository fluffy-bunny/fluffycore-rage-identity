import { LoadingButton } from '@mui/lab';
import { FormControl, Link, Stack, TextField, Typography } from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../api';
import { LoginModelsLoginPhaseOneRequest } from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { AuthSocialButtons } from '../components/auth/AuthSocialButtons/AuthSocialButtons';
import { RoutePaths } from '../constants/routes';
import { useNotification } from '../contexts/NotificationContext/NotificationContext';
import { useExternalLogin } from '../hooks/useExternalLogin';
import { PageProps } from '../types';

export const SignInPage: React.FC<PageProps> = ({ onNavigate }) => {
  const { showNotification } = useNotification();
  const [executeExternalLogin] = useExternalLogin();

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsLoginPhaseOneRequest>();

  const { mutateAsync, isLoading } = useMutation(
    async (values: LoginModelsLoginPhaseOneRequest) => {
      const { data: loginPhaseOneData } = await api.loginPhaseOneCreate(
        values,
        { withCredentials: true, withXSRFToken: true },
      );

      if (loginPhaseOneData.directive === 'displayPasswordPage') {
        return onNavigate(RoutePaths.SignInPassword, {
          email: loginPhaseOneData.email,
        });
      }

      if (loginPhaseOneData.directive === 'startExternalLogin') {
        return executeExternalLogin(
          loginPhaseOneData.directiveStartExternalLogin?.slug,
        );
      }
    },
  );

  async function onSubmit(values: { email: string; password?: string }) {
    try {
      await mutateAsync(values);
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  }

  return (
    <AuthLayout title="Sign in">
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormControl>
          <TextField
            {...register('email', { required: 'You must enter your email.' })}
            error={getFieldState('email').invalid}
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
              sx={{ marginLeft: 'auto', alignItems: 'center' }}
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
      </form>
    </AuthLayout>
  );
};
