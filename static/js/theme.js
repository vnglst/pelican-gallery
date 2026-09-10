(function () {
  const storageKey = "pelican-theme";
  const root = document.documentElement;
  const systemTheme = window.matchMedia("(prefers-color-scheme: dark)");
  const select = document.querySelector("[data-theme-select]");

  function applyTheme(preference) {
    const resolved = preference === "system" ? (systemTheme.matches ? "dark" : "light") : preference;
    root.dataset.theme = resolved;
    root.dataset.themePreference = preference;
    if (select) select.value = preference;
  }

  function savedPreference() {
    try {
      const preference = localStorage.getItem(storageKey);
      return preference === "light" || preference === "dark" ? preference : "system";
    } catch (_) {
      return "system";
    }
  }

  if (select) {
    select.addEventListener("change", function () {
      const preference = select.value;
      try {
        if (preference === "system") localStorage.removeItem(storageKey);
        else localStorage.setItem(storageKey, preference);
      } catch (_) {}
      applyTheme(preference);
    });
  }

  systemTheme.addEventListener("change", function () {
    if (savedPreference() === "system") applyTheme("system");
  });

  applyTheme(savedPreference());
})();
