import { NavLink } from 'react-router-dom'
import { cn } from '@/lib/utils'

export const links = [
  { to: '/upload', label: 'Upload' },
  { to: '/catalog', label: 'Catalog' },
]

export function AppSidebar() {
  return (
    <nav aria-label="Primary" className="flex flex-col gap-2">
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          className={({ isActive }) =>
            cn(
              'rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground',
              isActive && 'bg-muted text-foreground',
            )
          }
        >
          {link.label}
        </NavLink>
      ))}
    </nav>
  )
}
