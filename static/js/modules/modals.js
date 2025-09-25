import { html, useState, useEffect } from "https://esm.sh/htm/preact/standalone";
import { Modal } from "/static/js/modules/components.js";

export const ModelModal = ({ isOpen, onClose, onSelect, models, loading, error }) => {
  const [filterText, setFilterText] = useState("");

  const filteredModels = models.filter((model) => model.id.toLowerCase().includes(filterText.toLowerCase()));

  return html`
    <${Modal}
      isOpen=${isOpen}
      onClose=${onClose}
      title=${html`
        <div class="flex items-center gap-3">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 12l2 2 4-4M21 12c0 4.97-4.03 9-9 9s-9-4.03-9-9 4.03-9 9-9 9 4.03 9 9z" />
          </svg>
          Select AI Model
        </div>
      `}
      size="md"
    >
      ${error && html`<div class="mb-4 p-4 bg-fg text-bg border border-fg text-sm" role="alert">${error}</div>`}

      <div class="mb-4">
        <input
          type="text"
          placeholder="Filter models (e.g., 'openai', 'google')..."
          class="w-full px-3 py-2 border border-border bg-bg text-fg placeholder-fg/50 focus:outline-none focus:border-fg transition-colors duration-200"
          value=${filterText}
          onInput=${(e) => setFilterText(e.target.value)}
          aria-label="Filter models by ID"
        />
      </div>

      <div class="space-y-3" role="listbox" aria-label="Available AI models">
        ${
          loading
            ? html`<div class="text-center py-8 text-sm" role="status" aria-live="polite">Loading models...</div>`
            : filteredModels.length === 0
            ? html`<div class="text-center py-8 text-sm text-fg/60">No models match your filter.</div>`
            : filteredModels.map(
                (model) => html`
                  <div
                    key=${model.id}
                    class="model-card border border-border p-4 hover:bg-fg hover:text-bg transition-colors duration-200 cursor-pointer focus:outline-none focus:bg-fg focus:text-bg flex items-center justify-between"
                    role="option"
                    tabindex="0"
                    aria-label="Select ${model.name}"
                    onClick=${(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      onSelect(model);
                    }}
                    onKeyDown=${(e) => {
                      if (e.key === "Enter" || e.key === " ") {
                        e.preventDefault();
                        e.stopPropagation();
                        onSelect(model);
                      }
                    }}
                  >
                    <div class="font-semibold">${model.name}</div>
                    <div class="text-sm opacity-75">$${model.cost.toFixed(2)}/1M tokens</div>
                  </div>
                `
              )
        }
      </div>
    </${Modal}>
  `;
};

const DEFAULT_CONFIG = { temperature: 0.7, max_tokens: 50000 };

export const ConfigModal = ({ isOpen, onClose, onSave, artwork }) => {
  const [config, setConfig] = useState(DEFAULT_CONFIG);

  useEffect(() => {
    if (artwork) {
      try {
        setConfig({
          temperature: artwork.temperature,
          max_tokens: artwork.max_tokens,
        });
      } catch (e) {
        console.error("Failed to read artwork params:", e);
      }
    }
  }, [artwork]);

  const handleSave = () => {
    onSave(config);
    onClose();
  };

  const resetDefaults = () => {
    setConfig(DEFAULT_CONFIG);
  };

  return html`
    <${Modal}
      isOpen=${isOpen}
      onClose=${onClose}
      title="Generation Settings"
      size="lg"
    >
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div class="space-y-6">
          <h3 class="text-lg font-semibold">Model Parameters</h3>

          <div class="space-y-2">
            <label for="temperature-input" class="block text-sm font-medium">
              Temperature
              <span class="block text-xs text-fg/70 font-normal"
                >Controls creativity (0 = focused, 2 = creative)</span
              >
            </label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                id="temperature-input"
                class="flex-1 h-2 bg-border appearance-none cursor-pointer accent-fg"
                min="0"
                max="2"
                step="0.1"
                value=${config.temperature}
                onInput=${(e) => setConfig({ ...config, temperature: parseFloat(e.target.value) })}
              />
              <span class="text-sm font-mono min-w-8 text-right">${config.temperature}</span>
            </div>
          </div>

          <div class="space-y-2">
            <label for="max-tokens-input" class="block text-sm font-medium">Max Tokens</label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                id="max-tokens-input"
                class="flex-1 h-2 bg-border appearance-none cursor-pointer accent-fg"
                min="100"
                max="1000000"
                step="100"
                value=${config.max_tokens}
                onInput=${(e) => setConfig({ ...config, max_tokens: parseInt(e.target.value) })}
              />
              <span class="text-sm font-mono min-w-16 text-right">${config.max_tokens}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between pt-6 mt-8 border-t border-border">
        <button
          class="px-4 py-2 border border-border hover:bg-fg hover:text-bg transition-colors duration-200 text-sm font-medium"
          onClick=${resetDefaults}
        >
          Reset to Defaults
        </button>
        <button
          class="px-6 py-2 bg-fg text-bg hover:bg-opacity-80 transition-colors duration-200 text-sm font-medium"
          onClick=${handleSave}
        >
          Save Settings
        </button>
      </div>
    </${Modal}>
  `;
};
