class DetailsElement extends HTMLElement {
  quill;

  get description() {
    return this.quill.getSemanticHTML();
  }

  async connectedCallback() {
    await this.#loadQuill();

    const task_id = this.getAttribute("task-id");
    const close_btn = this.querySelector("button.close");
    const remove_btn = this.querySelector("button.remove");
    const dialog = this.querySelector("dialog.modal");
    const cancel_btn = this.querySelector("dialog button.default");
    const delete_btn = this.querySelector("dialog button.warning");

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

    close_btn.addEventListener("click", () => {
      const input = document.querySelector("input[name='task-detail']:checked");
      if (input) {
        input.checked = false;
      }
    });

    remove_btn.addEventListener("click", () => {
      dialog?.showModal();
    });
    cancel_btn.addEventListener("click", () => {
      dialog?.close();
    });
    delete_btn.addEventListener("click", () => {
      dialog?.close();
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
