class DetailsElement extends HTMLElement {
  quill;

  get description() {
    return this.quill.getSemanticHTML();
  }

  connectedCallback() {
    const task_id = this.getAttribute("task-id");
    const close_btn = this.querySelector("button.close");
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
  }
}

customElements.define("rm-task-details", DetailsElement);
