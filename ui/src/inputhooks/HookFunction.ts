/**
 * InputHooks are small functions that run when the value of an input X changes. During definition of a node
 * an input can declare itself as a hook paramter for a specific function.
 *
 * E.g. an input called "Url" can say that it wants to take on the role of the "path" parameter of the "Pathname" hook function.
 */
export type HookFunction = (arg: Record<string, unknown>) => string | number | boolean
