import { Box, Card, CardContent, Container } from "@mui/material";
import { Outlet } from "react-router-dom";

export const AuthLayout = () => {
  return (
    <Box sx={{ display: "flex", height: "100%" }}>
      <Container maxWidth="md" sx={{ margin: "auto" }}>
        <Card>
          <CardContent sx={{ padding: 4, "&:last-child": { padding: 4 } }}>
            <Outlet />
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};
