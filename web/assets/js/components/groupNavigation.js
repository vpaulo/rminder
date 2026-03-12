class GroupNavigationElement extends HTMLElement {
  /** @type string */
  groupId;
  /** @type string */
  groupName;
  /** @type HTMLSpanElement */
  settings;
  /** @type HTMLDialogElement */
  formContainer;

  #editGroupHandler;

  /**
   * @param {MouseEvent} e
   */
  edit(e) {
    e.preventDefault();

    const form = this.formContainer.querySelector("form");
    const nameInput = form.querySelector("input.new-group");
    const updateBtn = form.querySelector(".btn.add-new-group");
    const removeBtn = form.querySelector(".remove-group");

    form.removeAttribute("hx-post");
    form.setAttribute("hx-put", `/partials/groups/${this.groupId}`);
    removeBtn.setAttribute("hx-delete", `/partials/groups/${this.groupId}`);

    nameInput.value = this.groupName;
    updateBtn.innerHTML = "Update";
    removeBtn.classList.remove("hidden");

    htmx.process(form);
    htmx.process(removeBtn);

    this.formContainer.showModal();
    nameInput.focus();
  }

  connectedCallback() {
    this.formContainer = document.querySelector(".group-form-container");
    this.settings = this.querySelector(".group-settings");

    this.groupId = this.dataset.groupId;
    this.groupName = this.querySelector(".group-name")?.textContent?.trim() ?? "";

    this.#editGroupHandler = (e) => this.edit(e);

    this.settings?.addEventListener("click", this.#editGroupHandler);
  }

  disconnectedCallback() {
    this.settings?.removeEventListener("click", this.#editGroupHandler);
  }
}

customElements.define("rm-group-nav", GroupNavigationElement);
