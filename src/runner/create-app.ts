/* eslint-disable no-constant-condition */
import axios from 'axios'
import { BenchmarkApp } from '../apps'
import { CLUSTER_HOST, DEBUG, } from '../constants'
import { createAppK8s } from '../k8s/app'
import { sleep } from '../utils'

const verifyAppReady = async (nodePort: number) => {
  const url = `${CLUSTER_HOST}:${nodePort}`
  const doRequest = async () => {
    return axios.get(url)
  }

  while (true) {
    try {
      await doRequest()
      break
    } catch (err) {
      if (DEBUG && axios.isAxiosError(err)) {
        console.log(`   -> ${err.code} - ${JSON.stringify(err.response?.data)}`)
      }
      await sleep(500)
    }
  }
}

export const createApp = async ({ app, dbHost }: { app: BenchmarkApp, dbHost: string }) => {
  console.log('   Creating app...')
  const { deleteApp, appNodePort } = await createAppK8s({ app, dbHost })

  await verifyAppReady(appNodePort)

  return {
    deleteApp,
    appNodePort
  }
}
