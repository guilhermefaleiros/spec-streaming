import shaka from 'shaka-player'

shaka.polyfill.installAll()

type Options = {
  videoElement: HTMLVideoElement
  manifestUrl: string
  onError: (message: string) => void
}

type Result =
  | { ok: true; cleanup: () => Promise<void> }
  | { ok: false }

export async function createShakaPlayer({ videoElement, manifestUrl, onError }: Options): Promise<Result> {
  if (!shaka.Player.isBrowserSupported()) {
    console.error('Shaka browser support error', { manifestUrl })
    onError('Browser not supported')
    return { ok: false }
  }

  const player = new shaka.Player()
  ;(videoElement as HTMLVideoElement & { dataset?: DOMStringMap }).dataset.playbackState = 'buffering'

  player.addEventListener('error', (event: Event & { detail?: { message?: string } }) => {
    const error = event.detail ?? {}
    console.error('Shaka playback error', error)
    onError(error.message || 'Playback failed')
  })

  try {
    await player.attach(videoElement)
    await player.load(manifestUrl)
    ;(videoElement as HTMLVideoElement & { dataset?: DOMStringMap }).dataset.playbackState = 'ready'
    return {
      ok: true,
      cleanup: () => player.destroy(),
    }
  } catch (error) {
    ;(videoElement as HTMLVideoElement & { dataset?: DOMStringMap }).dataset.playbackState = 'failed'
    console.error('Shaka setup error', error)
    onError(error instanceof Error ? error.message : 'Playback failed')
    await player.destroy().catch(() => {})
    return { ok: false }
  }
}
