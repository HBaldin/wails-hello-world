import {useState} from 'react'
import {GreetForm} from './components/GreetForm'
import {greet} from './api/greetApi'
import styles from './GreetingPage.module.css'

export function GreetingPage() {
    const [name, setName] = useState('')
    const [result, setResult] = useState('Please enter your name below 👇')

    async function handleSubmit() {
        if (!name) return
        try {
            setResult(await greet(name))
        } catch (err) {
            setResult(err instanceof Error ? err.message : String(err))
        }
    }

    return (
        <div className={styles.page}>
            <div className={styles.result}>{result}</div>
            <GreetForm name={name} onNameChange={setName} onSubmit={handleSubmit} />
        </div>
    )
}
