import {
  KeyboardArrowDownOutlined,
  KeyboardArrowUpOutlined,
} from '@mui/icons-material';
import { Button, Menu, MenuItem } from '@mui/material';
import React, { useContext, useId } from 'react';

import { UserContext } from '../../../contexts/UserContext/UserContext';

export const ProfileDropdown = () => {
  const { user } = useContext(UserContext);
  const buttonId = useId();
  const menuId = useId();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

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
      <Button
        fullWidth
        id={buttonId}
        aria-controls={open ? menuId : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
        endIcon={
          open ? <KeyboardArrowUpOutlined /> : <KeyboardArrowDownOutlined />
        }
        onClick={handleClick}
      >
        {userName || user?.email}
      </Button>
      <Menu
        id={menuId}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'left',
        }}
        transformOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        MenuListProps={{
          'aria-labelledby': buttonId,
        }}
      >
        {/* TODO implement logout */}
        <MenuItem>Logout</MenuItem>
      </Menu>
    </>
  );
};
