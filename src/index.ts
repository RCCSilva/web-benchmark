// eslint-disable-next-line @typescript-eslint/no-var-requires
require('dotenv').config()

import { loadApps } from './apps'
import { buildAndUploadImage } from './docker'
import { cleanUp } from './k8s/cleanup'
import { runAllBenchmarks } from './runner/run-test'


const main = async () => {
  const apps = await loadApps()
  await cleanUp()
  await buildAndUploadImage(apps)
  await runAllBenchmarks(apps)
  await cleanUp()
}

main()
  .then(() => console.log('Finished!'))
  .catch((err) => console.error('Failed', err))

process.on('SIGINT', async () => {
  await cleanUp()
  process.exit()
})

process.on('SIGTERM', async () => {
  await cleanUp()
  process.exit()
})
