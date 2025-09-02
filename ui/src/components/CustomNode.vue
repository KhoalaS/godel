<script lang="ts" setup>
import { FunctionRegistry } from '@/types/InputHook'
import { HandleColors, type NodeIO, type PipelineNode } from '@/types/Node'
import {
  Handle,
  Position,
  useNodeConnections,
  useNodesData,
  useVueFlow,
  type GraphNode,
  type NodeProps,
} from '@vue-flow/core'
import { WAutocomplete, WindowBody, WindowComponent, WInput, type WindowControls } from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const { updateNodeData, removeNodes, findNode } = useVueFlow()

const targetConnections = useNodeConnections({
  handleType: 'source',
})

function updateSelf(value: string | number | boolean, io: NodeIO) {
  if (!io.readOnly && props.data.io != undefined) {
    const hookUpdates: Record<string, NodeIO> = {}

    if (io.hooks) {
      for (const [hookId, functionId] of Object.entries(io.hooks)) {
        const func = FunctionRegistry.get(functionId)
        if (func == undefined) {
          continue
        }

        const newValue = func(value)

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
  }

  updateTargetNodes(io.id, value)
}

function onUpdate(value: string | number | boolean, io: NodeIO) {
  updateSelf(value, io)
  updateTargetNodes(io.id, value)
}

function onValueChange(io: NodeIO) {
  if (io.value != undefined) {
    if (io.hooks) {
      hook(io.value, io.hooks)
    }
    updateTargetNodes(io.id, io.value)
  }
}

function updateTargetNodes(inputId: string, newValue: string | number | boolean) {
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

    updateNodeData<PipelineNode>(targetNode.id, {
      io: {
        ...targetNode.data.io,
        [conn.targetHandle]: {
          ...targetNode.data.io[conn.targetHandle],
          value: newValue,
        },
      },
    })
  }
}

function hook(input: string | number | boolean, overwrites: Record<string, string>) {
  const hookUpdates: Record<string, NodeIO> = {}

  for (const [hookId, functionId] of Object.entries(overwrites)) {
    const func = FunctionRegistry.get(functionId)
    if (func == undefined) {
      continue
    }

    const newValue = func(input)

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
    const data = useNodesData<GraphNode<PipelineNode>>(props.id)
    if (data.value) {
      updateNodeData<PipelineNode>(props.id, {
        io: {
          ...data.value.data.io,
          ...hookUpdates,
        },
      })
    }
  }
  for (const [id, io] of Object.entries(hookUpdates)) {
    if (io.value != undefined) {
      updateTargetNodes(id, io.value)
    }
  }
}

function onControlClick(ctrl: WindowControls) {
  if (ctrl == 'Close') {
    removeNodes(props.id)
  }
}
</script>

<template>
  <WindowComponent
    @click:control="onControlClick"
    :title="data.name"
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
            </div>
            <div v-else-if="input.type == 'selection'">
              <WAutocomplete
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
            <WInput
              v-else
              :value="input.value"
              @update="(v) => onUpdate(v, input)"
              @value-change="onValueChange(input)"
              :type="input.valueType"
              :id="props.data.id + input.id"
            />
            {{ input.value }}
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
