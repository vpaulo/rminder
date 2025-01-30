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
  /** @type string */
  filter;
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

    console.log(
      ">>> LIST: ",
      this.listId,
      this.listName,
      this.colour,
      this.icon,
      this.pinned,
      this.position,
      this.filter,
    );
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
    // filters
    const includeEl = form.querySelector('select[name="include"]');
    const completedEl = form.querySelector('select[name="completed"]');
    const importantEl = form.querySelector('select[name="important"]');
    const priorityEl = form.querySelector('select[name="priority"]');
    const dateEl = form.querySelector('select[name="date"]');
    const startRangeEl = form.querySelector('input[name="from"]');
    const endRangeEl = form.querySelector('input[name="to"]');
    const [include, completed, important, priority, date, startRange, endRange] = this.filter
      .split(";")
      .map((f) => f.split("=")[1]);

    form.removeAttribute("hx-post");
    form.setAttribute("hx-put", `/lists/${this.listId}`);
    removeBtn.setAttribute("hx-delete", `/lists/${this.listId}`);

    name.value = this.listName;
    pinned.toggleAttribute("checked", this.pinned === "");
    colour.setAttribute("checked", "");
    icon.setAttribute("checked", "");
    updateBtn.innerHTML = "Update";
    removeBtn.classList.remove("hidden");

    // Set filters
    includeEl.value = include;
    completedEl.value = completed;
    importantEl.value = important;
    priorityEl.value = priority;
    dateEl.value = date;
    startRangeEl.value = startRange;
    endRangeEl.value = endRange;

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
    this.filter = this.dataset.filter;

    this.#editListHandler = (e) => this.edit(e);

    this.style.setProperty("--list-colour", `var(${this.colour})`);

    this.settings?.addEventListener("click", this.#editListHandler);
  }

  disconnectedCallback() {
    this.settings?.removeEventListener("click", this.#editListHandler);
  }
}

customElements.define("rm-list-nav", ListNavigationElement);
