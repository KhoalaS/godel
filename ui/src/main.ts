import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

import 'vue-98/dist/main.css'
import './main.css'
import { initHooks } from './inputhooks'

const app = createApp(App)

app.config.errorHandler = (err, instance, info) => {
  console.error(err, instance, info)
}

initHooks()

app.use(createPinia())

app.mount('#app')
