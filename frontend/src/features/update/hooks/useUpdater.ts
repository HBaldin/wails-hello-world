import {useCallback, useState} from 'react'
import {checkForUpdate, installUpdate} from '../api/updateApi'
import type {UpdateStatusState} from '../types/update.types'

export function useUpdater() {
    const [status, setStatus] = useState<UpdateStatusState>({kind: 'idle'})

    const check = useCallback(async () => {
        setStatus({kind: 'checking'})
        try {
            const info = await checkForUpdate()
            setStatus(info.available ? {kind: 'available', info} : {kind: 'up-to-date'})
        } catch (err) {
            setStatus({kind: 'error', message: err instanceof Error ? err.message : String(err)})
        }
    }, [])

    const install = useCallback(async () => {
        setStatus({kind: 'installing'})
        try {
            // On success the app relaunches itself; there is nothing left to update here.
            await installUpdate()
        } catch (err) {
            setStatus({kind: 'error', message: err instanceof Error ? err.message : String(err)})
        }
    }, [])

    return {status, check, install}
}
