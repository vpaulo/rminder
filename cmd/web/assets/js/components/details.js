class DetailsElement extends HTMLElement {
  quill;

  get description() {
    return this.quill.getSemanticHTML();
  }

  connectedCallback() {
    const task_id = this.getAttribute("task-id");
    const close_btn = this.querySelector("button.close");
    const remove_btn = this.querySelector("button.remove");
    const dialog = this.querySelector("dialog.modal");
    const cancel_btn = this.querySelector("dialog button.default");
    const delete_btn = this.querySelector("dialog button.warning");
    console.log("task id : ", task_id);

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
}

customElements.define("rm-task-details", DetailsElement);
