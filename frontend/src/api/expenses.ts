import { api } from './client'
import type { Expense, CreateExpenseInput, UpdateExpenseInput, ExpenseListResponse, ExpenseFilter } from '../types'

export const listExpenses = (filter?: ExpenseFilter) =>
  api.get<{ data: ExpenseListResponse }>('/expenses', { params: filter })

export const getExpense = (id: number) =>
  api.get<{ data: Expense }>(`/expenses/${id}`)

export const createExpense = (input: CreateExpenseInput) =>
  api.post<{ data: Expense }>('/expenses', input)

export const updateExpense = (id: number, input: UpdateExpenseInput) =>
  api.put<{ data: Expense }>(`/expenses/${id}`, input)

export const deleteExpense = (id: number) =>
  api.delete(`/expenses/${id}`)
