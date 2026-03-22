export interface User {
  id: number
  name: string
  email: string
}

export interface AuthResponse {
  token: string
  user: User
}

export type Category =
  | 'Groceries'
  | 'Leisure'
  | 'Electronics'
  | 'Utilities'
  | 'Clothing'
  | 'Health'
  | 'Others'

export const CATEGORIES: Category[] = [
  'Groceries',
  'Leisure',
  'Electronics',
  'Utilities',
  'Clothing',
  'Health',
  'Others',
]

export interface Expense {
  id: number
  user_id: number
  title: string
  amount: number
  category: Category
  date: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface CreateExpenseInput {
  title: string
  amount: number
  category: Category
  date: string
  notes?: string
}

export interface UpdateExpenseInput {
  title?: string
  amount?: number
  category?: Category
  date?: string
  notes?: string
}

export interface ExpenseListResponse {
  expenses: Expense[]
  total: number
  page: number
  page_size: number
}

export interface ExpenseFilter {
  category?: Category
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

export interface ApiError {
  code: string
  message: string
}
