import { LoadingButton } from '@mui/lab';
import { Box, FormControl, Typography } from '@mui/material';
import { useMutation } from 'react-query';

import { MainLayout } from '../../components/MainLayout/MainLayout';
import { ProfileLayout } from '../../components/profile/ProfileLayout/ProfileLayout';
import { RoutePaths } from '../../constants/routes';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { usePublicKeyCredentialSupport } from '../../hooks/usePublicKeyCredentialSupport';
import { PageProps } from '../../types';
import { registerUser } from '../../utils/webauthn';

export const UserProfilePasskeysManagementPage: React.FC<PageProps> = ({
  onNavigate,
}) => {
  const isSupported = usePublicKeyCredentialSupport();
  const { showNotification } = useNotification();

  const { mutateAsync, isLoading } = useMutation(() =>
    registerUser(RoutePaths.ProfilePasskeysManagement),
  );

  const onRegister = async () => {
    try {
      await mutateAsync();
    } catch (error) {
      showNotification('Something went wrong. Please try again.', 'error');
    }
  };

  const renderSupportMessage = () => {
    if (isSupported === null) {
      return (
        <Typography paragraph>
          Checking support for PublicKeyCredential...
        </Typography>
      );
    }

    if (!isSupported) {
      return (
        <Typography paragraph color="red">
          PublicKeyCredential is not supported in this browser.
        </Typography>
      );
    }

    if (isSupported) {
      return (
        <Typography paragraph>
          PublicKeyCredential is supported in this browser.
        </Typography>
      );
    }
  };

  return (
    <MainLayout
      currentPage={RoutePaths.ProfilePersonalInformation}
      onNavigate={onNavigate}
    >
      <ProfileLayout
        currentPage={RoutePaths.ProfilePasskeysManagement}
        onNavigate={(route) => onNavigate(route)}
      >
        <Typography variant="h4" component="h1" gutterBottom>
          Manage Pass Keys
        </Typography>
        {renderSupportMessage()}
        <FormControl>
          <Box>
            <LoadingButton
              variant="contained"
              loading={isLoading}
              onClick={onRegister}
            >
              Register
            </LoadingButton>
          </Box>
        </FormControl>
      </ProfileLayout>
    </MainLayout>
  );
};
