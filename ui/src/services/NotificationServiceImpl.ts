import { h, createVNode, render } from 'vue'
import { type NotificationConfig, type NotificationService } from './NotificationService'
import ToastComponent from '@/components/notification/ToastComponent.vue'
import { Service } from '@n8n/di'

@Service()
export class NotificationServiceImpl implements NotificationService {
  readonly durationInSeconds: number

  constructor(durationInSeconds: number = 10) {
    this.durationInSeconds = durationInSeconds
  }

  showNotification(config: NotificationConfig): void {
    const wrapper = document.createElement('div')
    document.body.appendChild(wrapper)

    const node = createVNode(ToastComponent, {
      ...config,
      onClose: () => {
        wrapper.remove()
      },
    })

    render(node, wrapper)
  }
}
