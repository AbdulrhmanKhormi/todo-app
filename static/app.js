const apiUrl = "/todos";

async function loadTodos() {
  const res = await fetch(apiUrl);
  const todos = await res.json();
  const list = document.getElementById("todo-list");
  list.innerHTML = "";

  todos.forEach((todo) => {
    const li = document.createElement("li");
    li.className = todo.completed ? "done" : "";
    li.textContent = todo.title;

    li.onclick = async () => {
      await fetch(`${apiUrl}/${todo.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...todo, completed: !todo.completed }),
      });
      loadTodos();
    };

    const del = document.createElement("button");
    del.textContent = "X";
    del.onclick = async (e) => {
      e.stopPropagation();
      await fetch(`${apiUrl}/${todo.id}`, { method: "DELETE" });
      loadTodos();
    };

    li.appendChild(del);
    list.appendChild(li);
  });
}

document.getElementById("todo-form").onsubmit = async (e) => {
  e.preventDefault();
  const title = document.getElementById("title").value;
  await fetch(apiUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, completed: false }),
  });
  document.getElementById("title").value = "";
  loadTodos();
};

loadTodos();
