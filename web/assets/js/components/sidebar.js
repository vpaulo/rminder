class SidebarElement extends HTMLElement {
  /** @type HTMLDivElement */
  cancelBtn;
  /** @type HTMLButtonElement */
  addList;
  /** @type HTMLDivElement */
  formContainer;

  #closePopoverHandler;
  #formFocusHandler;

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

  addEvents() {
    this.#closePopoverHandler = () => this.closePopover();
    this.#formFocusHandler = () => this.formFocus();

    this.addList.addEventListener("click", this.#formFocusHandler);
    this.cancelBtn.addEventListener("click", this.#closePopoverHandler);
  }

  closePopover() {
    this.formContainer.close();
    this.formContainer.querySelector("form").reset();
  }

  formFocus() {
    this.formContainer.showModal();
    this.formContainer.querySelector(".new-list").focus();
  }

  connectedCallback() {
    this.init();

    this.cancelBtn = this.querySelector(".cancel-new-list");
    this.addList = this.querySelector(".add-list");
    this.formContainer = this.querySelector(".list-form-container");

    this.addEvents();
  }

  disconnectedCallback() {
    this.addList.removeEventListener("click", this.#formFocusHandler);
    this.cancelBtn.removeEventListener("click", this.#closePopoverHandler);
  }
}

customElements.define("rm-sidebar", SidebarElement);
