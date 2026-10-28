
  import { svgImageURL } from "../js/svg-image.js";

  const button = document.querySelector("#run");
  const output = document.querySelector("#result");
  button.addEventListener("click", async () => {
    button.disabled = true;
    const results = [];
    const failures = [];
    let tested = 0;
    async function check(source, label) {
      const url = svgImageURL(source);
      const xml = decodeURIComponent(url.slice(url.indexOf(",") + 1));
      const parsed = new DOMParser().parseFromString(xml, "image/svg+xml");
      if (!url || parsed.querySelector("parsererror")) throw new Error(`${label}: invalid SVG XML`);
      const image = new Image();
      image.src = url;
      await image.decode();
      if (!image.naturalWidth || !image.naturalHeight) throw new Error(`${label}: empty image`);
      tested++;
    }
    try {
      const fixtures = [
        '<svg viewBox="0 0 20 10"><rect width="20" height="10" fill="red"/></svg>',
        '<svg><text>Tom & Jerry</text></svg>',
        '<svg><rect fill="red"/><!-- unfinished comment',
        '<svg><circle cx="100"="120" r="4"/><circle r1.8"/><filter height%"></filter></svg>',
        '<svg viewBox="0 0 10 10"><defs><path id="shape" d="M0 0h10v10z"/></defs><use xlink:href="#shape"/></svg>',
      ];
      for (let i = 0; i < fixtures.length; i++) await check(fixtures[i], `Fixture ${i + 1}`);
      results.push({ category: "Malformed XML and namespaces", tested: fixtures.length, failures: 0 });
      for (const category of ["abstract", "classic", "fun", "modern", "photography", "sculpture"]) {
        output.textContent = `Checking ${category}…`;
        const response = await fetch(`/gallery/category/${category}`);
        if (!response.ok) throw new Error(`${category}: HTTP ${response.status}`);
        const page = new DOMParser().parseFromString(await response.text(), "text/html");
        const sources = [...page.querySelectorAll("[data-artwork-media] [data-svg-source]")];
        if (!sources.length) throw new Error(`${category}: no isolated SVG sources found`);
        let failed = 0;
        for (const source of sources) {
          const label = `${category}: ${source.closest("figure")?.querySelector("strong")?.textContent}`;
          try { await check(source.content.textContent, label); }
          catch (error) { failures.push(String(error)); failed++; }
        }
        results.push({ category, tested: sources.length, failures: failed });
      }
      output.textContent = JSON.stringify({ status: failures.length ? "FAIL" : "PASS", tested, results, failures }, null, 2);
    } catch (error) {
      output.textContent = `FAIL: ${error}`;
    } finally {
      button.disabled = false;
    }
  });
