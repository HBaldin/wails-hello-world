import type {ButtonProps} from './Button.types'
import styles from './Button.module.css'

export function Button({children, variant = 'primary', className, ...rest}: ButtonProps) {
    const classes = [styles.btn, styles[variant], className].filter(Boolean).join(' ')

    return (
        <button className={classes} {...rest}>
            {children}
        </button>
    )
}
