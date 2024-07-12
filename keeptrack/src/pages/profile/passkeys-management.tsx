import { LoadingButton } from '@mui/lab';
import { Box, FormControl, Typography } from '@mui/material';
import { useMutation } from 'react-query';

// import { RoutePaths } from '../../constants/routes';
import { useNotification } from '../../contexts/NotificationContext/NotificationContext';
import { usePublicKeyCredentialSupport } from '../../hooks/usePublicKeyCredentialSupport';
import { PageProps } from '../../types';
import { webauthn } from '../../utils/webauthn';

export const UserProfilePasskeysManagementPage: React.FC<PageProps> = () => {
  const isSupported = usePublicKeyCredentialSupport();
  const { showNotification } = useNotification();

  const { mutateAsync, isLoading } = useMutation(() => webauthn.registerUser());

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
    <>
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
    </>
  );
};
