class TasksAppElement extends HTMLElement {
  /** @type ResizeObserver */
  observer;

  connectedCallback() {
    this.observer = new ResizeObserver((entries) => {
      document.body.style.setProperty("--vh", `${document.body.clientHeight}px`);
    });

    this.observer.observe(document.body);
  }

  disconnectedCallback() {
    this.observer.disconnect();
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
