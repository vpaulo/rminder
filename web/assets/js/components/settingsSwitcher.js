import { tryCatch } from "../utils.js";

class SettingsSwitcherElement extends HTMLElement {
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
      // const [error] = await tryCatch(
      //   fetch("/v1/import", {
      //     method: "POST",
      //     headers: {
      //       "Content-Type": "application/json",
      //     },
      //     body: JSON.stringify(tasks),
      //   }),
      // );
      // if (error) {
      //   console.error("POST: Import lists: ", error);
      // }
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

  connectedCallback() {
    const rect = this.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();

    this.style.setProperty("--popover-top", `${rect.bottom}px`);
    this.style.setProperty("--popover-right", `${bodyRect.right - rect.right + 10}px`);

    this.handleClick = this.#handleClick.bind(this);
    this.addEventListener("click", this.handleClick, false);
  }

  disconnectedCallback() {
    this.removeEventListener("click", this.handleClick, false);
  }
}

customElements.define("rm-settings-switcher", SettingsSwitcherElement);
