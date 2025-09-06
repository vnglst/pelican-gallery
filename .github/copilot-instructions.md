# GitHub Copilot Instructions for Pelican Art Gallery

## Project Overview

This is **Pelican Art Gallery** - a Go-based web application for generating AI-powered SVG artwork. The project serves as an art gallery and workshop for creating custom SVG illustrations using various AI models, inspired by Simon Willison's famous experiment testing 16 different LLMs on generating SVG code for "a pelican riding a bicycle."

The project was sparked by the creative potential demonstrated in [Simon Willison's LLM SVG benchmark](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/), which showed how different AI models produce varying quality SVG outputs when given the same creative prompt.

## Architecture & Technologies

### Backend

- **Language**: Go 1.21+
- **Framework**: Standard library with custom routing
- **Templates**: Go HTML templates with component-based organization
- **Static Files**: CSS components, JavaScript, SVG assets
- **Build**: Makefile-based build system

### Frontend

- **CSS**: Component-based architecture with CSS custom properties (variables)
- **JavaScript**: **Modern ES6 modules architecture** with vanilla JavaScript
  - **Module System**: ES6 import/export with clear separation of concerns
  - **Architecture**: Functional composition, keeping things modular and simple
- **Templates**: Separated by purpose (homepage, workshop, gallery)

## Code Style & Conventions

### Go Code

- Follow standard Go conventions (gofmt, golint)
- Use explicit error handling
- Prefer composition over inheritance
- Keep handlers thin, business logic in separate packages
- Use struct embedding for configuration

### CSS Architecture

- **Component-based**: Each major UI component gets its own CSS file
- **CSS Custom Properties**: Use `--spacing-*`, `--color-*`, `--transition-*` variables
- **No CSS frameworks**: Keep dependencies minimal
- **Mobile-first**: Design for mobile, enhance for desktop
- **Dark mode support**: Use `prefers-color-scheme` media queries

### Template Organization

- `homepage.html`: Minimal landing page with project info and gallery link
- `workshop.html`: Full workshop interface for creating artwork
- `gallery.html`: Display generated artwork collection
- Keep templates focused on single purposes

### JavaScript

- **ES6 Modules**: Use import/export for all new code
- **Modern Patterns**: Private class fields (#private), async/await, destructuring
- **Architecture**: Functional composition over inheritance
- **Error Handling**: Comprehensive try/catch with user-friendly messages
- **Performance**: Lazy loading, request cancellation with AbortController
- **Accessibility**: Focus management, keyboard navigation, ARIA attributes

## Project Structure

```
/
├── main.go              # Entry point and routing
├── internal/
│   ├── api/            # HTTP handlers
│   ├── config/         # Configuration management
│   ├── database/       # Database operations
│   └── models/         # Data structures
├── templates/          # HTML templates (homepage, workshop, gallery)
├── static/
│   ├── js/             # Modern ES6 modules
│   │   ├── main.js     # Application entry point
│   │   ├── modules/    # Feature modules
│   │   │   ├── artwork.js    # Artwork generation and management
│   │   │   ├── config.js     # Configuration management
│   │   │   ├── modals.js     # Modal dialog management
│   │   │   ├── models.js     # Model selection and management
│   │   │   └── ui.js         # User interface utilities
│   │   └── utils/      # Pure utility functions
│   │       └── dom.js        # DOM manipulation utilities
│   ├── css/
│   │   ├── main.css    # Entry point imports all components
│   │   └── components/ # Individual component styles
│   └── style.css       # Base styles and utilities
├── config/             # YAML configuration files
└── bin/               # Built binaries
```

## Key Design Principles

### Simplicity First

- Prefer vanilla solutions over complex frameworks
- Keep dependencies minimal
- Choose readability over cleverness
- Avoid premature optimization

### Black & White Only - NO COLORS

**VERY IMPORTANT**: This project uses a strict black and white color scheme only. Never use any colors (red, blue, green, etc.) in the design. All styling should use only:

- `--fg` (black text)
- `--bg` (white background)
- `--bg-secondary` (light gray)
- `--border` (subtle borders)

Active states, highlights, and accents should use **bold typography, shadows, and transforms** instead of colors. For example:

- Use `font-weight: 600` for emphasis
- Use `box-shadow` for depth
- Use `transform: translateY(-1px)` for interactive feedback
- Use `background: var(--fg)` with `color: var(--bg)` for active states

### Component Separation

- Each template serves a distinct purpose
- CSS components match template structure
- JavaScript modules focus on specific functionality
- Clear separation of concerns

### User Experience

- **Homepage**: Extremely minimal - just project name, description, and gallery link
- **Workshop**: Full-featured interface for artwork creation
- **Gallery**: Clean display of generated artwork
- Mobile-responsive design throughout

### Performance

- Minimal JavaScript bundle size
- CSS custom properties for theming
- Efficient Go template rendering
- Static asset optimization

## Development Patterns

### Adding New Features

1. Start with Go handlers in `internal/api/`
2. Create necessary templates in `templates/`
3. Add component CSS in `static/css/components/`
4. Add JavaScript functionality if needed
5. Update routing in `main.go`

### CSS Development

- Create new component files for major UI sections
- Use existing CSS custom properties for consistency
- Test dark mode compatibility
- Ensure mobile responsiveness

### Template Updates

- Keep templates semantic and accessible
- Use Go template partials for shared components
- Maintain consistent header/footer structure
- Test across different content types

## Common Gotchas & Lessons Learned

### Template Routing

- Ensure route handlers use correct template files
- `index.html` vs `homepage.html` vs `workshop.html` confusion has happened
- Test routes after template changes

### CSS Variables

- Always use CSS custom properties for spacing, colors, transitions
- Check both light and dark mode when adding styles
- Component styles should be self-contained but use global variables

### JavaScript Complexity

- **Modular Architecture**: Refactored from monolithic script.js to ES6 modules for maintainability
- **Dependency Injection**: Use constructor injection instead of dynamic imports for better testability
- **Context Binding**: Arrow functions solve event handler context issues in class methods
- **Private Fields**: Use `#private` syntax for encapsulation instead of closures
- **Event Delegation**: Centralized event handling prevents memory leaks and improves performance

### Build Process

- Use `make dev` for development with hot reload
- Use `make clean && make build` for production builds
- Test builds before deployment

## File Naming Conventions

- CSS components: `component-name.css` in `static/css/components/`
- Templates: `purpose.html` in `templates/`
- Go packages: lowercase, descriptive names
- JavaScript: Semantic function and variable names

## Testing & Quality

- Test responsive design on multiple screen sizes
- Verify dark mode functionality
- Check accessibility with semantic HTML
- Validate Go code with standard tools
- Test build process regularly

## Common Tasks

### Adding a New Page

1. Create template in `templates/new-page.html`
2. Add route handler in `internal/api/handlers.go`
3. Create CSS component in `static/css/components/new-page.css`
4. Import CSS component in `static/css/main.css`
5. Add route in `main.go`

### Styling Updates

1. Use existing CSS custom properties when possible
2. Add new component CSS file if needed
3. Test both light and dark modes
4. Ensure mobile compatibility
5. Keep styles semantic and maintainable

### JavaScript Features

1. **ES6 Modules**: Use import/export for all new code with clear module boundaries
2. **Class-based Architecture**: Each feature module exports a class with private fields (#private)
3. **Dependency Injection**: Managers receive DOM elements and modal instances at construction
4. **Modern JavaScript**: async/await, destructuring, optional chaining, arrow functions
5. **Event Delegation**: Centralized event handling with proper context binding
6. **Error Boundaries**: Comprehensive try/catch with user-friendly error messages
7. **Request Cancellation**: AbortController for cancelling in-flight requests
8. **Accessibility**: Focus management, keyboard navigation, ARIA attributes

### JavaScript Development Tasks

1. **Adding New Modules**: Create class-based modules in `static/js/modules/` with dependency injection
2. **Updating Existing Features**: Modify specific module files rather than monolithic script
3. **DOM Utilities**: Add pure functions to `static/js/utils/dom.js` for reusable DOM operations
4. **Event Handling**: Use event delegation in main.js for new interactive elements
5. **Error Handling**: Implement try/catch with UI.showError() for user-friendly feedback

## ES6 Module Architecture

### Module Organization

- **`main.js`**: Application bootstrap, event setup, and global state management
- **`modules/artwork.js`**: Artwork generation, saving, regeneration, and gallery management
- **`modules/config.js`**: Configuration management with server synchronization
- **`modules/modals.js`**: Modal dialog management (model selector, examples, config)
- **`modules/models.js`**: Model selection, management, and persistence
- **`modules/ui.js`**: User interface utilities, error handling, and feedback
- **`utils/dom.js`**: Pure DOM manipulation utilities and element queries

### Architecture Patterns

- **Class-based modules**: Each feature module exports a class with private fields (#private)
- **Dependency injection**: Managers receive DOM elements and modal instances
- **Event delegation**: Centralized event handling with proper context binding
- **Async/await patterns**: Modern asynchronous programming throughout
- **Global state management**: Centralized state in main.js with selective exports

### Modern Patterns

```javascript
// Class-based module with dependency injection
export class ArtworkManager {
  #elements;
  #modals;

  constructor(elements, modals) {
    this.#elements = elements;
    this.#modals = modals;
  }

  async handleGenerate(prompt, selectedModels) {
    // Implementation with proper error handling
    try {
      const response = await fetch("/api/generate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt, models: selectedModels }),
      });

      if (!response.ok) throw new Error("Generation failed");

      const result = await response.json();
      this.#displayArtwork(result);
    } catch (error) {
      UI.showError(`Generation failed: ${error.message}`);
    }
  }
}
```

### Event Handling

- Use event delegation instead of scattered listeners
- Implement custom events for module communication
- Provide keyboard navigation and accessibility

### Error Management

- Comprehensive error boundaries
- User-friendly error messages
- Console logging for debugging
- Optional error tracking integration

## Inspiration & Context

This project was directly inspired by [Simon Willison's LLM SVG benchmark](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/), where he tested 16 different AI models on their ability to generate SVG code for "a pelican riding a bicycle." Key insights from this experiment:

- **Creative Prompt Selection**: The pelican bicycle prompt was chosen because it's unlikely to exist in training data
- **Model Diversity**: Tested models from OpenAI, Anthropic, Google Gemini, Meta's Llama, and others
- **Quality Variations**: Different models produce significantly different SVG quality and styles
- **Accessibility**: Demonstrated that AI can create unique, creative SVG art beyond what's in training data

Pelican Art Gallery extends this concept by providing a user-friendly interface for exploring AI-generated SVG art with any creative prompt, not just pelicans on bicycles.
