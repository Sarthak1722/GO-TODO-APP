const DEFAULT_API_BASE_URL = '/api'

export const API_BASE_URL = (
  import.meta.env.VITE_API_URL || DEFAULT_API_BASE_URL
).replace(/\/$/, '')

export class ApiError extends Error {
  constructor(message, { status, details, payload } = {}) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.details = details
    this.payload = payload
  }
}

const normalizeTodo = (todo) => ({
  id: Number(todo.id),
  body: String(todo.body ?? ''),
  completed: Boolean(todo.completed),
})

const parsePayload = async (response) => {
  if (response.status === 204) return null

  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return response.json()
  }

  const text = await response.text()
  return text ? { error: text } : null
}

const request = async (path, options = {}) => {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  })

  const payload = await parsePayload(response)

  if (!response.ok) {
    throw new ApiError(
      payload?.error || `Request failed with status ${response.status}`,
      {
        status: response.status,
        details: payload?.details,
        payload,
      },
    )
  }

  return payload?.success && Object.hasOwn(payload, 'data') ? payload.data : payload
}

export const todoApi = {
  async listTodos() {
    const todos = await request('/todos')
    return Array.isArray(todos) ? todos.map(normalizeTodo) : []
  },

  async getTodo(id) {
    const todo = await request(`/todos/${id}`)
    return normalizeTodo(todo)
  },

  async createTodo(body) {
    const todo = await request('/todos', {
      method: 'POST',
      body: JSON.stringify({ body }),
    })
    return normalizeTodo(todo)
  },

  async updateTodo(id, todo) {
    const updatedTodo = await request(`/todos/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(todo),
    })
    return normalizeTodo(updatedTodo)
  },

  async deleteTodo(id) {
    await request(`/todos/${id}`, {
      method: 'DELETE',
      headers: {
        Accept: 'application/json',
      },
    })
    return id
  },
}
