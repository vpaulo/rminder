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
  /** @type {HTMLUListElement} */
  tasksList;
  /** @type {HTMLLIElement} */
  draggedTask = null;

  /** @type {function(KeyboardEvent): void} */
  #keybidings(e) {
    if (this.editing) return;

    switch (true) {
      case e.shiftKey === false && (e.code === "ArrowDown" || e.code === "KeyJ"):
        this.navigation(navigationElement.Task, navigationType.Next);
        break;
      case e.shiftKey === false && (e.code === "ArrowUp" || e.code === "KeyK"):
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
      case e.shiftKey === true && (e.code === "ArrowDown" || e.code === "KeyJ"):
        this.reorderTask(navigationType.Next);
        break;
      case e.shiftKey === true && (e.code === "ArrowUp" || e.code === "KeyK"):
        this.reorderTask(navigationType.Previous);
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

  /** @type {function(DragEvent): void} */
  #handleDragStart(e) {
    this.draggedTask = e.target;
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/html", this.draggedTask.innerHTML);
    e.target.style.opacity = "0.5";
  }

  /** @type {function(DragEvent): void} */
  #handleDragOver(e) {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    const targetItem = e.target.classList.contains("drag-item")
      ? e.target
      : e.target.parentNode.classList.contains("drag-item")
        ? e.target.parentNode
        : null;

    if (targetItem !== this.draggedTask && targetItem) {
      const boundingRect = targetItem.getBoundingClientRect();
      const offset = boundingRect.y + boundingRect.height / 2;
      if (e.clientY - offset > 0) {
        targetItem.style.borderBottom = "solid 2px transparent";
        targetItem.style.borderTop = "";
      } else {
        targetItem.style.borderTop = "solid 2px transparent";
        targetItem.style.borderBottom = "";
      }
    }
  }

  /** @type {function(DragEvent): void} */
  #handleDrop(e) {
    e.preventDefault();

    const targetItem = e.target.classList.contains("drag-item")
      ? e.target
      : e.target.parentNode.classList.contains("drag-item")
        ? e.target.parentNode
        : null;

    if (targetItem !== this.draggedTask && targetItem) {
      if (e.clientY > targetItem.getBoundingClientRect().top + targetItem.offsetHeight / 2) {
        targetItem.parentNode.insertBefore(this.draggedTask, targetItem.nextSibling);
      } else {
        targetItem.parentNode.insertBefore(this.draggedTask, targetItem);
      }
    }
    this.querySelectorAll(".drag-item").forEach((el) => {
      el.style.borderTop = "";
      el.style.borderBottom = "";
    });
    this.draggedTask.style.opacity = "";
    this.draggedTask = null;

    // TODO: POST order changes
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

  /** @type {function(NavigationType): void} */
  reorderTask(type) {
    window.getSelection().removeAllRanges();
    const items = this.querySelectorAll(".tasks__list > li input");

    if (items.length === 0) return;

    const selected = [...items].find((e) => e.checked)?.closest(".drag-item");

    if (!selected) return;

    if (type === navigationType.Next && selected.nextSibling) {
      selected.parentNode.insertBefore(selected, selected.nextSibling.nextSibling);
    }

    if (type === navigationType.Previous && selected.previousSibling) {
      selected.parentNode.insertBefore(selected, selected.previousSibling);
    }

    // TODO: POST order changes
  }

  connectedCallback() {
    this.tasksList = this.querySelector(".tasks__list");

    this.observer = new ResizeObserver((entries) => {
      document.body.style.setProperty("--vh", `${document.body.clientHeight}px`);
    });

    this.observer.observe(document.body);

    this.keybidings = this.#keybidings.bind(this);
    this.focusInHandler = this.#focusInHandler.bind(this);
    this.focusOutHandler = this.#focusOutHandler.bind(this);

    this.handleDragStart = this.#handleDragStart.bind(this);
    this.handleDragOver = this.#handleDragOver.bind(this);
    this.handleDrop = this.#handleDrop.bind(this);

    document.addEventListener("keydown", this.keybidings, false);
    document.addEventListener("focusin", this.focusInHandler, false);
    document.addEventListener("focusout", this.focusOutHandler, false);

    // Add event listeners for drag and drop events
    this.tasksList.addEventListener("dragstart", this.handleDragStart, false);
    this.tasksList.addEventListener("dragover", this.handleDragOver, false);
    this.tasksList.addEventListener("drop", this.handleDrop, false);
  }

  disconnectedCallback() {
    this.observer.disconnect();
    document.removeEventListener("keydown", this.keybidings, false);
    document.removeEventListener("focusin", this.focusInHandler, false);
    document.removeEventListener("focusout", this.focusOutHandler, false);

    this.tasksList.removeEventListener("dragstart", this.handleDragStart, false);
    this.tasksList.removeEventListener("dragover", this.handleDragOver, false);
    this.tasksList.removeEventListener("drop", this.handleDrop, false);
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
