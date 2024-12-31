class TasksAppElement extends HTMLElement {
  connectedCallback() {
    console.log("Hello tasks app");
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
