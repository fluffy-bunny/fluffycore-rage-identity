import { useForm } from "react-hook-form";
import {
  Box,
  Button,
  FormControl,
  Link,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { AuthSocialButtons } from "../components/auth/AuthSocialButtons/AuthSocialButtons";
import { useMutation } from "react-query";
import { api } from "../api";
import { LoadingButton } from "@mui/lab";
import { LoginModelsLoginPhaseOneRequest } from "../api/Api";
import { RoutePaths } from "../constants/routes";

export const Page = () => {
  const { register, handleSubmit } = useForm<LoginModelsLoginPhaseOneRequest>();
  const { mutateAsync, isLoading } = useMutation(
    (values: LoginModelsLoginPhaseOneRequest) =>
      api.loginPhaseOneCreate(values, {
        withXSRFToken: true,
        withCredentials: true,
      })
  );

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Sign in
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
          <Box>
            <Link href={RoutePaths.ForgotPassword}>Forgot password?</Link>
          </Box>
        </FormControl>
        <FormControl fullWidth sx={{ marginTop: 3 }}>
          <Stack direction="row">
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography>Sign in with socials</Typography>
              <AuthSocialButtons />
            </Stack>
            <Stack direction="row" spacing={1} sx={{ marginLeft: "auto" }}>
              <Button href={RoutePaths.SignUp}>Sign up</Button>
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
    </>
  );
};

// export function Page() {
//   const [email, setEmail] = useState("");
//   const [redirectToPassword, setRedirectToPassword] = useState(false);
//   const [socialIdps, setSocialIdps] = useState<{ slug: string }[]>([]);

//   useEffect(() => {
//     const fetchData = async () => {
//       const response = await fetch("http://localhost1.com:9044/api/manifest", {
//         method: "GET",
//         headers: {
//           Accept: "application/json",
//           "Content-Type": "application/json",
//         },
//       });

//       if (response.ok) {
//         const data = await response.json();
//         setSocialIdps(data.social_idps);
//         console.log(data.social_idps);
//       } else {
//         // Handle error
//       }
//     };

//     fetchData();
//   }, []);

//   const handleSubmit = async (event: React.FormEvent) => {
//     event.preventDefault();
//     let csrf = myCommon.getCSRF();
//     console.log(csrf);
//     // Fetch call to validate the email
//     const response = await fetch(
//       "http://localhost1.com:9044/api/login-phase-one",
//       {
//         method: "POST",
//         credentials: "include",
//         headers: {
//           Accept: "application/json",
//           "Content-Type": "application/json",
//           "X-Csrf-Token": csrf,
//         },
//         body: JSON.stringify({
//           email: email,
//         }), // send the email as part of the request body
//       }
//     );

//     if (response.ok) {
//       const data = await response.json();
//       console.log(data);

//       if (data.isValid) {
//         setRedirectToPassword(true);
//       } else {
//         // Handle invalid email
//       }
//     } else {
//       // Handle error
//     }
//   };

//   return (
//     <form onSubmit={handleSubmit}>
//       <h1>Login</h1>
//       <label>
//         Email:
//         <input
//           type="email"
//           value={email}
//           onChange={(e) => setEmail(e.target.value)}
//           required
//         />
//       </label>
//       <button type="submit">Next</button>
//       <div>
//         <a href="/signup">Sign Up</a> |{" "}
//         <a href="/forgot-password">Forgot Password?</a>
//       </div>
//       <div
//         style={{
//           display: "flex",
//           flexDirection: "row",
//           justifyContent: "space-between",
//         }}
//       >
//         {socialIdps.map((idp) => (
//           <button key={idp.slug}>{idp.slug}</button>
//         ))}
//       </div>
//     </form>
//   );
// }
