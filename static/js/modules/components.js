import { html, useEffect, useRef } from "https://esm.sh/htm/preact/standalone";

// SVG Display component for safe HTML rendering
const SVGDisplay = ({ svgContent }) => {
  const containerRef = useRef();
  useEffect(() => {
    if (containerRef.current && svgContent) {
      containerRef.current.innerHTML = svgContent;
      const svgElements = containerRef.current.querySelectorAll("svg");
      svgElements.forEach((svg) => {
        svg.style.maxWidth = "100%";
        svg.style.maxHeight = "100%";
        svg.style.width = "auto";
        svg.style.height = "auto";
        svg.style.display = "block";
        svg.style.margin = "0 auto";
      });
    }
  }, [svgContent]);
  return html`<div
    ref=${containerRef}
    class="w-full h-full flex items-center justify-center overflow-hidden"
    style="max-width: 100%; max-height: 100%;"
  ></div>`;
};

// Toast system
const Toast = ({ message, type = "info", onClose }) => {
  useEffect(() => {
    const timer = setTimeout(onClose, 3000);
    return () => clearTimeout(timer);
  }, [onClose]);
  const icons = {
    success: "✓",
    error: "✕",
    warning: "⚠",
    info: "ℹ",
  };
  const typeStyles = {
    success: "bg-fg text-bg border-fg",
    error: "bg-fg text-bg border-fg",
    warning: "bg-fg text-bg border-fg",
    info: "bg-bg text-fg border-border",
  };
  return html`
    <div
      class="fixed right-4 top-4 z-50 min-w-80 max-w-md ${typeStyles[
        type
      ]} border shadow-lg animate-slide-in flex items-center gap-3 p-4"
    >
      <div class="flex-shrink-0 font-bold text-lg">${icons[type]}</div>
      <div class="flex-1 text-sm font-medium">${message}</div>
      <button
        class="flex-shrink-0 w-6 h-6 flex items-center justify-center hover:bg-opacity-80 transition-colors duration-200 text-lg font-bold leading-none"
        onClick=${onClose}
      >
        ×
      </button>
    </div>
  `;
};

// Toast container component
export const ToastContainer = ({ toasts, removeToast }) => {
  return html`
    <div id="toast-container" class="fixed top-4 right-4 z-50 space-y-2">
      ${toasts.map(
        (toast, index) =>
          html`<${Toast}
            key=${index}
            message=${toast.message}
            type=${toast.type}
            onClose=${() => removeToast(index)}
          />`
      )}
    </div>
  `;
};

// Loading overlay component
export const LoadingOverlay = ({ message, visible }) => {
  if (!visible) return null;
  return html`
    <div id="global-loading-overlay" class="fixed inset-0 z-[2000] bg-fg/50 flex items-center justify-center">
      <div class="bg-bg border border-border p-8 shadow-xl max-w-sm w-full mx-4 flex flex-col items-center gap-4">
        <svg class="w-8 h-8 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" />
          <path d="M12 6v6l4 2" />
        </svg>
        <div class="text-sm font-medium text-center">${message}</div>
      </div>
    </div>
  `;
};

// Artwork card component
export const ArtworkCard = ({ artwork, onRegenerate, onConfigure, onRemove, isGenerating }) => {
  const hasContent = artwork.svg !== "";

  return html`
    <div
      class="border border-border bg-bg ${isGenerating ? "opacity-80 cursor-loading" : ""}"
      data-artwork-id=${artwork.id}
    >
      <div class="flex items-center justify-between p-4 border-b border-border">
        <div class="flex-1 min-w-0">
          <h3 class="font-semibold text-sm truncate">${artwork.model}</h3>
        </div>
        <div class="flex items-center gap-1 ml-4">
          <button
            class="w-8 h-8 flex items-center justify-center hover:bg-fg hover:text-bg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            title="Regenerate"
            onClick=${() => onRegenerate(artwork.id)}
            disabled=${isGenerating}
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" />
              <path d="M21 3v5h-5" />
              <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" />
              <path d="M3 21v-5h5" />
            </svg>
          </button>
          <button
            class="w-8 h-8 flex items-center justify-center hover:bg-fg hover:text-bg transition-colors duration-200"
            title="Model Settings"
            onClick=${() => onConfigure(artwork)}
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3" />
              <path
                d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
              />
            </svg>
          </button>
          <button
            class="w-8 h-8 flex items-center justify-center hover:bg-fg hover:text-bg transition-colors duration-200"
            title="Remove Model"
            onClick=${() => onRemove(artwork.id)}
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
      </div>

      <div class="aspect-square relative bg-bg min-h-[200px]">
        ${hasContent
          ? html` <${SVGDisplay} svgContent=${artwork.svg} /> `
          : html`
              <div class="w-full h-full flex flex-col items-center justify-center text-fg/50 gap-3">
                <svg class="w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                  <circle cx="8.5" cy="8.5" r="1.5" />
                  <polyline points="21,15 16,10 5,21" />
                </svg>
                <span class="text-sm font-medium">Ready to generate</span>
              </div>
            `}
        ${isGenerating &&
        html`
          <div class="absolute inset-0 z-10 bg-fg/60 text-bg flex flex-col items-center justify-center gap-3">
            <svg class="w-8 h-8 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <path d="M12 6v6l4 2" />
            </svg>
            <span class="text-sm font-semibold">Generating…</span>
          </div>
        `}
      </div>
    </div>
  `;
};

// Generic Modal component
export const Modal = ({ isOpen, onClose, title, children, size = "md", showCloseButton = true, className = "" }) => {
  const dialogRef = useRef(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    const handleEscape = (e) => {
      if (e.key === "Escape") onClose();
    };

    const handleClickOutside = (e) => {
      if (e.target === dialog) onClose();
    };

    if (isOpen) {
      try {
        dialog.showModal();
      } catch (e) {
        dialog.setAttribute("open", "");
      }
      dialog.addEventListener("click", handleClickOutside);
      document.addEventListener("keydown", handleEscape);
    } else {
      try {
        dialog.close();
      } catch (e) {
        dialog.removeAttribute("open");
      }
    }

    return () => {
      dialog.removeEventListener("click", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose]);

  const sizeClasses = {
    sm: "max-w-md",
    md: "max-w-2xl",
    lg: "max-w-4xl",
    xl: "max-w-6xl",
  };

  return html`
    <dialog
      ref=${dialogRef}
      class="bg-transparent border-0 p-4 w-full max-h-[90vh] overflow-visible fixed inset-0 m-auto ${sizeClasses[
        size
      ]} ${className}"
    >
      <div class="bg-bg border border-border shadow-xl max-h-[90vh] overflow-hidden flex flex-col">
        ${title &&
        html`
          <div class="flex items-center justify-between p-6 border-b border-border">
            <h2 class="text-xl font-bold">${title}</h2>
            ${showCloseButton &&
            html`
              <button
                class="w-8 h-8 flex items-center justify-center hover:bg-fg hover:text-bg transition-colors duration-200 text-xl font-bold leading-none"
                type="button"
                onClick=${onClose}
                aria-label="Close modal"
              >
                ×
              </button>
            `}
          </div>
        `}
        <div class="flex-1 p-6 overflow-y-auto">${children}</div>
      </div>
    </dialog>
  `;
};
