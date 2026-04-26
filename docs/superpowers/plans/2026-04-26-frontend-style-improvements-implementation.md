# Frontend Style Improvements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Modernize the frontend with a Netflix-inspired UI, shadcn/ui components, TanStack Query, and Shaka Player while keeping Vite + React Router.

**Architecture:** Keep the app as a Vite SPA with a shared app shell, three routes (`/upload`, `/catalog`, `/videos/:id`), and React Query as the source of truth for server state. Use shadcn/ui for composition, a single Netflix palette defined in CSS variables, and a dedicated Shaka player wrapper so playback lifecycle stays isolated from page code.

**Tech Stack:** React 19, Vite, TypeScript, React Router, `@tanstack/react-query`, `shaka-player`, shadcn/ui, Tailwind CSS, Vitest, Testing Library, Playwright.

---

## File Structure

### Frontend foundation

- Create: `frontend/components.json`
- Create: `frontend/postcss.config.cjs`
- Create: `frontend/tailwind.config.ts`
- Create: `frontend/src/styles/globals.css`
- Create: `frontend/src/lib/utils.ts`
- Create: `frontend/src/lib/utils.test.ts`
- Modify: `frontend/package.json`
- Modify: `frontend/vite.config.ts`
- Modify: `frontend/src/main.tsx`

### shadcn/ui components

- Create: `frontend/src/components/ui/button.tsx`
- Create: `frontend/src/components/ui/card.tsx`
- Create: `frontend/src/components/ui/input.tsx`
- Create: `frontend/src/components/ui/label.tsx`
- Create: `frontend/src/components/ui/badge.tsx`
- Create: `frontend/src/components/ui/sheet.tsx`
- Create: `frontend/src/components/ui/separator.tsx`
- Create: `frontend/src/components/ui/skeleton.tsx`
- Create: `frontend/src/components/ui/alert.tsx`
- Create: `frontend/src/components/ui/scroll-area.tsx`

### Routing and shell

- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/pages/home-page.tsx`
- Modify: `frontend/src/pages/video-page.tsx`
- Create: `frontend/src/pages/catalog-page.tsx`
- Create: `frontend/src/components/app-shell.tsx`
- Create: `frontend/src/components/app-sidebar.tsx`
- Create: `frontend/src/components/mobile-nav-drawer.tsx`
- Create: `frontend/src/components/page-header.tsx`
- Create: `frontend/src/components/app-shell.test.tsx`

### Data layer

- Create: `frontend/src/lib/query-client.ts`
- Create: `frontend/src/lib/queries.ts`
- Create: `frontend/src/lib/polling.ts`
- Create: `frontend/src/lib/polling.test.ts`
- Create: `frontend/src/lib/queries.test.tsx`
- Modify: `frontend/src/lib/api.ts`
- Modify: `frontend/src/lib/types.ts`
- Modify: `frontend/src/main.tsx`

### Upload and catalog UI

- Modify: `frontend/src/components/upload-form.tsx`
- Modify: `frontend/src/components/video-list.tsx`
- Modify: `frontend/src/lib/format.ts`
- Create: `frontend/src/components/upload-panel.tsx`
- Create: `frontend/src/components/video-card.tsx`
- Create: `frontend/src/components/video-grid.tsx`
- Create: `frontend/src/components/video-status-badge.tsx`
- Create: `frontend/src/components/empty-state.tsx`
- Create: `frontend/src/components/loading-state.tsx`
- Create: `frontend/src/components/error-state.tsx`
- Create: `frontend/src/components/video-card.test.tsx`

### Shaka player

- Modify: `frontend/src/components/video-player.tsx`
- Modify: `frontend/src/pages/video-page.tsx`
- Create: `frontend/src/lib/shaka-player.ts`
- Create: `frontend/src/lib/shaka-player.test.ts`
- Create: `frontend/src/components/player-status-panel.tsx`
- Create: `frontend/src/components/video-player.test.tsx`

### E2E and verification

- Modify: `frontend/playwright.config.ts`
- Modify: `frontend/tests/e2e/upload-and-play.spec.ts`
- Create: `frontend/tests/fixtures/sample.mp4`

---

## Working Rules

- Before editing library-specific code, check current docs with Context7 for the exact package/version being used.
- Use `frontend-design` for layout, composition, spacing, motion, and the Netflix palette.
- Use `shadcn` for component setup and composition instead of hand-rolled UI when a component already exists.
- Use `vercel-react-best-practices` when deciding where state lives, how to structure hooks, and how to avoid rerender churn.
- Use `tanstack-start-best-practices` only for reusable TanStack patterns that still fit a Vite + React Router app; do not migrate to TanStack Start.

## Task 1: Bootstrap the UI foundation and Netflix theme

**Files:**
- Modify: `frontend/package.json`
- Create: `frontend/components.json`
- Create: `frontend/postcss.config.cjs`
- Create: `frontend/tailwind.config.ts`
- Create: `frontend/src/styles/globals.css`
- Create: `frontend/src/lib/utils.ts`
- Create: `frontend/src/lib/utils.test.ts`
- Modify: `frontend/vite.config.ts`
- Modify: `frontend/src/main.tsx`

- [ ] **Step 1: Write the failing utility test**

```ts
import { describe, expect, it } from 'vitest'
import { cn } from './utils'

