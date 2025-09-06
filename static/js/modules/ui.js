// UI utility functions module
import { DOM } from "../utils/dom.js";

export const UI = {
  updateSliders() {
    const tempInput = document.getElementById("temperature-input");
    const tempValue = document.getElementById("temperature-value");
    const maxTokensInput = document.getElementById("max-tokens-input");
    const maxTokensValue = document.getElementById("max-tokens-value");

    tempInput.addEventListener("input", (e) => {
      DOM.updateSliderValue(e.target, tempValue);
    });

    maxTokensInput.addEventListener("input", (e) => {
      DOM.updateSliderValue(e.target, maxTokensValue);
    });
  },

  showError(message) {
    const errorMessage = document.getElementById("error-message");
    errorMessage.textContent = message;
    DOM.showElement(errorMessage);
    errorMessage.scrollIntoView({ behavior: "smooth", block: "nearest" });
  },

  hideError() {
    const errorMessage = document.getElementById("error-message");
    DOM.hideElement(errorMessage);
  },

  showSuccessMessage(message) {
    const successDiv = document.createElement("div");
    successDiv.className = "success-message";
    successDiv.textContent = message;
    successDiv.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: var(--fg);
      color: var(--bg);
      padding: 0.75rem 1rem;
      border-radius: var(--radius);
      z-index: 1000;
      font-family: var(--font);
      font-size: 0.9rem;
    `;

    document.body.appendChild(successDiv);

    setTimeout(() => {
      if (successDiv.parentNode) {
        successDiv.parentNode.removeChild(successDiv);
      }
    }, 3000);
  },

  autoResizeTextarea(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + "px";
  },
};
