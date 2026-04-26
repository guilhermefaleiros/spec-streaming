import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

const files = [
  'src/components/empty-state.tsx',
  'src/components/loading-state.tsx',
  'src/components/player-status-panel.tsx',
  'src/components/upload-panel.tsx',
  'src/components/video-card.tsx',
  'src/components/video-player.tsx',
  'src/pages/video-page.tsx',
]

describe('theme class usage', () => {
  it('avoids non-existent custom utility classes', () => {
    for (const file of files) {
      const source = readFileSync(resolve(process.cwd(), file), 'utf8')

      expect(source).not.toContain('bg-surface')
      expect(source).not.toContain('bg-surface-2')
      expect(source).not.toContain('text-danger')
    }
  })
})
