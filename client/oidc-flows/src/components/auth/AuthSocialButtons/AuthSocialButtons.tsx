import { GitHub, Google, Microsoft } from '@mui/icons-material';
import { IconButton, Stack } from '@mui/material';
import { useMutation } from 'react-query';

import { api } from '../../../api';
import {
  ExternalIdpStartExternalIDPLoginRequest,
  ManifestIDP,
} from '../../../api/Api';
import { useManifest } from '../../../contexts/ManifestContext/ManifestContext';

export const AuthSocialButtons = () => {
  const { data } = useManifest();
  const { mutateAsync } = useMutation(
    (values: ExternalIdpStartExternalIDPLoginRequest) =>
      api.startExternalLoginCreate(values),
    {
      onSuccess: (data) => {
        if (data.data.redirectUri) {
          window.location.href = data.data.redirectUri;
        }
      },
    },
  );

  return (
    <Stack direction="row" spacing={1}>
      {data?.social_idps
        ?.filter((item): item is Required<ManifestIDP> => !!item.slug)
        .map((item) => {
          const Icon = IconsMap[item.slug as SocialIdps];

          return (
            <IconButton
              key={item.slug}
              onClick={() =>
                mutateAsync({ slug: item.slug, directive: 'login' })
              }
            >
              <Icon />
            </IconButton>
          );
        })}
    </Stack>
  );
};

enum SocialIdps {
  Google = 'google-social',
  Github = 'github-social',
  Microsoft = 'microsoft-social',
}

const IconsMap: Record<SocialIdps, React.ElementType> = {
  [SocialIdps.Google]: Google,
  [SocialIdps.Github]: GitHub,
  [SocialIdps.Microsoft]: Microsoft,
};
