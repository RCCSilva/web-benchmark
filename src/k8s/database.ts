
import { k8sApi } from './api'
import { K8S_NAMESPACE, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER } from '../constants'


const createDatabasePod = async (name: string) => {
  return k8sApi.createNamespacedPod(K8S_NAMESPACE, {
    apiVersion: 'v1',
    kind: 'Pod',
    metadata: {
      name,
      labels: {
        app: name
      }
    },
    spec: {
      containers: [
        {
          name: 'postgres',
          image: 'postgres:14-alpine',
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
              name: 'POSTGRES_USER',
              value: POSTGRES_USER,
            },
            {
              name: 'POSTGRES_PASSWORD',
              value: POSTGRES_PASSWORD
            },
            {
              name: 'POSTGRES_DB',
              value: POSTGRES_DB,
            }
          ]
        }
      ]
    }
  })
}

const createDatabaseNodePort = async (
  { name, podName }: { name: string, podName: string }
) => {
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
          name: 'db-port',
          protocol: 'TCP',
          port: 5432,
          targetPort: 5432,
        }
      ]
    },
  })


  const databaseNodePort = service.body.spec?.ports?.[0].nodePort
  if (!databaseNodePort) {
    throw new Error('Failed to setup external-app-srv')
  }

  return {
    databaseNodePort
  }
}

const createDatabaseService = async (
  { name, podName }: { name: string, podName: string }
) => {
  return k8sApi.createNamespacedService(K8S_NAMESPACE, {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      name
    },
    spec: {
      selector: {
        app: podName
      },
      ports: [
        {
          name: 'database-port',
          protocol: 'TCP',
          port: 5432,
          targetPort: 5432
        }
      ]
    },
  })
}


const makeDeleteDatabase = (
  { podName, serviceName, externalServiceName }:
    { podName: string, serviceName: string, externalServiceName: string }
) => {
  return () => {
    return Promise.all([
      k8sApi.deleteNamespacedPod(podName, K8S_NAMESPACE),
      k8sApi.deleteNamespacedService(serviceName, K8S_NAMESPACE),
      k8sApi.deleteNamespacedService(externalServiceName, K8S_NAMESPACE)
    ])
  }
}

export const createDatabaseK8s = async (name: string) => {
  const podName = `${name}-db`
  const serviceName = `${name}-db-srv`
  const externalServiceName = `${name}-db-external-srv`

  await Promise.all([
    createDatabasePod(podName),
    createDatabaseService({ name: serviceName, podName }),
  ])

  const { databaseNodePort } = await createDatabaseNodePort({ name: externalServiceName, podName })

  return {
    databaseNodePort,
    deleteDatabase: makeDeleteDatabase({ podName, serviceName, externalServiceName }),
    serviceName
  }
}

