import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

import 'vue-98/dist/main.css'
import './main.css'
import type { ErrorService } from './error/ErrorService'
import { DefaultErrorService } from './error/DefaultErrorService'
import { ErrorServiceKey, HookFunctionServiceKey, NotificationServiceKey } from './InjectionKeys'
import { router } from './router'
import type { HookFunctionService } from './inputhooks/HookFunctionService'
import { HookFunctionServiceImpl } from './inputhooks/HookFunctionServiceImpl'
import { BaseNameHook } from './inputhooks/Basename'
import { DisplayHook } from './inputhooks/Display'
import { SuffixHook } from './inputhooks/Suffix'
import { ToBytesHook } from './inputhooks/ToBytes'
import { NotificationServiceImpl } from './services/NotificationServiceImpl'
import { Container } from '@n8n/di'

const app = createApp(App)

app.config.errorHandler = (err, instance, info) => {
  console.error(err, instance, info)
}

const errorService: ErrorService = DefaultErrorService
const hookFunctionService: HookFunctionService = Container.get(HookFunctionServiceImpl)
hookFunctionService.register('basename', BaseNameHook)
hookFunctionService.register('display', DisplayHook)
hookFunctionService.register('suffix', SuffixHook)
hookFunctionService.register('toBytes', ToBytesHook)

const notificationService = Container.get(NotificationServiceImpl)

app.config.errorHandler = (err, vm, info) => {
  console.error(err, vm, info)
  notificationService.showNotification({
    message: (err as Error).message,
    title: 'Error',
    level: 'error',
    durationInSeconds: 15,
  })
}

app.use(createPinia())
app.use(router)

app.mount('#app')
