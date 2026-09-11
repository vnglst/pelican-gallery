(() => {
  document.querySelectorAll("[data-chronology-scope]").forEach((scope) => {
    scope.querySelectorAll("[data-chronology-scroll]").forEach((button) => {
      button.addEventListener("click", () => {
        const chronology = scope.querySelector("[data-chronology]:not([hidden])");
        if (!chronology) return;
        const direction = button.dataset.chronologyScroll === "previous" ? -1 : 1;
        chronology.scrollBy({
          left: direction * Math.max(240, chronology.clientWidth * 0.75),
          behavior: "smooth",
        });
      });
    });
  });

  document.querySelectorAll("[data-chronology]").forEach((chronology) => {
    let startX = 0;
    let startScrollLeft = 0;
    let pointerId = null;
    let dragging = false;
    let suppressClick = false;
    const dragThreshold = 6;

    chronology.addEventListener("pointerdown", (event) => {
      if (event.pointerType !== "mouse" || event.button !== 0 || event.target.closest("button, a")) return;
      startX = event.clientX;
      startScrollLeft = chronology.scrollLeft;
      pointerId = event.pointerId;
      dragging = false;
    });

    chronology.addEventListener("pointermove", (event) => {
      if (event.pointerId !== pointerId) return;
      if (!dragging && Math.abs(event.clientX - startX) >= dragThreshold) {
        dragging = true;
        suppressClick = true;
        chronology.classList.add("is-dragging");
        chronology.setPointerCapture(event.pointerId);
      }
      if (!dragging) return;
      chronology.scrollLeft = startScrollLeft - (event.clientX - startX);
    });

    const stopDragging = (event) => {
      if (event.pointerId !== pointerId) return;
      if (dragging) chronology.classList.remove("is-dragging");
      if (chronology.hasPointerCapture(event.pointerId)) chronology.releasePointerCapture(event.pointerId);
      pointerId = null;
      dragging = false;
    };

    chronology.addEventListener("pointerup", stopDragging);
    chronology.addEventListener("pointercancel", stopDragging);
    chronology.addEventListener("click", (event) => {
      if (!suppressClick) return;
      suppressClick = false;
      event.preventDefault();
      event.stopPropagation();
    }, true);
  });
})();
