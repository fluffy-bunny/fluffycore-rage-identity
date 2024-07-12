import { LoadingButton } from '@mui/lab';
import { Typography } from '@mui/material';
import { useContext } from 'react';
import { useMutation } from 'react-query';

import { api } from '../../api';
import { RoutePaths } from '../../constants/routes';
import { AppContext } from '../../contexts/AppContext/AppContext';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { UserContext } from '../../contexts/UserContext/UserContext';
import { AppType, PageProps } from '../../types';

export const UserProfileSecuritySettingsPage: React.FC<PageProps> = ({
  onNavigate,
}) => {
  const { user } = useContext(UserContext);
  const { setApp } = useContext(AppContext);
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
          onNavigate(RoutePaths.VerifyCode, {
            code: data.data.directiveEmailCodeChallenge?.code,
          });
          setApp(AppType.Auth);
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
    <>
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
    </>
  );
};
