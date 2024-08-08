import { AppBar, Box, Link, Stack, Toolbar } from '@mui/material';
import React from 'react';

import { RoutePaths } from '../../constants/routes';
import { ProfileDropdown } from '../profile/ProfileDropdown/ProfileDropdown';

const navItems = [
  {
    label: 'Mapped Account',
    path: RoutePaths.Root,
  },
  {
    label: 'Profile',
    path: RoutePaths.ProfilePersonalInformation,
  },
];

export const MainLayout = ({
  children,
  currentPage,
  onNavigate,
}: {
  children: React.ReactNode;
  currentPage?: string | 'default';
  onNavigate: (path: string) => void;
}) => {
  return (
    <>
      <Stack sx={{ height: '100%' }}>
        <AppBar
          position="static"
          sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
        >
          <Toolbar sx={{ width: '100%' }}>
            <Stack component="nav" direction="row" spacing={2}>
              {navItems.map((nav, index) => {
                const isActive =
                  currentPage === 'default'
                    ? index === 0
                    : nav.path === currentPage;

                return (
                  <Link
                    key={nav.label}
                    color="inherit"
                    component="button"
                    sx={{ textDecoration: isActive ? 'underline' : 'none' }}
                    onClick={() => onNavigate(nav.path)}
                  >
                    {nav.label}
                  </Link>
                );
              })}
            </Stack>
            <Box sx={{ marginLeft: 'auto' }}>
              <ProfileDropdown />
            </Box>
          </Toolbar>
        </AppBar>
        <Box
          component="main"
          sx={{ flexGrow: 1, background: 'background.default', p: 3 }}
        >
          {children}
        </Box>
      </Stack>
    </>
  );
};
