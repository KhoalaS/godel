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
import { WindowBody, WindowComponent, WInput, type WindowControls } from 'vue-98'
const props = defineProps<NodeProps<PipelineNode>>()

const { updateNodeData, removeNodes } = useVueFlow()

const sourceConnections = useNodeConnections({
  // type target means all connections where *this* node is the target
  // that means we go backwards in the graph to find the source of the connection(s)
  handleType: 'target',
})

const sourceData = useNodesData<GraphNode<PipelineNode>>(() =>
  sourceConnections.value.map((connection) => connection.source),
)

function onUpdate(value: string | number | boolean, io: NodeIO) {
  if (!io.readOnly && props.data.io) {
    updateNodeData<PipelineNode>(props.id, {
      io: {
        ...props.data.io,
        [io.id]: { ...props.data.io?.[io.id], value: value },
      },
    })

    if (io.hooks) {
      hook(value, io.hooks)
    }
  }
}

function hasIncoming(inputId: string) {
  const targetConn = sourceConnections.value.find((conn) => {
    return conn.targetHandle == inputId
  })

  return Boolean(targetConn)
}

function getIncomingData(inputId: string) {
  const targetConn = sourceConnections.value.find((conn) => {
    return conn.targetHandle == inputId
  })

  if (targetConn && targetConn.targetHandle) {
    const node = sourceData.value?.find((node) => {
      return node.id == targetConn.source
    })
    if (node) {
      return node.data.io?.[targetConn.sourceHandle!].value!
    }
  }

  return ''
}

function hook(input: string | number | boolean, overwrites: Record<string, string>) {
  for (const [inputId, functionId] of Object.entries(overwrites)) {
    const func = FunctionRegistry.get(functionId)
    if (func == undefined) {
      continue
    }

    if (props.data.io?.[inputId].value != undefined) {
      return
    }

    const newValue = func(input)

    if (props.data.io) {
      updateNodeData<PipelineNode>(props.id, {
        io: {
          ...props.data.io,
          [inputId]: { ...props.data.io?.[inputId], value: newValue },
        },
      })
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
    class="w-64 text-center p-3"
    :controls="['Close']"
  >
    <template #body>
      <WindowBody class="m-2">
        <div class="flex flex-col items-start mt-1" :key="input.id" v-for="input in data.io">
          <label
            :for="props.data.id + input.id"
            class="text-xs"
            v-if="input.label && input.type != 'generated'"
            >{{ input.label }}</label
          >
          <div class="relative w-full">
            <Handle
              v-if="input.type == 'input' || input.type == 'passthrough'"
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
            <WInput
              v-else-if="!hasIncoming(input.id)"
              :initial="input.value"
              @update="(v) => onUpdate(v, input)"
              :type="input.valueType"
              :id="props.data.id + input.id"
            />
            <WInput
              v-else
              :value="getIncomingData(input.id)"
              :initial="getIncomingData(input.id)"
              @update="(v) => onUpdate(v, input)"
              :disabled="true"
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
