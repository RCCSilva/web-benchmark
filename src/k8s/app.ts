import { BenchmarkApp } from '../apps'
import { K8S_NAMESPACE, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER } from '../constants'
import { k8sApi } from './api'


const createAppPod = async (
  { name, image, dbHost }: { name: string, image: string, dbHost: string }
) => {
  await k8sApi.createNamespacedPod(K8S_NAMESPACE, {
    apiVersion: 'v1',
    kind: 'Pod',
    metadata: {
      name: name,
      labels: {
        app: name
      }
    },
    spec: {
      containers: [
        {
          name: name,
          image: image,
          resources: {
            requests: {
              memory: '1Gi',
              cpu: '1000m'
            },
            limits: {
              memory: '1Gi',
              cpu: '1000m'
            }
          },
          env: [
            {
              name: 'DB_USER',
              value: POSTGRES_USER,
            },
            {
              name: 'DB_PASSWORD',
              value: POSTGRES_PASSWORD,
            },
            {
              name: 'DB_NAME',
              value: POSTGRES_DB,
            },
            {
              name: 'DB_HOST',
              value: dbHost
            },
            {
              name: 'DB_PORT',
              value: '5432'
            }
          ]
        }
      ]
    }
  })
}


const createAppNodePort = async ({ name, podName }: { name: string, podName: string }) => {
  const service = await k8sApi.createNamespacedService(K8S_NAMESPACE, {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      name
    },
    spec: {
      selector: {
        app: podName
      },
      type: 'NodePort',
      ports: [
        {
          name: 'app-port',
          protocol: 'TCP',
          port: 3000,
          targetPort: 3000,
        }
      ]
    },
  })


  const appNodePort = service.body.spec?.ports?.[0]?.nodePort
  if (!appNodePort) {
    throw new Error('Failed to setup external-app-srv')
  }

  return {
    appNodePort
  }
}

const makeDeleteApp = (
  { podName, serviceName }: { podName: string, serviceName: string }
) => {
  return () => {
    return Promise.all([
      k8sApi.deleteNamespacedPod(podName, K8S_NAMESPACE),
      k8sApi.deleteNamespacedService(serviceName, K8S_NAMESPACE)
    ])
  }
}

export const createAppK8s = async ({ app, dbHost }: { app: BenchmarkApp, dbHost: string }) => {
  const podName = `${app.slug}-app`
  const serviceName = `${app.slug}-srv`

  await createAppPod({ name: podName, image: app.image, dbHost })
  const { appNodePort } = await createAppNodePort({ name: serviceName, podName })

  return {
    appNodePort,
    deleteApp: makeDeleteApp({ podName, serviceName })
  }
}
