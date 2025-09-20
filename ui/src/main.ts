import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

import 'vue-98/dist/main.css'
import './main.css'
import { initHooks } from './inputhooks'
import type { ErrorService } from './error/ErrorService'
import { DefaultErrorService } from './error/DefaultErrorService'
import { ErrorServiceKey } from './InjectionKeys'
import { router } from './router'

const app = createApp(App)

app.config.errorHandler = (err, instance, info) => {
  console.error(err, instance, info)
}

initHooks()

const errorService: ErrorService = DefaultErrorService
app.provide(ErrorServiceKey, errorService)

app.config.errorHandler = (err, vm, info) => console.error(err, vm, info)

app.use(createPinia())
app.use(router)

app.mount('#app')
