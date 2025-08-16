<script setup lang="ts">
import { ref } from 'vue'
import JobsTable from './components/JobsTable.vue'
import { useJobStore } from './stores/job'
import {
  TaskbarGroupheader,
  TitlebarIcon,
  WAutocomplete,
  WButton,
  WindowBody,
  WindowComponent,
  WInput,
} from 'vue-98'
import type { Config } from './types/Config'

const url = ref('')
const jobStore = useJobStore()
const config = ref<Config | undefined>()
const noneOption: Config = {
  deleteOnCancel: false,
  destPath: '',
  id: 'NONE',
  name: 'None',
  transformer: [],
}

function download() {
  const parsedUrl = URL.parse(url.value)
  if (!parsedUrl) {
    console.warn('invalid url')
    return
  }

  const configId = config.value?.id != noneOption.id ? config.value : undefined
  jobStore.addJob(url.value, configId?.id)
  url.value = ''
}
</script>

<template>
  <Suspense>
    <!-- component with nested async dependencies -->
    <WindowComponent
      title="Downloads"
      :controls="['Minimize', 'Maximize', 'Close']"
      style="width: 800px; height: 600px"
    >
      <template #title-icon>
        <TitlebarIcon icon="document"></TitlebarIcon>
      </template>
      <template #body>
        <WindowBody>
          <template #toolbars>
            <div style="display: flex; gap: 2px; align-items: center">
              <TaskbarGroupheader></TaskbarGroupheader>
              <label for="url" aria-label="Url">Url</label>
              <WInput v-model="url" id="url" />
              <WButton @click="download">Download</WButton>
            </div>
            <div>
              <WAutocomplete
                style="width: 200px"
                v-model="config"
                :options="jobStore.configs"
                :none-option="noneOption"
              >
              </WAutocomplete>
            </div>
          </template>
          <JobsTable />
        </WindowBody>
      </template>
    </WindowComponent>

    <!-- loading state via #fallback slot -->
    <template #fallback> Loading... </template>
  </Suspense>
</template>

<style scoped></style>
