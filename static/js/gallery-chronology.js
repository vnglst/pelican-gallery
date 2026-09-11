document.querySelectorAll("[data-gallery-study]").forEach((study) => {
  const buttons = study.querySelectorAll("[data-gallery-provider]");
  const panels = study.querySelectorAll("[data-gallery-panel]");

  study.querySelectorAll("[data-gallery-edit]").forEach((link) => {
    link.addEventListener("click", (event) => event.stopPropagation());
  });

  buttons.forEach((button) => {
    button.addEventListener("click", () => {
      const provider = button.dataset.galleryProvider;
      buttons.forEach((candidate) => candidate.removeAttribute("aria-current"));
      button.setAttribute("aria-current", "page");
      panels.forEach((panel) => {
        panel.hidden = panel.dataset.galleryPanel !== provider;
        if (!panel.hidden) panel.scrollLeft = 0;
      });
    });
  });

  if (`#${study.id}` === window.location.hash) study.open = true;
  study.addEventListener("toggle", () => {
    if (study.open) {
      history.replaceState(null, "", `#${study.id}`);
    } else if (`#${study.id}` === window.location.hash) {
      history.replaceState(null, "", `${window.location.pathname}${window.location.search}`);
    }
  });
});