describe('cn', () => {
  it('merges conditional classes', () => {
    expect(cn('px-2', false && 'hidden', 'px-4')).toBe('px-4')
  })
})
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd frontend && npm test -- --run src/lib/utils.test.ts`
Expected: FAIL with `Cannot find module './utils'`.

- [ ] **Step 3: Initialize shadcn, add base components, and define the theme**

Run these commands from `frontend/`:

```bash
npx shadcn@latest init -t vite
npx shadcn@latest add button card input label badge sheet separator skeleton alert scroll-area
```

Create `frontend/src/lib/utils.ts`:

```ts
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

Create `frontend/src/styles/globals.css` with the Netflix palette and semantic tokens:

```css
:root {
  --background: 20 20 20;
  --surface: 24 24 24;
  --surface-2: 35 35 35;
  --primary: 229 9 20;
  --primary-foreground: 255 255 255;
  --text: 255 255 255;
  --text-muted: 179 179 179;
  --border: 47 47 47;
  --success: 70 211 105;
  --warning: 245 197 24;
  --danger: 255 77 79;
  --ring: 229 9 20;
}
```

Update `frontend/src/main.tsx` to import the global stylesheet.

- [ ] **Step 4: Run the test and a frontend build**

Run: `cd frontend && npm test -- --run src/lib/utils.test.ts`
Expected: PASS.

Run: `cd frontend && npm run build`
Expected: PASS.

- [ ] **Step 5: Commit the foundation**

```bash
git add frontend/package.json frontend/components.json frontend/postcss.config.cjs frontend/tailwind.config.ts frontend/src/styles/globals.css frontend/src/lib/utils.ts frontend/src/lib/utils.test.ts frontend/src/main.tsx frontend/vite.config.ts frontend/src/components/ui
git commit -m "chore(frontend): bootstrap shadcn theme foundation"
```

## Task 2: Add the app shell and route structure

**Files:**
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/pages/home-page.tsx`
- Modify: `frontend/src/pages/video-page.tsx`
- Create: `frontend/src/pages/catalog-page.tsx`
- Create: `frontend/src/components/app-shell.tsx`
- Create: `frontend/src/components/app-sidebar.tsx`
- Create: `frontend/src/components/mobile-nav-drawer.tsx`
- Create: `frontend/src/components/page-header.tsx`
- Create: `frontend/src/components/app-shell.test.tsx`

- [ ] **Step 1: Write the failing shell test**

```tsx
import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AppShell } from './app-shell'

