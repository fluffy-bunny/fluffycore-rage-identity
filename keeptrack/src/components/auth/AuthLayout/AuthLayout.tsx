import { Box, Card, CardContent, Container, Typography } from '@mui/material';
import React from 'react';

export const AuthLayout = ({
  children,
  title,
  description,
}: {
  children: React.ReactNode;
  title: string;
  description?: string;
}) => {
  return (
    <Box sx={{ display: 'flex', height: '100%' }}>
      <Container maxWidth="md" sx={{ margin: 'auto' }}>
        <Card>
          <CardContent sx={{ padding: 4, '&:last-child': { padding: 4 } }}>
            <Typography variant="h4" component="h1" gutterBottom>
              {title}
            </Typography>
            {description && <Typography gutterBottom>{description}</Typography>}
            {children}
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};
