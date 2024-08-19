import { LoadingButton } from '@mui/lab';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from '@mui/material';
import { useMutation } from 'react-query';

import { api } from '../../../api';
import { useNotification } from '../../../contexts/NotificationContext/NotificationContext';

export const ProfileDeletePasskeyModal = ({
  passkeyId,
  onClose,
  onSuccess,
}: {
  passkeyId: string;
  onClose(): void;
  onSuccess(): void;
}) => {
  const { showNotification } = useNotification();

  const { isLoading, mutateAsync } = useMutation(api.userRemovePasskeyCreate, {
    onSuccess: () => {
      onSuccess();
      onClose();
      showNotification('Passkey has been deleted successfully.', 'success');
    },
    onError: () => {
      showNotification('Something went wrong. Please try again.', 'error');
    },
  });

  return (
    <Dialog
      open={!!passkeyId}
      onClose={() => {
        if (isLoading) {
          return;
        }

        onClose();
      }}
    >
      <DialogTitle>Delete passkey?</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Are you sure you want to delete this passkey?
        </DialogContentText>
      </DialogContent>
      <DialogActions sx={{ paddingTop: 0, paddingX: 3, paddingBottom: 3 }}>
        <Button onClick={onClose}>Cancel</Button>
        <LoadingButton
          loading={isLoading}
          variant="contained"
          onClick={() =>
            mutateAsync({
              aaguid: passkeyId,
            })
          }
        >
          Delete
        </LoadingButton>
      </DialogActions>
    </Dialog>
  );
};
