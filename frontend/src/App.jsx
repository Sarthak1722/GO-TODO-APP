import { useEffect, useMemo, useState } from 'react'
import { AnimatePresence, Reorder, motion, useReducedMotion } from 'framer-motion'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Activity,
  CheckCircle2,
  Circle,
  ClipboardList,
  ListFilter,
  Loader2,
  Plus,
  Search,
  Sparkles,
  TriangleAlert,
  X,
} from 'lucide-react'
import { Toaster, toast } from 'sonner'
import { TodoItem } from './components/TodoItem'
import { API_BASE_URL, todoApi } from './lib/api'

const TODOS_QUERY_KEY = ['todos']
const TODO_ORDER_STORAGE_KEY = 'kaizen.todoOrder'

const filters = [
  { id: 'all', label: 'All', icon: ListFilter },
  { id: 'active', label: 'Active', icon: Circle },
  { id: 'done', label: 'Done', icon: CheckCircle2 },
]

const panelMotion = {
  hidden: { opacity: 0, y: 16, filter: 'blur(8px)' },
  visible: {
    opacity: 1,
    y: 0,
    filter: 'blur(0px)',
    transition: { duration: 0.55, ease: [0.22, 1, 0.36, 1] },
  },
}

const getStoredOrder = () => {
  try {
    const order = JSON.parse(localStorage.getItem(TODO_ORDER_STORAGE_KEY) || '[]')
    return Array.isArray(order) ? order.map(Number).filter(Number.isFinite) : []
  } catch {
    return []
  }
}

const getOrderedTodos = (todos, orderIds) => {
  const order = new Map(orderIds.map((id, index) => [id, index]))

  return [...todos].sort((a, b) => {
    const aOrder = order.has(a.id) ? order.get(a.id) : Number.MAX_SAFE_INTEGER
    const bOrder = order.has(b.id) ? order.get(b.id) : Number.MAX_SAFE_INTEGER

    return aOrder - bOrder || b.id - a.id
  })
}

const getVisibleTodos = (todos, filter, search) => {
  const query = search.trim().toLowerCase()

  return todos.filter((todo) => {
    if (filter === 'active' && todo.completed) return false
    if (filter === 'done' && !todo.completed) return false
    if (!query) return true
    return todo.body.toLowerCase().includes(query)
  })
}

const friendlyError = (error, action) => {
  if (error?.status === 404) {
    return {
      title: 'Todo not found',
      description: `The backend could not find that todo while trying to ${action}.`,
    }
  }

  if (error?.status === 500) {
    return {
      title: 'Backend had a moment',
      description: 'Your UI was rolled back. Try again once the API settles.',
    }
  }

  return {
    title: error?.message || 'Request failed',
    description: `Could not ${action}. Please try again.`,
  }
}

const showErrorToast = (error, action) => {
  const message = friendlyError(error, action)
  toast.error(message.title, { description: message.description })
}

