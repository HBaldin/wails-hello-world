import type {UpdateStatusState} from '../../types/update.types'

export interface UpdateStatusProps {
    status: UpdateStatusState
    onInstall: () => void
}
