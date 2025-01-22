class SidebarElement extends HTMLElement {
  /** @type HTMLDivElement */
  cancelBtn;
  /** @type HTMLButtonElement */
  addList;
  /** @type HTMLDivElement */
  formContainer;

  #closePopoverHandler;
  #formFocusHandler;

  addEvents() {
    this.#closePopoverHandler = () => this.closePopover();
    this.#formFocusHandler = () => this.formFocus();

    this.addList.addEventListener("click", this.#formFocusHandler);
    this.cancelBtn.addEventListener("click", this.#closePopoverHandler);
  }

  closePopover() {
    this.formContainer.close();
    const form = this.formContainer.querySelector("form");
    const btn = form.querySelector(".btn.add-new-list");
    const removeBtn = form.querySelector(".remove-list");
    form.reset();

    form.removeAttribute("hx-put");
    form.setAttribute("hx-post", "/lists/create");

    form.querySelectorAll("input:checked")?.forEach((el) => {
      el.removeAttribute("checked");
    });
    form.querySelector('input[name="swatch"]')?.setAttribute("checked", "");
    form.querySelector('input[name="icon"]')?.setAttribute("checked", "");
    btn.innerHTML = "Add";
    removeBtn?.classList.add("hidden");
    removeBtn?.removeAttribute("hx-delete");

    // if you edit htmx attributes through js you need to run this
    htmx.process(form);
  }

  formFocus() {
    this.formContainer.showModal();
    this.formContainer.querySelector(".new-list").focus();
  }

  connectedCallback() {
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
