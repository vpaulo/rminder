class PremiumElement extends HTMLElement {
  quill;

  get description() {
    return this.quill.getSemanticHTML();
  }

  connectedCallback() {
    const dialog = this.querySelector("dialog.modal");
    const get_premium_btn = this.querySelector("button.primary");
    const buy_btn = this.querySelector("dialog button.primary");
    const cancel_btn = this.querySelector("dialog button.default");
 
    get_premium_btn.addEventListener("click", () => {
      dialog?.showModal();
    });
    buy_btn.addEventListener("click", () => {
      dialog?.close();
    });
    cancel_btn.addEventListener("click", () => {
      dialog?.close();
    });
  }
}

customElements.define("rm-get-premium", PremiumElement);
