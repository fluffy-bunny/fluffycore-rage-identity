import { AccountCircle } from '@mui/icons-material';
import { Box, IconButton, Menu, MenuItem, Tooltip } from '@mui/material';
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
        <Tooltip title={userName || user?.email}>
          <IconButton
            id={buttonId}
            aria-controls={open ? menuId : undefined}
            aria-haspopup="true"
            aria-expanded={open ? 'true' : undefined}
            color="inherit"
            onClick={handleClick}
          >
            <AccountCircle />
          </IconButton>
        </Tooltip>
      </Box>
      <Menu
        id={menuId}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
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
