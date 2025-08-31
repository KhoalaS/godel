<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue'
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
import PipelineBuilder from './components/PipelineBuilder.vue'

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

const noneTransformer = {
  id: 'NONE',
  name: 'None',
}
const transformer = ref<{ id: string; name: string }>(noneTransformer)
const pipelineBuilder = useTemplateRef('pipelineBuilder')

function download() {
  const parsedUrl = URL.parse(url.value)
  if (!parsedUrl) {
    console.warn('invalid url')
    return
  }

  const configId = config.value?.id != noneOption.id ? config.value : undefined
  const tr: string[] = []
  if (transformer.value.id != 'NONE') {
    tr.push(transformer.value.id)
  }

  jobStore.addJob(url.value, configId?.id, tr)
  url.value = ''
}

onMounted(async () => {
  await jobStore.init()
})
</script>

<template>
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
            <WAutocomplete
              style="width: 200px"
              v-model="transformer"
              :options="
                jobStore.transformers.map((tr) => {
                  return {
                    id: tr,
                    name: tr,
                  }
                })
              "
              :none-option="noneTransformer"
            >
            </WAutocomplete>
          </div>
        </template>
        <JobsTable />
      </WindowBody>
    </template>
  </WindowComponent>
  <WindowComponent title="Pipeline Builder" :controls="['Minimize', 'Maximize']">
    <template #body>
      <WindowBody>
        <template #toolbars>
          <div class="flex gap-1 items-center py-1 px-[2px] justify-start">
            <WButton @click="pipelineBuilder?.saveGraph()">Start</WButton>
          </div>
        </template>
        <PipelineBuilder ref="pipelineBuilder"></PipelineBuilder>
      </WindowBody>
    </template>
  </WindowComponent>
</template>

<style scoped></style>
