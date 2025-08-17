import { ref, readonly } from 'vue'
import { defineStore } from 'pinia'
import { DownloadJob } from '@/types/DownloadJob'
import z from 'zod'
import { Config } from '@/types/Config'

export const useJobStore = defineStore('job', () => {
  const baseUrl = 'localhost:9095'
  const jobs = ref<DownloadJob[]>([])
  const configs = ref<Config[]>([])
  const transformers = ref<string[]>([])

  async function init() {
    try {
      const configResponse = await fetch(`http://${baseUrl}/configs`)
      if (configResponse.status != 200) {
        console.warn('could not get configs, got status code', configResponse.status)
      } else {
        const data = await configResponse.json()
        configs.value = z.array(Config).parse(data)
      }

      const trResponse = await fetch(`http://${baseUrl}/transformers`)
      if (trResponse.status != 200) {
        console.warn('could not get transformers, got status code', trResponse.status)
        return
      } else {
        const data = await trResponse.json()
        transformers.value = z.array(z.string()).parse(data)
      }

      const response = await fetch(`http://${baseUrl}/jobs`)
      if (response.status != 200) {
        console.warn('could not get jobs, got status code', response.status)
        return
      }

      const data = await response.json()
      jobs.value = z.array(DownloadJob).parse(data)
    } catch (e: unknown) {
      console.warn(e)
    }

    initWs()
  }

  async function initWs() {
    const socket = new WebSocket(`ws://${baseUrl}/updates/jobinfo`)
    // Connection opened
    socket.addEventListener('open', () => {
      console.log('Connection opened')
    })

    // Listen for messages
    socket.addEventListener('message', (event) => {
      try {
        const jobData = JSON.parse(event.data)
        const job = DownloadJob.parse(jobData)

        const targetIdx = jobs.value.findIndex((j) => j.id == job.id)
        if (targetIdx != -1) {
          jobs.value[targetIdx] = job
        } else {
          jobs.value.push(job)
        }
      } catch (e: unknown) {
        console.warn(e)
      }
    })
  }

  async function addJob(url: string, configId?: string, transformer?: string[]) {
    const job: DownloadJob = {
      url: url,
      id: '-1',
      configId: configId,
      transformer: transformer
    }

    try {
      const response = await fetch(`http://${baseUrl}/add`, {
        method: 'POST',
        body: JSON.stringify(job),
      })

      if (response.status != 200) {
        console.warn('could not get jobs, got status code', response.status)
        return
      }
    } catch (e: unknown) {
      console.warn(e)
    }
  }

  async function pauseJob(id: string) {
    try {
      const response = await fetch(`http://${baseUrl}/pause`, {
        method: 'POST',
        body: id,
      })

      if (response.status != 200) {
        console.warn('could not pause job, got status code', response.status)
        return
      }
    } catch (e: unknown) {
      console.warn(e)
    }
  }

  return { init, addJob, pauseJob, jobs: readonly(jobs), configs: configs, transformers }
})
