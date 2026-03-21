import { useState, useEffect, type FormEvent } from 'react'
import { CATEGORIES, type Category, type CreateExpenseInput, type Expense } from '../types'

interface Props {
  initial?: Expense
  onSubmit: (input: CreateExpenseInput) => Promise<void>
  onCancel: () => void
  loading: boolean
}

const CATEGORY_COLORS: Record<Category, string> = {
  Groceries: 'bg-green-100 text-green-800',
  Leisure: 'bg-purple-100 text-purple-800',
  Electronics: 'bg-blue-100 text-blue-800',
  Utilities: 'bg-yellow-100 text-yellow-800',
  Clothing: 'bg-pink-100 text-pink-800',
  Health: 'bg-red-100 text-red-800',
  Others: 'bg-gray-100 text-gray-800',
}

export { CATEGORY_COLORS }

export default function ExpenseForm({ initial, onSubmit, onCancel, loading }: Props) {
  const [title, setTitle] = useState(initial?.title ?? '')
  const [amount, setAmount] = useState(initial?.amount?.toString() ?? '')
  const [category, setCategory] = useState<Category>(initial?.category ?? 'Others')
  const [date, setDate] = useState(initial?.date?.slice(0, 10) ?? new Date().toISOString().slice(0, 10))
  const [notes, setNotes] = useState(initial?.notes ?? '')
  const [error, setError] = useState('')

  useEffect(() => {
    if (initial) {
      setTitle(initial.title)
      setAmount(initial.amount.toString())
      setCategory(initial.category)
      setDate(initial.date.slice(0, 10))
      setNotes(initial.notes ?? '')
    }
  }, [initial])

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    const parsedAmount = parseFloat(amount)
    if (isNaN(parsedAmount) || parsedAmount <= 0) {
      setError('Amount must be a positive number.')
      return
    }
    try {
      await onSubmit({ title, amount: parsedAmount, category, date, notes: notes || undefined })
    } catch (err: unknown) {
      if (err && typeof err === 'object' && 'response' in err) {
        const res = (err as { response?: { data?: { error?: { message?: string } } } }).response
        setError(res?.data?.error?.message ?? 'Failed to save expense.')
      } else {
        setError('Failed to save expense.')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="px-4 py-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
          {error}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
        <input
          type="text"
          required
          value={title}
          onChange={e => setTitle(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          placeholder="e.g. Weekly groceries"
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Amount ($)</label>
          <input
            type="number"
            required
            min="0.01"
            step="0.01"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="0.00"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input
            type="date"
            required
            value={date}
            onChange={e => setDate(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
        <select
          value={category}
          onChange={e => setCategory(e.target.value as Category)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white"
        >
          {CATEGORIES.map(c => (
            <option key={c} value={c}>{c}</option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Notes <span className="text-gray-400">(optional)</span></label>
        <textarea
          value={notes}
          onChange={e => setNotes(e.target.value)}
          rows={2}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
          placeholder="Any additional notes…"
        />
      </div>

      <div className="flex justify-end gap-3 pt-2">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm text-gray-600 hover:text-gray-900 transition-colors"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={loading}
          className="px-5 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? 'Saving…' : initial ? 'Update' : 'Add Expense'}
        </button>
      </div>
    </form>
  )
}
