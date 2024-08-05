import { LoadingButton } from '@mui/lab';
import { Box, FormControl, Link, Stack, TextField } from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../api';
import { LoginModelsPasswordResetFinishRequest } from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { RoutePaths } from '../constants/routes';
import { useNotification } from '../contexts/NotificationContext/NotificationContext';
import { PageProps } from '../types';

export const ResetPasswordPage: React.FC<PageProps> = ({ onNavigate }) => {
  const { showNotification } = useNotification();

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsPasswordResetFinishRequest>();

  const { mutateAsync, isLoading } = useMutation(
    (values: LoginModelsPasswordResetFinishRequest) =>
      api.passwordResetFinishCreate(values, {
        withCredentials: true,
        withXSRFToken: true,
      }),
    {
      onSuccess: () => {
        window.location.href = RoutePaths.Profile;
      },
    },
  );

  async function onSubmit(values: LoginModelsPasswordResetFinishRequest) {
    try {
      await mutateAsync(values);
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  }

  return (
    <AuthLayout title="Reset password">
      <Box component="form" onSubmit={handleSubmit(onSubmit)}>
        <FormControl>
          <TextField
            {...register('password', {
              required: 'You must enter your password.',
            })}
            type="password"
            error={getFieldState('password').invalid}
            helperText={errors.password?.message}
            label="Password"
            placeholder="Enter your password"
          />
        </FormControl>
        <FormControl>
          <TextField
            {...register('passwordConfirm', {
              required: 'You must enter your password confirmation.',
            })}
            type="password"
            error={getFieldState('passwordConfirm').invalid}
            helperText={errors.passwordConfirm?.message}
            label="Confirm Password"
            placeholder="Enter your password"
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
