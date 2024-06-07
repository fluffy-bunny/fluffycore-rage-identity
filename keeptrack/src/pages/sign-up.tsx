import { useForm } from "react-hook-form";
import {
  Box,
  Button,
  FormControl,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation } from "react-query";
import { api } from "../api";
import { LoadingButton } from "@mui/lab";
import { LoginModelsSignupRequest } from "../api/Api";
import { RoutePaths } from "../constants/routes";

export const Page = () => {
  const { register, handleSubmit } = useForm<LoginModelsSignupRequest>();
  const { mutateAsync, isLoading } = useMutation(api.signupCreate);

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Sign up
      </Typography>
      <Box
        component="form"
        onSubmit={handleSubmit((values) => mutateAsync(values))}
      >
        <FormControl>
          <TextField
            {...register("email")}
            label="Email address"
            placeholder="Enter your email"
          />
        </FormControl>
        <FormControl>
          <TextField
            {...register("password")}
            type="password"
            label="Password"
            placeholder="Enter your password"
          />
        </FormControl>
        <FormControl sx={{ display: "flex", marginTop: 3 }}>
          <Stack direction="row" spacing={1} sx={{ marginLeft: "auto" }}>
            <Button href={RoutePaths.SignIn}>Sign in</Button>
            <LoadingButton
              loading={isLoading}
              type="submit"
              variant="contained"
            >
              Next
            </LoadingButton>
          </Stack>
        </FormControl>
      </Box>
    </>
  );
};
