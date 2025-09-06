// Configuration management module
export class ConfigManager {
  constructor(DOM, UI, modals) {
    this.DOM = DOM;
    this.UI = UI;
    this.modals = modals;
    this.default = null;
  }

  save() {
    this.default = {
      temperature: parseFloat(document.getElementById("temperature-input").value),
      maxTokens: parseInt(document.getElementById("max-tokens-input").value),
      reasoningEnabled: document.getElementById("reasoning-enabled").checked,
      reasoningEffort: document.getElementById("reasoning-effort").value,
    };
  }

  reset() {
    if (this.default) {
      const tempInput = document.getElementById("temperature-input");
      const maxTokensInput = document.getElementById("max-tokens-input");
      const reasoningEnabled = document.getElementById("reasoning-enabled");
      const reasoningEffort = document.getElementById("reasoning-effort");
      const reasoningOptions = document.getElementById("reasoning-options");

      tempInput.value = this.default.temperature;
      document.getElementById("temperature-value").textContent = this.default.temperature;
      maxTokensInput.value = this.default.maxTokens;
      document.getElementById("max-tokens-value").textContent = this.default.maxTokens;
      reasoningEnabled.checked = this.default.reasoningEnabled;
      reasoningEffort.value = this.default.reasoningEffort;

      this.DOM.toggleElement(reasoningOptions, this.default.reasoningEnabled);
      this.UI.showSuccessMessage("Configuration reset to defaults");
    }
  }

  saveToServer() {
    this.UI.showSuccessMessage("Settings saved successfully");
    this.modals.config.close();
  }
}
