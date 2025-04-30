import { tryCatch } from "../utils.js";

class SettingsSwitcherElement extends HTMLElement {
  /** @type HTMLDialogElement */
  formContainer;
  /** @type HTMLInputElement */
  fileElem;
  /** @type HTMLDivElement */
  messages;

  maxFileSize = 10 * 1024 * 1024; // 10MB
  allowedTypes = ["application/json"];

  /** @type {function(MouseEvent): void} */
  async #handleClick(e) {
    if (!e.target.dataset.action) return;

    if (e.target.dataset.action === "export") {
      const [error, data] = await tryCatch(fetch("/v1/export"));

      if (error) {
        console.error("POST: Export lists: ", error);
      }

      const [err, resp] = await tryCatch(data.json());

      if (err) {
        console.error("JSON: Export lists response: ", err);
      }

      this.saveJSON(resp, "rminder_data");
    }

    if (e.target.dataset.action === "import") {
      this.formContainer.showModal();
    }
  }

  /**
   *
   * @param {[T]} data
   * @param {string} saveAs
   */
  saveJSON(data, saveAs) {
    var stringified = JSON.stringify(data, null, 2);
    var blob = new Blob([stringified], { type: "application/json" });
    var url = URL.createObjectURL(blob);

    var a = document.createElement("a");
    a.download = saveAs + ".json";
    a.href = url;
    a.id = saveAs;
    document.body.appendChild(a);
    a.click();
    document.querySelector("#" + a.id).remove();
  }

  #handleFiles() {
    this.validateFile(this.inputElem.files[0]);
  }

  validateFile(file) {
    this.messages.textContent = ""; // clear messages
    if (!this.allowedTypes.includes(file.type)) {
      this.showError(`${file.name} is not an allowed file type.`);
      return false;
    }

    if (file.size > this.maxFileSize) {
      this.showError(`${file.name} exceeds the maximum file size of 10MB.`);
      return false;
    }

    return true;
  }

  showError(message) {
    this.messages.textContent = message;
  }

  connectedCallback() {
    const rect = this.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();
    this.formContainer = this.querySelector(".file-form-container");
    this.inputElem = this.querySelector("#file");
    this.messages = this.querySelector("#messages");

    this.style.setProperty("--popover-top", `${rect.bottom}px`);
    this.style.setProperty("--popover-right", `${bodyRect.right - rect.right + 10}px`);

    this.handleClick = this.#handleClick.bind(this);
    this.handleFiles = this.#handleFiles.bind(this);

    this.addEventListener("click", this.handleClick, false);
    this.inputElem.addEventListener("change", this.handleFiles, false);
  }

  disconnectedCallback() {
    this.removeEventListener("click", this.handleClick, false);
    this.inputElem.removeEventListener("change", this.handleFiles, false);
  }
}

customElements.define("rm-settings-switcher", SettingsSwitcherElement);
