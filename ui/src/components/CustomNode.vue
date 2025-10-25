<script lang="ts" setup>
import { HandleColors, type PipelineNode } from '@/models/Node'
import { Handle, Position, type NodeProps } from '@vue-flow/core'
import { computed } from 'vue'
import {
  WAutocomplete,
  WButton,
  WindowBody,
  WindowComponent,
  WInput,
  type WindowControls,
} from 'vue-98'
import { usePipelineStore } from '@/stores/pipeline'
import { useNodeUpdates } from '@/composables/useNodeUpdates'
import CodeArea from './CodeArea.vue'

// TODO split different cases for input components in sub-components

const props = defineProps<NodeProps<PipelineNode>>()
const store = usePipelineStore()
const vueFlow = store.vueFlow
const { onUpdate, onValueChange } = useNodeUpdates(props, store)

function onControlClick(ctrl: WindowControls) {
  if (ctrl == 'Close') {
    vueFlow.removeNodes(props.id)
  }
}

const _title = computed(() => {
  if (props.data.status == 'running' && props.data.progress) {
    return `${props.data.name} ${Math.floor(props.data.progress * 100)}%`
  }
  return props.data.name
})
</script>

<template>
  <WindowComponent
    @click:control="onControlClick"
    :title="_title"
    class="w-64 p-3"
    :controls="['Close']"
  >
    <template #body>
      <WindowBody class="m-2">
        <div class="flex flex-col items-start mt-1" :key="input.id" v-for="input in data.io">
          <label
            :for="props.data.id + input.id"
            class="text-xs"
            v-if="input.label && input.type != 'generated' && input.type != 'connected_only'"
            >{{ input.label }}</label
          >
          <div class="relative w-full">
            <Handle
              v-if="
                input.type == 'input' ||
                input.type == 'passthrough' ||
                input.type == 'connected_only'
              "
              class="handle handle-input"
              :class="{ disabled: input.disabled }"
              :key="input.id"
              type="target"
              :id="input.id"
              :position="Position.Left"
              :connectable-start="!input.disabled && false"
              :connectable-end="!input.disabled && true"
              :style="{ background: HandleColors[input.valueType] }"
            />
            <Handle
              v-if="
                input.type == 'output' || input.type == 'passthrough' || input.type == 'generated'
              "
              class="handle handle-output"
              :class="{ disabled: input.disabled }"
              :key="input.id"
              type="source"
              :id="input.id"
              :position="Position.Right"
              :connectable-start="!input.disabled && true"
              :connectable-end="!input.disabled && false"
              :style="{ background: HandleColors[input.valueType] }"
            />
            <div class="text-xs text-right" v-if="input.type == 'generated'">{{ input.label }}</div>
            <div class="text-xs text-left" v-else-if="input.type == 'connected_only'">
              {{ input.label }}
              <WInput
                style="display: none"
                :value="input.value"
                @update="(v) => onUpdate(v, input)"
                @value-change="onValueChange(input)"
                :type="input.valueType"
                :id="props.data.id + input.id"
              />
            </div>
            <div v-else-if="input.type == 'selection'">
              <WAutocomplete
                @update:model-value="(value) => onUpdate(value.name, input)"
                :initial="{
                  name: String(input.value),
                  id: String(input.value),
                }"
                :options="
                  input.options?.map((opt) => {
                    return { name: opt, id: opt }
                  })
                "
              ></WAutocomplete>
            </div>
            <div v-else-if="input.valueType == 'code'">
              <CodeArea
                :value="String(input.value ?? 'function abc() {\n\n}\n')"
                @update="(v) => onUpdate(v, input)"
                @value-change="onValueChange(input)"
                :type="input.valueType"
                :id="props.data.id + input.id"
              >
              </CodeArea>
            </div>
            <div v-else-if="input.valueType == 'directory' && input.type === 'output'">
              <div class="flex gap-1">
                <WInput
                  :value="input.value"
                  @update="(v) => onUpdate(v, input)"
                  @value-change="onValueChange(input)"
                  :id="props.data.id + input.id"
                >
                </WInput>
                <WButton>TODO</WButton>
              </div>
            </div>
            <WInput
              v-else
              :value="input.value"
              @update="(v) => onUpdate(v, input)"
              @value-change="onValueChange(input)"
              :type="input.valueType"
              :id="props.data.id + input.id"
            />
          </div>
        </div>
      </WindowBody>
    </template>
  </WindowComponent>
</template>

<style scoped>
.handle-input {
  left: -12px;
  position: absolute;
}

.handle-output {
  right: -12px;
  position: absolute;
}

.handle.disabled {
  filter: grayscale(1);
  opacity: 0.5;
  transition-duration: 100ms;
}

.vue-flow__handle {
  border-radius: unset;
  width: 8px;
  height: 8px;
}
</style>
