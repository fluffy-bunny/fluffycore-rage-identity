import {
  KeyboardArrowDownOutlined,
  KeyboardArrowUpOutlined,
} from '@mui/icons-material';
import { Box, Button, Menu, MenuItem } from '@mui/material';
import React, { useContext, useId } from 'react';
import { useMutation } from 'react-query';

import { api } from '../../../api';
import { UserContext } from '../../../contexts/UserContext/UserContext';

export const ProfileDropdown = () => {
  const { user } = useContext(UserContext);
  const buttonId = useId();
  const menuId = useId();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const logoutMutation = useMutation('logout', api.logoutCreate, {
    onSuccess: (data) => {
      if (data.data.directive === 'redirect' && data.data.redirectURL) {
        window.location.href = data.data.redirectURL;
      }
    },
  });

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const userName = [user?.givenName, user?.familyName]
    .filter(Boolean)
    .join(' ');

  return (
    <>
      <Box sx={{ display: 'flex' }}>
        <Button
          id={buttonId}
          aria-controls={open ? menuId : undefined}
          aria-haspopup="true"
          aria-expanded={open ? 'true' : undefined}
          endIcon={
            open ? <KeyboardArrowUpOutlined /> : <KeyboardArrowDownOutlined />
          }
          sx={{ marginX: 'auto' }}
          onClick={handleClick}
        >
          {userName || user?.email}
        </Button>
      </Box>
      <Menu
        id={menuId}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        MenuListProps={{
          'aria-labelledby': buttonId,
        }}
        slotProps={{
          paper: {
            sx: { borderRadius: 4 },
          },
        }}
      >
        <MenuItem
          disabled={logoutMutation.isLoading}
          onClick={() => logoutMutation.mutateAsync({})}
        >
          Logout
        </MenuItem>
      </Menu>
    </>
  );
};
