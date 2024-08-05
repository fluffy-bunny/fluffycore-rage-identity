import { Typography } from '@mui/material';

import { MainLayout } from '../components/MainLayout/MainLayout';
import { PageProps } from '../types';

export const HomePage: React.FC<PageProps> = ({ currentPage, onNavigate }) => {
  return (
    <MainLayout currentPage={currentPage} onNavigate={onNavigate}>
      <Typography variant="h4" component="h1" gutterBottom>
        Home page
      </Typography>
    </MainLayout>
  );
};
