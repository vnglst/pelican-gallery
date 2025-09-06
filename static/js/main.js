// Main entry point for the application
import { DOM } from "./utils/dom.js";
import { Modal } from "./modules/modals.js";
import { ConfigManager } from "./modules/config.js";
import { UI } from "./modules/ui.js";
import { ArtworkManager } from "./modules/artwork.js";
import { ModelManager } from "./modules/models.js";

// Global state
let isGenerating = false;
let selectedModels = [];

// DOM Elements
const elements = {
  modelSelectorBtn: DOM.get("model-selector-btn"),
  modelSelectorText: DOM.get("model-selector-text"),
  selectedModelsList: DOM.get("selected-models-list"),
  temperatureInput: DOM.get("temperature-input"),
  temperatureValue: DOM.get("temperature-value"),
  maxTokensInput: DOM.get("max-tokens-input"),
  maxTokensValue: DOM.get("max-tokens-value"),
  promptInput: DOM.get("prompt-input"),
  titleInput: DOM.get("title-input"),
  categoryInput: DOM.get("category-input"),
  generateBtn: DOM.get("generate-btn"),
  configBtn: DOM.get("config-btn"),
  svgPreview: DOM.get("svg-preview"),
  errorMessage: DOM.get("error-message"),
  btnText: DOM.query("#generate-btn .btn-text"),
  btnLoading: DOM.query("#generate-btn .btn-loading"),

  // Modal buttons
  clearModelsBtn: DOM.get("clear-models-btn"),
  selectAllModelsBtn: DOM.get("select-all-models-btn"),
  saveModelsBtn: DOM.get("save-models-btn"),
  examplesBtn: DOM.get("examples-btn"),
  resetConfigBtn: DOM.get("reset-config-btn"),
  saveConfigBtn: DOM.get("save-config-btn"),

  // Reasoning controls
  reasoningEnabled: DOM.get("reasoning-enabled"),
  reasoningOptions: DOM.get("reasoning-options"),
  reasoningEffort: DOM.get("reasoning-effort"),
};

// Modal instances
const modals = {
  model: new Modal("model-modal", "close-model-btn"),
  examples: new Modal("examples-modal", "close-examples-btn"),
  config: new Modal("config-modal", "close-config-btn"),
};

// Initialize managers
const artworkManager = new ArtworkManager(elements, modals);
const modelManager = new ModelManager(elements, modals, selectedModels);
const configManager = new ConfigManager(DOM, UI, modals);

// Initialize the app
document.addEventListener("DOMContentLoaded", function () {
  setupEventListeners();
  modelManager.initializeSelectedModels();
  updateGenerateButton();
  configManager.save();
  artworkManager.displayExistingArtworksOnLoad();
});

// Setup event listeners
function setupEventListeners() {
  // Main UI event listeners
  elements.modelSelectorBtn.addEventListener("click", () => modals.model.open());
  elements.promptInput.addEventListener("input", updateGenerateButton);
  elements.generateBtn.addEventListener("click", handleGenerate);
  elements.configBtn.addEventListener("click", () => modals.config.open());

  // Setup slider listeners
  UI.updateSliders();

  // Model modal event listeners
  elements.clearModelsBtn.addEventListener("click", () => modelManager.clearAllModels());
  elements.selectAllModelsBtn.addEventListener("click", () => modelManager.selectAllModels());
  elements.saveModelsBtn.addEventListener("click", () => modelManager.saveModelSelection());

  // Other modal event listeners
  elements.examplesBtn.addEventListener("click", () => modals.examples.open());
  elements.resetConfigBtn.addEventListener("click", () => configManager.reset());
  elements.saveConfigBtn.addEventListener("click", () => configManager.saveToServer());

  // Model card click handlers
  document.addEventListener("click", (e) => modelManager.handleModelCardClick(e));

  // Example card click handlers
  document.addEventListener("click", (e) => artworkManager.handleExampleCardClick(e));

  // Keyboard shortcuts
  elements.promptInput.addEventListener("keydown", handleKeyboardShortcuts);
  document.addEventListener("keydown", handleGlobalKeyboardShortcuts);

  // Reasoning controls event listeners
  elements.reasoningEnabled.addEventListener("change", function (e) {
    DOM.toggleElement(elements.reasoningOptions, e.target.checked);
  });

  // Auto-resize textarea
  elements.promptInput.addEventListener("input", function () {
    UI.autoResizeTextarea(this);
    artworkManager.generateTitleFromPrompt();
  });

  // Initialize textarea height
  UI.autoResizeTextarea(elements.promptInput);
}