describe('AppShell', () => {
  it('renders upload and catalog navigation', () => {
    render(
      <MemoryRouter initialEntries={['/upload']}>
        <AppShell />
      </MemoryRouter>,
    )

    expect(screen.getByRole('link', { name: /upload/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /catalog/i })).toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd frontend && npm test -- --run src/components/app-shell.test.tsx`
Expected: FAIL with missing `AppShell`.

- [ ] **Step 3: Build the shell and routes with shadcn composition**

Implement `AppShell` around `Outlet` and use `NavLink` for active states. Keep the sidebar fixed on desktop and open the mobile drawer with `Sheet`.

Target route tree in `frontend/src/App.tsx`:

```tsx
<BrowserRouter>
  <Routes>
    <Route element={<AppShell />}>
      <Route index element={<Navigate to="/upload" replace />} />
      <Route path="/upload" element={<UploadPage />} />
      <Route path="/catalog" element={<CatalogPage />} />
      <Route path="/videos/:id" element={<VideoPage />} />
    </Route>
  </Routes>
</BrowserRouter>
```

Keep `home-page.tsx` as the upload page implementation so the current code can be upgraded incrementally, then introduce `catalog-page.tsx` for the catalog view.

- [ ] **Step 4: Run the shell test and the app build**

Run: `cd frontend && npm test -- --run src/components/app-shell.test.tsx`
Expected: PASS.

Run: `cd frontend && npm run build`
Expected: PASS.

- [ ] **Step 5: Commit the shell and routes**

```bash
git add frontend/src/App.tsx frontend/src/pages/home-page.tsx frontend/src/pages/video-page.tsx frontend/src/pages/catalog-page.tsx frontend/src/components/app-shell.tsx frontend/src/components/app-sidebar.tsx frontend/src/components/mobile-nav-drawer.tsx frontend/src/components/page-header.tsx frontend/src/components/app-shell.test.tsx
git commit -m "feat(frontend): add shell and route structure"
```

## Task 3: Add the TanStack Query data layer and polling helpers

**Files:**
- Modify: `frontend/src/lib/api.ts`
- Modify: `frontend/src/lib/types.ts`
- Modify: `frontend/src/main.tsx`
- Create: `frontend/src/lib/query-client.ts`
- Create: `frontend/src/lib/queries.ts`
- Create: `frontend/src/lib/polling.ts`
- Create: `frontend/src/lib/polling.test.ts`
- Create: `frontend/src/lib/queries.test.tsx`
- Modify: `frontend/src/pages/home-page.tsx`
- Modify: `frontend/src/pages/catalog-page.tsx`
- Modify: `frontend/src/pages/video-page.tsx`

- [ ] **Step 1: Write the failing polling test**

```ts
import { describe, expect, it } from 'vitest'
import { shouldPollVideos } from './polling'

describe('shouldPollVideos', () => {
  it('keeps polling while any video is processing', () => {
    expect(
      shouldPollVideos([
        { id: '1', title: 'Trailer', status: 'ready' },
        { id: '2', title: 'Interview', status: 'processing' },
      ]),
    ).toBe(3000)
  })

  it('stops polling when every video is ready', () => {
    expect(shouldPollVideos([{ id: '1', title: 'Trailer', status: 'ready' }])).toBe(false)
  })
})
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd frontend && npm test -- --run src/lib/polling.test.ts`
Expected: FAIL with missing `shouldPollVideos`.

- [ ] **Step 3: Add query keys, the singleton client, and the hooks**

Use a query key factory so the keys stay consistent as the app grows:

```ts
export const videoKeys = {
  all: ['videos'] as const,
  detail: (id: string) => ['video', id] as const,
}
```

Create `query-client.ts` with one exported `queryClient` and `defaultOptions`.

Update `api.ts` so failures throw instead of silently returning empty data.

Create hooks that use `useQuery`, `useMutation`, and `useQueryClient`:

```ts
export function useVideosQuery() { ... }
export function useVideoQuery(id: string) { ... }
export function useUploadVideoMutation() { ... }
```

Wire `refetchInterval` through `shouldPollVideos` and `shouldPollVideo` so polling is declarative, not manual `setInterval` logic.

- [ ] **Step 4: Run the polling and query tests**

Run: `cd frontend && npm test -- --run src/lib/polling.test.ts src/lib/queries.test.tsx`
Expected: PASS.

- [ ] **Step 5: Commit the data layer**

```bash
git add frontend/src/lib/api.ts frontend/src/lib/types.ts frontend/src/lib/query-client.ts frontend/src/lib/queries.ts frontend/src/lib/polling.ts frontend/src/lib/polling.test.ts frontend/src/lib/queries.test.tsx frontend/src/main.tsx frontend/src/pages/home-page.tsx frontend/src/pages/catalog-page.tsx frontend/src/pages/video-page.tsx
git commit -m "feat(frontend): add react query data layer"
```

## Task 4: Rebuild upload and catalog screens with shadcn and Netflix styling

**Files:**
- Modify: `frontend/src/pages/home-page.tsx`
- Create: `frontend/src/components/upload-panel.tsx`
- Create: `frontend/src/components/video-card.tsx`
- Create: `frontend/src/components/video-grid.tsx`
- Create: `frontend/src/components/video-status-badge.tsx`
- Create: `frontend/src/components/empty-state.tsx`
- Create: `frontend/src/components/loading-state.tsx`
- Create: `frontend/src/components/error-state.tsx`
- Modify: `frontend/src/components/upload-form.tsx`
- Modify: `frontend/src/components/video-list.tsx`
- Modify: `frontend/src/lib/format.ts`
- Create: `frontend/src/components/video-card.test.tsx`
- Modify: `frontend/src/components/upload-form.test.tsx`
- Modify: `frontend/src/components/video-list.test.tsx`

- [ ] **Step 1: Write the failing card test**

```tsx
import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { VideoCard } from './video-card'

describe('VideoCard', () => {
  it('renders title, status, and player link', () => {
    render(
      <MemoryRouter>
        <VideoCard video={{ id: '1', title: 'Trailer', status: 'ready' }} />
      </MemoryRouter>,
    )

    expect(screen.getByText('Trailer')).toBeInTheDocument()
    expect(screen.getByText(/ready/i)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /open/i })).toHaveAttribute('href', '/videos/1')
  })
})
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd frontend && npm test -- --run src/components/video-card.test.tsx`
Expected: FAIL with missing `VideoCard`.

- [ ] **Step 3: Compose the upload and catalog UI using shadcn**

Use `Card` composition for upload panels and catalog cards, `Badge` for status, `Skeleton` for loading, and `Alert` for errors.

Keep the layout dark and high-contrast with the Netflix tokens from `globals.css`:

```tsx
<Card>
  <CardHeader>
    <CardTitle>Upload</CardTitle>
    <CardDescription>Send MP4 videos and watch their status update.</CardDescription>
  </CardHeader>
  <CardContent>{/* upload form */}</CardContent>
  <CardFooter>{/* actions */}</CardFooter>
</Card>
```

Make `video-list.tsx` render a `VideoGrid` of `VideoCard` components instead of a plain list.

- [ ] **Step 4: Run the frontend component tests**

Run: `cd frontend && npm test -- --run src/components/upload-form.test.tsx src/components/video-list.test.tsx src/components/video-card.test.tsx`
Expected: PASS.

- [ ] **Step 5: Commit the UI refresh**

```bash
git add frontend/src/pages/home-page.tsx frontend/src/components/upload-panel.tsx frontend/src/components/video-card.tsx frontend/src/components/video-grid.tsx frontend/src/components/video-status-badge.tsx frontend/src/components/empty-state.tsx frontend/src/components/loading-state.tsx frontend/src/components/error-state.tsx frontend/src/components/upload-form.tsx frontend/src/components/video-list.tsx frontend/src/lib/format.ts frontend/src/components/upload-form.test.tsx frontend/src/components/video-list.test.tsx frontend/src/components/video-card.test.tsx
git commit -m "feat(frontend): rebuild upload and catalog ui"
```

## Task 5: Replace dashjs with Shaka Player and improve the playback page

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/src/components/video-player.tsx`
- Modify: `frontend/src/pages/video-page.tsx`
- Create: `frontend/src/lib/shaka-player.ts`
- Create: `frontend/src/lib/shaka-player.test.ts`
- Create: `frontend/src/components/player-status-panel.tsx`
- Create: `frontend/src/components/video-player.test.tsx`

- [ ] **Step 1: Write the failing player lifecycle test**

```ts
import { describe, expect, it, vi } from 'vitest'
import { render, cleanup } from '@testing-library/react'
import { VideoPlayer } from './video-player'

vi.mock('shaka-player', () => ({
  default: {
    polyfill: { installAll: vi.fn() },
    Player: {
      isBrowserSupported: () => true,
    },
  },
}))

describe('VideoPlayer', () => {
  it('mounts and tears down the player', async () => {
    render(<VideoPlayer manifestUrl="/videos/1/stream/manifest.mpd" />)
    cleanup()
    expect(true).toBe(true)
  })
})
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd frontend && npm test -- --run src/components/video-player.test.tsx`
Expected: FAIL with missing Shaka wrapper.

- [ ] **Step 3: Build a dedicated Shaka wrapper with cleanup**

Use the npm package `shaka-player`, not `dashjs`.

The helper should:

- call `shaka.polyfill.installAll()` once on initialization,
- check `shaka.Player.isBrowserSupported()`,
- create a player instance,
- attach it to the `video` element,
- load the manifest URL,
- register an error handler,
- destroy the player on cleanup.

If the player should expose quality controls, use the programmatic Shaka UI overlay and enable the quality selector in the overflow menu.

- [ ] **Step 4: Run the player test and the app build**

Run: `cd frontend && npm test -- --run src/components/video-player.test.tsx src/lib/shaka-player.test.ts`
Expected: PASS.

Run: `cd frontend && npm run build`
Expected: PASS.

- [ ] **Step 5: Commit the playback upgrade**

```bash
git add frontend/package.json frontend/src/components/video-player.tsx frontend/src/pages/video-page.tsx frontend/src/lib/shaka-player.ts frontend/src/lib/shaka-player.test.ts frontend/src/components/player-status-panel.tsx frontend/src/components/video-player.test.tsx
git commit -m "feat(frontend): replace dashjs with shaka player"
```

## Task 6: Tighten E2E coverage for the 3-page flow

**Files:**
- Modify: `frontend/playwright.config.ts`
- Modify: `frontend/tests/e2e/upload-and-play.spec.ts`
- Create: `frontend/tests/fixtures/sample.mp4`

- [ ] **Step 1: Write the updated E2E spec**

```ts
import { expect, test } from '@playwright/test'
import path from 'node:path'

test('can upload, browse catalog, and open the player page', async ({ page }) => {
  await page.route('**/videos', async (route) => {
    const method = route.request().method()

    if (method === 'GET') {
      await route.fulfill({
        json: [
          { id: '1', title: 'Trailer', status: 'ready' },
        ],
      })
      return
    }

    await route.fulfill({
      json: { id: '1', title: 'Trailer', status: 'uploaded' },
    })
  })

  await page.goto('/upload')
  await page.getByLabel('Title').fill('Trailer')
  await page.getByLabel('File').setInputFiles(path.join(process.cwd(), 'tests/fixtures/sample.mp4'))
  await page.getByRole('button', { name: /upload/i }).click()

  await expect(page.getByText('Trailer')).toBeVisible()
  await page.getByRole('link', { name: /trailer/i }).click()
  await expect(page).toHaveURL(/\/videos\/1/)
})
```

- [ ] **Step 2: Run the E2E test and confirm it fails**

Run: `cd frontend && npx playwright test`
Expected: FAIL until the shell, routes, and API hooks are wired.

- [ ] **Step 3: Keep the Playwright config aligned with the Vite dev server**

Keep `webServer` pointed at the local frontend dev server and keep tracing, screenshots, and video capture enabled for failed runs.

- [ ] **Step 4: Run the E2E suite**

Run: `cd frontend && npx playwright test`
Expected: PASS.

- [ ] **Step 5: Commit the coverage**

```bash
git add frontend/playwright.config.ts frontend/tests/e2e/upload-and-play.spec.ts frontend/tests/fixtures/sample.mp4
git commit -m "test(e2e): cover upload catalog and player flow"
```

## Task 7: Final polish and verification

**Files:**
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/main.tsx`
- Modify: `frontend/src/lib/api.ts`
- Modify: `frontend/src/pages/home-page.tsx`
- Modify: `frontend/src/pages/catalog-page.tsx`
- Modify: `frontend/src/pages/video-page.tsx`
- Modify: `frontend/package.json`

- [ ] **Step 1: Remove obsolete frontend assumptions**

Remove `dashjs` imports and any old manual polling logic. Make sure the query hooks are the only source of polling behavior.

- [ ] **Step 2: Run the full frontend verification matrix**

Run:

```bash
cd frontend && npm test -- --run
cd frontend && npm run build
cd frontend && npx playwright test
```

Expected: all commands PASS.

- [ ] **Step 3: Commit the finished frontend refresh**

```bash
git add frontend/package.json frontend/src/App.tsx frontend/src/main.tsx frontend/src/lib/api.ts frontend/src/pages/home-page.tsx frontend/src/pages/catalog-page.tsx frontend/src/pages/video-page.tsx
git commit -m "chore(frontend): finish style overhaul"
```

## Self-Review

- Netflix palette: covered by Task 1 and Task 4.
- Three pages plus sidebar/drawer navigation: covered by Task 2.
- React Query as the server-state layer: covered by Task 3.
- Shaka Player as the playback engine: covered by Task 5.
- shadcn/ui composition: covered by Tasks 1, 2, and 4.
- Context7 and library skills usage: covered in the Working Rules and required before each library-specific step.
- Performance and rerender hygiene: covered by the query key factory, declarative polling, singleton query client, and isolated player lifecycle.
- Testing: unit tests, component tests, and E2E are each assigned their own task.
