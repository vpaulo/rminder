/** @typedef {{ Next: number, Previous: number }} NavigationType */
const navigationType = {
  Next: 0,
  Previous: 1,
};
/** @typedef {{ Task: number, List: number }} NavigationElement */
const navigationElement = {
  Task: 0,
  List: 1,
};

class TasksAppElement extends HTMLElement {
  /** @type {ResizeObserver} */
  observer;
  /** @type {boolean} */
  isEditing = false;

  /** @type {function(KeyboardEvent): void} */
  #keybidings(e) {
    if (this.editing) return;

    switch (true) {
      case e.code === "ArrowDown" || e.code === "KeyJ":
        this.navigation(navigationElement.Task, navigationType.Next);
        break;
      case e.code === "ArrowUp" || e.code === "KeyK":
        this.navigation(navigationElement.Task, navigationType.Previous);
        break;
      case e.code === "BracketLeft":
        this.navigation(navigationElement.List, navigationType.Next);
        break;
      case e.code === "BracketRight":
        this.navigation(navigationElement.List, navigationType.Previous);
        break;
      case e.code === "KeyA":
        e.preventDefault();
        this.querySelector("#task")?.focus();
        break;
      case e.code === "KeyN":
        e.preventDefault();
        this.querySelector(".add-list")?.click();
        break;
      case e.code === "KeyF":
        e.preventDefault();
        this.querySelector(".searchbox")?.focus();
        break;
      case e.code === "KeyX":
        this.toggleTaskCompletion();
        break;
      case e.code === "KeyI":
        this.toggleTaskImportance();
        break;
      case e.code === "KeyP":
        this.switchPriority();
        break;
    }
  }

  /** @type {function(FocusEvent): void} */
  #focusInHandler(e) {
    this.setEdit(e.target, true);
  }

  /** @type {function(FocusEvent): void} */
  #focusOutHandler(e) {
    this.setEdit(e.target, false);
  }

  /** @type {function(HTMLElement, boolean): void} */
  setEdit(element, state) {
    if (
      element.tagName.toLowerCase() === "input" ||
      (element.tagName.toLowerCase() === "div" && element.classList.hasClass("ql-editor"))
    ) {
      this.editing = state;
    }
  }

  /** @type {function(NavigationElement, NavigationType): void} */
  navigation(element, type) {
    const items =
      element === navigationElement.Task
        ? this.querySelectorAll(".tasks__list > li input")
        : this.querySelectorAll("rm-list-nav input");

    if (items.length === 0) return;

    const lastIndex = items.length - 1;
    const selected = [...items].findIndex((e) => e.checked);

    if (type === navigationType.Next && selected === lastIndex) return;
    if (type === navigationType.Previous && selected === 0) return;

    if (selected < 0) {
      items[type === navigationType.Next ? 0 : lastIndex].click();
    } else {
      items[type === navigationType.Next ? selected + 1 : selected - 1].click();
    }
  }

  toggleTaskCompletion() {
    const selected = [...this.querySelectorAll(".tasks__list > li input")]?.findIndex((e) => e.checked);

    if (selected >= 0) {
      this.querySelectorAll(".tasks__list > li .completed-ckeck")[selected].click();
    }
  }

  toggleTaskImportance() {
    const selected = [...this.querySelectorAll(".tasks__list > li input")]?.findIndex((e) => e.checked);

    if (selected >= 0) {
      this.querySelectorAll(".tasks__list > li .importance-check")[selected].click();
    }
  }

  switchPriority() {
    const selected = [...this.querySelectorAll(".tasks__list > li input")]?.findIndex((e) => e.checked);

    if (selected >= 0) {
      const element = this.querySelector(".details__priority select[name='priority']");

      element.value = element.value === "4" ? 0 : Number(element.value) + 1;
      element.dispatchEvent(new Event("change"));
      // TODO: when first option is selected it does not show the name
    }
  }

  connectedCallback() {
    this.observer = new ResizeObserver((entries) => {
      document.body.style.setProperty("--vh", `${document.body.clientHeight}px`);
    });

    this.observer.observe(document.body);

    this.keybidings = this.#keybidings.bind(this);
    this.focusInHandler = this.#focusInHandler.bind(this);
    this.focusOutHandler = this.#focusOutHandler.bind(this);

    document.addEventListener("keydown", this.keybidings, false);
    document.addEventListener("focusin", this.focusInHandler, false);
    document.addEventListener("focusout", this.focusOutHandler, false);
  }

  disconnectedCallback() {
    this.observer.disconnect();
    document.removeEventListener("keydown", this.keybidings, false);
    document.removeEventListener("focusin", this.focusInHandler, false);
    document.removeEventListener("focusout", this.focusOutHandler, false);
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
