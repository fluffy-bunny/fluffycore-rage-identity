import { Delete } from '@mui/icons-material';
import {
  IconButton,
  List,
  ListItem,
  ListItemText,
  Tooltip,
} from '@mui/material';
import { useContext, useState } from 'react';

import { UserContext } from '../../../contexts/UserContext/UserContext';
import { ProfileDeletePasskeyModal } from '../ProfileDeletePasskeyModal/ProfileDeletePasskeyModal';

export const ProfilePasskeysGrid = () => {
  const [deletePasskeyId, setDeletePasskyId] = useState<string>();
  const { user, refetch } = useContext(UserContext);

  if (!user?.passkeys) {
    return null;
  }

  return (
    <>
      <List>
        {user.passkeys.map((passkey) => (
          <ListItem
            key={passkey.aaguid}
            secondaryAction={
              <Tooltip title={`Delete the ${passkey.name} passkey`}>
                <IconButton
                  color="error"
                  size="small"
                  onClick={() => setDeletePasskyId(passkey.aaguid)}
                >
                  <Delete fontSize="inherit" />
                </IconButton>
              </Tooltip>
            }
            disablePadding
          >
            <ListItemText>{passkey.name}</ListItemText>
          </ListItem>
        ))}
      </List>
      {deletePasskeyId && (
        <ProfileDeletePasskeyModal
          passkeyId={deletePasskeyId}
          onClose={() => setDeletePasskyId(undefined)}
          onSuccess={() => refetch?.()}
        />
      )}
    </>
  );
};
