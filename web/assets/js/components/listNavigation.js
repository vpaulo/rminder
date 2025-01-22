class ListNavigationElement extends HTMLElement {
  /** @type string */
  listId;
  /** @type string */
  listName;
  /** @type string */
  colour;
  /** @type string */
  icon;
  /** @type string */
  position;
  /** @type string */
  pinned;
  /** @type HTMLSpanElement */
  settings;
  /** @type HTMLDivElement */
  formContainer;

  #editListHandler;

  /**
   *
   * @param {MouseEvent} e
   */
  edit(e) {
    e.preventDefault();

    console.log(">>> LIST: ", this.listId, this.listName, this.colour, this.icon, this.pinned, this.position);
    const form = this.formContainer.querySelector("form");

    form.querySelectorAll("input:checked")?.forEach((el) => {
      el.removeAttribute("checked");
    });

    const name = form.querySelector("input.new-list");
    const pinned = form.querySelector('input[name="pin"]');
    const colour = form.querySelector(`input[name="swatch"][value="${this.colour}"]`);
    const icon = form.querySelector(`input[name="icon"][value="${this.icon}"]`);
    const updateBtn = form.querySelector(".btn.add-new-list");
    const removeBtn = form.querySelector(".remove-list");

    form.removeAttribute("hx-post");
    form.setAttribute("hx-put", `/lists/${this.listId}`);
    removeBtn.setAttribute("hx-delete", `/lists/${this.listId}`);

    name.value = this.listName;
    pinned.toggleAttribute("checked", this.pinned === "");
    colour.setAttribute("checked", "");
    icon.setAttribute("checked", "");
    updateBtn.innerHTML = "Update";
    removeBtn.classList.remove("hidden");

    // if you edit htmx attributes through js you need to run this
    htmx.process(form);

    this.formContainer.showModal();
  }

  connectedCallback() {
    this.formContainer = document.querySelector(".list-form-container");
    this.settings = this.querySelector(".settings");

    this.listId = this.dataset.id;
    this.listName = this.dataset.name;
    this.colour = this.dataset.colour;
    this.icon = this.dataset.icon;
    this.pinned = this.dataset.pinned;
    this.position = this.dataset.position;

    this.#editListHandler = (e) => this.edit(e);

    this.style.setProperty("--list-colour", `var(${this.colour})`);

    this.settings?.addEventListener("click", this.#editListHandler);
  }

  disconnectedCallback() {
    this.settings?.removeEventListener("click", this.#editListHandler);
  }
}

customElements.define("rm-list-nav", ListNavigationElement);
