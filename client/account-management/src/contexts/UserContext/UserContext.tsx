import { Box, CircularProgress } from '@mui/material';
import { createContext } from 'react';

import {
  ApiUserIdentityInfoUserIdentityInfo,
  ApiUserProfileProfile,
} from '../../api/Api';
import { useUserIdentity } from '../../hooks/useUserIdentity';
import { useUserProfile } from '../../hooks/useUserProfile';

interface UserContextValue {
  user?: ApiUserProfileProfile & ApiUserIdentityInfoUserIdentityInfo;
  refetch?: () => Promise<void>;
}

const defaultUserContextValue: UserContextValue = {
  user: undefined,
};

export const UserContext = createContext<UserContextValue>(
  defaultUserContextValue,
);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const userProfile = useUserProfile();
  const userIdentity = useUserIdentity();

  if (userProfile.isLoading || userIdentity.isLoading) {
    return (
      <Box sx={{ display: 'flex', height: '100%' }}>
        <CircularProgress sx={{ margin: 'auto' }} />
      </Box>
    );
  }

  return (
    <UserContext.Provider
      value={{
        user: {
          ...userProfile.data,
          ...userIdentity.data,
        },
        refetch: async () => {
          await Promise.all([userProfile.refetch(), userIdentity.refetch()]);
        },
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
