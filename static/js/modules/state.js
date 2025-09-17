export function createInitialState(windowObj = window) {
  const initialArtworks = (() => {
    try {
      const raw = windowObj.existingArtworks;
      let list = [];
      if (Array.isArray(raw)) {
        list = raw;
      } else if (typeof raw === "string" && raw.trim()) {
        try {
          const parsed = JSON.parse(raw);
          if (Array.isArray(parsed)) list = parsed;
        } catch (e) {
          console.warn("Failed to parse window.existingArtworks string:", e);
        }
      }
      const map = new Map();
      list.forEach((art) => {
        const id = art.id;
        if (id != null) map.set(Number(id), art);
      });
      return map;
    } catch (e) {
      console.warn("Error initializing artworks from window:", e);
      return new Map();
    }
  })();

  return {
    currentGroup: null,
    artworks: initialArtworks,
    models: [],
    toasts: [],
    loading: { visible: false, message: "" },
    modals: { model: false, config: false },
    configArtwork: null,
    modelsLoading: false,
    modelsError: "",
    generatingArtworks: new Set(),
    formData: {
      title: "Sunflowers",
      prompt: "Sunflowers by Vincent van Gogh.",
      category: "Art",
    },
  };
}

export const reducer = (state, action) => {
  switch (action.type) {
    case "SET_CURRENT_GROUP":
      return { ...state, currentGroup: action.payload };
    case "SET_ARTWORKS":
      return { ...state, artworks: action.payload };
    case "ADD_ARTWORK": {
      const map = new Map(state.artworks);
      const art = action.payload;
      if (art?.id != null) map.set(Number(art.id), art);
      return { ...state, artworks: map };
    }
    case "UPDATE_ARTWORK": {
      const map = new Map(state.artworks);
      const art = action.payload;
      if (art?.id != null) map.set(Number(art.id), art);
      return { ...state, artworks: map };
    }
    case "REMOVE_ARTWORK": {
      const map = new Map(state.artworks);
      map.delete(Number(action.payload));
      return { ...state, artworks: map };
    }
    case "SET_MODELS":
      return { ...state, models: action.payload };
    case "PUSH_TOAST":
      return { ...state, toasts: [...state.toasts, action.payload] };
    case "REMOVE_TOAST": {
      const i = action.payload;
      return { ...state, toasts: state.toasts.filter((_, idx) => idx !== i) };
    }
    case "SHOW_LOADING":
      return { ...state, loading: { visible: true, message: action.payload || "Loading..." } };
    case "HIDE_LOADING":
      return { ...state, loading: { visible: false, message: "" } };
    case "SET_MODAL":
      return { ...state, modals: { ...state.modals, [action.payload.modal]: action.payload.value } };
    case "SET_CONFIG_ARTWORK":
      return { ...state, configArtwork: action.payload };
    case "SET_MODELS_LOADING":
      return { ...state, modelsLoading: action.payload };
    case "SET_MODELS_ERROR":
      return { ...state, modelsError: action.payload };
    case "ADD_GENERATING": {
      const s = new Set(state.generatingArtworks);
      s.add(Number(action.payload));
      return { ...state, generatingArtworks: s };
    }
    case "REMOVE_GENERATING": {
      const s = new Set(state.generatingArtworks);
      s.delete(Number(action.payload));
      return { ...state, generatingArtworks: s };
    }
    case "SET_FORM_DATA":
      return { ...state, formData: action.payload };
    default:
      return state;
  }
};
