import { createMemoryHistory, createRouter } from 'vue-router'

import PipelineView from './views/PipelineView.vue'
import NodeBuilderView from './views/NodeBuilderView.vue'

const routes = [
  { path: '/', component: PipelineView },
  { path: '/builder', component: NodeBuilderView },
]

export const router = createRouter({
  history: createMemoryHistory(),
  routes,
})
