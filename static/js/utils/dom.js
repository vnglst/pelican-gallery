// DOM utility functions
export const DOM = {
  get: (id) => document.getElementById(id),
  query: (selector) => document.querySelector(selector),
  queryAll: (selector) => document.querySelectorAll(selector),

  updateSliderValue: (input, display) => {
    display.textContent = input.value;
  },

  showElement: (element) => {
    element.style.display = "block";
  },

  hideElement: (element) => {
    element.style.display = "none";
  },

  toggleElement: (element, show) => {
    element.style.display = show ? "block" : "none";
  },
};
