import { useState, useEffect, useCallback } from 'react'
import Navbar from '../components/Navbar'
import ExpenseList from '../components/ExpenseList'
import ExpenseModal from '../components/ExpenseModal'
import { listExpenses, createExpense, updateExpense, deleteExpense } from '../api/expenses'
import type { Expense, CreateExpenseInput, ExpenseFilter, Category } from '../types'
import { CATEGORIES } from '../types'

export default function DashboardPage() {
  const [expenses, setExpenses] = useState<Expense[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const pageSize = 20

  const [filter, setFilter] = useState<ExpenseFilter>({})
  const [filterCategory, setFilterCategory] = useState<Category | ''>('')
  const [filterStart, setFilterStart] = useState('')
  const [filterEnd, setFilterEnd] = useState('')

  const [modalOpen, setModalOpen] = useState(false)
  const [editingExpense, setEditingExpense] = useState<Expense | undefined>()
  const [modalLoading, setModalLoading] = useState(false)

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const fetchExpenses = useCallback(async (f: ExpenseFilter, p: number) => {
    setLoading(true)
    setError('')
    try {
      const res = await listExpenses({ ...f, page: p, page_size: pageSize })
      setExpenses(res.data.data.expenses ?? [])
      setTotal(res.data.data.total)
    } catch {
      setError('Failed to load expenses.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchExpenses(filter, page)
  }, [filter, page, fetchExpenses])

  function applyFilters() {
    const newFilter: ExpenseFilter = {}
    if (filterCategory) newFilter.category = filterCategory
    if (filterStart) newFilter.start_date = filterStart
    if (filterEnd) newFilter.end_date = filterEnd
    setPage(1)
    setFilter(newFilter)
  }

  function clearFilters() {
    setFilterCategory('')
    setFilterStart('')
    setFilterEnd('')
    setPage(1)
    setFilter({})
  }

  function openCreate() {
    setEditingExpense(undefined)
    setModalOpen(true)
  }

  function openEdit(expense: Expense) {
    setEditingExpense(expense)
    setModalOpen(true)
  }

  async function handleSubmit(input: CreateExpenseInput) {
    setModalLoading(true)
    try {
      if (editingExpense) {
        await updateExpense(editingExpense.id, input)
      } else {
        await createExpense(input)
      }
      setModalOpen(false)
      fetchExpenses(filter, page)
    } finally {
      setModalLoading(false)
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this expense?')) return
    try {
      await deleteExpense(id)
      fetchExpenses(filter, page)
    } catch {
      alert('Failed to delete expense.')
    }
  }

  const totalPages = Math.ceil(total / pageSize)

  const thisMonthTotal = expenses
    .filter(e => {
      const d = new Date(e.date)
      const now = new Date()
      return d.getUTCMonth() === now.getMonth() && d.getUTCFullYear() === now.getFullYear()
    })
    .reduce((sum, e) => sum + e.amount, 0)

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-5xl mx-auto px-4 sm:px-6 py-6 space-y-6">
        {/* Stats */}
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <StatCard label="Total expenses" value={total.toString()} />
          <StatCard label="Shown on page" value={expenses.length.toString()} />
          <StatCard label="This month" value={`$${thisMonthTotal.toFixed(2)}`} highlight />
        </div>

        {/* Filters */}
        <div className="bg-white rounded-2xl border border-gray-200 p-4">
          <div className="flex flex-wrap gap-3 items-end">
            <div>
              <label className="block text-xs text-gray-500 mb-1">Category</label>
              <select
                value={filterCategory}
                onChange={e => setFilterCategory(e.target.value as Category | '')}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">All categories</option>
                {CATEGORIES.map(c => <option key={c} value={c}>{c}</option>)}
              </select>
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">From</label>
              <input
                type="date"
                value={filterStart}
                onChange={e => setFilterStart(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">To</label>
              <input
                type="date"
                value={filterEnd}
                onChange={e => setFilterEnd(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
            <button
              onClick={applyFilters}
              className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors"
            >
              Apply
            </button>
            {(filterCategory || filterStart || filterEnd) && (
              <button
                onClick={clearFilters}
                className="px-4 py-2 text-sm text-gray-500 hover:text-gray-700 transition-colors"
              >
                Clear
              </button>
            )}
            <div className="ml-auto">
              <button
                onClick={openCreate}
                className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors"
              >
                + Add Expense
              </button>
            </div>
          </div>
        </div>

        {/* Expense list */}
        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden">
          {error && (
            <div className="px-4 py-3 bg-red-50 border-b border-red-200 text-red-700 text-sm">
              {error}
            </div>
          )}
          {loading ? (
            <div className="text-center py-16 text-gray-400 text-sm">Loading…</div>
          ) : (
            <ExpenseList expenses={expenses} onEdit={openEdit} onDelete={handleDelete} />
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-2">
            <button
              disabled={page === 1}
              onClick={() => setPage(p => p - 1)}
              className="px-3 py-1.5 text-sm border border-gray-300 rounded-lg disabled:opacity-40 hover:bg-gray-50 transition-colors"
            >
              Previous
            </button>
            <span className="text-sm text-gray-500">Page {page} of {totalPages}</span>
            <button
              disabled={page === totalPages}
              onClick={() => setPage(p => p + 1)}
              className="px-3 py-1.5 text-sm border border-gray-300 rounded-lg disabled:opacity-40 hover:bg-gray-50 transition-colors"
            >
              Next
            </button>
          </div>
        )}
      </main>

      <ExpenseModal
        isOpen={modalOpen}
        title={editingExpense ? 'Edit Expense' : 'New Expense'}
        expense={editingExpense}
        loading={modalLoading}
        onSubmit={handleSubmit}
        onClose={() => setModalOpen(false)}
      />
    </div>
  )
}

function StatCard({ label, value, highlight }: { label: string; value: string; highlight?: boolean }) {
  return (
    <div className={`rounded-2xl border p-4 ${highlight ? 'bg-indigo-600 border-indigo-600 text-white' : 'bg-white border-gray-200'}`}>
      <p className={`text-xs font-medium mb-1 ${highlight ? 'text-indigo-200' : 'text-gray-500'}`}>{label}</p>
      <p className={`text-2xl font-bold ${highlight ? 'text-white' : 'text-gray-900'}`}>{value}</p>
    </div>
  )
}
