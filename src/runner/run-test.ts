import autocannon from 'autocannon'
import { BenchmarkApp } from '../apps'
import { CLUSTER_HOST, TEST, TEST_DURATION } from '../constants'
import { buildResultHeader, copyResultToResults, processResult } from '../result'
import { sleep } from '../utils'
import { createApp } from './create-app'
import { createDatabase } from './create-database'

const warmUp = async (nodePort: number) => {
  console.log('   Warming up...')
  await autocannon({
    url: `${CLUSTER_HOST}:${nodePort}`,
    connections: 50,
    pipelining: 10,
    duration: TEST_DURATION
  })
}

const benchmark = async (nodePort: number) => {
  console.log('   Running benchmark...')
  return autocannon({
    url: `${CLUSTER_HOST}:${nodePort}`,
    connections: 50,
    pipelining: 10,
    duration: TEST_DURATION,
  })
}

const runAppBenchmark = async (app: BenchmarkApp) => {
  const { deleteDatabase, serviceName } = await createDatabase(app)
  const { deleteApp, appNodePort } = await createApp({ app, dbHost: serviceName })

  if (!TEST) {
    await warmUp(appNodePort)
  }

  const benchmarkResult = await benchmark(appNodePort)

  processResult(app, benchmarkResult)
  return Promise.all([deleteApp(), deleteDatabase()])

}

export const runAllBenchmarks = async (apps: BenchmarkApp[]) => {
  buildResultHeader()

  for await (const app of apps) {
    console.log(`Testing ${app.name}`)
    await runAppBenchmark(app)
    await sleep(10000)
  }
  copyResultToResults()
}
