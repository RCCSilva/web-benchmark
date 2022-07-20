/* eslint-disable @typescript-eslint/no-var-requires */
import util from 'util'

export const execAsync = util.promisify(require('child_process').exec)

export const sleep = (ms: number) => new Promise((resolve) => setTimeout(() => resolve(undefined), ms))
