// Artwork management module
import { DOM } from "../utils/dom.js";
import { UI } from "./ui.js";

export class ArtworkManager {
  constructor(elements, modals) {
    this.elements = elements;
    this.modals = modals;
  }

  // Generate a title from the current prompt
  generateTitleFromPrompt() {
    const prompt = this.elements.promptInput.value.trim();
    if (!prompt || !this.elements.titleInput) return;

    // Take first few words and convert to title case
    const title = prompt
      .split(/\s+/)
      .slice(0, 6) // Take first 6 words for title
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(" ");

    this.elements.titleInput.value = title;
  }

  // Handle saving artwork
  async handleSaveArtwork(svgContent, modelId, modelName, saveButton) {
    let title = this.elements.titleInput ? this.elements.titleInput.value.trim() : "";
    const category = this.elements.categoryInput ? this.elements.categoryInput.value.trim() : "";
    const prompt = this.elements.promptInput ? this.elements.promptInput.value.trim() : "";

    // Provide default title if empty
    if (!title) {
      title = "Untitled Artwork";
      if (this.elements.titleInput) {
        this.elements.titleInput.value = title;
      }
    }

    if (!category) {
      UI.showError("Please enter a category before saving");
      return;
    }

    // Disable save button and show loading state
    const originalContent = saveButton.innerHTML;
    saveButton.disabled = true;
    saveButton.innerHTML = `
      <svg class="icon spinning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 12a9 9 0 11-6.219-8.56"/>
      </svg>
      Saving...
    `;

    try {
      const requestBody = {
        title: title,
        category: category,
        prompt: prompt,
        model: modelId,
        svg_content: svgContent,
        temperature: this.elements.temperatureInput ? parseFloat(this.elements.temperatureInput.value) : 0.7,
        max_tokens: this.elements.maxTokensInput ? parseInt(this.elements.maxTokensInput.value) : 1000,
      };

      const response = await fetch("/api/save-artwork", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      const data = await response.json();

      if (!response.ok || data.error) {
        throw new Error(data.error || "Failed to save artwork");
      }

      // Success - replace save button with regenerate and delete buttons
      const saveContainer = saveButton.parentElement;
      const svgItem = saveContainer.closest(".svg-item");

      // Create action buttons container
      const actionButtons = document.createElement("div");
      actionButtons.className = "svg-action-buttons";

      // Create regenerate button
      const regenerateBtn = document.createElement("button");
      regenerateBtn.className = "regenerate-btn";
      regenerateBtn.title = "Regenerate this artwork";
      regenerateBtn.innerHTML = `
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
          <path d="M21 3v5h-5"/>
          <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
          <path d="M3 21v-5h5"/>
        </svg>
      `;

      // Create delete button
      const deleteBtn = document.createElement("button");
      deleteBtn.className = "delete-btn";
      deleteBtn.title = "Delete this artwork";
      deleteBtn.innerHTML = `
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 6h18"/>
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
        </svg>
      `;

      // Add data attributes for functionality
      if (svgItem) {
        svgItem.setAttribute("data-id", data.id || "");
        svgItem.setAttribute("data-slug", data.slug || "");
        svgItem.setAttribute("data-model", modelId);
      }

      actionButtons.appendChild(regenerateBtn);
      actionButtons.appendChild(deleteBtn);

      // Replace save container with action buttons
      saveContainer.parentElement.replaceChild(actionButtons, saveContainer);

      UI.showSuccessMessage(`Artwork "${title}" saved successfully!`);
    } catch (error) {
      console.error("Error saving artwork:", error);
      UI.showError(`Failed to save artwork: ${error.message}`);

      // Reset button on error
      saveButton.disabled = false;
      saveButton.innerHTML = originalContent;
    }
  }

  // Handle saving all artworks at once
  async handleSaveAllArtworks() {
    let title = this.elements.titleInput ? this.elements.titleInput.value.trim() : "";
    const category = this.elements.categoryInput ? this.elements.categoryInput.value.trim() : "";
    const prompt = this.elements.promptInput ? this.elements.promptInput.value.trim() : "";

    // Provide default title if empty
    if (!title) {
      title = "Untitled Artwork";
      if (this.elements.titleInput) {
        this.elements.titleInput.value = title;
      }
    }

    if (!category) {
      UI.showError("Please enter a category before saving");
      return;
    }

    // Find all generated SVGs
    const svgContainers = document.querySelectorAll(".svg-container");
    if (svgContainers.length === 0) {
      UI.showError("No artworks to save");
      return;
    }

    // Collect all artwork data
    const artworks = [];
    svgContainers.forEach((container) => {
      const svgContent = container.innerHTML;
      const svgItem = container.closest(".svg-item");
      const modelName = svgItem ? svgItem.querySelector(".svg-item-header").textContent : "Unknown Model";
      const modelId = svgItem ? svgItem.getAttribute("data-model") : "unknown";

      artworks.push({
        svgContent,
        modelId,
        modelName,
      });
    });

    // Get Save All button
    const saveAllButton = document.querySelector(".save-all-btn");
    const originalContent = saveAllButton.innerHTML;

    // Disable button and show loading state
    saveAllButton.disabled = true;
    saveAllButton.innerHTML = `
      <svg class="icon spinning" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 12a9 9 0 11-6.219-8.56"/>
      </svg>
      Saving All...
    `;

    try {
      // Save all artworks
      for (let i = 0; i < artworks.length; i++) {
        const artwork = artworks[i];
        const response = await fetch("/api/save-artwork", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            title: title,
            category: category,
            prompt: prompt,
            model: artwork.modelId,
            svg_content: artwork.svgContent,
            temperature: this.elements.temperatureInput ? parseFloat(this.elements.temperatureInput.value) : 0.7,
            max_tokens: this.elements.maxTokensInput ? parseInt(this.elements.maxTokensInput.value) : 1000,
          }),
        });

        const data = await response.json();
        if (!response.ok || data.error) {
          throw new Error(data.error || `Failed to save artwork for ${artwork.modelName}`);
        }
      }

