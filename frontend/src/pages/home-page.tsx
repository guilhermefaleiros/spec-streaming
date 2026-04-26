import { UploadPanel } from '@/components/upload-panel'
import { useUploadVideoMutation } from '@/lib/queries'

export function HomePage() {
  const uploadVideo = useUploadVideoMutation()

  return (
    <div className="mx-auto w-full max-w-2xl p-6">
      <UploadPanel
        onSubmit={async (input) => {
          await uploadVideo.mutateAsync(input)
        }}
      />
    </div>
  )
}
