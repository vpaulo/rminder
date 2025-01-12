class ListNavigationElement extends HTMLElement {
  /** @type string */
  colour;

  // init() {
  //   this.insertAdjacentHTML("afterbegin", ``);
  // }

  connectedCallback() {
    this.colour = this.dataset.colour;

    this.style.setProperty("--list-colour", `var(${this.colour})`);

    // this.init();
  }

  disconnectedCallback() {}
}

customElements.define("rm-list-nav", ListNavigationElement);
