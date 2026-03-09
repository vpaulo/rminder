import { apiHTML, swapHTML } from "../utils/fetch.js";

class SidebarElement extends HTMLElement {
  /** @type HTMLDivElement */
  cancelBtn;
  /** @type HTMLButtonElement */
  addList;
  /** @type HTMLDialogElement */
  formContainer;
  /** @type HTMLIElement */
  clearIcon;
  /** @type HTMLInputElement */
  searchInput;

  #closePopoverHandler;
  #formFocusHandler;
  #searchClearHandler;
  #formSubmitHandler;
  #removeListHandler;
  #searchSubmitHandler;

  addEvents() {
    this.#closePopoverHandler = () => this.closePopover();
    this.#formFocusHandler = () => this.formFocus();
    this.#searchClearHandler = () => this.clearSearch();
    this.#formSubmitHandler = (e) => this.handleListFormSubmit(e);
    this.#removeListHandler = (e) => this.handleRemoveList(e);
    this.#searchSubmitHandler = (e) => this.handleSearchSubmit(e);

    this.addList.addEventListener("click", this.#formFocusHandler);
    this.cancelBtn.addEventListener("click", this.#closePopoverHandler);
    this.clearIcon.addEventListener("click", this.#searchClearHandler);
    this.formContainer.querySelector("form")?.addEventListener("submit", this.#formSubmitHandler);
    this.formContainer.querySelector(".remove-list")?.addEventListener("click", this.#removeListHandler);
    this.querySelector(".searchbox__form")?.addEventListener("submit", this.#searchSubmitHandler);
  }

  async handleListFormSubmit(e) {
    e.preventDefault();
    const form = this.formContainer.querySelector("form");
    const listId = form.dataset.listId;
    const method = listId ? "PUT" : "POST";
    const url = listId ? `/partials/lists/${listId}` : "/partials/lists/create";
    const html = await apiHTML(method, url, new URLSearchParams(new FormData(form)));
    swapHTML(".lists__container", html, "outerHTML");
    this.closePopover();
  }

  async handleRemoveList(e) {
    e.preventDefault();
    const form = this.formContainer.querySelector("form");
    const listId = form.dataset.listId;
    if (!listId) return;
    const html = await apiHTML("DELETE", `/partials/lists/${listId}`);
    swapHTML(".lists__container", html, "outerHTML");
    this.closePopover();
  }

  async handleSearchSubmit(e) {
    e.preventDefault();
    const form = e.target;
    const html = await apiHTML("POST", "/partials/lists/search", new URLSearchParams(new FormData(form)));
    swapHTML(".main", html, "innerHTML");
    // refresh sidebar list counts after delay
    setTimeout(async () => {
      const listsHtml = await apiHTML("GET", "/partials/lists/all");
      swapHTML(".lists__container", listsHtml, "outerHTML");
    }, 20);
  }

  closePopover() {
    this.formContainer.close();
    const form = this.formContainer.querySelector("form");
    const btn = form.querySelector(".btn.add-new-list");
    const removeBtn = form.querySelector(".remove-list");
    form.reset();

    delete form.dataset.listId;

    form.querySelectorAll("input:checked")?.forEach((el) => {
      el.removeAttribute("checked");
    });
    form.querySelectorAll("option")?.forEach((el) => {
      el.removeAttribute("selected");
    });
    form.querySelector('input[name="swatch"]')?.setAttribute("checked", "");
    form.querySelector('input[name="icon"]')?.setAttribute("checked", "");
    btn.innerHTML = "Add";
    removeBtn?.classList.add("hidden");
    delete removeBtn?.dataset.listId;
  }

  formFocus() {
    this.formContainer.showModal();
    this.formContainer.querySelector(".new-list").focus();
  }

  clearSearch() {
    this.searchInput.value = "";
  }

  connectedCallback() {
    this.cancelBtn = this.querySelector(".cancel-new-list");
    this.addList = this.querySelector(".add-list");
    this.formContainer = this.querySelector(".list-form-container");
    this.clearIcon = this.querySelector(".searchbox__container i.clear-icon");
    this.searchInput = this.querySelector(".searchbox");

    this.addEvents();
  }

  disconnectedCallback() {
    this.addList.removeEventListener("click", this.#formFocusHandler);
    this.cancelBtn.removeEventListener("click", this.#closePopoverHandler);
    this.clearIcon.removeEventListener("click", this.#searchClearHandler);
    this.formContainer.querySelector("form")?.removeEventListener("submit", this.#formSubmitHandler);
    this.formContainer.querySelector(".remove-list")?.removeEventListener("click", this.#removeListHandler);
    this.querySelector(".searchbox__form")?.removeEventListener("submit", this.#searchSubmitHandler);
  }
}

customElements.define("rm-sidebar", SidebarElement);
