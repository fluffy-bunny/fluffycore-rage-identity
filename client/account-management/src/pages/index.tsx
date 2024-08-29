import logo from '../assets/logo.svg';
import { MainLayout } from '../components/MainLayout/MainLayout';
import { PageProps } from '../types';

export const HomePage: React.FC<PageProps> = ({ currentPage, onNavigate }) => {
  return (
    <MainLayout currentPage={currentPage} onNavigate={onNavigate}>
      <img src={logo} alt="Logo" />
    </MainLayout>
  );
};
