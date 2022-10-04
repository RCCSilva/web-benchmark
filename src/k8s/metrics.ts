import { K8S_NAMESPACE } from '../constants'
import { metricsClient } from './api'

export const getMetrics = () => {
  metricsClient.getPodMetrics(K8S_NAMESPACE)
}
