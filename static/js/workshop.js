import { html, render, useReducer, useEffect } from "https://esm.sh/htm/preact/standalone";
import api from "/static/js/modules/api.js";
import { ToastContainer, LoadingOverlay, ArtworkCard } from "/static/js/modules/components.js";
import { ModelModal, ConfigModal } from "/static/js/modules/modals.js";
import { createInitialState, reducer } from "/static/js/modules/state.js";

const WorkshopApp = () => {
  const [state, dispatch] = useReducer(reducer, createInitialState(window));

  const showToast = (message, type = "info") => dispatch({ type: "PUSH_TOAST", payload: { message, type } });
  const removeToast = (index) => dispatch({ type: "REMOVE_TOAST", payload: index });
  const showLoading = (message = "Loading...") => dispatch({ type: "SHOW_LOADING", payload: message });
  const hideLoading = () => dispatch({ type: "HIDE_LOADING" });

  // Load group data for editing
  const loadGroupForEdit = async (groupId) => {
    try {
      showLoading("Loading group...");
      const data = await api.getGroup(groupId);
      const group = data.group;
      const artworks = data.artworks || [];

      dispatch({ type: "SET_CURRENT_GROUP", payload: group });

      const artworkMap = new Map();

      artworks.forEach((artwork) => {
        const id = artwork.id;
        if (id != null) {
          artworkMap.set(Number(id), artwork);
        }
      });
      dispatch({ type: "SET_ARTWORKS", payload: artworkMap });

      showToast("Group loaded for editing", "success");
    } catch (error) {
      console.error("Failed to load group:", error);
      showToast(`Failed to load group: ${error.message}`, "error");
    } finally {
      hideLoading();
    }
  };

  // Initialize from window data and URL parameters
  useEffect(() => {
    const initializeWorkshop = async () => {
      const groupData = window.currentGroup;

      let parsed = null;
      if (groupData) {
        if (typeof groupData === "string") {
          try {
            parsed = JSON.parse(groupData);
          } catch (e) {
            console.error("Failed to parse window.currentGroup:", e);
            parsed = null;
          }
        }
      }

      // Check for edit parameter in URL
      const urlParams = new URLSearchParams(window.location.search);
      const editId = urlParams.get("edit");

      if (editId && !parsed) {
        // Load group data from API - this will set formData internally
        await loadGroupForEdit(editId);
      } else if (parsed) {
        // Initialize from existing group
        dispatch({ type: "SET_CURRENT_GROUP", payload: parsed });
      }
      // If neither editId nor group, keep default form data
      const raw = window.existingArtworks;
      const existingArtworks = JSON.parse(raw);

      const artworkMap = new Map();
      existingArtworks.forEach((artwork) => {
        const id = artwork.id;
        if (id != null) {
          artworkMap.set(Number(id), artwork);
        }
      });

      if (artworkMap.size > 0) {
        dispatch({ type: "SET_ARTWORKS", payload: artworkMap });
      }
    };

    initializeWorkshop();
  }, []);

  // Update form data when currentGroup changes
  useEffect(() => {
    if (state.currentGroup) {
      const groupFormData = {
        title: state.currentGroup.title,
        prompt: state.currentGroup.prompt,
        category: state.currentGroup.category,
        original_url: state.currentGroup.original_url || "",
        artist_name: state.currentGroup.artist_name,
      };
      dispatch({ type: "SET_FORM_DATA", payload: groupFormData });
    }
  }, [state.currentGroup]);

  // Load models when modal opens
  const loadModels = async () => {
    dispatch({ type: "SET_MODELS_LOADING", payload: true });
    dispatch({ type: "SET_MODELS_ERROR", payload: "" });
    try {
      const data = await api.getModels();
      dispatch({ type: "SET_MODELS", payload: data.models || [] });
    } catch (error) {
      console.error("Failed to load models:", error);
      dispatch({ type: "SET_MODELS_ERROR", payload: `Failed to load models: ${error.message}` });
    } finally {
      dispatch({ type: "SET_MODELS_LOADING", payload: false });
    }
  };

  // API functions
  const saveGroup = async () => {
    const { title, prompt, category, original_url, artist_name } = state.formData;

    if (!title?.trim() || !prompt?.trim() || !category?.trim()) {
      showToast("Title, prompt and category are required", "error");
      return;
    }

    try {
      showLoading("Saving group...");
      const payload = {
        title: title.trim(),
        prompt: prompt.trim(),
        category: category.trim(),
        original_url: original_url?.trim(),
        artist_name: artist_name?.trim(),
      };
      const groupId = state.currentGroup?.id;
      const group = await (groupId ? api.updateGroup(groupId, payload) : api.createGroup(payload));
      dispatch({ type: "SET_CURRENT_GROUP", payload: group });
      window.currentGroup = group;
      showToast("Group saved", "success");
    } catch (error) {
      console.error("Save group error:", error);
      showToast(`Failed to save group: ${error.message}`, "error");
    } finally {
      hideLoading();
    }
  };

  const deleteGroup = async () => {
    const groupId = state.currentGroup?.id;

    if (!groupId) {
      showToast("No group to delete", "error");
      return;
    }

    const groupTitle = state.currentGroup?.title;
    if (!confirm('Delete group "' + groupTitle + '" and all artworks?')) {
      return;
    }

    try {
      showLoading("Deleting group...");
      await api.deleteGroup(groupId);
      showToast("Group deleted", "success");
      setTimeout(() => {
        window.location.href = "/gallery";
      }, 500);
    } catch (error) {
      showToast("Failed to delete group: " + error.message, "error");
    } finally {
      hideLoading();
    }
  };

  const addModel = async (model) => {
    let groupId = state.currentGroup?.id;

    // If no group exists, save the group first
    if (!groupId) {
      const { title, prompt, category, original_url, artist_name } = state.formData;

      if (!title?.trim() || !prompt?.trim()) {
        showToast("Please enter a title and prompt before adding models", "error");
        return;
      }

      try {
        showLoading("Saving group...");
        const groupPayload = {
          title: title.trim(),
          prompt: prompt.trim(),
          category: category.trim(),
          original_url: original_url?.trim() || "",
          artist_name: artist_name?.trim() || "",
        };

        const newGroup = await api.createGroup(groupPayload);
        dispatch({ type: "SET_CURRENT_GROUP", payload: newGroup });
        window.currentGroup = newGroup;
        groupId = newGroup.id;
        showToast("Group saved", "success");
      } catch (error) {
        console.error("Save group error:", error);
        const errorMessage =
          error.name === "TypeError" && error.message && error.message.indexOf("fetch") !== -1
            ? "Network error: Unable to connect to server"
            : "Failed to save group: " + error.message;
        showToast(errorMessage, "error");
        return;
      } finally {
        hideLoading();
      }
    }

    try {
      showLoading("Adding model...");
      const payload = {
        group_id: groupId,
        model: model.id,
        temperature: 0.7,
        max_tokens: 50000,
      };

      const artwork = await api.createArtwork(payload);
      dispatch({ type: "ADD_ARTWORK", payload: artwork });
      showToast("Added " + model.name, "success");
    } catch (error) {
      console.error("Add model error:", error);
      showToast("Failed to add model: " + error.message, "error");
    } finally {
      hideLoading();
    }
  };

  const generateArtwork = async (artworkId) => {
    const numericId = Number(artworkId);
    const artwork = state.artworks.get(numericId);
    if (!artwork) {
      showToast("Artwork " + numericId + " not found", "error");
      return;
    }

    dispatch({ type: "ADD_GENERATING", payload: numericId });

    try {
      const result = await api.generateArtwork(numericId);
      const updatedArtwork = { ...artwork, svg: result.svg };
      dispatch({ type: "UPDATE_ARTWORK", payload: updatedArtwork });
      showToast("Artwork generated", "success");
    } catch (error) {
      console.error("Generate artwork error:", error);
      showToast("Generation failed: " + error.message, "error");
    } finally {
      dispatch({ type: "REMOVE_GENERATING", payload: numericId });
    }
  };

  const updateArtworkParams = async (artworkId, params) => {
    try {
      const updated = await api.updateArtwork(artworkId, {
        temperature: params.temperature,
        max_tokens: params.max_tokens,
      });
      dispatch({ type: "UPDATE_ARTWORK", payload: updated });
      showToast("Parameters updated", "success");
    } catch (error) {
      showToast("Failed to update parameters: " + error.message, "error");
    }
  };

  const removeArtwork = async (artworkId) => {
    if (!confirm("Remove this model?")) return;

    try {
      await api.deleteArtwork(artworkId);
      dispatch({ type: "REMOVE_ARTWORK", payload: artworkId });
      showToast("Model removed", "success");
    } catch (error) {
      console.error("Remove artwork error:", error);
      var errorMessage = "Failed to remove model";

      if (error.name === "TypeError" && error.message && error.message.indexOf("fetch") !== -1) {
        errorMessage = "Network error: Unable to connect to server";
      } else if (error.message) {
        errorMessage = "Failed to remove model: " + error.message;
      }

      showToast(errorMessage, "error");
    }
  };

  // Event handlers
  const handleAddModel = () => {
    loadModels();
    dispatch({ type: "SET_MODAL", payload: { modal: "model", value: true } });
  };

  const handleModelSelect = (model) => {
    dispatch({ type: "SET_MODAL", payload: { modal: "model", value: false } });
    addModel(model);
  };

  const handleConfigure = (artwork) => {
    dispatch({ type: "SET_CONFIG_ARTWORK", payload: artwork });
    dispatch({ type: "SET_MODAL", payload: { modal: "config", value: true } });
  };

  const handleConfigSave = (params) => {
    if (state.configArtwork) {
      updateArtworkParams(state.configArtwork.id, params);
    }
  };

  // Form values
  const isEditing = !!state.currentGroup?.id;

  return html`
    <div>
      <!-- Group Form -->
      <div class="lg:grid lg:grid-cols-3 lg:gap-8 space-y-8 lg:space-y-0">
        <div class="space-y-6">
          <div class="space-y-6">
            <div class="space-y-2">
              <label for="prompt-input" class="block text-sm font-medium">Describe your artwork</label>
              <textarea
                id="prompt-input"
                class="w-full p-3 border border-border bg-bg text-fg text-sm focus:outline-none focus:border-fg resize-none"
                placeholder="A serene mountain landscape with geometric patterns, flowing rivers, and abstract shapes in harmonious colors..."
                rows="6"
                value=${state.formData.prompt}
                onInput=${(e) =>
                  dispatch({ type: "SET_FORM_DATA", payload: { ...state.formData, prompt: e.target.value } })}
              ></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-2">
                <label for="title-input" class="block text-sm font-medium">Title</label>
                <input
                  type="text"
                  id="title-input"
                  class="w-full p-3 border border-border bg-bg text-fg text-sm focus:outline-none focus:border-fg"
                  placeholder="Mountain Landscape"
                  title="Enter a descriptive title for your artwork group"
                  value=${state.formData.title}
                  onInput=${(e) =>
                    dispatch({ type: "SET_FORM_DATA", payload: { ...state.formData, title: e.target.value } })}
                />
              </div>

              <div class="space-y-2">
                <label for="category-input" class="block text-sm font-medium">Category</label>
                <input
                  type="text"
                  id="category-input"
                  class="w-full p-3 border border-border bg-bg text-fg text-sm focus:outline-none focus:border-fg"
                  placeholder="abstract, nature, geometric, etc."
                  value=${state.formData.category}
                  onInput=${(e) =>
                    dispatch({ type: "SET_FORM_DATA", payload: { ...state.formData, category: e.target.value } })}
                />
              </div>
            </div>

            <div class="space-y-2">
              <label for="original-url-input" class="block text-sm font-medium">Original Artwork URL</label>
              <input
                type="url"
                id="original-url-input"
                class="w-full p-3 border border-border bg-bg text-fg text-sm focus:outline-none focus:border-fg"
                placeholder="https://example.com/original-artwork.jpg"
                value=${state.formData.original_url || ""}
                onInput=${(e) =>
                  dispatch({ type: "SET_FORM_DATA", payload: { ...state.formData, original_url: e.target.value } })}
              />
            </div>

            <div class="space-y-2">
              <label for="artist-name-input" class="block text-sm font-medium">Artist Name</label>
              <input
                type="text"
                id="artist-name-input"
                class="w-full p-3 border border-border bg-bg text-fg text-sm focus:outline-none focus:border-fg"
                placeholder="Jane Doe"
                value=${state.formData.artist_name || ""}
                onInput=${(e) =>
                  dispatch({ type: "SET_FORM_DATA", payload: { ...state.formData, artist_name: e.target.value } })}
              />
            </div>

            <div class="flex items-center gap-3">
              <button
                class="px-6 py-2 bg-fg text-bg hover:bg-opacity-80 transition-colors duration-200 text-sm font-medium"
                onClick=${saveGroup}
              >
                ${isEditing ? "Update Group" : "Save Group"}
              </button>
              ${isEditing &&
              html`
                <button
                  class="px-4 py-2 border border-border hover:bg-fg hover:text-bg transition-colors duration-200 text-sm font-medium"
                  onClick=${deleteGroup}
                >
                  Delete Group
                </button>
              `}
            </div>
          </div>
        </div>

        <div class="lg:col-span-2 space-y-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
            ${Array.from(state.artworks.entries()).map(
              ([id, artwork]) =>
                html`
                  <${ArtworkCard}
                    key=${id}
                    artwork=${artwork}
                    onRegenerate=${generateArtwork}
                    onConfigure=${handleConfigure}
                    onRemove=${removeArtwork}
                    isGenerating=${state.generatingArtworks.has(Number(id))}
                  />
                `
            )}
            <!-- Add Model Card -->
            <div class="border border-border p-8 text-center space-y-6">
              <h3 class="font-semibold mt-auto">Add AI Model</h3>
              <p class="text-sm text-fg/70">Select AI models to generate artwork variations</p>
              <button
                class="w-full px-4 py-2 border border-border hover:bg-fg hover:text-bg transition-colors duration-200 text-sm font-medium flex items-center justify-center gap-2"
                onClick=${handleAddModel}
              >
                <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 5v14M5 12h14" />
                </svg>
                Add Model
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modals -->
      <${ModelModal}
        isOpen=${state.modals.model}
        onClose=${() => dispatch({ type: "SET_MODAL", payload: { modal: "model", value: false } })}
        onSelect=${handleModelSelect}
        models=${state.models}
        loading=${state.modelsLoading}
        error=${state.modelsError}
      />

      <${ConfigModal}
        isOpen=${state.modals.config}
        onClose=${() => dispatch({ type: "SET_MODAL", payload: { modal: "config", value: false } })}
        onSave=${handleConfigSave}
        artwork=${state.configArtwork}
      />

      <!-- Loading overlay -->
      <${LoadingOverlay} message=${state.loading.message} visible=${state.loading.visible} />

      <!-- Toast container -->
      <${ToastContainer} toasts=${state.toasts} removeToast=${removeToast} />
    </div>
  `;
};

document.addEventListener("DOMContentLoaded", () => {
  try {
    const container = document.querySelector("main");
    if (container) {
      render(html`<${WorkshopApp} />`, container);
    } else {
      console.error("Workshop container not found");
    }
  } catch (error) {
    console.error("Failed to initialize Workshop app:", error);
    // Fallback: show error message
    const container = document.querySelector("main");
    if (container) {
      container.innerHTML = `
          <div class="error-message" style="display: block; text-align: center; padding: 2rem;">
            <h3>Workshop Initialization Error</h3>
            <p>Failed to load the workshop interface. Please refresh the page.</p>
            <p><small>Error: ${error.message}</small></p>
          </div>
        `;
    }
  }
});
