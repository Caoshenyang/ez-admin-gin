import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const rootDir = dirname(dirname(fileURLToPath(import.meta.url)))
const migrationsDir = join(rootDir, 'server/migrations')
const fullInitializationVersion = 1

for (const dialect of ['mysql', 'postgres']) {
  const outputFile = join(migrationsDir, dialect, 'full_schema_and_seed.sql')
  writeFileSync(outputFile, normalizeFullSql(readFileSync(outputFile, 'utf8'), dialect))
  console.log(`normalized ${outputFile}`)
}

function normalizeFullSql(sql, dialect) {
  const baseline =
    dialect === 'mysql'
      ? /VALUES \(\d+, 0\)/
      : /VALUES \(\d+, false\)/
  const replacement =
    dialect === 'mysql'
      ? `VALUES (${fullInitializationVersion}, 0)`
      : `VALUES (${fullInitializationVersion}, false)`

  return `${sql.replace(baseline, replacement).trim()}\n`
}
