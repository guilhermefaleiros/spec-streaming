import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AppShell } from './components/app-shell'
import { CatalogPage } from './pages/catalog-page'
import { HomePage } from './pages/home-page'
import { VideoPage } from './pages/video-page'

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/upload" replace />} />
        <Route element={<AppShell />}>
          <Route path="/upload" element={<HomePage />} />
          <Route path="/catalog" element={<CatalogPage />} />
          <Route path="/videos/:id" element={<VideoPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
