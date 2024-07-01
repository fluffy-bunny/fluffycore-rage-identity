import { LoadingButton } from '@mui/lab';
import { Box, FormControl, Link, Stack, TextField } from '@mui/material';
import { useForm } from 'react-hook-form';
import { useMutation } from 'react-query';

import { api } from '../api';
import { LoginModelsVerifyCodeRequest } from '../api/Api';
import { AuthLayout } from '../components/auth/AuthLayout/AuthLayout';
import { RoutePaths } from '../constants/routes';
import { PageProps } from '../types';

export const VerifyCodePage: React.FC<PageProps<{ email: string }>> = ({
  pageProps,
  onNavigate,
}) => {
  const {
    formState: { errors },
    register,
    handleSubmit,
    getFieldState,
  } = useForm<LoginModelsVerifyCodeRequest>();
  const { mutateAsync, isLoading } = useMutation(api.verifyCodeCreate, {
    onSuccess: (data) => {
      if (data.data.directiveRedirect?.redirectUri) {
        window.location.href = data.data.directiveRedirect.redirectUri;
      }

      if (data.data.directive === 'displayLoginPhaseOnePage') {
        return onNavigate(RoutePaths.SignIn);
      }

      if (data.data.directive === 'displayPasswordResetPage') {
        return onNavigate(RoutePaths.ResetPassword);
      }
    },
  });

  return (
    <AuthLayout
      title="Verify Code"
      description={`A verification code has be emailed to ${pageProps?.email} If an account exists.`}
    >
      <Box
        component="form"
        onSubmit={handleSubmit((values) => mutateAsync(values))}
      >
        <FormControl>
          <TextField
            {...register('code', { required: 'You must enter the code.' })}
            error={getFieldState('code').invalid}
            helperText={errors.code?.message}
            label="Verification code"
            placeholder="Enter verification code"
          />
        </FormControl>
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack
              spacing={2}
              direction="row"
              sx={{ marginLeft: 'auto', alignItems: 'center' }}
            >
              <Link
                component="button"
                onClick={() => onNavigate(RoutePaths.SignIn)}
              >
                Cancel
              </Link>
              <LoadingButton
                loading={isLoading}
                type="submit"
                variant="contained"
              >
                Next
              </LoadingButton>
            </Stack>
          </Stack>
        </FormControl>
      </Box>
    </AuthLayout>
  );
};
