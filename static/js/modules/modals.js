// Modal management module
export class Modal {
  constructor(modalId, closeButtonId = null) {
    this.modal = document.getElementById(modalId);
    this.closeButton = closeButtonId ? document.getElementById(closeButtonId) : null;
    this.setupEventListeners();
  }

  setupEventListeners() {
    // Close button listener
    if (this.closeButton) {
      this.closeButton.addEventListener("click", () => this.close());
    }

    // Backdrop click listener - for dialog elements, we need to check for clicks outside content
    this.modal.addEventListener("click", (e) => {
      const rect = this.modal.getBoundingClientRect();
      if (e.clientX < rect.left || e.clientX > rect.right || e.clientY < rect.top || e.clientY > rect.bottom) {
        this.close();
      }
    });
  }

  open() {
    if (this.modal.tagName.toLowerCase() === "dialog") {
      this.modal.showModal();
    } else {
      this.modal.style.display = "flex";
      document.body.style.overflow = "hidden";
    }
  }

  close() {
    if (this.modal.tagName.toLowerCase() === "dialog") {
      this.modal.close();
    } else {
      this.modal.style.display = "none";
      document.body.style.overflow = "";
    }
  }
}
