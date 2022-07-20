/* eslint-disable no-constant-condition */
import { Sequelize } from 'sequelize'
import { DEBUG } from './constants'
import { sleep } from './utils'

interface Props {
  user: string
  password: string
  host: string
  port: number
  name: string
}

export const createDatabaseTable = async ({ user, password, host, port, name }: Props) => {
  const url = `postgres://${user}:${password}@${host}:${port}/${name}`
  while (true) {
    try {
      const sequelize = new Sequelize(url)

      await sequelize.query(`
CREATE TABLE "users" (
  ID    BIGSERIAL         PRIMARY KEY,
  NAME  VARCHAR(64)                  ,
  AGE   INT                     
)`, { logging: false })
      await sequelize.connectionManager.close()
      break
    } catch (err) {
      if (DEBUG) {
        console.log(`   -> ${err}`)
      }
      await sleep(500)
    }
  }
}
