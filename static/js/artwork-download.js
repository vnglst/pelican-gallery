import { parseSVG, svgDimensions, svgImageURL } from "./svg-image.js";

function safeFilename(title) {
  const filename = title
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");

  return filename || "artwork";
}

async function downloadArtwork(button) {
  const result = button.closest("figure");
  const media = result && result.querySelector("[data-artwork-media]");
  const source = media && media.querySelector("[data-svg-source]");
  const sourceSvg = source && parseSVG(source.content.textContent);
  const sourceImage = media && media.querySelector("img");
  if (!sourceSvg && !sourceImage) return;

  button.disabled = true;
  button.setAttribute("aria-busy", "true");

  let image = sourceImage;
  let width = sourceImage && sourceImage.naturalWidth;
  let height = sourceImage && sourceImage.naturalHeight;

  try {
    if (sourceSvg) {
      const dimensions = svgDimensions(sourceSvg);
      const scale = Math.max(1, 2048 / Math.max(dimensions.width, dimensions.height));
      width = Math.round(dimensions.width * scale);
      height = Math.round(dimensions.height * scale);
      image = new Image();
      image.decoding = "async";
      image.src = svgImageURL(sourceSvg, { width, height });
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
    button.disabled = false;
    button.removeAttribute("aria-busy");
  }
}

document.addEventListener("click", (event) => {
  const button = event.target.closest(".download-png");
  if (!button) return;

  downloadArtwork(button);
});
