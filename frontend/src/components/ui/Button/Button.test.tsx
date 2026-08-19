import {describe, expect, it, vi} from 'vitest'
import {fireEvent, render, screen} from '@testing-library/react'
import {Button} from './Button'

describe('Button', () => {
    it('renders its children', () => {
        render(<Button>Click me</Button>)
        expect(screen.getByText('Click me')).toBeInTheDocument()
    })

    it('calls onClick when clicked', () => {
        const onClick = vi.fn()
        render(<Button onClick={onClick}>Click me</Button>)
        fireEvent.click(screen.getByText('Click me'))
        expect(onClick).toHaveBeenCalledTimes(1)
    })

    it('does not call onClick when disabled', () => {
        const onClick = vi.fn()
        render(
            <Button onClick={onClick} disabled>
                Click me
            </Button>,
        )
        fireEvent.click(screen.getByText('Click me'))
        expect(onClick).not.toHaveBeenCalled()
    })
})
