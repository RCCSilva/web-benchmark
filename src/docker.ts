import { BenchmarkApp } from './apps'
import { execAsync } from './utils'

export const buildAndUploadImage = async (apps: BenchmarkApp[]) => {
  for await (const app of apps) {
    try {
      console.log(`Building ${app.name}`)
      await execAsync(`docker build -t ${app.tag} ${app.path}`)
    } catch (err) {
      console.error(`Failed to build ${app.tag}`)
      throw err
    }
  }
}
