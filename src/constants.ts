if (!process.env.CLUSTER_HOST) {
  throw new Error('CLUSTER_HOST not defined')
}
export const CLUSTER_HOST = process.env.CLUSTER_HOST
export const CLUSTER_IP = CLUSTER_HOST.replace(/https:\/\/|http:\/\//, '')

export const DEBUG = process.env.DEBUG === '1'
export const TEST = process.env.TEST === '1'

export const TEST_DURATION = TEST ? '5' : '40'

export const POSTGRES_USER = 'user'
export const POSTGRES_PASSWORD = 'Password!234'
export const POSTGRES_DB = 'mydb'
export const K8S_NAMESPACE = 'default'

export const IMAGE_PREFIX = process.env.IMAGE_PREFIX || ''
