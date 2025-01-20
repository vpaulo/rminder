class ListNavigationElement extends HTMLElement {
  /** @type string */
  colour;

  connectedCallback() {
    this.colour = this.dataset.colour;

    this.style.setProperty("--list-colour", `var(${this.colour})`);
  }

  disconnectedCallback() {}
}

customElements.define("rm-list-nav", ListNavigationElement);
