import { Stack, Typography } from '@mui/material';

import { MainLayout } from '../components/MainLayout/MainLayout';
import { PageProps } from '../types';

export const HomePage: React.FC<PageProps> = ({ currentPage, onNavigate }) => {
  return (
    <MainLayout currentPage={currentPage} onNavigate={onNavigate}>
      <Stack direction="row" spacing={4} sx={{ alignItems: 'center' }}>
        <img src={`${process.env.PUBLIC_URL}/assets/logo.svg`} alt="Logo" />
        <Typography variant="h4">Data AI for the real world</Typography>
      </Stack>
    </MainLayout>
  );
};