      // Success - update button to show saved state
      saveAllButton.innerHTML = `
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 6L9 17l-5-5"/>
        </svg>
        All Saved!
      `;
      saveAllButton.classList.add("saved");

      UI.showSuccessMessage(`All ${artworks.length} artworks saved successfully!`);

      // Reset button after delay
      setTimeout(() => {
        saveAllButton.disabled = false;
        saveAllButton.innerHTML = originalContent;
        saveAllButton.classList.remove("saved");
      }, 3000);
    } catch (error) {
      console.error("Error saving all artworks:", error);
      UI.showError(`Failed to save all artworks: ${error.message}`);

      // Reset button on error
      saveAllButton.disabled = false;
      saveAllButton.innerHTML = originalContent;
    }
  }

  // Check if all artworks are generated and show Save All button
  checkAndShowSaveAllButton() {
    const svgItems = document.querySelectorAll(".svg-item");
    const svgContainers = document.querySelectorAll(".svg-container");
    const saveAllContainer = document.querySelector(".save-all-container");

    // Check if we have at least one successful generation and no loading states
    const hasLoadingStates = document.querySelectorAll(".svg-item-loading").length > 0;

    if (svgContainers.length > 0 && !hasLoadingStates && saveAllContainer) {
      saveAllContainer.style.display = "block";
    }
  }

  // Display existing artworks when in edit mode
  displayExistingArtworksOnLoad() {
    // Wait a bit to ensure all scripts are loaded
    setTimeout(() => {
      let artworks = window.existingArtworks;

      // If it's a string, try to parse it
      if (typeof artworks === "string") {
        try {
          artworks = JSON.parse(artworks);
          // Update the global variable so other functions can use it
          window.existingArtworks = artworks;
        } catch (e) {
          console.error("Failed to parse existing artworks:", e);
          return;
        }
      }

      if (artworks && Array.isArray(artworks) && artworks.length > 0) {
        // Clear the placeholder
        this.elements.svgPreview.innerHTML = "";
        this.elements.svgPreview.classList.add("has-content");

        // Create gallery container
        const gallery = document.createElement("div");
        gallery.className = "svg-gallery";
        this.elements.svgPreview.appendChild(gallery);

        // Display each existing artwork
        artworks.forEach((artwork) => {
          const svgItem = document.createElement("div");
          svgItem.className = "svg-item";
          svgItem.setAttribute("data-id", artwork.ID);
          svgItem.setAttribute("data-slug", artwork.Slug);
          svgItem.setAttribute("data-model", artwork.Model);

          svgItem.innerHTML = `
          <div class="svg-item-header">
            <span>${artwork.Model}</span>
            <div class="svg-action-buttons">
              <button class="regenerate-btn" title="Regenerate this artwork">
                <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
                  <path d="M21 3v5h-5"/>
                  <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
                  <path d="M3 21v-5h5"/>
                </svg>
              </button>
              <button class="delete-btn" title="Delete this artwork">
                <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 6h18"/>
                  <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                  <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="svg-item-content">
            <div class="svg-container">
              ${artwork.SVGContent}
            </div>
          </div>
        `;

          gallery.appendChild(svgItem);
        });

        // Add "Save All" button container if there are artworks
        const saveAllContainer = document.createElement("div");
        saveAllContainer.className = "save-all-container";
        saveAllContainer.style.display = "block"; // Show since we have artworks

        const saveAllButton = document.createElement("button");
        saveAllButton.className = "save-all-btn";
        saveAllButton.title = "Save all artworks";
        saveAllButton.innerHTML = `
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
          <polyline points="17,21 17,13 7,13 7,21"/>
          <polyline points="7,3 7,8 15,8"/>
        </svg>
        Save all
      `;

        saveAllButton.addEventListener("click", () => this.handleSaveAllArtworks());
        saveAllContainer.appendChild(saveAllButton);

        // Append to edit banner instead of svg preview
        const editBanner = document.querySelector(".edit-banner-content");
        if (editBanner) {
          editBanner.appendChild(saveAllContainer);
        } else {
          this.elements.svgPreview.appendChild(saveAllContainer);
        }
      }
    }, 100); // Wait 100ms for all scripts to load
  }

  // Handle SVG generation for multiple models
  async handleGenerate(prompt, selectedModels) {
    // Filter out models that already have existing artworks (when in edit mode)
    let modelsToGenerate = selectedModels;
    if (window.existingArtworks && window.existingArtworks.length > 0) {
      const existingModelIds = new Set(window.existingArtworks.map((artwork) => artwork.Model));

      modelsToGenerate = selectedModels.filter((model) => {
        const hasExisting = existingModelIds.has(model.id);
        return !hasExisting;
      });

      if (modelsToGenerate.length === 0) {
        UI.showError(
          "All selected models already have artworks for this collection. Please select different models or create a new collection."
        );
        return;
      }
    }

    // Setup the gallery layout
    this.setupSVGGallery(modelsToGenerate);

    // Generate SVGs for models that don't have existing artworks
    const promises = modelsToGenerate.map((model) => this.generateSVGForModel(model.id, model.name, prompt));

    try {
      await Promise.allSettled(promises);
    } catch (error) {
      console.error("Error in generation:", error);
    }
  }

  // Setup the SVG gallery layout
  setupSVGGallery(models) {
    this.elements.svgPreview.innerHTML = "";
    this.elements.svgPreview.classList.add("has-content");

    const gallery = document.createElement("div");
    gallery.className = "svg-gallery";
    this.elements.svgPreview.appendChild(gallery);

    // Create containers for each model
    models.forEach((model) => {
      const svgItem = document.createElement("div");
      svgItem.className = "svg-item";
      svgItem.id = `svg-item-${model.id.replace(/[^a-zA-Z0-9]/g, "_")}`;
      svgItem.setAttribute("data-model", model.id);

      svgItem.innerHTML = `
        <div class="svg-item-header">${model.name}</div>
        <div class="svg-item-content">
          <div class="svg-item-loading">Generating...</div>
        </div>
      `;

      gallery.appendChild(svgItem);
    });

    // Add "Save All" button container
    const saveAllContainer = document.createElement("div");
    saveAllContainer.className = "save-all-container";
    saveAllContainer.style.display = "none"; // Hidden until artworks are generated

    const saveAllButton = document.createElement("button");
    saveAllButton.className = "save-all-btn";
    saveAllButton.title = "Save all artworks";
    saveAllButton.innerHTML = `
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
        <polyline points="17,21 17,13 7,13 7,21"/>
        <polyline points="7,3 7,8 15,8"/>
    </svg>
    Save all
  `;

    saveAllButton.addEventListener("click", () => this.handleSaveAllArtworks());
    saveAllContainer.appendChild(saveAllButton);

    // Append to edit banner instead of svg preview
    const editBanner = document.querySelector(".edit-banner-content");
    if (editBanner) {
      editBanner.appendChild(saveAllContainer);
    } else {
      this.elements.svgPreview.appendChild(saveAllContainer);
    }
  }

  // Generate SVG for a specific model
  async generateSVGForModel(modelId, modelName, prompt) {
    const itemId = `svg-item-${modelId.replace(/[^a-zA-Z0-9]/g, "_")}`;
    const svgItem = DOM.get(itemId);
    const contentDiv = svgItem.querySelector(".svg-item-content");

    try {
      const titleInput = document.getElementById("title-input");
      const requestBody = {
        model: modelId,
        prompt: prompt,
        title: titleInput ? titleInput.value.trim() : "",
        category: this.elements.categoryInput.value,
        temperature: parseFloat(this.elements.temperatureInput.value),
        max_tokens: parseInt(this.elements.maxTokensInput.value),
      };

      // Add reasoning parameters if enabled
      if (this.elements.reasoningEnabled.checked) {
        requestBody.reasoning = {
          enabled: true,
          effort: this.elements.reasoningEffort.value,
          exclude: false,
        };
      }

      const response = await fetch("/api/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      const data = await response.json();

      if (!response.ok || data.error) {
        throw new Error(data.error || "Failed to generate SVG");
      }

      // Display the SVG
      this.displaySVGForModel(contentDiv, data.svg, modelId, modelName);

      // Check if all artworks are generated and show Save All button
      this.checkAndShowSaveAllButton();
    } catch (error) {
      console.error(`Error generating SVG for ${modelName}:`, error);
      this.displayErrorForModel(contentDiv, error.message);
    }
  }

  // Display SVG for a specific model
  displaySVGForModel(contentDiv, svgContent, modelId, modelName) {
    contentDiv.innerHTML = "";

    // Create main container
    const container = document.createElement("div");
    container.className = "svg-display-container";

    // Create SVG container
    const svgContainer = document.createElement("div");
    svgContainer.className = "svg-container";
    svgContainer.innerHTML = svgContent;

    // Create save button container
    const saveContainer = document.createElement("div");
    saveContainer.className = "svg-save-container";

    const saveButton = document.createElement("button");
    saveButton.className = "save-artwork-btn";
    saveButton.title = "Save Artwork";
    saveButton.innerHTML = `
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
        <polyline points="17,21 17,13 7,13 7,21"/>
        <polyline points="7,3 7,8 15,8"/>
      </svg>
    `;

    saveButton.addEventListener("click", () => {
      this.handleSaveArtwork(svgContent, modelId, modelName, saveButton);
    });

    saveContainer.appendChild(saveButton);
    container.appendChild(svgContainer);
    container.appendChild(saveContainer);
    contentDiv.appendChild(container);
  }

  // Display error for a specific model
  displayErrorForModel(contentDiv, errorMessage) {
    contentDiv.innerHTML = `<div class="svg-item-error">${errorMessage}</div>`;
  }

  // Handle example card click
  handleExampleCardClick(e) {
    const exampleCard = e.target.closest(".example-card");
    if (exampleCard) {
      const prompt = exampleCard.dataset.prompt;
      const title = exampleCard.dataset.title;
      const category = exampleCard.dataset.category;
      if (prompt) {
        this.selectExample(prompt, title, category);
      }
    }
  }

  // Examples functionality
  selectExample(prompt, title, category) {
    this.elements.promptInput.value = prompt;
    UI.autoResizeTextarea(this.elements.promptInput);

    // Fill in title if provided
    if (title) {
      const titleInput = document.getElementById("title-input");
      if (titleInput) {
        titleInput.value = title;
      }
    }

    // Fill in category if provided
    if (category && this.elements.categoryInput) {
      this.elements.categoryInput.value = category;
    }

    // Update generate button
    const hasModels = window.selectedModels && window.selectedModels.length > 0;
    const hasPrompt = prompt.trim() !== "";
    this.elements.generateBtn.disabled = !hasModels || !hasPrompt;

    this.modals.examples.close();

    // Focus on the prompt input for any edits
    this.elements.promptInput.focus();

    // Auto-generate title from prompt if no title was provided
    if (!title) {
      this.generateTitleFromPrompt();
    }

    // Optional: Show a brief success message
    UI.showSuccessMessage("Example loaded successfully!");
  }

  // Regenerate artwork functionality
  async regenerateArtwork(slug, model, btn, svgItem) {
    try {
      // Disable button and show loading state
      btn.disabled = true;
      btn.style.transform = "rotate(180deg)";
      btn.style.opacity = "0.7";

      // Add loading indicator to SVG container
      const svgContainer = svgItem.querySelector(".svg-container");
      const originalOpacity = svgContainer.style.opacity;
      svgContainer.style.opacity = "0.5";
      svgContainer.style.position = "relative";

      // Create loading overlay
      const loadingOverlay = document.createElement("div");
      loadingOverlay.style.cssText = `
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.7);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 10;
      `;

      const spinner = document.createElement("div");
      spinner.style.cssText = `
        width: 24px;
        height: 24px;
        border: 2px solid #374151;
        border-top: 2px solid #4ade80;
        border-radius: 50%;
        animation: spin 1s linear infinite;
      `;

      loadingOverlay.appendChild(spinner);
      svgContainer.appendChild(loadingOverlay);

      const response = await fetch("/api/regenerate-artwork", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          slug: slug,
          model: model,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to regenerate artwork");
      }

      const result = await response.json();

      // Update the SVG content
      const svgElement = svgContainer.querySelector("svg");
      if (svgElement) {
        // Create a temporary container to parse the new SVG
        const tempDiv = document.createElement("div");
        tempDiv.innerHTML = result.svgContent;
        const newSvg = tempDiv.querySelector("svg");

        if (newSvg) {
          svgElement.replaceWith(newSvg);
        }
      }

      // Show success message
      UI.showSuccessMessage("Artwork regenerated successfully!");
    } catch (error) {
      console.error("Error regenerating artwork:", error);
      UI.showError("Failed to regenerate artwork: " + error.message);
    } finally {
      // Re-enable button and restore styles
      btn.disabled = false;
      btn.style.transform = "";
      btn.style.opacity = "";

      // Remove loading overlay and restore container
      const svgContainer = svgItem.querySelector(".svg-container");
      const loadingOverlay = svgContainer.querySelector('div[style*="position: absolute"]');
      if (loadingOverlay) {
        loadingOverlay.remove();
      }
      svgContainer.style.opacity = "";
      svgContainer.style.position = "";
    }
  }

  // Delete artwork functionality
  async deleteArtwork(artworkId, btn, svgItem) {
    try {
      // Disable button and show loading state
      btn.disabled = true;
      btn.innerHTML =
        '<svg class="icon spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="m16 12-4-4-4 4"/><path d="m12 16 4-4"/></svg>';

      const response = await fetch("/api/delete-artwork", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: parseInt(artworkId),
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to delete artwork");
      }

      const data = await response.json();

      if (data.success) {
        // Remove the artwork from the UI with animation
        svgItem.style.transition = "opacity 0.3s ease, transform 0.3s ease";
        svgItem.style.opacity = "0";
        svgItem.style.transform = "scale(0.8)";

        setTimeout(() => {
          svgItem.remove();

          // Check if this was the last artwork in the gallery
          const remainingArtworks = document.querySelectorAll(".svg-item");
          if (remainingArtworks.length === 0) {
            // Reload the page to show the empty state or redirect
            window.location.reload();
          }
        }, 300);

        UI.showSuccessMessage("Artwork deleted successfully!");
      } else {
        throw new Error(data.message || "Unknown error");
      }
    } catch (error) {
      console.error("Error deleting artwork:", error);

      // Restore button state
      btn.disabled = false;
      btn.innerHTML =
        '<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>';

      UI.showError("Failed to delete artwork: " + error.message);
    }
  }
}
