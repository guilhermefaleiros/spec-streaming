import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AppShell } from './app-shell'

describe('AppShell', () => {
  it('renders upload and catalog navigation plus the drawer trigger', () => {
    render(
      <MemoryRouter initialEntries={['/upload']}>
        <AppShell />
      </MemoryRouter>,
    )

    expect(screen.getByRole('link', { name: /upload/i })).toHaveAttribute('href', '/upload')
    expect(screen.getByRole('link', { name: /catalog/i })).toHaveAttribute('href', '/catalog')

    expect(screen.getByRole('button', { name: /open navigation/i })).toBeInTheDocument()
  })
})
