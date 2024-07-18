import {
  Box,
  Card,
  CardContent,
  List,
  ListItemButton,
  ListItemText,
  Stack,
} from '@mui/material';
import React from 'react';

import { RoutePaths } from '../../../constants/routes';

const navItems = [
  {
    value: RoutePaths.ProfilePersonalInformation,
    label: 'Personal Information',
  },
  {
    value: RoutePaths.ProfileSecuritySettings,
    label: 'Security settings',
  },
  {
    value: RoutePaths.ProfilePasskeysManagement,
    label: 'Passkeys management',
  },
];

const SidebarWidth = 250;

export const ProfileLayout = ({
  children,
  currentPage,
  onNavigate,
}: {
  children: React.ReactNode;
  currentPage?: string | 'default';
  onNavigate: (path: string) => void;
}) => {
  return (
    <Card sx={{ height: '100%' }}>
      <CardContent sx={{ height: '100%', padding: 3 }}>
        <Stack direction="row" spacing={3} sx={{ height: '100%' }}>
          <List sx={{ width: SidebarWidth, flexBasis: SidebarWidth }}>
            {navItems.map((item, index) => (
              <ListItemButton
                key={item.label}
                selected={
                  currentPage === 'default'
                    ? index === 0
                    : item.value === currentPage
                }
                onClick={() => onNavigate(item.value)}
              >
                <ListItemText>{item.label}</ListItemText>
              </ListItemButton>
            ))}
          </List>
          <Box component="main" sx={{ flexGrow: 1, paddingY: 1.5 }}>
            {children}
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
};
