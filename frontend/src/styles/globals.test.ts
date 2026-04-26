import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

describe('globals.css', () => {
  it('imports Tailwind CSS so utility classes are generated', () => {
    const css = readFileSync(resolve(process.cwd(), 'src/styles/globals.css'), 'utf8')

    expect(css).toContain('@import "tailwindcss";')
  })

  it('defines Tailwind v4 semantic theme mappings', () => {
    const css = readFileSync(resolve(process.cwd(), 'src/styles/globals.css'), 'utf8')

    expect(css).toContain('@theme inline')
    expect(css).toContain('--color-card: rgb(var(--card))')
    expect(css).toContain('--color-muted: rgb(var(--muted))')
    expect(css).toContain('--color-destructive: rgb(var(--destructive))')
  })
})
