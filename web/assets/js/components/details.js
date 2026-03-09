import { apiHTML, swapHTML } from "../utils/fetch.js";

class DetailsElement extends HTMLElement {
  quill;

  get description() {
    return this.quill.getSemanticHTML();
  }

  async connectedCallback() {
    await this.#loadQuill();

    const taskId = this.getAttribute("task-id");
    const taskSelector = `.tasks__list > li[data-id='${taskId}']`;

    this.quill = new Quill("#note-editor", {
      modules: {
        toolbar: [
          [{ header: [1, 2, 3, false] }],
          ["bold", "italic", "underline", "strike"],
          ["blockquote", "code-block"],
          ["link", "image"],
          [{ list: "ordered" }, { list: "bullet" }],
        ],
      },
      bounds: this,
      placeholder: "Add description...",
      theme: "snow",
    });

    // Title form
    const titleForm = this.querySelector(".detail__title form");
    titleForm?.addEventListener("submit", async (e) => {
      e.preventDefault();
      const html = await apiHTML("PUT", `/partials/tasks/${taskId}/title`, new URLSearchParams(new FormData(titleForm)));
      swapHTML(taskSelector, html, "outerHTML");
    });

    // Priority select
    const prioritySelect = this.querySelector("select[name='priority']");
    prioritySelect?.addEventListener("change", async () => {
      const body = new URLSearchParams({ priority: prioritySelect.value });
      const html = await apiHTML("PUT", `/partials/tasks/${taskId}/priority`, body);
      swapHTML(taskSelector, html, "outerHTML");
    });

    // List select
    const listSelect = this.querySelector("select[name='list']");
    listSelect?.addEventListener("change", async () => {
      const body = new URLSearchParams({ list: listSelect.value });
      const html = await apiHTML("PUT", `/partials/tasks/${taskId}/list`, body);
      swapHTML(taskSelector, html, "outerHTML");
    });

    // Start date
    const startDate = this.querySelector("input[name='from']");
    startDate?.addEventListener("change", async () => {
      const body = new URLSearchParams({ from: startDate.value });
      const html = await apiHTML("PUT", `/partials/tasks/${taskId}/date-start`, body);
      swapHTML(taskSelector, html, "outerHTML");
    });

    // End date
    const endDate = this.querySelector("input[name='to']");
    endDate?.addEventListener("change", async () => {
      const body = new URLSearchParams({ to: endDate.value });
      const html = await apiHTML("PUT", `/partials/tasks/${taskId}/date-end`, body);
      swapHTML(taskSelector, html, "outerHTML");
    });

    // Subtask form
    const subtaskForm = this.querySelector("form.add-subtask");
    subtaskForm?.addEventListener("submit", async (e) => {
      e.preventDefault();
      const html = await apiHTML("POST", `/partials/tasks/${taskId}/subtask`, new URLSearchParams(new FormData(subtaskForm)));
      swapHTML(`#subtasks-${taskId}`, html, "innerHTML");
      subtaskForm.reset();
    });

    // Description form
    const descForm = this.querySelector(".detail__note form");
    descForm?.addEventListener("submit", async (e) => {
      e.preventDefault();
      const body = new URLSearchParams({ description: this.description });
      await apiHTML("PUT", `/partials/tasks/${taskId}/description`, body);
    });

    // Subtask interactions (delegated — subtasks are swapped in dynamically)
    const subtasksContainer = this.querySelector(".subtasks");
    subtasksContainer?.addEventListener("click", async (e) => {
      const check = e.target.closest(".completed-check");
      const del = e.target.closest(".subtask-delete");

      if (check) {
        const subtask = check.closest(".subtask");
        const subtaskId = subtask?.dataset.id;
        if (!subtaskId) return;
        const html = await apiHTML("PUT", `/partials/tasks/${subtaskId}/completed`);
        swapHTML(`.subtask[data-id='${subtaskId}'] .completed-check`, html, "outerHTML");
      }

      if (del) {
        const subtask = del.closest(".subtask");
        const subtaskId = subtask?.dataset.id;
        if (!subtaskId) return;
        const html = await apiHTML("DELETE", `/partials/tasks/${subtaskId}`);
        swapHTML(`.subtask[data-id='${subtaskId}']`, html, "outerHTML");
      }
    });

    // Close button
    const closeBtn = this.querySelector("button.close");
    closeBtn?.addEventListener("click", async () => {
      const input = document.querySelector("input[name='task-detail']:checked");
      if (input) input.checked = false;
      await apiHTML("PUT", `/partials/tasks/${taskId}/remove-persistence`);
    });

    // Delete task (modal)
    const removeBtn = this.querySelector("button.remove");
    const dialog = this.querySelector("dialog.modal");
    const cancelBtn = this.querySelector("dialog button.default");
    const deleteBtn = this.querySelector("dialog button.warning");

    removeBtn?.addEventListener("click", () => dialog?.showModal());
    cancelBtn?.addEventListener("click", () => dialog?.close());
    deleteBtn?.addEventListener("click", async () => {
      dialog?.close();
      const html = await apiHTML("DELETE", `/partials/tasks/${taskId}`);
      swapHTML(taskSelector, html, "outerHTML");
    });
  }

  async #loadQuill() {
    if (window.Quill) return;

    if (!document.querySelector('link[href="/assets/css/libs/quill.snow.css"]')) {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = "/assets/css/libs/quill.snow.css";
      document.head.appendChild(link);
    }

    return new Promise((resolve, reject) => {
      const script = document.createElement("script");
      script.src = "/assets/js/libs/quill.js";
      script.onload = resolve;
      script.onerror = reject;
      document.head.appendChild(script);
    });
  }
}

customElements.define("rm-task-details", DetailsElement);
