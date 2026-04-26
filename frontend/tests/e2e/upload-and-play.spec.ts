import { expect, test } from '@playwright/test'
import path from 'node:path'

test('redirects to upload, submits a video, and opens the video route', async ({ page }) => {
  let uploadCount = 0

  await page.route('http://localhost:8080/videos', async (route) => {
    if (route.request().method() !== 'POST') {
      await route.continue()
      return
    }

    uploadCount += 1

    await route.fulfill({
      json: {
        id: 'video-123',
        title: 'Trailer',
        status: 'uploaded',
      },
    })
  })

  await page.route('http://localhost:8080/videos/video-123', async (route) => {
    await route.fulfill({
      json: {
        id: 'video-123',
        title: 'Trailer',
        status: 'processing',
      },
    })
  })

  await page.goto('/')
  await expect(page).toHaveURL('/upload')
  await expect(page.getByRole('link', { name: 'Upload' })).toBeVisible()

  await page.getByLabel('Title').fill('Trailer')
  await page.getByLabel('File').setInputFiles(
    path.join(process.cwd(), 'tests/fixtures/sample.mp4'),
  )
  await page.getByRole('button', { name: 'Upload' }).click()

  await expect.poll(() => uploadCount).toBe(1)
  await expect(page.getByLabel('Title')).toHaveValue('')
  await expect(page.getByText(/Selected:/)).toHaveCount(0)

  await page.goto('/videos/video-123')
  await expect(page).toHaveURL('/videos/video-123')
  await expect(page.getByRole('heading', { name: 'Trailer' })).toBeVisible()
  await expect(page.getByText('Current status: processing')).toBeVisible()
  await expect(page.getByText('Video is not ready for playback.')).toBeVisible()

  await page.getByRole('link', { name: 'Upload' }).click()
  await expect(page).toHaveURL('/upload')
})
