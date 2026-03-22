import type { Expense, CreateExpenseInput } from '../types'
import ExpenseForm from './ExpenseForm'

interface Props {
  isOpen: boolean
  title: string
  expense?: Expense
  loading: boolean
  onSubmit: (input: CreateExpenseInput) => Promise<void>
  onClose: () => void
}

export default function ExpenseModal({ isOpen, title, expense, loading, onSubmit, onClose }: Props) {
  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center px-4"
      onClick={e => { if (e.target === e.currentTarget) onClose() }}
    >
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" />
      <div className="relative w-full max-w-md bg-white rounded-2xl shadow-xl p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">{title}</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors text-xl leading-none"
          >
            ×
          </button>
        </div>
        <ExpenseForm
          initial={expense}
          onSubmit={onSubmit}
          onCancel={onClose}
          loading={loading}
        />
      </div>
    </div>
  )
}
