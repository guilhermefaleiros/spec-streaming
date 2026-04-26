import { Outlet } from 'react-router-dom'
import { AppSidebar } from './app-sidebar'
import { PageHeader } from './page-header'

export function AppShell() {
  return (
    <div className="min-h-screen bg-background text-foreground lg:grid lg:grid-cols-[240px_1fr]">
      <aside className="hidden border-r border-border px-4 py-6 lg:block">
        <div className="mb-6">
          <p className="text-sm text-muted-foreground">Spec Streaming</p>
        </div>
        <AppSidebar />
      </aside>
      <div className="flex min-w-0 flex-col">
        <PageHeader />
        <main className="flex-1 px-4 py-6 lg:px-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
