import { useMutation, useQuery } from "react-query";
import { IconButton, Skeleton, Stack } from "@mui/material";
import { GitHub, Google, Microsoft } from "@mui/icons-material";
import { api } from "../../../api";
import { ManifestIDP } from "../../../api/Api";

export const AuthSocialButtons = () => {
  const { isLoading, data } = useQuery("manifest", api.manifestList);
  const { mutateAsync } = useMutation(api.startExternalLoginCreate);

  if (!data || isLoading) {
    return (
      <Stack direction="row" spacing={1}>
        {Array.from({ length: 3 }).map((_, index) => (
          <Skeleton key={index} variant="circular" width="40px" height="40px">
            <IconButton />
          </Skeleton>
        ))}
      </Stack>
    );
  }

  return (
    <Stack direction="row" spacing={1}>
      {data.data.social_idps
        ?.filter((item): item is Required<ManifestIDP> => !!item.slug)
        .map((item) => {
          const Icon = IconsMap[item.slug as SocialIdps];

          return (
            <IconButton
              key={item.slug}
              onClick={() =>
                mutateAsync({ slug: item.slug, directive: "login" })
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
  Google = "google-social",
  Github = "github-social",
  Microsoft = "microsoft-social",
}

const IconsMap: Record<SocialIdps, React.ElementType> = {
  [SocialIdps.Google]: Google,
  [SocialIdps.Github]: GitHub,
  [SocialIdps.Microsoft]: Microsoft,
};
