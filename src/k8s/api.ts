import { CoreV1Api, KubeConfig, Metrics } from '@kubernetes/client-node'

const kc = new KubeConfig()
kc.loadFromDefault()

export const metricsClient = new Metrics(kc)

export const k8sApi = kc.makeApiClient(CoreV1Api)
