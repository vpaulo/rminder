import { tryCatch } from "../../helpers/tryCatch.js";

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

/** @typedef {{ id: number, position: number }} Reorder */

function csrfToken() {
  return document.querySelector('meta[name="csrf-token"]')?.content || "";
}

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
      case e.shiftKey === false && e.code === "BracketLeft":
        this.navigation(navigationElement.List, navigationType.Next);
        break;
      case e.shiftKey === false && e.code === "BracketRight":
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
      case e.shiftKey === true && e.code === "BracketLeft":
        this.reorderList(navigationType.Next);
        break;
      case e.shiftKey === true && e.code === "BracketRight":
        this.reorderList(navigationType.Previous);
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
    const targetItem = e.target.classList.contains("drag-item") ? e.target : e.target.closest(".drag-item");

    if (targetItem !== this.draggedTask && targetItem) {
      const boundingRect = targetItem.getBoundingClientRect();
      const offset = boundingRect.top + boundingRect.height / 2;
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

    const targetItem = e.target.classList.contains("drag-item") ? e.target : e.target.closest(".drag-item");

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

    // TODO: maybe add a delay or debounce the call to server
    this.updateTaskPosition();
  }

  /** @type {function(HTMLElement, boolean): void} */
  setEdit(element, state) {
    if (
      element.tagName.toLowerCase() === "input" ||
      (element.tagName.toLowerCase() === "div" && element.classList.contains("ql-editor"))
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

    if (element === navigationElement.Task) {
      this.querySelector(".tasks")?.scrollTo({
        top: selected >= lastIndex / 2 ? (selected || 1) * 44 : 0, // 44 is the min height for a task
        behavior: "smooth",
      });
    }
  }

  toggleTaskCompletion() {
    const selected = [...this.querySelectorAll(".tasks__list > li input")]?.findIndex((e) => e.checked);

    if (selected >= 0) {
      this.querySelectorAll(".tasks__list > li .completed-check")[selected].click();
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
    // TODO: maybe add a delay or debounce the call to server
    this.updateTaskPosition();
  }

  async updateTaskPosition() {
    /** @type {Reorder[]} */
    const tasks = [];

    this.querySelectorAll(".tasks__list > li")?.forEach((task, i) => {
      if (+task.dataset.position !== i) {
        tasks.push({
          id: +task.dataset.id,
          position: i,
        });
        task.dataset.position = i;
      }
    });

    if (tasks.length === 0) return;

    const [error] = await tryCatch(
      fetch("/api/tasks/reorder", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": csrfToken(),
        },
        body: JSON.stringify(tasks),
      }),
    );

    if (error) {
      console.error("POST: Tasks reorder: ", error);
    }
  }

  /** @type {function(NavigationType): void} */
  reorderList(type) {
    window.getSelection().removeAllRanges();
    const items = this.querySelectorAll(".lists-holder > rm-list-nav input");

    if (items.length === 0) return;

    const selected = [...items].find((e) => e.checked)?.closest("rm-list-nav");

    if (!selected) return;

    if (type === navigationType.Next && selected.nextSibling) {
      selected.parentNode.insertBefore(selected, selected.nextSibling.nextSibling);
    }

    if (type === navigationType.Previous && selected.previousSibling) {
      selected.parentNode.insertBefore(selected, selected.previousSibling);
    }
    // TODO: maybe add a delay or debounce the call to server
    this.updateListPosition();
  }

  // TODO: maybe combine updateListPosition with updateTaskPosition, because it's very similar
  async updateListPosition() {
    /** @type {Reorder[]} */
    const lists = [];

    this.querySelectorAll(".lists-holder > rm-list-nav")?.forEach((list, i) => {
      if (+list.dataset.position !== i) {
        lists.push({
          id: +list.dataset.id,
          position: i,
        });
        list.dataset.position = i;
      }
    });

    if (lists.length === 0) return;

    const [error] = await tryCatch(
      fetch("/api/lists/reorder", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": csrfToken(),
        },
        body: JSON.stringify(lists),
      }),
    );

    if (error) {
      console.error("POST: Lists reorder: ", error);
    }
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
    this.tasksList?.addEventListener("dragstart", this.handleDragStart, false);
    this.tasksList?.addEventListener("dragover", this.handleDragOver, false);
    this.tasksList?.addEventListener("drop", this.handleDrop, false);
  }

  disconnectedCallback() {
    this.observer.disconnect();
    document.removeEventListener("keydown", this.keybidings, false);
    document.removeEventListener("focusin", this.focusInHandler, false);
    document.removeEventListener("focusout", this.focusOutHandler, false);

    this.tasksList?.removeEventListener("dragstart", this.handleDragStart, false);
    this.tasksList?.removeEventListener("dragover", this.handleDragOver, false);
    this.tasksList?.removeEventListener("drop", this.handleDrop, false);
  }
}

customElements.define("rm-tasks-app", TasksAppElement);
