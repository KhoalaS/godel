import registerBasename from './Basename'
import registerToBytes from './ToBytes'
import registerSuffix from './Suffix'

export function initHooks() {
  registerBasename()
  registerToBytes()
  registerSuffix()
}