function App() {
  const queryClient = useQueryClient()
  const prefersReducedMotion = useReducedMotion()
  const [draft, setDraft] = useState('')
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState('all')
  const [orderIds, setOrderIds] = useState(getStoredOrder)

  const todosQuery = useQuery({
    queryKey: TODOS_QUERY_KEY,
    queryFn: todoApi.listTodos,
  })

  const todos = useMemo(
    () => getOrderedTodos(todosQuery.data ?? [], orderIds),
    [orderIds, todosQuery.data],
  )
  const stats = useMemo(() => {
    const done = todos.filter((todo) => todo.completed).length
    return {
      total: todos.length,
      active: todos.length - done,
      done,
      rate: todos.length ? Math.round((done / todos.length) * 100) : 0,
    }
  }, [todos])

  const visibleTodos = useMemo(
    () => getVisibleTodos(todos, filter, search),
    [todos, filter, search],
  )

  useEffect(() => {
    if (todosQuery.isError) {
      showErrorToast(todosQuery.error, 'load todos')
    }
  }, [todosQuery.error, todosQuery.isError])

  const createTodo = useMutation({
    mutationFn: todoApi.createTodo,
    onMutate: async (body) => {
      await queryClient.cancelQueries({ queryKey: TODOS_QUERY_KEY })
      const previousTodos = queryClient.getQueryData(TODOS_QUERY_KEY) ?? []
      const tempId = -Date.now()

      queryClient.setQueryData(TODOS_QUERY_KEY, [
        { id: tempId, body, completed: false, optimistic: true },
        ...previousTodos,
      ])

      setOrderIds((currentOrder) => [tempId, ...currentOrder])

      setDraft('')
      return { previousTodos, tempId, body }
    },
    onError: (error, _body, context) => {
      queryClient.setQueryData(TODOS_QUERY_KEY, context?.previousTodos ?? [])
      setOrderIds((currentOrder) =>
        currentOrder.filter((id) => id !== context?.tempId),
      )
      setDraft(context?.body ?? '')
      showErrorToast(error, 'create the todo')
    },
    onSuccess: (createdTodo, _body, context) => {
      queryClient.setQueryData(TODOS_QUERY_KEY, (currentTodos = []) =>
        currentTodos.map((todo) =>
          todo.id === context.tempId ? createdTodo : todo,
        ),
      )
      setOrderIds((currentOrder) =>
        currentOrder.map((id) => (id === context.tempId ? createdTodo.id : id)),
      )
      toast.success('Todo added', {
        description: 'Created instantly and confirmed by the Go API.',
      })
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: TODOS_QUERY_KEY })
    },
  })

  const toggleTodo = useMutation({
    mutationFn: ({ todo, completed, body = todo.body }) =>
      todoApi.updateTodo(todo.id, {
        id: todo.id,
        body,
        completed,
      }),
    onMutate: async ({ todo, completed, body = todo.body }) => {
      await queryClient.cancelQueries({ queryKey: TODOS_QUERY_KEY })
      const previousTodos = queryClient.getQueryData(TODOS_QUERY_KEY) ?? []

      queryClient.setQueryData(TODOS_QUERY_KEY, (currentTodos = []) =>
        currentTodos.map((item) =>
          item.id === todo.id ? { ...item, body, completed } : item,
        ),
      )

      return { previousTodos }
    },
    onError: (error, _variables, context) => {
      queryClient.setQueryData(TODOS_QUERY_KEY, context?.previousTodos ?? [])
      showErrorToast(error, 'update the todo')
    },
    onSuccess: (updatedTodo) => {
      queryClient.setQueryData(TODOS_QUERY_KEY, (currentTodos = []) =>
        currentTodos.map((todo) =>
          todo.id === updatedTodo.id ? updatedTodo : todo,
        ),
      )
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: TODOS_QUERY_KEY })
    },
  })

  const reorderTodos = (nextVisibleTodos) => {
    const visibleIds = new Set(visibleTodos.map((todo) => todo.id))
    const reorderedVisibleTodos = [...nextVisibleTodos]

    const nextTodos = todos.map((todo) =>
      visibleIds.has(todo.id) ? reorderedVisibleTodos.shift() : todo,
    )
    const nextOrder = nextTodos.map((todo) => todo.id)

    setOrderIds(nextOrder)
    localStorage.setItem(TODO_ORDER_STORAGE_KEY, JSON.stringify(nextOrder))
    queryClient.setQueryData(TODOS_QUERY_KEY, nextTodos)
  }

  const deleteTodo = useMutation({
    mutationFn: (todo) => todoApi.deleteTodo(todo.id),
    onMutate: async (todo) => {
      await queryClient.cancelQueries({ queryKey: TODOS_QUERY_KEY })
      const previousTodos = queryClient.getQueryData(TODOS_QUERY_KEY) ?? []

      queryClient.setQueryData(TODOS_QUERY_KEY, (currentTodos = []) =>
        currentTodos.filter((item) => item.id !== todo.id),
      )

      return { previousTodos }
    },
    onError: (error, _todo, context) => {
      queryClient.setQueryData(TODOS_QUERY_KEY, context?.previousTodos ?? [])
      showErrorToast(error, 'delete the todo')
    },
    onSuccess: () => {
      toast('Todo deleted', {
        description: 'Removed from the list and confirmed by the API.',
      })
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: TODOS_QUERY_KEY })
    },
  })

  const handleSubmit = (event) => {
    event.preventDefault()
    const body = draft.trim()

    if (body.length < 3) {
      toast.warning('Add a little more detail', {
        description: 'The API requires todo bodies to be at least 3 characters.',
      })
      return
    }

    if (body.length > 100) {
      toast.warning('Todo is too long', {
        description: 'The API accepts a maximum of 100 characters.',
      })
      return
    }

    createTodo.mutate(body)
  }

  const weekday = new Intl.DateTimeFormat(undefined, {
    weekday: 'long',
  }).format(new Date())

  return (
    <main className="relative min-h-screen overflow-hidden bg-zinc-950 px-4 py-6 text-zinc-100 sm:px-6 lg:px-8">
      <div className="noise-mask pointer-events-none absolute inset-0 opacity-60" />
      <div className="pointer-events-none absolute left-1/2 top-0 h-[34rem] w-[34rem] -translate-x-1/2 rounded-full bg-violet-500/10 blur-3xl" />
      <div className="pointer-events-none absolute right-0 top-20 h-80 w-80 rounded-full bg-teal-300/10 blur-3xl" />

      <Toaster
        theme="dark"
        richColors
        closeButton
        position="top-right"
        toastOptions={{
          classNames: {
            toast:
              'border border-white/10 bg-zinc-950/85 text-zinc-100 backdrop-blur-xl',
          },
        }}
      />

      <motion.section
        initial={prefersReducedMotion ? false : 'hidden'}
        animate="visible"
        variants={panelMotion}
        className="relative mx-auto flex w-full max-w-5xl flex-col gap-7"
      >
        <header className="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
          <div className="flex items-center gap-4">
            <motion.div
              aria-hidden="true"
              animate={
                prefersReducedMotion
                  ? undefined
                  : {
                      boxShadow: [
                        '0 0 20px rgba(94,234,212,0.18)',
                        '0 0 38px rgba(167,139,250,0.28)',
                        '0 0 20px rgba(94,234,212,0.18)',
                      ],
                    }
              }
              transition={{ duration: 3.2, repeat: Infinity }}
              className="grid size-12 place-items-center rounded-2xl border border-teal-200/20 bg-zinc-900/80 text-teal-200 shadow-[0_0_40px_rgba(45,212,191,0.14)] backdrop-blur-xl"
            >
              <Sparkles className="size-6" />
            </motion.div>
            <div>
              {/* <p className="text-[0.65rem] font-semibold uppercase tracking-[0.42em] text-zinc-600">
                Go · Fiber · v1.0
              </p> */}
              <h1 className="mt-1 flex flex-wrap items-baseline gap-x-3 gap-y-1 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
                Kaizen
                <span className="font-serif text-2xl font-light text-teal-200/90 sm:text-3xl">
                  改善
                </span>
              </h1>
              <p className="mt-2 max-w-xl text-sm leading-6 text-zinc-500">
                From the Japanese idea of continuous improvement: small,
                deliberate changes that compound into meaningful progress.
              </p>
            </div>
          </div>

          <div className="inline-flex max-w-full items-center gap-2 self-start rounded-full border border-white/10 bg-white/[0.04] px-4 py-2 text-sm text-zinc-400 shadow-2xl shadow-black/20 backdrop-blur-xl">
            <Activity className="size-4 text-teal-200" />
            <span className="truncate">{API_BASE_URL}</span>
            <span className="size-2 rounded-full bg-emerald-400 shadow-[0_0_18px_rgba(52,211,153,0.8)]" />
          </div>
        </header>

        <section className="max-w-3xl pt-4">
          <p className="text-xs font-semibold uppercase tracking-[0.44em] text-zinc-600">
            Inbox · {weekday}
          </p>
          <h2 className="mt-4 text-5xl font-semibold tracking-tight text-white sm:text-6xl">
            What will you{' '}
            <span className="font-serif font-extralight italic text-teal-200">finish</span>{' '}
            today?
          </h2>
          {/* <p className="mt-5 max-w-2xl text-lg leading-8 text-zinc-500">
            A small surface for big intent. Optimistic updates, fluid motion,
            and your Go API doing the heavy lifting.
          </p> */}
        </section>

        <form
          onSubmit={handleSubmit}
          className="flex flex-col gap-3 rounded-3xl border border-white/10 bg-zinc-900/45 p-2 shadow-[0_24px_80px_rgba(0,0,0,0.35)] backdrop-blur-2xl sm:flex-row sm:items-center"
        >
          <div className="flex min-h-14 flex-1 items-center gap-3 rounded-2xl bg-black/10 px-3">
            <Plus className="size-5 shrink-0 text-zinc-600" />
            <input
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              maxLength={100}
              placeholder="What needs doing?"
              className="h-12 min-w-0 flex-1 border-0 bg-transparent text-base text-white outline-none placeholder:text-zinc-600"
            />
            {draft ? (
              <button
                type="button"
                onClick={() => setDraft('')}
                className="grid size-8 place-items-center rounded-xl text-zinc-600 transition hover:bg-white/5 hover:text-zinc-300"
                aria-label="Clear todo input"
              >
                <X className="size-4" />
              </button>
            ) : null}
            <span className="hidden text-xs tabular-nums text-zinc-600 sm:block">
              {draft.trim().length}/100
            </span>
          </div>
          <button
            type="submit"
            disabled={createTodo.isPending}
            className="inline-flex h-12 items-center justify-center gap-2 rounded-2xl border border-teal-200/20 bg-teal-200 px-5 text-sm font-semibold text-zinc-950 shadow-[0_0_34px_rgba(94,234,212,0.32)] transition hover:bg-teal-100 disabled:opacity-60"
          >
            {createTodo.isPending ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <Plus className="size-4" />
            )}
            Add
          </button>
        </form>

        <section className="grid gap-3 sm:grid-cols-3">
          <StatCard label="Total" value={stats.total} />
          <StatCard label="Active" value={stats.active} tone="text-teal-200" />
          <StatCard
            label="Done"
            value={stats.done}
            suffix={`${stats.rate}%`}
            tone="text-violet-200"
          />
        </section>

        <section className="flex flex-col gap-3 md:flex-row md:items-center">
          <label className="flex h-[3.25rem] flex-1 items-center gap-3 rounded-2xl border border-white/10 bg-zinc-950/70 px-4 text-zinc-500 shadow-inner shadow-black/20 backdrop-blur-xl">
            <Search className="size-5" />
            <input
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder="Search todos..."
              className="min-w-0 flex-1 border-0 bg-transparent text-sm text-zinc-100 outline-none placeholder:text-zinc-600"
            />
          </label>

          <div className="grid grid-cols-3 gap-1 rounded-2xl border border-white/10 bg-zinc-950/70 p-1 shadow-inner shadow-black/20 backdrop-blur-xl">
            {filters.map((option) => {
              const Icon = option.icon
              const active = filter === option.id
              const count =
                option.id === 'active'
                  ? stats.active
                  : option.id === 'done'
                    ? stats.done
                    : stats.total

              return (
                <button
                  key={option.id}
                  type="button"
                  onClick={() => setFilter(option.id)}
                  className={[
                    'relative inline-flex h-11 items-center justify-center gap-2 rounded-xl px-3 text-sm font-medium transition',
                    active ? 'text-zinc-950' : 'text-zinc-400 hover:text-white',
                  ].join(' ')}
                >
                  {active ? (
                    <motion.span
                      layoutId="active-filter"
                      className="absolute inset-0 rounded-xl bg-teal-200 shadow-[0_0_28px_rgba(94,234,212,0.34)]"
                      transition={{
                        type: 'spring',
                        stiffness: 500,
                        damping: 36,
                      }}
                    />
                  ) : null}
                  <span className="relative inline-flex items-center gap-2">
                    <Icon className="size-4" />
                    {option.label}
                    <span
                      className={[
                        'rounded-md px-1.5 py-0.5 text-xs tabular-nums',
                        active ? 'bg-zinc-950/10' : 'bg-white/5 text-zinc-500',
                      ].join(' ')}
                    >
                      {count}
                    </span>
                  </span>
                </button>
              )
            })}
          </div>
        </section>

        <section className="pb-12">
          <div className="mb-3 flex items-center justify-between px-1">
            <p className="text-[0.65rem] font-semibold uppercase tracking-[0.44em] text-zinc-600">
              Earlier
            </p>
            <p className="text-xs text-zinc-600">
              {visibleTodos.length}{' '}
              {visibleTodos.length === 1 ? 'item' : 'items'}
            </p>
          </div>

          {todosQuery.isLoading ? (
            <TodoSkeletons />
          ) : todosQuery.isError ? (
            <ErrorState onRetry={() => todosQuery.refetch()} />
          ) : (
            <AnimatePresence mode="popLayout" initial={false}>
              {visibleTodos.length ? (
                <Reorder.Group
                  as="ul"
                  axis="y"
                  values={visibleTodos}
                  onReorder={reorderTodos}
                  className="space-y-3"
                >
                  {visibleTodos.map((todo) => (
                    <TodoItem
                      key={todo.id}
                      todo={todo}
                      onToggle={(item) =>
                        toggleTodo.mutate({
                          todo: item,
                          completed: !item.completed,
                        })
                      }
                      onDelete={(item) => deleteTodo.mutate(item)}
                      onEdit={(item, body) =>
                        toggleTodo.mutate({
                          todo: item,
                          body,
                          completed: item.completed,
                          action: 'edit',
                        })
                      }
                      isToggling={
                        toggleTodo.isPending &&
                        toggleTodo.variables?.todo?.id === todo.id &&
                        toggleTodo.variables?.action !== 'edit'
                      }
                      isEditing={
                        toggleTodo.isPending &&
                        toggleTodo.variables?.todo?.id === todo.id &&
                        toggleTodo.variables?.action === 'edit'
                      }
                      isDeleting={
                        deleteTodo.isPending &&
                        deleteTodo.variables?.id === todo.id
                      }
                    />
                  ))}
                </Reorder.Group>
              ) : (
                <EmptyState
                  hasTodos={todos.length > 0}
                  onClearFilters={() => {
                    setSearch('')
                    setFilter('all')
                  }}
                />
              )}
            </AnimatePresence>
          )}
        </section>
      </motion.section>
    </main>
  )
}

