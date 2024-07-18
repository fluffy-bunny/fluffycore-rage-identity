import { LoadingButton } from '@mui/lab';
import { FormControl, TextField, Typography } from '@mui/material';
import { useContext } from 'react';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../../api';
import { ApiUserProfileProfile } from '../../api/Api';
import { MainLayout } from '../../components/MainLayout/MainLayout';
import { ProfileLayout } from '../../components/profile/ProfileLayout/ProfileLayout';
import { RoutePaths } from '../../constants/routes';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { UserContext } from '../../contexts/UserContext/UserContext';
import { PageProps } from '../../types';

export const UserProfilePersonalInformationPage: React.FC<PageProps> = ({
  onNavigate,
}) => {
  const { user, refetch } = useContext(UserContext);
  const { showNotification } = useNotification();

  const { mutateAsync, isLoading } = useMutation(api.userProfileCreate, {
    onSuccess: () => {
      refetch?.();
      showNotification('Personal information has been updated.', 'success');
    },
  });

  const {
    formState: { errors },
    register,
    getFieldState,
    handleSubmit,
  } = useForm<ApiUserProfileProfile>({
    defaultValues: {
      givenName: user?.givenName ?? '',
      familyName: user?.familyName ?? '',
      phoneNumber: user?.phoneNumber ?? '',
    },
  });

  const onSubmit = async (values: ApiUserProfileProfile) => {
    try {
      await mutateAsync(values);
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  };

  return (
    <MainLayout>
      <ProfileLayout
        currentPage={RoutePaths.ProfilePersonalInformation}
        onNavigate={(route) => onNavigate(route)}
      >
        <Typography variant="h4" component="h1" gutterBottom>
          Personal information
        </Typography>
        <form onSubmit={handleSubmit(onSubmit)}>
          <FormControl>
            <TextField
              label="Email"
              defaultValue={user?.email}
              InputProps={{ readOnly: true }}
            />
          </FormControl>
          <FormControl>
            <TextField
              {...register('givenName')}
              error={getFieldState('givenName').invalid}
              helperText={errors.givenName?.message}
              label="Given name"
              placeholder="Enter your given name"
            />
          </FormControl>
          <FormControl>
            <TextField
              {...register('familyName')}
              error={getFieldState('familyName').invalid}
              helperText={errors.familyName?.message}
              label="Family name"
              placeholder="Enter your family name"
            />
          </FormControl>
          <FormControl>
            <TextField
              {...register('phoneNumber')}
              error={getFieldState('phoneNumber').invalid}
              helperText={errors.phoneNumber?.message}
              label="Phone number"
              placeholder="Enter your phone number"
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
      </ProfileLayout>
    </MainLayout>
  );
};
