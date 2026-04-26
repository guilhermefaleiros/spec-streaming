import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createShakaPlayer } from './shaka-player'

const mocks = vi.hoisted(() => {
  const destroy = vi.fn().mockResolvedValue(undefined)
  const attach = vi.fn().mockResolvedValue(undefined)
  const load = vi.fn().mockResolvedValue(undefined)
  const addEventListener = vi.fn()
  return { destroy, attach, load, addEventListener }
})

vi.mock('shaka-player', () => ({
  default: {
    polyfill: { installAll: vi.fn() },
    Player: class {
      static isBrowserSupported = vi.fn(() => true)
      attach = mocks.attach
      load = mocks.load
      destroy = mocks.destroy
      addEventListener = mocks.addEventListener
    },
  },
}))

beforeEach(() => {
  mocks.destroy.mockClear()
  mocks.attach.mockClear()
  mocks.load.mockClear()
  mocks.addEventListener.mockClear()
})

describe('createShakaPlayer', () => {
  it('attaches, loads, and returns a cleanup function', async () => {
    const result = await createShakaPlayer({
      videoElement: document.createElement('video'),
      manifestUrl: '/manifest.mpd',
      onError: () => {},
    })

    expect(result.ok).toBe(true)
    if (result.ok) {
      await result.cleanup()
    }
    expect(mocks.attach).toHaveBeenCalled()
    expect(mocks.load).toHaveBeenCalledWith('/manifest.mpd')
    expect(mocks.destroy).toHaveBeenCalled()
  })

  it('reports load failures through onError and destroys the player', async () => {
    mocks.load.mockRejectedValueOnce(new Error('load failed'))
    const messages: string[] = []
    const result = await createShakaPlayer({
      videoElement: document.createElement('video'),
      manifestUrl: '/manifest.mpd',
      onError: (message) => messages.push(message),
    })

    expect(result.ok).toBe(false)
    expect(messages).toContain('load failed')
    expect(mocks.destroy).toHaveBeenCalled()
  })
})
