import { AccountBoxOutlined } from '@mui/icons-material';
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
} from '@mui/material';
import React from 'react';

import { ProfileDropdown } from '../profile/ProfileDropdown/ProfileDropdown';

const navItems = [
  {
    label: 'Profile',
    icon: AccountBoxOutlined,
  },
];

const DrawerWidth = 320;

export const MainLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <Stack direction="row" sx={{ height: '100%' }}>
      <Drawer
        sx={{
          width: DrawerWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: DrawerWidth,
            boxSizing: 'border-box',
          },
        }}
        variant="permanent"
        anchor="left"
      >
        <Stack sx={{ height: '100%' }}>
          <List>
            {navItems.map((nav) => (
              <ListItem key={nav.label}>
                <ListItemButton selected={true}>
                  <ListItemIcon>
                    <nav.icon fontSize="small" />
                  </ListItemIcon>
                  <ListItemText primary={nav.label} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Box sx={{ marginTop: 'auto', padding: 2 }}>
            <ProfileDropdown />
          </Box>
        </Stack>
      </Drawer>
      <Box
        component="main"
        sx={{ flexGrow: 1, background: 'background.default', p: 3 }}
      >
        {children}
      </Box>
    </Stack>
  );
};
