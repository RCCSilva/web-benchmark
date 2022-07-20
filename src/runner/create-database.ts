import { BenchmarkApp } from '../apps'
import { CLUSTER_IP, POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER } from '../constants'
import { createDatabaseTable } from '../db'
import { createDatabaseK8s } from '../k8s/database'

export const createDatabase = async (app: BenchmarkApp) => {
  console.log('   Creating database...')
  const { databaseNodePort, deleteDatabase, serviceName } = await createDatabaseK8s(app.slug)

  await createDatabaseTable(
    {
      user: POSTGRES_USER,
      password: POSTGRES_PASSWORD,
      host: CLUSTER_IP,
      port: databaseNodePort,
      name: POSTGRES_DB,
    }
  )

  return {
    deleteDatabase,
    serviceName
  }
}
