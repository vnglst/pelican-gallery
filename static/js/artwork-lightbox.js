(() => {
  const dialog = document.createElement("dialog");
  dialog.className = "artwork-lightbox";
  dialog.setAttribute("aria-modal", "true");
  dialog.setAttribute("aria-labelledby", "artwork-lightbox-title");
  dialog.innerHTML = `
    <header class="artwork-lightbox-header">
      <div class="artwork-lightbox-details">
        <strong id="artwork-lightbox-title"></strong>
        <span class="artwork-lightbox-meta"></span>
      </div>
      <div class="artwork-lightbox-controls">
        <button type="button" class="artwork-lightbox-previous" aria-label="Previous artwork" title="Previous artwork (Left arrow)">←</button>
        <span class="artwork-lightbox-position" aria-live="polite"></span>
        <button type="button" class="artwork-lightbox-next" aria-label="Next artwork" title="Next artwork (Right arrow)">→</button>
        <button type="button" class="artwork-lightbox-close" aria-label="Close enlarged artwork" title="Close (Escape)">×</button>
      </div>
    </header>
    <div class="artwork-lightbox-content"></div>
  `;
  document.body.append(dialog);

  const content = dialog.querySelector(".artwork-lightbox-content");
  const title = dialog.querySelector("#artwork-lightbox-title");
  const meta = dialog.querySelector(".artwork-lightbox-meta");
  const position = dialog.querySelector(".artwork-lightbox-position");
  const previousButton = dialog.querySelector(".artwork-lightbox-previous");
  const nextButton = dialog.querySelector(".artwork-lightbox-next");
  const closeButton = dialog.querySelector(".artwork-lightbox-close");
  let closeFallback = null;
  let currentMedia = null;
  let navigationItems = [];

  const lockPageScroll = () => {
    document.documentElement.classList.add("lightbox-open");
  };

  const unlockPageScroll = () => {
    document.documentElement.classList.remove("lightbox-open");
  };

  const finishClose = () => {
    if (!dialog.open) return;
    window.clearTimeout(closeFallback);
    closeFallback = null;
    dialog.close();
  };

  const close = () => {
    if (!dialog.open || dialog.classList.contains("is-closing")) return;
    dialog.classList.add("is-closing");
    closeFallback = window.setTimeout(finishClose, 220);
  };

  const getNavigationItems = (media) => {
    const scope = media.closest("[data-chronology-scope]") || document;
    return [...scope.querySelectorAll("[data-artwork-media]")].filter((item) =>
      !item.closest("[hidden]") && item.getClientRects().length > 0
    );
  };

  const showArtwork = (media) => {
    const artwork = media.querySelector("img, svg");
    if (!artwork) return false;

    const header = media.closest("figure")?.querySelector(".chronology-card-header");
    const metaGroup = header?.querySelector(".chronology-card-meta");
    const metaParts = metaGroup
      ? [...metaGroup.children].map((item) => item.textContent.trim()).filter(Boolean)
      : [...(header?.querySelectorAll(":scope > span") || [])].map((item) => item.textContent.trim()).filter(Boolean);

    currentMedia = media;
    navigationItems = getNavigationItems(media);
    const currentIndex = navigationItems.indexOf(media);

    title.textContent = header?.querySelector("strong")?.textContent.trim() || "Artwork";
    meta.textContent = metaParts.join(" · ");
    meta.hidden = metaParts.length === 0;
    position.textContent = navigationItems.length > 1 ? `${currentIndex + 1} / ${navigationItems.length}` : "";
    previousButton.hidden = navigationItems.length < 2;
    nextButton.hidden = navigationItems.length < 2;
    content.replaceChildren(artwork.cloneNode(true));
    return true;
  };

  const navigate = (offset) => {
    if (!currentMedia || navigationItems.length < 2) return;
    const currentIndex = navigationItems.indexOf(currentMedia);
    const nextIndex = (currentIndex + offset + navigationItems.length) % navigationItems.length;
    showArtwork(navigationItems[nextIndex]);
  };

  closeButton.addEventListener("click", close);
  previousButton.addEventListener("click", () => navigate(-1));
  nextButton.addEventListener("click", () => navigate(1));
  dialog.addEventListener("click", (event) => {
    if (event.target === dialog) close();
  });
  dialog.addEventListener("keydown", (event) => {
    if (event.key === "ArrowLeft") {
      event.preventDefault();
      navigate(-1);
    } else if (event.key === "ArrowRight") {
      event.preventDefault();
      navigate(1);
    }
  });
  dialog.addEventListener("cancel", (event) => {
    event.preventDefault();
    close();
  });
  dialog.addEventListener("animationend", (event) => {
    if (event.animationName === "artwork-lightbox-out") finishClose();
  });
  dialog.addEventListener("close", () => {
    window.clearTimeout(closeFallback);
    closeFallback = null;
    dialog.classList.remove("is-closing");
    content.replaceChildren();
    currentMedia = null;
    navigationItems = [];
    unlockPageScroll();
  });

  document.querySelectorAll("[data-artwork-media]").forEach((media) => {
    let pointerStartX = null;
    media.setAttribute("tabindex", "0");
    media.setAttribute("role", "button");
    media.setAttribute("aria-label", `Enlarge ${media.closest("figure")?.querySelector("strong")?.textContent || "artwork"}`);

    const open = (event) => {
      if (event.target.closest("button, a")) return;
      if (event.type === "click" && pointerStartX !== null && Math.abs(event.clientX - pointerStartX) > 5) return;
      if (!showArtwork(media)) return;
      dialog.classList.remove("is-closing");
      dialog.showModal();
      lockPageScroll();
      closeButton.focus();
    };

    media.addEventListener("pointerdown", (event) => {
      pointerStartX = event.clientX;
    });
    media.addEventListener("click", open);
    media.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        open(event);
      }
    });
  });
})();
