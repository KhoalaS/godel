<script setup lang="ts">
import { WindowBody, WindowComponent } from 'vue-98'
</script>

<template>
  <WindowComponent class="h-full w-full" title="Node Builder" :controls="['Minimize', 'Maximize']">
    <template #body>
      <WindowBody>
        <template #toolbars> </template>
      </WindowBody>
    </template>
  </WindowComponent>
</template>

<style scoped></style>
