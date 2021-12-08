/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
import { promises as fs } from 'fs'
import { resolve } from 'path'
import { exec } from 'child_process'
import { promisify } from 'util'
const execPromise = promisify(exec)
import addApis from './add-apis.js'
const __dirname = new URL(import.meta.url).pathname

const SWAGGER_PATH =
  '../docs/openapi/openapi.yaml'
const GENERATED_DIR = 'src/lib/apis/generated'

const npx = process.platform === 'win32' ? 'npx.cmd' : 'npx'

const generateCmd = [
  npx,
  'openapi-generator-cli',
  'generate',
  '-g',
  'typescript-axios',
  '-i',
  SWAGGER_PATH,
  '-o',
  GENERATED_DIR
]

  ; (async () => {
    await fs.mkdir(resolve(__dirname, '../', GENERATED_DIR), {
      recursive: true
    })

    await execPromise(generateCmd.join(' '))

    // generate Apis class
    await addApis(GENERATED_DIR)
  })()