// Update generate button state
function updateGenerateButton() {
  const hasModels = selectedModels.length > 0;
  const hasPrompt = elements.promptInput.value.trim() !== "";
  elements.generateBtn.disabled = !hasModels || !hasPrompt || isGenerating;
}

// Handle SVG generation
async function handleGenerate() {
  if (isGenerating) return;

  const prompt = elements.promptInput.value.trim();

  if (selectedModels.length === 0 || !prompt) {
    UI.showError("Please select at least one model and enter a prompt.");
    return;
  }

  setGenerating(true);
  UI.hideError();

  try {
    await artworkManager.handleGenerate(prompt, selectedModels);
  } finally {
    setGenerating(false);
  }
}

// Set generating state
function setGenerating(generating) {
  isGenerating = generating;

  if (generating) {
    DOM.hideElement(elements.btnText);
    DOM.showElement(elements.btnLoading);
    elements.generateBtn.disabled = true;
  } else {
    DOM.showElement(elements.btnText);
    DOM.hideElement(elements.btnLoading);
    updateGenerateButton();
  }
}

// Keyboard shortcuts
function handleKeyboardShortcuts(e) {
  if (e.ctrlKey && e.key === "Enter") {
    e.preventDefault();
    if (!elements.generateBtn.disabled) {
      handleGenerate();
    }
  }
}

function handleGlobalKeyboardShortcuts(e) {
  if (e.key === "F1") {
    e.preventDefault();
    alert("Keyboard Shortcuts:\n\nCtrl + Enter: Generate SVG\nF1: Show this help");
  }
}

// Export for global access if needed
window.selectedModels = selectedModels;
window.isGenerating = isGenerating;

// Set up global model removal function
window.removeModel = function (modelId) {
  modelManager.removeModel(modelId);
  updateGenerateButton();
};

// Event delegation for regenerate buttons
document.addEventListener("click", function (e) {
  if (e.target.closest(".regenerate-btn")) {
    e.preventDefault();

    const btn = e.target.closest(".regenerate-btn");
    const svgItem = btn.closest(".svg-item");
    const slug = svgItem.dataset.slug;
    const model = svgItem.dataset.model;

    if (slug && model) {
      artworkManager.regenerateArtwork(slug, model, btn, svgItem);
    } else {
      console.error("Missing slug or model data attributes");
      UI.showError("Unable to regenerate: missing artwork information");
    }
  }
});

// Event delegation for delete buttons
document.addEventListener("click", function (e) {
  if (e.target.closest(".delete-btn")) {
    e.preventDefault();

    const btn = e.target.closest(".delete-btn");
    const svgItem = btn.closest(".svg-item");
    const artworkId = svgItem.dataset.id;

    if (artworkId) {
      // Confirm deletion
      if (confirm("Are you sure you want to delete this artwork? This action cannot be undone.")) {
        artworkManager.deleteArtwork(artworkId, btn, svgItem);
      }
    } else {
      console.error("Missing artwork ID data attribute");
      UI.showError("Unable to delete: missing artwork information");
    }
  }
});

// Add spin animation for loading spinner
const style = document.createElement("style");
style.textContent = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  .spin {
    animation: spin 1s linear infinite;
  }
`;
document.head.appendChild(style);
