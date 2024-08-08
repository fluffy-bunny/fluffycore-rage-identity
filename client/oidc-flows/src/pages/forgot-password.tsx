import { LoadingButton } from '@mui/lab';
import { Box, FormControl, Link, Stack, TextField } from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../api';
import { LoginModelsPasswordResetStartRequest } from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { RoutePaths } from '../constants/routes';
import { useNotification } from '../contexts/NotificationContext/NotificationContext';
import { PageProps } from '../types';
import { withPreventDefault } from '../utils/links';

export const ForgotPasswordPage: React.FC<PageProps> = ({ onNavigate }) => {
  const { showNotification } = useNotification();

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsPasswordResetStartRequest>();

  const { mutateAsync, isLoading } = useMutation(
    (values: LoginModelsPasswordResetStartRequest) =>
      api.passwordResetStartCreate(values),
    {
      onSuccess: (data) => {
        if (data.data.directive === 'displayVerifyCodePage') {
          return onNavigate(RoutePaths.VerifyCode, {
            email: data.data.email,
            code: data.data.directiveEmailCodeChallenge?.code,
          });
        }
      },
    },
  );

  async function onSubmit(values: LoginModelsPasswordResetStartRequest) {
    try {
      await mutateAsync(values);
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  }

  return (
    <AuthLayout title="Forgot password">
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
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack
              spacing={2}
              direction="row"
              sx={{ marginLeft: 'auto', alignItems: 'center' }}
            >
              <Link
                href="#"
                onClick={withPreventDefault(() =>
                  onNavigate(RoutePaths.SignIn),
                )}
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
