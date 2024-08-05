import { LoadingButton } from '@mui/lab';
import {
  Button,
  FormControl,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useContext, useState } from 'react';
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
  const [isEditEnabled, setIsEditEnabled] = useState(false);

  return (
    <MainLayout
      currentPage={RoutePaths.ProfilePersonalInformation}
      onNavigate={onNavigate}
    >
      <ProfileLayout
        currentPage={RoutePaths.ProfilePersonalInformation}
        onNavigate={(route) => onNavigate(route)}
      >
        <Stack direction="row" sx={{ alignItems: 'center', marginBottom: 1.5 }}>
          <Typography variant="h4" component="h1" gutterBottom={false}>
            Personal information
          </Typography>
          {!isEditEnabled && (
            <Button
              sx={{ marginLeft: 'auto' }}
              onClick={() => setIsEditEnabled(true)}
            >
              Edit
            </Button>
          )}
        </Stack>
        {isEditEnabled ? (
          <PersonalInformationForm onCancel={() => setIsEditEnabled(false)} />
        ) : (
          <PersonalInformation />
        )}
      </ProfileLayout>
    </MainLayout>
  );
};

const PersonalInformationForm = ({ onCancel }: { onCancel(): void }) => {
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
        <Stack direction="row" spacing={2} sx={{ marginLeft: 'auto' }}>
          <Button disabled={isLoading} onClick={onCancel}>
            Cancel
          </Button>
          <LoadingButton type="submit" variant="contained" loading={isLoading}>
            Save
          </LoadingButton>
        </Stack>
      </FormControl>
    </form>
  );
};

const PersonalInformation = () => {
  const { user } = useContext(UserContext);

  const items = [
    {
      label: 'Given name',
      value: user?.givenName,
    },
    {
      label: 'Family name',
      value: user?.familyName,
    },
    {
      label: 'Email',
      value: user?.email,
    },
    {
      label: 'Phone number',
      value: user?.phoneNumber,
    },
  ];

  return (
    <>
      {items.map((item) => (
        <FormControl key={item.label}>
          <Typography color="GrayText">{item.label}</Typography>
          <Typography>{item.value || '-'}</Typography>
        </FormControl>
      ))}
    </>
  );
};
