import { expect, test } from '@playwright/test'
import path from 'node:path'

test('uploads a video and opens the player page', async ({ page }) => {
  await page.goto('/')

  await page.getByLabel('Title').fill('Trailer')
  await page.getByLabel('File').setInputFiles(
    path.join(process.cwd(), 'tests/fixtures/sample.mp4'),
  )
  await page.getByRole('button', { name: 'Upload' }).click()

  await expect(page.getByText('Trailer')).toBeVisible()
  await expect(page.getByText('ready')).toBeVisible({ timeout: 30000 })

  await page.getByRole('link', { name: 'Trailer' }).click()
  await expect(page).toHaveURL(/\/videos\//)
  await expect(page.locator('video')).toBeVisible()
})
