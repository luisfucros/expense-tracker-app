import type { Expense } from '../types'
import { CATEGORY_COLORS } from './ExpenseForm'

interface Props {
  expenses: Expense[]
  onEdit: (expense: Expense) => void
  onDelete: (id: number) => void
}

export default function ExpenseList({ expenses, onEdit, onDelete }: Props) {
  if (expenses.length === 0) {
    return (
      <div className="text-center py-16 text-gray-400">
        <p className="text-4xl mb-3">🧾</p>
        <p className="text-sm">No expenses yet. Add your first one!</p>
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-100">
            <th className="text-left py-3 px-4 text-xs font-semibold text-gray-500 uppercase tracking-wide">Date</th>
            <th className="text-left py-3 px-4 text-xs font-semibold text-gray-500 uppercase tracking-wide">Title</th>
            <th className="text-left py-3 px-4 text-xs font-semibold text-gray-500 uppercase tracking-wide">Category</th>
            <th className="text-right py-3 px-4 text-xs font-semibold text-gray-500 uppercase tracking-wide">Amount</th>
            <th className="text-left py-3 px-4 text-xs font-semibold text-gray-500 uppercase tracking-wide hidden md:table-cell">Notes</th>
            <th className="py-3 px-4" />
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-50">
          {expenses.map(expense => (
            <tr key={expense.id} className="hover:bg-gray-50 transition-colors group">
              <td className="py-3 px-4 text-gray-500 whitespace-nowrap">
                {formatDate(expense.date)}
              </td>
              <td className="py-3 px-4 font-medium text-gray-900">{expense.title}</td>
              <td className="py-3 px-4">
                <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${CATEGORY_COLORS[expense.category]}`}>
                  {expense.category}
                </span>
              </td>
              <td className="py-3 px-4 text-right font-semibold text-gray-900 tabular-nums">
                ${expense.amount.toFixed(2)}
              </td>
              <td className="py-3 px-4 text-gray-400 truncate max-w-[180px] hidden md:table-cell">
                {expense.notes || '—'}
              </td>
              <td className="py-3 px-4">
                <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity justify-end">
                  <button
                    onClick={() => onEdit(expense)}
                    className="text-indigo-500 hover:text-indigo-700 text-xs font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(expense.id)}
                    className="text-red-400 hover:text-red-600 text-xs font-medium"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', timeZone: 'UTC' })
}
