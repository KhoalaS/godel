import { createWebHistory, createRouter } from 'vue-router'

import PipelineView from './views/PipelineView.vue'
import NodeBuilderView from './views/NodeBuilderView.vue'

const routes = [
  { path: '/', component: PipelineView },
  { path: '/builder', component: NodeBuilderView, name: 'nodeBuilder' },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
