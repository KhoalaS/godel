import registerBasename from './Basename'
import registerToBytes from './ToBytes'
import registerSuffix from './Suffix'
import registerDisplay from './Display'

export function initHooks() {
  registerBasename()
  registerToBytes()
  registerSuffix()
  registerDisplay()
}
