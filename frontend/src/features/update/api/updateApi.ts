import {AppVersion, CheckForUpdate, InstallUpdate} from '@wailsjs/go/main/App'
import {callBackend} from '@/services/wailsClient'
import type {UpdateInfo} from '../types/update.types'

export function getAppVersion(): Promise<string> {
    return callBackend('AppVersion', () => AppVersion())
}

export function checkForUpdate(): Promise<UpdateInfo> {
    return callBackend('CheckForUpdate', () => CheckForUpdate())
}

export function installUpdate(): Promise<void> {
    return callBackend('InstallUpdate', () => InstallUpdate())
}
