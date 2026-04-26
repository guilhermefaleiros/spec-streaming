import { useMutation, useQuery } from '@tanstack/react-query'
import { getVideo, listVideos, uploadVideo } from './api'
import { shouldPollVideo, shouldPollVideos } from './polling'
import { queryClient } from './query-client'

export const videoKeys = {
  all: ['videos'] as const,
  detail: (id: string) => ['videos', id] as const,
}

export function useVideosQuery() {
  return useQuery({
    queryKey: videoKeys.all,
    queryFn: listVideos,
    refetchInterval: (query) => shouldPollVideos(query.state.data ?? []),
  })
}

export function useVideoQuery(id?: string) {
  return useQuery({
    queryKey: videoKeys.detail(id ?? ''),
    queryFn: () => getVideo(id ?? ''),
    enabled: Boolean(id),
    refetchInterval: (query) => shouldPollVideo(query.state.data),
  })
}

export function useUploadVideoMutation() {
  return useMutation({
    mutationFn: uploadVideo,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: videoKeys.all })
    },
  })
}
