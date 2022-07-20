import { K8S_NAMESPACE } from '../constants'
import { k8sApi } from './api'

export const cleanUp = async () => {
  const pods = await k8sApi.listNamespacedPod(K8S_NAMESPACE)
  const services = await k8sApi.listNamespacedService(K8S_NAMESPACE)

  console.log(`Cleaning up before exiting (${pods.body.items.length} pods, ${services.body.items.length} services)...`)

  const mapDeletePromises = pods.body.items.map(pod => {
    console.log(` Deleting pod ${pod.metadata?.name}...`)
    !!pod.metadata?.name &&
      k8sApi.deleteNamespacedPod(pod.metadata.name, K8S_NAMESPACE)
        .catch(() => console.log(` Failed to delete pod ${pod.metadata?.name}`))
  })

  const serviceDeletePromises = services.body.items.map(service => {
    console.log(` Deleting service ${service.metadata?.name}...`)
    !!service.metadata?.name &&
      k8sApi.deleteNamespacedService(service.metadata.name, K8S_NAMESPACE)
        .catch(() => console.log(`  Failed to delete service ${service.metadata?.name}`))
  })

  await Promise.all(mapDeletePromises)
  await Promise.all(serviceDeletePromises)
}
