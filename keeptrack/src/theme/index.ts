import { createTheme } from '@mui/material';

export const theme = createTheme({
  palette: {
    background: {
      default: '#f0f4f9',
    },
  },
  components: {
    MuiButton: {
      defaultProps: {
        size: 'large',
        disableElevation: true,
      },
      styleOverrides: {
        root: {
          borderRadius: 32,
          textTransform: 'initial',
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        fullWidth: true,
      },
    },
    MuiFormControl: {
      defaultProps: {
        margin: 'dense',
        fullWidth: true,
      },
    },
    MuiCard: {
      defaultProps: {
        elevation: 0,
      },
      styleOverrides: {
        root: {
          borderRadius: 32,
        },
      },
    },
  },
});
