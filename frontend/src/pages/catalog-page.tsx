import { ErrorState } from '@/components/error-state'
import { LoadingState } from '@/components/loading-state'
import { VideoList } from '@/components/video-list'
import { useVideosQuery } from '@/lib/queries'

export function CatalogPage() {
  const videos = useVideosQuery()

  if (videos.isLoading) {
    return (
      <div className="mx-auto w-full max-w-6xl p-6">
        <LoadingState />
      </div>
    )
  }

  if (videos.isError) {
    return (
      <div className="mx-auto w-full max-w-6xl p-6">
        <ErrorState title="Catalog unavailable" description="We could not load the video catalog right now." />
      </div>
    )
  }

  return (
    <div className="mx-auto w-full max-w-6xl p-6">
      <div className="mb-6 flex flex-col gap-2">
        <h1 className="text-3xl font-semibold tracking-tight">Catalog</h1>
        <p className="text-sm text-muted-foreground">Browse processed videos and open them in the player.</p>
      </div>
      <VideoList videos={videos.data ?? []} />
    </div>
  )
}
