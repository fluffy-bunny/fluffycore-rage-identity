import { createContext } from 'react';

import { AppType } from '../../types';

interface AppContextValue {
  app: AppType;
  setApp: (appType: AppType) => void;
}

const defaultAppContextValue: AppContextValue = {
  app: AppType.Auth,
  setApp: () => {},
};

export const AppContext = createContext<AppContextValue>(
  defaultAppContextValue,
);
