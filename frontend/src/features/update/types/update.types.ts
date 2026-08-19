export interface UpdateInfo {
    available: boolean
    version: string
    notes: string
}

export type UpdateStatusState =
    | {kind: 'idle'}
    | {kind: 'checking'}
    | {kind: 'up-to-date'}
    | {kind: 'available'; info: UpdateInfo}
    | {kind: 'installing'}
    | {kind: 'error'; message: string}
