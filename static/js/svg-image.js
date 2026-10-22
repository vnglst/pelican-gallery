// Parse generated SVG as HTML, as the gallery did before, then serialize the
// browser's repaired SVG tree. Raw model output is not always well-formed XML.
const validAttributes = new Map();

function isXMLAttribute(attribute) {
  if (/[\s<>"'=/%]/u.test(attribute.name)) return false;
  const key = `${attribute.namespaceURI || ""}:${attribute.name}`;
  if (!validAttributes.has(key)) {
    try {
      document.createAttributeNS(attribute.namespaceURI, attribute.name);
      validAttributes.set(key, true);
    } catch {
      validAttributes.set(key, false);
    }
  }
  return validAttributes.get(key);
}

export function parseSVG(source) {
  if (typeof source === "string") {
    const parsed = new DOMParser().parseFromString(source, "text/html");
    return parsed.querySelector("svg");
  }
  return source;
}

export function svgImageURL(source, dimensions) {
  const svg = parseSVG(source);
  if (!svg) return "";

  const image = svg.cloneNode(true);
  // HTML accepts typo attributes such as cx="100"="120". XML rejects the
  // resulting attribute named ="120", so discard only invalid XML names.
  [image, ...image.querySelectorAll("*")].forEach((element) => {
    [...element.attributes].forEach((attribute) => {
      if (!isXMLAttribute(attribute)) element.removeAttributeNode(attribute);
    });
  });
  image.setAttribute("xmlns", "http://www.w3.org/2000/svg");
  const { width, height } = dimensions || svgDimensions(svg);
  image.setAttribute("width", String(width));
  image.setAttribute("height", String(height));
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(new XMLSerializer().serializeToString(image))}`;
}

export function svgDimensions(svg) {
  const viewBox = svg.viewBox && svg.viewBox.baseVal;
  if (viewBox && viewBox.width > 0 && viewBox.height > 0) {
    return { width: viewBox.width, height: viewBox.height };
  }
  const width = Number.parseFloat(svg.getAttribute("width"));
  const height = Number.parseFloat(svg.getAttribute("height"));
  return {
    width: Number.isFinite(width) && width > 0 ? width : 1024,
    height: Number.isFinite(height) && height > 0 ? height : 1024,
  };
}

document.querySelectorAll("[data-svg-frame]").forEach((frame) => {
  const svg = parseSVG(frame.querySelector("template").content.textContent);
  const image = frame.querySelector("img");
  if (!svg) {
    image.hidden = true;
    return;
  }
  image.alt = svg?.getAttribute("aria-label") || svg?.querySelector("title")?.textContent || "";
  image.src = svgImageURL(svg);
});
