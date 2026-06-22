import { useEffect, useRef, useState } from 'react'
import { Reorder, motion } from 'framer-motion'
import { Check, GripVertical, Loader2, Pencil, Save, Trash2, X } from 'lucide-react'
import { toast } from 'sonner'

const spring = {
  type: 'spring',
  stiffness: 520,
  damping: 38,
  mass: 0.72,
}

export function TodoItem({
  todo,
  onToggle,
  onDelete,
  onEdit,
  isToggling = false,
  isEditing = false,
  isDeleting = false,
}) {
  const [isEditingBody, setIsEditingBody] = useState(false)
  const [draft, setDraft] = useState(todo.body)
  const inputRef = useRef(null)
  const isBusy = isToggling || isEditing || isDeleting || todo.optimistic
  const completed = todo.completed

  useEffect(() => {
    if (isEditingBody) {
      inputRef.current?.focus()
      inputRef.current?.select()
    }
  }, [isEditingBody])

  const cancelEdit = () => {
    setDraft(todo.body)
    setIsEditingBody(false)
  }

  const submitEdit = () => {
    const body = draft.trim()

    if (body === todo.body) {
      setIsEditingBody(false)
      return
    }

    if (body.length < 3) {
      toast.warning('Add a little more detail', {
        description: 'Todo bodies need at least 3 characters.',
      })
      return
    }

    if (body.length > 100) {
      toast.warning('Todo is too long', {
        description: 'The API accepts a maximum of 100 characters.',
      })
      return
    }

    onEdit(todo, body)
    setIsEditingBody(false)
  }

  const handleEditKeyDown = (event) => {
    if (event.key === 'Enter') {
      event.preventDefault()
      submitEdit()
    }

    if (event.key === 'Escape') {
      cancelEdit()
    }
  }

  return (
    <Reorder.Item
      value={todo}
      layout
      initial={{ opacity: 0, y: 18, scale: 0.98, filter: 'blur(10px)' }}
      animate={{ opacity: 1, y: 0, scale: 1, filter: 'blur(0px)' }}
      exit={{
        opacity: 0,
        height: 0,
        y: -10,
        scale: 0.96,
        filter: 'blur(8px)',
        marginTop: 0,
        marginBottom: 0,
      }}
      transition={spring}
      className="group relative overflow-hidden rounded-2xl border border-white/10 bg-zinc-950/70 shadow-[0_18px_70px_rgba(0,0,0,0.22)] backdrop-blur-xl"
    >
      <motion.div
        layout
        className="absolute inset-0 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
        style={{
          background:
            'linear-gradient(90deg, rgba(167,139,250,0.11), rgba(45,212,191,0.07), transparent)',
        }}
      />

      <div className="relative flex min-h-16 items-center gap-3 px-4 py-3 sm:px-5">
        <button
          type="button"
          className="hidden size-8 shrink-0 cursor-grab place-items-center rounded-xl text-zinc-700 transition active:cursor-grabbing group-hover:text-zinc-500 hover:bg-white/5 sm:grid"
          aria-label="Drag to reorder todo"
        >
          <GripVertical aria-hidden="true" className="size-4" />
        </button>

        <motion.button
          type="button"
          whileTap={{ scale: 0.9 }}
          whileHover={{ scale: isBusy ? 1 : 1.06 }}
          disabled={isBusy}
          onClick={() => onToggle(todo)}
          className={[
            'relative grid size-8 shrink-0 place-items-center rounded-xl border transition-all duration-300',
            completed
              ? 'border-teal-200/70 bg-teal-300 text-zinc-950 shadow-[0_0_28px_rgba(94,234,212,0.45)]'
              : 'border-zinc-700/90 bg-zinc-900/80 text-transparent hover:border-violet-300/70 hover:shadow-[0_0_28px_rgba(167,139,250,0.25)]',
          ].join(' ')}
          aria-label={completed ? 'Mark todo active' : 'Mark todo completed'}
        >
          {isToggling ? (
            <Loader2 className="size-4 animate-spin text-zinc-950" />
          ) : (
            <motion.span
              initial={false}
              animate={{ scale: completed ? 1 : 0.35, opacity: completed ? 1 : 0 }}
              transition={spring}
            >
              <Check className="size-4" strokeWidth={3} />
            </motion.span>
          )}
        </motion.button>

        {isEditingBody ? (
          <div className="flex min-w-0 flex-1 items-center gap-2">
            <input
              ref={inputRef}
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              onKeyDown={handleEditKeyDown}
              maxLength={100}
              disabled={isEditing}
              className="h-10 min-w-0 flex-1 rounded-xl border border-teal-200/20 bg-black/20 px-3 text-sm font-medium text-white outline-none transition placeholder:text-zinc-600 focus:border-teal-200/60 focus:shadow-[0_0_24px_rgba(94,234,212,0.18)] sm:text-base"
              aria-label="Edit todo"
            />
            <span className="hidden text-xs tabular-nums text-zinc-600 sm:block">
              {draft.trim().length}/100
            </span>
          </div>
        ) : (
          <motion.p
            layout
            animate={{
              color: completed ? 'rgb(113 113 122)' : 'rgb(244 244 245)',
            }}
            transition={{ duration: 0.22 }}
            className={[
              'min-w-0 flex-1 text-left text-sm font-medium leading-6 sm:text-base',
              completed ? 'line-through decoration-zinc-500 decoration-2' : '',
            ].join(' ')}
          >
            {todo.body}
          </motion.p>
        )}

        {todo.optimistic ? (
          <span className="rounded-full border border-violet-300/20 bg-violet-300/10 px-2.5 py-1 text-[0.65rem] font-semibold uppercase tracking-[0.24em] text-violet-200">
            Syncing
          </span>
        ) : null}

        {isEditingBody ? (
          <>
            <motion.button
              type="button"
              whileTap={{ scale: 0.92 }}
              disabled={isEditing}
              onClick={submitEdit}
              className="grid size-9 shrink-0 place-items-center rounded-xl border border-teal-200/20 bg-teal-200/10 text-teal-100 transition hover:bg-teal-200 hover:text-zinc-950 disabled:opacity-40"
              aria-label="Save todo edit"
            >
              {isEditing ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <Save className="size-4" />
              )}
            </motion.button>
            <motion.button
              type="button"
              whileTap={{ scale: 0.92 }}
              disabled={isEditing}
              onClick={cancelEdit}
              className="grid size-9 shrink-0 place-items-center rounded-xl border border-transparent text-zinc-600 transition hover:border-white/10 hover:bg-white/5 hover:text-zinc-200 disabled:opacity-40"
              aria-label="Cancel todo edit"
            >
              <X className="size-4" />
            </motion.button>
          </>
        ) : (
          <motion.button
            type="button"
            whileTap={{ scale: 0.92 }}
            disabled={isBusy}
            onClick={() => {
              setDraft(todo.body)
              setIsEditingBody(true)
            }}
            className="grid size-9 shrink-0 place-items-center rounded-xl border border-transparent text-zinc-600 opacity-100 transition hover:border-teal-200/20 hover:bg-teal-200/10 hover:text-teal-100 disabled:opacity-40 sm:opacity-0 sm:group-hover:opacity-100"
            aria-label="Edit todo"
          >
            <Pencil className="size-4" />
          </motion.button>
        )}

        <motion.button
          type="button"
          whileTap={{ scale: 0.92 }}
          disabled={isDeleting || todo.optimistic}
          onClick={() => onDelete(todo)}
          className="grid size-9 shrink-0 place-items-center rounded-xl border border-transparent text-zinc-600 opacity-100 transition hover:border-rose-400/25 hover:bg-rose-400/10 hover:text-rose-200 disabled:opacity-40 sm:opacity-0 sm:group-hover:opacity-100"
          aria-label="Delete todo"
        >
          {isDeleting ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Trash2 className="size-4" />
          )}
        </motion.button>
      </div>
    </Reorder.Item>
  )
}
