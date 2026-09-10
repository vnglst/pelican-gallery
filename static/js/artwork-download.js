function safeFilename(title) {
  const filename = title
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");

  return filename || "artwork";
}

function svgDimensions(svg) {
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

async function downloadArtwork(button) {
  const result = button.closest("figure");
  const sourceSvg = result && result.querySelector("div > div > svg");
  const sourceImage = result && result.querySelector("img");
  if (!sourceSvg && !sourceImage) return;

  button.disabled = true;
  button.setAttribute("aria-busy", "true");

  let image = sourceImage;
  let imageUrl;
  let width = sourceImage && sourceImage.naturalWidth;
  let height = sourceImage && sourceImage.naturalHeight;

  try {
    if (sourceSvg) {
      const svg = sourceSvg.cloneNode(true);
      svg.setAttribute("xmlns", "http://www.w3.org/2000/svg");

      const dimensions = svgDimensions(svg);
      const scale = Math.max(1, 2048 / Math.max(dimensions.width, dimensions.height));
      width = Math.round(dimensions.width * scale);
      height = Math.round(dimensions.height * scale);
      svg.setAttribute("width", String(width));
      svg.setAttribute("height", String(height));

      const svgBlob = new Blob([new XMLSerializer().serializeToString(svg)], {
        type: "image/svg+xml;charset=utf-8",
      });
      imageUrl = URL.createObjectURL(svgBlob);
      image = new Image();
      image.decoding = "async";
      image.src = imageUrl;
      await image.decode();
    } else if (!sourceImage.complete) {
      await sourceImage.decode();
      width = sourceImage.naturalWidth;
      height = sourceImage.naturalHeight;
    }

    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;
    const context = canvas.getContext("2d");
    context.drawImage(image, 0, 0, width, height);

    const pngBlob = await new Promise((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (blob) resolve(blob);
        else reject(new Error("PNG encoding failed"));
      }, "image/png");
    });

    const downloadUrl = URL.createObjectURL(pngBlob);
    const link = document.createElement("a");
    link.href = downloadUrl;
    link.download = `${safeFilename(button.dataset.title || "artwork")}.png`;
    link.click();
    URL.revokeObjectURL(downloadUrl);
  } catch (error) {
    console.error("Could not download artwork as PNG", error);
    window.alert("This artwork could not be downloaded as a PNG.");
  } finally {
    if (imageUrl) URL.revokeObjectURL(imageUrl);
    button.disabled = false;
    button.removeAttribute("aria-busy");
  }
}

document.addEventListener("click", (event) => {
  const button = event.target.closest(".download-png");
  if (!button) return;

  downloadArtwork(button);
});
