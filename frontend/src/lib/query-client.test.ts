import { describe, expect, it } from 'vitest'
import { queryClient } from './query-client'

describe('queryClient', () => {
  it('uses stable default query options', () => {
    const defaults = queryClient.getDefaultOptions()

    expect(defaults.queries?.refetchOnWindowFocus).toBe(false)
    expect(defaults.queries?.retry).toBe(1)
    expect(defaults.queries?.staleTime).toBe(5_000)
    expect(defaults.mutations?.retry).toBe(0)
  })
})
