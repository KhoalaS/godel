<script lang="ts" setup>
import { FunctionRegistry } from '@/types/InputHook'
import { HandleColors, type NodeIO, type PipelineNode } from '@/types/Node'
import { Handle, Position, useNodeConnections, useVueFlow, type NodeProps } from '@vue-flow/core'
import { computed } from 'vue'
import {
  WAutocomplete,
  WButton,
  WindowBody,
  WindowComponent,
  WInput,
  type WindowControls,
} from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const { updateNodeData, removeNodes, findNode } = useVueFlow()

const targetConnections = useNodeConnections({
  handleType: 'source',
})

function onUpdate(value: string | number | boolean, io: NodeIO) {
  console.log('CustomNode:onUpdate', io.id)
  if (!io.readOnly && props.data.io != undefined) {
    const hookUpdates: Record<string, NodeIO> = {}

    if (io.hooks) {
      for (const [hookId, functionId] of Object.entries(io.hooks)) {
        const func = FunctionRegistry.get(functionId)
        if (func == undefined) {
          continue
        }

        const hookValues: Record<string, unknown> = {}
        for (const _io of Object.values(props.data.io ?? {})) {
          if (_io.hookMapping?.[hookId] != undefined) {
            hookValues[_io.hookMapping[hookId]] = _io.id == io.id ? value : _io.value
          }
        }

        const newValue = func(hookValues)

        if (props.data.io?.[hookId].value == newValue) {
          continue
        }

        if (props.data.io) {
          hookUpdates[hookId] = {
            ...props.data.io?.[hookId],
            value: newValue,
          }
        }
      }
    }

    updateNodeData<PipelineNode>(props.id, {
      io: {
        ...props.data.io,
        [io.id]: { ...io, value: value },
        ...hookUpdates,
      },
    })

    const inputs: {
      inputId: string
      newValue: string | number | boolean | undefined
    }[] = [{ inputId: io.id, newValue: value }]

    for (const [id, io] of Object.entries(hookUpdates)) {
      inputs.push({ inputId: id, newValue: io.value })
    }
    updateTargetNodes(inputs)
  }
}

function onValueChange(io: NodeIO) {
  console.log('CustomNode:onValueChange', io.id)

  updateTargetNodes([
    {
      inputId: io.id,
      newValue: io.value,
    },
  ])
  if (io.value != undefined && io.hooks) {
    hook(io)
  }
}

function updateTargetNodes(
  inputs: {
    inputId: string
    newValue: string | number | boolean | undefined
  }[],
) {
  console.log('CustomNode:updateTargetNodes')

  const data = new Map<string, Record<string, NodeIO>>()

  for (const { inputId, newValue } of inputs) {
    if (newValue == undefined) {
      continue
    }

    const conns = targetConnections.value.filter((conn) => {
      return conn.sourceHandle == inputId
    })

    for (const conn of conns) {
      if (!conn.targetHandle) {
        return
      }

      const targetNode = findNode<PipelineNode>(conn.target)
      if (!targetNode || !targetNode.data.io) {
        continue
      }

      if (!data.has(conn.target)) {
        data.set(conn.target, {})
      }

      data.set(conn.target, {
        ...data.get(conn.target),
        [conn.targetHandle]: {
          ...targetNode.data.io[conn.targetHandle],
          value: newValue,
        },
      })
    }
  }

  for (const [nodeId, updates] of data.entries()) {
    const targetNode = findNode<PipelineNode>(nodeId)

    if (targetNode?.data.io) {
      updateNodeData<PipelineNode>(nodeId, {
        io: {
          ...targetNode.data.io,
          ...updates,
        },
      })
    }
  }
}

function hook(io: NodeIO) {
  console.log('CustomNode:hook', io.value)

  const hookUpdates: Record<string, NodeIO> = {}

  for (const [hookId, functionId] of Object.entries(io.hooks ?? {})) {
    const func = FunctionRegistry.get(functionId)
    if (func == undefined) {
      continue
    }

    const hookValues: Record<string, unknown> = {}
    for (const io of Object.values(props.data.io ?? {})) {
      if (io.hookMapping?.[hookId] != undefined) {
        hookValues[io.hookMapping[hookId]] = io.value
      }
    }

    const newValue = func(hookValues)

    if (props.data.io?.[hookId].value == newValue) {
      continue
    }

    if (props.data.io) {
      hookUpdates[hookId] = {
        ...props.data.io?.[hookId],
        value: newValue,
      }
    }
  }

  if (props.data.io) {
    updateNodeData<PipelineNode>(props.id, {
      io: {
        ...props.data.io,
        ...hookUpdates,
      },
    })
  }

  const inputs: {
    inputId: string
    newValue: string | number | boolean | undefined
  }[] = []
  for (const [id, io] of Object.entries(hookUpdates)) {
    inputs.push({ inputId: id, newValue: io.value })
  }
  updateTargetNodes(inputs)
}

function onControlClick(ctrl: WindowControls) {
  if (ctrl == 'Close') {
    removeNodes(props.id)
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
