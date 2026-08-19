import {Button} from '@/components/ui/Button'
import type {GreetFormProps} from './GreetForm.types'
import styles from './GreetForm.module.css'

export function GreetForm({name, onNameChange, onSubmit}: GreetFormProps) {
    return (
        <div className={styles.inputBox}>
            <input
                className={styles.input}
                autoComplete="off"
                type="text"
                value={name}
                onChange={(e) => onNameChange(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && onSubmit()}
            />
            <Button onClick={onSubmit}>Greet</Button>
        </div>
    )
}
