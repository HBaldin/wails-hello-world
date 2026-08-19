// Thin bridge over the generated Wails/Go bindings (see frontend/wailsjs).
// Feature-level `api/` modules call `callBackend` instead of the raw
// generated functions so every IPC failure surfaces as a typed error.

export class BackendError extends Error {
    constructor(operation: string, public readonly cause: unknown) {
        super(`${operation} failed: ${cause instanceof Error ? cause.message : String(cause)}`)
        this.name = 'BackendError'
    }
}

export async function callBackend<T>(operation: string, fn: () => Promise<T>): Promise<T> {
    try {
        return await fn()
    } catch (cause) {
        throw new BackendError(operation, cause)
    }
}
