import { CssBaseline, GlobalStyles, ThemeProvider } from '@mui/material';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from 'react-query';

import { App } from './App';
import { NotificationProvider } from './contexts/NotificationContext/NotificationContext';
import { theme } from './theme';
import { AppType } from './types';

const rootElement = document.getElementById('root') as HTMLElement;
const rootAppAttribute = rootElement.getAttribute('data-app') as AppType;

const root = ReactDOM.createRoot(rootElement);

const queryClient = new QueryClient();

root.render(
  <ThemeProvider theme={theme}>
    <CssBaseline>
      <GlobalStyles
        styles={{
          html: { height: '100%' },
          body: { height: '100%', margin: 0 },
          '#root': { height: '100%', width: '100%' },
        }}
      />
      <QueryClientProvider client={queryClient}>
        <NotificationProvider>
          <App app={rootAppAttribute} />
        </NotificationProvider>
      </QueryClientProvider>
    </CssBaseline>
  </ThemeProvider>,
);
