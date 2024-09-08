class DetailsElement extends HTMLElement {
  connectedCallback() {
    const task_id = this.getAttribute("task-id");
    console.log("task id : ", task_id);

    const quill = new Quill("#note-editor", {
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
  }
}

customElements.define("rm-task-details", DetailsElement);
