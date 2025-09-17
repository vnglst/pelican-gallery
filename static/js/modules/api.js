// API helper module for Workshop
// Uses fetch with JSON handling, consistent errors, and timeouts.

const DEFAULT_TIMEOUT = 10 * 60 * 1_000; // 10 minutes

const withTimeout = (promise, ms = DEFAULT_TIMEOUT, signal) => {
  if (signal) return promise; // external AbortController handles cancel
  let timeoutId;
  const timeout = new Promise((_, reject) => {
    timeoutId = setTimeout(() => reject(new Error(`Request timed out after ${ms}ms`)), ms);
  });
  return Promise.race([promise.finally(() => clearTimeout(timeoutId)), timeout]);
};

const parseError = async (response) => {
  const contentType = response.headers.get("content-type") || "";
  try {
    if (contentType.includes("application/json")) {
      const data = await response.json();
      return data?.message || JSON.stringify(data);
    }
    return await response.text();
  } catch (e) {
    return response.statusText || `HTTP ${response.status} error`;
  }
};

const request = async (url, options = {}) => {
  try {
    const res = await withTimeout(fetch(url, options), options.timeout, options.signal);
    if (!res.ok) {
      throw new Error(await parseError(res));
    }
    const ct = res.headers.get("content-type") || "";
    if (ct.includes("application/json")) return res.json();
    return res.text();
  } catch (err) {
    if (err.name === "TypeError" && String(err.message).includes("fetch")) {
      throw new Error("Network error: Unable to connect to server");
    }
    throw err;
  }
};

// Endpoints
const getModels = () => request("/api/models");

const getGroup = (groupId) => request(`/api/groups/${groupId}`);

const createGroup = (payload) =>
  request("/api/groups", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

const updateGroup = (groupId, payload) =>
  request(`/api/groups/${groupId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

const deleteGroup = (groupId) => request(`/api/groups/${groupId}`, { method: "DELETE" });

const createArtwork = (payload) =>
  request("/api/artworks", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

const updateArtwork = (artworkId, payload) =>
  request(`/api/artworks/${artworkId}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

const deleteArtwork = (artworkId) => request(`/api/delete-artwork/${artworkId}`, { method: "DELETE" });

const generateArtwork = (artworkId) =>
  request("/api/generate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ artwork_id: Number(artworkId) }),
  });

// Default aggregated API object for convenient imports
const api = {
  getModels,
  getGroup,
  createGroup,
  updateGroup,
  deleteGroup,
  createArtwork,
  updateArtwork,
  deleteArtwork,
  generateArtwork,
};

export default api;
