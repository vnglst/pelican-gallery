// Model selection management module
import { DOM } from "../utils/dom.js";
import { UI } from "./ui.js";

export class ModelManager {
  constructor(elements, modals, selectedModels) {
    this.elements = elements;
    this.modals = modals;
    this.selectedModels = selectedModels;
  }

  // Initialize with default selected models
  initializeSelectedModels() {
    // Find all default model checkboxes and select them
    const defaultModelCheckboxes = DOM.queryAll('input[type="checkbox"][checked]');
    this.selectedModels.length = 0; // Clear array
    Array.from(defaultModelCheckboxes).forEach((checkbox) => {
      const modelCard = checkbox.closest(".model-card");
      // Apply visual selected state to pre-checked models
      this.updateModelCardSelection(modelCard, true);
      this.selectedModels.push({
        id: checkbox.value,
        name: modelCard.dataset.modelName,
      });
    });
    this.updateSelectedModelsDisplay();
  }

  // Setup event listeners for model cards
  setupModelCardListeners() {
    document.addEventListener("click", this.handleModelCardClick.bind(this));
  }

  // Handle model card click
  handleModelCardClick(e) {
    const modelCard = e.target.closest(".model-card");
    if (modelCard) {
      const checkbox = modelCard.querySelector('input[type="checkbox"]');

      // If the click was directly on the checkbox, let it handle itself naturally
      if (e.target === checkbox) {
        // Update visual state after the checkbox change
        setTimeout(() => {
          this.updateModelCardSelection(modelCard, checkbox.checked);
        }, 0);
        return;
      }

      // For any other click within the card, toggle the checkbox
      checkbox.checked = !checkbox.checked;
      this.updateModelCardSelection(modelCard, checkbox.checked);
    }
  }

  // Update model card selection visual state
  updateModelCardSelection(modelCard, selected) {
    if (selected) {
      modelCard.classList.add("selected");
    } else {
      modelCard.classList.remove("selected");
    }
  }

  // Clear all models
  clearAllModels() {
    const checkboxes = this.modals.model.modal.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((checkbox) => {
      checkbox.checked = false;
      this.updateModelCardSelection(checkbox.closest(".model-card"), false);
    });
  }

  // Select all models
  selectAllModels() {
    const checkboxes = this.modals.model.modal.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((checkbox) => {
      checkbox.checked = true;
      this.updateModelCardSelection(checkbox.closest(".model-card"), true);
    });
  }

  // Save model selection
  saveModelSelection() {
    const checkboxes = this.modals.model.modal.querySelectorAll('input[type="checkbox"]:checked');
    this.selectedModels.length = 0; // Clear array
    Array.from(checkboxes).forEach((checkbox) => {
      this.selectedModels.push({
        id: checkbox.value,
        name: checkbox.closest(".model-card").dataset.modelName,
      });
    });

    this.updateSelectedModelsDisplay();
    this.modals.model.close();

    if (this.selectedModels.length > 0) {
      UI.showSuccessMessage(`${this.selectedModels.length} model(s) selected`);
    }
  }

  // Update selected models display
  updateSelectedModelsDisplay() {
    const modelCount = DOM.get("model-count");

    if (this.selectedModels.length === 0) {
      this.elements.modelSelectorText.textContent = "Select Models";
      DOM.hideElement(this.elements.selectedModelsList);
      DOM.hideElement(modelCount);
    } else {
      this.elements.modelSelectorText.textContent = `${this.selectedModels.length} Model${
        this.selectedModels.length === 1 ? "" : "s"
      } Selected`;
      modelCount.textContent = `(${this.selectedModels.length})`;
      modelCount.style.display = "inline";

      // Show selected models list
      DOM.showElement(this.elements.selectedModelsList);
      this.elements.selectedModelsList.innerHTML = this.selectedModels
        .map(
          (model) => `
        <div class="selected-model-item">
          <span>${model.name}</span>
          <button class="remove-model-btn" onclick="removeModel('${model.id}')" title="Remove">×</button>
        </div>
      `
        )
        .join("");
    }
  }

  // Remove model
  removeModel(modelId) {
    const index = this.selectedModels.findIndex((model) => model.id === modelId);
    if (index > -1) {
      this.selectedModels.splice(index, 1);
      this.updateSelectedModelsDisplay();
    }
  }
}

// Global function for removing models (needed for onclick handlers)
window.removeModel = function (modelId) {
  // This will be set up by the main script
};
