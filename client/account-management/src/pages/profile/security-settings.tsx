import { LoadingButton } from '@mui/lab';
import { Typography } from '@mui/material';
import { useContext } from 'react';
import { useMutation } from 'react-query';

import { api } from '../../api';
import { MainLayout } from '../../components/MainLayout/MainLayout';
import { ProfileLayout } from '../../components/profile/ProfileLayout/ProfileLayout';
import { RoutePaths } from '../../constants/routes';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { UserContext } from '../../contexts/UserContext/UserContext';
import { PageProps } from '../../types';

export const UserProfileSecuritySettingsPage: React.FC<PageProps> = ({
  onNavigate,
}) => {
  const { user } = useContext(UserContext);
  const { showNotification } = useNotification();

  const { mutateAsync, isLoading } = useMutation(
    () =>
      api.passwordResetStartCreate(
        {
          email: user?.email!,
        },
        {
          withCredentials: true,
          withXSRFToken: true,
        },
      ),
    {
      onSuccess: (data) => {
        if (data.data.directive === 'displayVerifyCodePage') {
          return onNavigate(RoutePaths.VerifyCode, {
            code: data.data.directiveEmailCodeChallenge?.code,
          });
        }
      },
    },
  );

  const onResetPassowrd = async () => {
    try {
      await mutateAsync();
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  };

  return (
    <MainLayout>
      <ProfileLayout
        currentPage={RoutePaths.ProfileSecuritySettings}
        onNavigate={(route) => onNavigate(route)}
      >
        <Typography variant="h4" component="h1" gutterBottom>
          Reset Password
        </Typography>
        <LoadingButton
          variant="contained"
          loading={isLoading}
          onClick={onResetPassowrd}
        >
          Reset
        </LoadingButton>
      </ProfileLayout>
    </MainLayout>
  );
};
