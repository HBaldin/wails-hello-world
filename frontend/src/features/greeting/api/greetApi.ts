import {Greet} from '@wailsjs/go/main/App'
import {callBackend} from '@/services/wailsClient'

export function greet(name: string): Promise<string> {
    return callBackend('Greet', () => Greet(name))
}
