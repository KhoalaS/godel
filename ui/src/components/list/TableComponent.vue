<script setup lang="ts" generic="U extends object">
import type { Column, Row, TypedColumn } from './Table'

interface TableProps<T> {
  rows: Row<T>[]
  columns: (Column | TypedColumn<T>)[]
}

const props = defineProps<TableProps<U>>()
</script>
<template>
  <table>
    <tr>
      <th :style="{ width: `${column.width}px` }" :key="column.key" v-for="column in props.columns">
        {{ column.headerName }}
      </th>
    </tr>
    <tbody>
      <tr :key="row.id" v-for="row in props.rows">
        <td v-for="value in row.data">
          {{ value }}
        </td>
      </tr>
    </tbody>
  </table>
</template>
