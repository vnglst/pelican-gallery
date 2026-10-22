# SVG isolation validation

Baseline: `51f6a39c0b075950d86537ab3dea0f91bddf8c53`.

Chrome was used for both measurements. The production baseline and HTML
rendered by the unmodified local Go server produced the same counts. The
before/after previews used server-rendered HTML, the same SQLite snapshot,
and each version's actual JavaScript and CSS. A static preview served these
pages to Chrome because its browser runs separately from the Go process.

An affected result means that at least one `url(#id)` or fragment `href`
resolves to an element outside its own SVG, despite having a local definition.
It does **not** mean the entire image was invisible. This measure excludes
additional visual changes caused by CSS leaking between SVGs.

| Category | Generated images | Affected before | Wrong references before | Affected after | Loaded after |
| --- | ---: | ---: | ---: | ---: | ---: |
| abstract | 78 | 25 | 115 | 0 | 78 |
| classic | 116 | 79 | 439 | 0 | 116 |
| fun | 30 | 18 | 239 | 0 | 30 |
| modern | 95 | 56 | 549 | 0 | 95 |
| photography | 52 | 33 | 125 | 0 | 52 |
| sculpture | 39 | 18 | 85 | 0 | 39 |
| Total | 410 | 229 | 1,552 | 0 | 410 |

After isolation, generated SVGs live in separate image documents. None of
their IDs or styles exist in the gallery document. All 410 `<img>` elements
reported `complete && naturalWidth > 0` in Chrome.

Visual check: in Delicate Tension No. 85, GPT-5 Nano previously referenced
the `bg` rectangle belonging to Composition X instead of its own gradient.
Its background was absent. The isolated result shows the intended gradient.

Additional validation:

- `go test ./...` passes, including preservation/escaping of malformed source
  in the frame template.
- All 433 nonempty database SVG sources serialize to valid XML in a DOM
  compatibility check. Eleven raw sources originally failed XML parsing;
  three contain invalid attribute names retained by HTML parsing.
- `/static/tests/svg-image.html` checks the real six-category corpus with
  Chrome's XML parser and image decoder, plus five fixtures for malformed
  markup, unescaped text, missing namespace declarations, and invalid
  attributes: 415 checks passed.
- Lightbox opening and next-image navigation load isolated images and add
  no SVG definitions to the page.
- Downloading the GPT-5 Nano result as PNG from `/group/76?provider=openai`
  successfully produces a download in Chrome.

The regression page can be rerun against the application by opening
`/static/tests/svg-image.html` and clicking **Run checks**. Safari was not
tested. No generated artwork or SQLite data was changed.
