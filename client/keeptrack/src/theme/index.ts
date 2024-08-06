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
          borderRadius: 24,
        },
      },
    },
    MuiLink: {
      styleOverrides: {
        root: {
          fontSize: 15,
          textDecoration: 'none',
          '&:hover': {
            textDecoration: 'underline',
          },
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 4,
        },
      },
    },
    MuiListItemIcon: {
      styleOverrides: {
        root: {
          minWidth: 'auto',
          marginRight: 8,
        },
      },
    },
  },
});
