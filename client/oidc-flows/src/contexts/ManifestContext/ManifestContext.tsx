import { Box, CircularProgress } from '@mui/material';
import React, { createContext, ReactNode, useContext } from 'react';
import { useQuery } from 'react-query';

import { api } from '../../api';
import { ManifestManifest } from '../../api/Api';

interface ManifestContextProps {
  data?: ManifestManifest;
}

const ManifestContext = createContext<ManifestContextProps | undefined>({
  data: undefined,
});

interface ManifestProviderProps {
  children: ReactNode;
}

export const ManifestProvider: React.FC<ManifestProviderProps> = ({
  children,
}) => {
  const { isLoading, data } = useQuery('manifest', api.manifestList);

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', height: '100%' }}>
        <CircularProgress sx={{ margin: 'auto' }} />
      </Box>
    );
  }

  return (
    <ManifestContext.Provider value={{ data: data?.data }}>
      {children}
    </ManifestContext.Provider>
  );
};

export const useManifest = (): ManifestContextProps => {
  const context = useContext(ManifestContext);

  if (context === undefined) {
    throw new Error('useManifest must be used within a ManifestProvider');
  }

  return context;
};
