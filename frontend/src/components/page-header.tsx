import { MobileNavDrawer } from './mobile-nav-drawer'

export function PageHeader() {
  return (
    <header className="flex items-center justify-between gap-4 border-b border-border px-4 py-4 lg:px-6">
      <div>
        <p className="text-sm text-muted-foreground">Spec Streaming</p>
      </div>
      <div className="lg:hidden">
        <MobileNavDrawer />
      </div>
    </header>
  )
}
