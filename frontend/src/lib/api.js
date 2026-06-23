const DEFAULT_API_BASE_URL = "/api";

export const API_BASE_URL = (
  `${import.meta.env.VITE_API_URL || ""}${DEFAULT_API_BASE_URL}`
).replace(/\/$/, "");

export class ApiError extends Error {
  constructor(message, { status, details, payload } = {}) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.details = details;
    this.payload = payload;
  }
}

const normalizeTodo = (todo) => ({
  id: Number(todo.id),
  body: String(todo.body ?? ""),
  completed: Boolean(todo.completed),
});

const parsePayload = async (response) => {
  if (response.status === 204) return null;

  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    return response.json();
  }

  const text = await response.text();
  return text ? { error: text } : null;
};

const request = async (path, options = {}, getToken) => {
  const token = typeof getToken === "function" ? await getToken() : getToken;

  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
    ...options,
  });

  const payload = await parsePayload(response);

  if (!response.ok) {
    throw new ApiError(
      payload?.error || `Request failed with status ${response.status}`,
      {
        status: response.status,
        details: payload?.details,
        payload,
      },
    );
  }

  return payload?.success && Object.hasOwn(payload, "data")
    ? payload.data
    : payload;
};

export const todoApi = {
  async listTodos(getToken) {
    const todos = await request("/todos", {}, getToken);
    return Array.isArray(todos) ? todos.map(normalizeTodo) : [];
  },

  async getTodo(id, getToken) {
    const todo = await request(`/todos/${id}`, {}, getToken);
    return normalizeTodo(todo);
  },

  async createTodo(body, getToken) {
    const todo = await request(
      "/todos",
      {
        method: "POST",
        body: JSON.stringify({ body }),
      },
      getToken,
    );
    return normalizeTodo(todo);
  },

  async updateTodo(id, todo, getToken) {
    const updatedTodo = await request(
      `/todos/${id}`,
      {
        method: "PATCH",
        body: JSON.stringify(todo),
      },
      getToken,
    );
    return normalizeTodo(updatedTodo);
  },

  async deleteTodo(id, getToken) {
    await request(
      `/todos/${id}`,
      {
        method: "DELETE",
      },
      getToken,
    );
    return id;
  },
};
