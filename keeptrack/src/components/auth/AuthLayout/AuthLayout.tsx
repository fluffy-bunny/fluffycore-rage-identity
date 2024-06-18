import React from 'react';
import { Box, Card, CardContent, Container } from '@mui/material';

export const AuthLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <Box sx={{ display: 'flex', height: '100%' }}>
      <Container maxWidth="md" sx={{ margin: 'auto' }}>
        <Card>
          <CardContent sx={{ padding: 4, '&:last-child': { padding: 4 } }}>
            {children}
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};
