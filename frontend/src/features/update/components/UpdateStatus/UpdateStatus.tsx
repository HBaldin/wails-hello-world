import {Button} from '@/components/ui/Button'
import type {UpdateStatusProps} from './UpdateStatus.types'
import styles from './UpdateStatus.module.css'

export function UpdateStatus({status, onInstall}: UpdateStatusProps) {
    switch (status.kind) {
        case 'idle':
            return null
        case 'checking':
            return <p className={styles.message}>Checking for updates...</p>
        case 'up-to-date':
            return <p className={styles.message}>You're up to date.</p>
        case 'available':
            return (
                <p className={styles.message}>
                    Update v{status.info.version} available.
                    <Button variant="secondary" onClick={onInstall}>
                        Install and restart
                    </Button>
                </p>
            )
        case 'installing':
            return <p className={styles.message}>Downloading and installing update...</p>
        case 'error':
            return <p className={styles.error}>Update failed: {status.message}</p>
    }
}
