import { useQuery } from 'react-query';

import { api } from '../api';

export function useUserProfile() {
  const { isLoading, data, refetch } = useQuery(
    'user-profile',
    api.userProfileList,
  );

  return {
    isLoading,
    data: data?.data,
    refetch,
  };
}
