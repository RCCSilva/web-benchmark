import autocannon from 'autocannon'
import fs from 'fs'
import { BenchmarkApp } from './apps'

export const processResult = (app: BenchmarkApp, result: autocannon.Result) => {
  const language = `${app.language.name} (${app.language.version})`
  const packages = app.packages.map(p => p.name).join(' / ')
  const hasDatabase = `${app.database ? '✓' : '✗'}`
  const averageLatency = result.latency.average.toFixed(2)
  const medianLatency = result.latency.p50.toFixed(2)
  const averageRequestsPerMinute = result.requests.average.toFixed(2)
  const medianRequestsPerMinute = result.requests.p50.toFixed(2)
  const successRate = `${((1 - (result.non2xx / result.requests.total)) * 100).toFixed(2)}%`
  const text =
    `| ${language} | ${packages} | ${hasDatabase} | ${averageLatency} | ${medianLatency} | ${averageRequestsPerMinute} | ${medianRequestsPerMinute} | ${successRate} | \n`

  fs.appendFileSync('RESULT.md', text)

}

export const buildResultHeader = () => {
  fs.writeFileSync('RESULT.md', '')
  fs.appendFileSync('RESULT.md', `| Language | Packages | Database | Average Latency (ms) | Median Latency (ms) | Average Requests/s | Median Requests/s |Success Rate |
|----------|-----------------|----------|------------:|------------:|-------------:| -----------:| -------------:|
`)
}

export const copyResultToResults = () => {
  fs.copyFileSync('RESULT.md', `./results/${new Date().toISOString().slice(0,19)}_RESULT.md`)
}
