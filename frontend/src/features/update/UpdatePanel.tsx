import {useEffect, useState} from 'react'
import {Button} from '@/components/ui/Button'
import {UpdateStatus} from './components/UpdateStatus'
import {getAppVersion} from './api/updateApi'
import {useUpdater} from './hooks/useUpdater'
import styles from './UpdatePanel.module.css'

export function UpdatePanel() {
    const [version, setVersion] = useState('')
    const {status, check, install} = useUpdater()

    useEffect(() => {
        getAppVersion().then(setVersion)
    }, [])

    return (
        <div className={styles.panel}>
            <span className={styles.version}>v{version}</span>
            <Button variant="secondary" onClick={check} disabled={status.kind === 'checking'}>
                Check for updates
            </Button>
            <UpdateStatus status={status} onInstall={install} />
        </div>
    )
}
