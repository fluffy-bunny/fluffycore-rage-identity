import { LoadingButton } from '@mui/lab';
import { FormControl, TextField, Typography } from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../../api';
import { LoginModelsPasswordResetFinishRequest } from '../../api/Api';
import { RoutePaths } from '../../constants/routes';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { PageProps } from '../../types';

export const UserProfileSecuritySettingsPage: React.FC<PageProps> = ({
  onNavigate,
}) => {
  const { showNotification } = useNotification();

  const { mutateAsync, isLoading } = useMutation(
    (values: LoginModelsPasswordResetFinishRequest) =>
      api.passwordResetFinishCreate(values, {
        withCredentials: true,
        withXSRFToken: true,
      }),
    {
      onSuccess: (data) => {
        if (data.data.directive === 'displayLoginPhaseOnePage') {
          return onNavigate(RoutePaths.SignIn);
        }
      },
    },
  );

  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsPasswordResetFinishRequest>();

  const onSubmit = async (values: LoginModelsPasswordResetFinishRequest) => {
    try {
      await mutateAsync(values);
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  };

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Reset Password
      </Typography>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FormControl>
          <TextField
            {...register('password', {
              required: 'You must enter your password.',
            })}
            error={getFieldState('password').invalid}
            helperText={errors.password?.message}
            label="Password"
            type="password"
          />
        </FormControl>
        <FormControl>
          <TextField
            {...register('passwordConfirm', {
              required: 'You must enter your password confirmation.',
            })}
            error={getFieldState('passwordConfirm').invalid}
            helperText={errors.passwordConfirm?.message}
            label="Confirm password"
            type="password"
          />
        </FormControl>
        <FormControl>
          <LoadingButton
            type="submit"
            variant="contained"
            loading={isLoading}
            sx={{ marginLeft: 'auto' }}
          >
            Save
          </LoadingButton>
        </FormControl>
      </form>
    </>
  );
};
