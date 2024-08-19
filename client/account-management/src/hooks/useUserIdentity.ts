import { useQuery } from 'react-query';

import { api } from '../api';

export function useUserIdentity() {
  const { isLoading, data, refetch } = useQuery(
    'user-identity',
    api.userIdentityInfoList,
  );

  return {
    isLoading,
    data: data?.data,
    refetch,
  };
}
