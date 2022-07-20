import crypto from 'crypto'
import * as yup from 'yup'
import { readdir, readFile } from 'fs/promises'
import path from 'path'
import { IMAGE_PREFIX } from './constants'

const tagVersion = crypto.randomBytes(5).toString('hex')

export type BenchmarkApp = {
  tag: string
  image: string
  path: string
  slug: string
  name: string
  database: boolean
  language: {
    name: string
    version: string
  },
  packages: {
    name: string
    version: string
  }[],
}

type MaybeBenchmarkApp = {
  ok: true;
  data: BenchmarkApp;
} | {
  ok: false;
  data?: undefined
}

const benchmark = yup.object({
  language: yup.object({
    name: yup.string(),
    version: yup.string()
  }),
  packages: yup.array(
    yup.object({
      name: yup.string(),
      version: yup.string()
    })
  ),
  database: yup.boolean()
})

const loadApps = async () => {
  const maybeBenchmarks = await loadAllBenchmarkFiles()
  return maybeBenchmarks
    .filter(benchmark => benchmark.ok === true)
    // .filter(benchmark => benchmark.data?.name.includes('Entity'))
    .map(benchmark => benchmark.data) as unknown as BenchmarkApp[]
}


const loadAllBenchmarkFiles = async (): Promise<MaybeBenchmarkApp[]> => {

  const directoryPath = path.join(__dirname, '../apps')
  const appsList = await readdir(directoryPath)

  const promises = appsList.map(loadOneBenchmarkFile)

  return Promise.all(promises)
}

const loadOneBenchmarkFile = async (app: string): Promise<MaybeBenchmarkApp> => {
  const appPath = path.join(__dirname, '../apps', app)
  const benchmarkSpecPath = path.join(__dirname, '../apps', app, 'benchmark.json')
  const file = await readFile(benchmarkSpecPath, { encoding: 'utf-8' })
  const attempt = JSON.parse(file)

  try {
    const validatedApp = await benchmark.validate(attempt)
    const slug = validatedApp.packages?.map(p => p.name).join('').toLowerCase()

    const packageNames = validatedApp.packages?.map(p => `${p.name} (${p.version})`).join(' / ')
    const tag = `web-benchmark-ecr:${slug}${tagVersion}`.toLowerCase()

    const benchmarkApp = {
      ...validatedApp,
      path: appPath,
      slug,
      name: `${validatedApp.language.name} (${validatedApp.language.version}) / ${packageNames}`,
      tag,
      image: `${IMAGE_PREFIX}${tag}`
    } as BenchmarkApp

    return {
      ok: true,
      data: benchmarkApp
    }
  } catch (err) {
    return {
      ok: false
    }
  }
}

export {
  loadApps
}
