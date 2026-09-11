(() => {
  const chronology = document.querySelector("[data-chronology]");
  if (!chronology) return;

  document.querySelectorAll("[data-chronology-scroll]").forEach((button) => {
    button.addEventListener("click", () => {
      const direction = button.dataset.chronologyScroll === "previous" ? -1 : 1;
      chronology.scrollBy({
        left: direction * Math.max(240, chronology.clientWidth * 0.75),
        behavior: "smooth",
      });
    });
  });

  let startX = 0;
  let startScrollLeft = 0;

  chronology.addEventListener("pointerdown", (event) => {
    if (event.pointerType !== "mouse" || event.button !== 0 || event.target.closest("button, a")) return;
    startX = event.clientX;
    startScrollLeft = chronology.scrollLeft;
    chronology.classList.add("is-dragging");
    chronology.setPointerCapture(event.pointerId);
  });

  chronology.addEventListener("pointermove", (event) => {
    if (!chronology.classList.contains("is-dragging")) return;
    chronology.scrollLeft = startScrollLeft - (event.clientX - startX);
  });

  const stopDragging = (event) => {
    if (!chronology.classList.contains("is-dragging")) return;
    chronology.classList.remove("is-dragging");
    if (chronology.hasPointerCapture(event.pointerId)) chronology.releasePointerCapture(event.pointerId);
  };

  chronology.addEventListener("pointerup", stopDragging);
  chronology.addEventListener("pointercancel", stopDragging);
})();
