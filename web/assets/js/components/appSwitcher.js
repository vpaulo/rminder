class AppSwitcherElement extends HTMLElement {
  connectedCallback() {
    const rect = this.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();

    this.style.setProperty("--popover-top", `${rect.bottom}px`);
    this.style.setProperty("--popover-right", `${bodyRect.right - rect.right + 10}px`);
  }
}

customElements.define("rm-app-switcher", AppSwitcherElement);
