class SidebarElement extends HTMLElement {
  connectedCallback() {
    console.log("Hello sidebar");

    this.init();
  }

  init() {
    this.insertAdjacentHTML(
      "afterbegin",
      `
    <div class="header">
      <label>
        <span class="menu" aria-label="Toggle sidebar"></span>
        <input class="hidden" type="checkbox" value="1" checked/>
      </label>
    </div>`,
    );
  }
}

customElements.define("rm-sidebar", SidebarElement);
