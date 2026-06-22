import { motion } from 'framer-motion'
import { Check, GripVertical, Loader2, Trash2 } from 'lucide-react'

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
  isToggling = false,
  isDeleting = false,
}) {
  const isBusy = isToggling || isDeleting || todo.optimistic
  const completed = todo.completed

  return (
    <motion.li
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
        <GripVertical
          aria-hidden="true"
          className="hidden size-4 shrink-0 text-zinc-700 transition-colors group-hover:text-zinc-500 sm:block"
        />

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

        {todo.optimistic ? (
          <span className="rounded-full border border-violet-300/20 bg-violet-300/10 px-2.5 py-1 text-[0.65rem] font-semibold uppercase tracking-[0.24em] text-violet-200">
            Syncing
          </span>
        ) : null}

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
    </motion.li>
  )
}
