import { ref, readonly } from 'vue'
import { defineStore } from 'pinia'
import { DownloadJob } from '@/types/DownloadJob'
import z from 'zod'

export const useJobStore = defineStore('job', () => {
  const jobs = ref<DownloadJob[]>([])
  const baseUrl = 'http://localhost:9095'

  async function init() {
    try {
      const response = await fetch(`${baseUrl}/jobs`)
      if (response.status != 200) {
        console.warn('could not get jobs, got status code', response.status)
        return
      }

      const data = await response.json()

      const _jobs = z.array(DownloadJob).parse(data)
      jobs.value = _jobs
    } catch (e: unknown) {
      console.warn(e)
    }
  }

  return { init, jobs: readonly(jobs) }
})
