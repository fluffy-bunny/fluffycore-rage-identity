import { Box, CircularProgress } from '@mui/material';
import { useEffect } from 'react';

import { RoutePaths } from '../constants/routes';
import { PageProps } from '../types';
import { loginUser } from '../utils/webauth';

export const SignInPasskeysPage: React.FC<
  PageProps<{ email: string; hasPasskey: boolean }>
> = ({ onNavigate, pageProps }) => {
  useEffect(() => {
    const run = async () => {
      try {
        await loginUser(() => (window.location.href = RoutePaths.Root));
      } catch (error) {
        return onNavigate(RoutePaths.SignInPassword, pageProps);
      }
    };

    run();
  }, []);

  return (
    <Box sx={{ height: '100%', display: 'flex' }}>
      <CircularProgress sx={{ margin: 'auto' }} />
    </Box>
  );
};
