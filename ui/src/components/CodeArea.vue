<script setup lang="ts">
import { useTemplateRef, watch } from 'vue'

const props = defineProps<{
  value?: string
  id?: string
}>()

const textAreaRef = useTemplateRef('textAreaRef')

const model = defineModel<string>({
  default: '',
})

const emit = defineEmits<{
  update: [value: string]
  valueChange: []
}>()

if (props.value != undefined) {
  model.value = props.value
}

watch(
  () => props.value,
  (newVal) => {
    if (newVal !== undefined && newVal !== model.value) {
      model.value = newVal
      emit('valueChange')
    }
  },
)

defineExpose({
  textAreaRef,
})
</script>
<template>
  <textarea
    @input="emit('update', model)"
    v-bind="$attrs"
    :id="id"
    ref="textAreaRef"
    v-model="model"
  ></textarea>
</template>
<style scoped="scoped">
textarea {
  width: 100%;
  height: 150px;
  padding: 0.75rem;
  box-sizing: border-box;
  border: 2px solid #ccc;
  border-radius: 4px;
  background-color: white;
  font-size: 1rem;
  resize: none;
}
</style>