function StatCard({ label, value, suffix, tone = 'text-white' }) {
  return (
    <motion.div
      layout
      className="rounded-3xl border border-white/10 bg-white/[0.035] p-5 shadow-[0_24px_70px_rgba(0,0,0,0.22)] backdrop-blur-xl"
    >
      <p className="text-[0.65rem] font-semibold uppercase tracking-[0.38em] text-zinc-600">
        {label}
      </p>
      <div className="mt-3 flex items-end gap-2">
        <motion.span
          key={value}
          initial={{ y: 10, opacity: 0, filter: 'blur(6px)' }}
          animate={{ y: 0, opacity: 1, filter: 'blur(0px)' }}
          className={`text-4xl font-semibold tracking-tight ${tone}`}
        >
          {value}
        </motion.span>
        {suffix ? (
          <span className="pb-1 text-sm font-medium text-zinc-600">
            {suffix}
          </span>
        ) : null}
      </div>
    </motion.div>
  )
}

function TodoSkeletons() {
  return (
    <ul className="space-y-3">
      {[0, 1, 2].map((item) => (
        <li
          key={item}
          className="h-16 animate-pulse rounded-2xl border border-white/10 bg-white/[0.035]"
        />
      ))}
    </ul>
  )
}

function EmptyState({ hasTodos, onClearFilters }) {
  return (
    <motion.div
      key="empty-state"
      initial={{ opacity: 0, y: 18, scale: 0.98, filter: 'blur(10px)' }}
      animate={{ opacity: 1, y: 0, scale: 1, filter: 'blur(0px)' }}
      exit={{ opacity: 0, y: -12, scale: 0.97, filter: 'blur(8px)' }}
      transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }}
      className="relative overflow-hidden rounded-3xl border border-white/10 bg-zinc-950/70 px-6 py-12 text-center shadow-[0_24px_80px_rgba(0,0,0,0.28)] backdrop-blur-xl"
    >
      <motion.div
        aria-hidden="true"
        animate={{
          rotate: [0, 8, -8, 0],
          y: [0, -8, 0],
        }}
        transition={{ duration: 4, repeat: Infinity, ease: 'easeInOut' }}
        className="mx-auto grid size-16 place-items-center rounded-3xl border border-teal-200/20 bg-teal-200/10 text-teal-200 shadow-[0_0_50px_rgba(94,234,212,0.18)]"
      >
        <ClipboardList className="size-8" />
      </motion.div>
      <h3 className="mt-6 text-xl font-semibold text-white">
        {hasTodos ? 'No todos match that view' : 'Nothing here yet'}
      </h3>
      <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-zinc-500">
        {hasTodos
          ? 'Clear the search or switch filters to bring your tasks back.'
          : 'Add one tiny next action and the list will glide into place.'}
      </p>
      {hasTodos ? (
        <button
          type="button"
          onClick={onClearFilters}
          className="mt-6 rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-2 text-sm font-medium text-zinc-200 transition hover:bg-white/[0.08]"
        >
          Clear filters
        </button>
      ) : null}
    </motion.div>
  )
}

function ErrorState({ onRetry }) {
  return (
    <div className="rounded-3xl border border-rose-400/20 bg-rose-400/10 p-6 text-center text-rose-100">
      <TriangleAlert className="mx-auto size-8" />
      <h3 className="mt-3 text-lg font-semibold">Could not load todos</h3>
      <p className="mt-2 text-sm text-rose-100/70">
        Check that the Go API is reachable at {API_BASE_URL}.
      </p>
      <button
        type="button"
        onClick={onRetry}
        className="mt-5 rounded-2xl bg-rose-100 px-4 py-2 text-sm font-semibold text-rose-950"
      >
        Retry
      </button>
    </div>
  )
}

export default App
