class TasksAppElement extends HTMLElement {
  connectedCallback() {
    const observer = new ResizeObserver((entries) => {
      document.body.style.setProperty("--vh", `${document.body.clientHeight}px`);
    });

    observer.observe(document.body);
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
