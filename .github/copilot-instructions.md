# GitHub Copilot Instructions for Pelican Art Gallery

## Project Overview├── static/

│ ├── js/ # Modern ES6 modules
│ │ └── workshop.js # Preact-based workshop interface
│ ├── css/
│ │ ├── input.css # Tailwind CSS v4 configuration with @theme
│ │ └── output.css # Generated Tailwind CSS

│ └── favicon.svg # Site icon

**Pelican Art Gallery** is a Go-based web application for generating AI-powered SVG artwork. What happens when you ask an AI to draw famous paintings as simple computer graphics? This project shows how language models—usually built for text—try to turn classic artworks into SVG images. Even though these AIs aren't trained to make vector graphics, it's fun to see how they mix what they know about art with their ability to write SVG code. This project is inspired by [Simon Willison's article](https://simonwillison.net/2025/Jun/6/six-months-in-llms/) about creative LLM benchmarks.

## Architecture & Technologies

### Backend

- **Language**: Go 1.21+
- **Framework**: Standard library with custom routing
- **Templates**: Go HTML templates with component-based organization
- **Static Files**: CSS components, JavaScript, SVG assets
- **Build**: Makefile-based build system

### Frontend

- **CSS**: Tailwind CSS v4 with CSS-first configuration using `@theme` directive
- **JavaScript**: **Preact-based workshop interface** with modern ES6 modules
  - **Workshop App**: Preact with HTM for component-based rendering in `workshop.js`
  - **Module System**: ES6 import/export for utility functions and components
  - **Architecture**: Functional composition with Preact hooks (useState, useEffect, useRef)
- **Templates**: Go HTML templates for server-side rendering (homepage, workshop, gallery)

## Code Style & Conventions

### Go Code

- Follow standard Go conventions (gofmt, golint)
- Use explicit error handling
- Prefer composition over inheritance
- Keep handlers thin, business logic in separate packages
- Use struct embedding for configuration

### CSS Architecture

- **Tailwind CSS v4**: CSS-first configuration with `@theme` directive in `static/css/input.css`
- **Custom Theme**: Black/white color scheme with CSS custom properties
- **Utility-First**: All styling uses Tailwind utility classes, no custom CSS components
- **Mobile-first**: Responsive design with Tailwind breakpoints
- **Build Process**: Standalone Tailwind binary processes classes from templates and JavaScript

### Template Organization

- `homepage.html`: Minimal landing page with project info and gallery link
- `workshop.html`: Full workshop interface for creating artwork
- `gallery.html`: Display generated artwork collection
- Keep templates focused on single purposes

### JavaScript

- **ES6 Modules**: Use import/export for all new code
- **Preact Components**: Workshop interface built with Preact functional components
- **HTM**: Template literals for JSX-like syntax without build step
- **Modern Patterns**: Hooks (useState, useEffect, useRef), async/await, destructuring
- **Architecture**: Functional composition with component-based UI
- **Error Handling**: Comprehensive try/catch with user-friendly messages
- **Performance**: Lightweight Preact bundle, efficient re-rendering
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

### Design Aesthetic & Visual Language

A modern, minimalist design aesthetic with a strict high-contrast, black-and-white color palette. The style is clean and gallery-inspired, prioritizing readability and uncluttered space.

It uses a professional, legible sans-serif font with a clear typographic scale; headings are bold and significantly larger than the body text. The layout is spacious and centered, with generous white space that creates a sense of focus and calm.

Interactive elements like buttons and links have subtle, non-color-based hover and active states. Interactivity is communicated through gentle shadows, slight upward shifts, or by inverting the black-and-white scheme for a striking effect on active selection. Borders are minimal, thin, and light gray, used only to softly define containers without distracting from the content.

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
3. Add Tailwind styling using utility classes
4. Add Preact components in `workshop.js` if needed
5. Update routing in `main.go`

### CSS Development

- Use Tailwind utility classes for all styling
- Follow the established black/white design tokens
- Test responsive behavior with Tailwind breakpoints
- Ensure mobile responsiveness with mobile-first approach

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

- Always use Tailwind utility classes for styling
- Custom theme values are defined in `static/css/input.css` with `@theme` directive
- Use semantic color names: `bg-bg`, `text-fg`, `border-border`
- Leverage Tailwind's spacing, typography, and responsive utilities

### JavaScript Complexity

- **Preact Architecture**: Workshop interface uses Preact functional components with hooks
- **HTM Templates**: JSX-like syntax using htm template literals, no build step required
- **Component State**: useState, useEffect, useRef hooks for component lifecycle
- **Event Handling**: Preact event handlers with proper context binding
- **Modular Design**: Single workshop.js file with well-organized component structure

### Build Process

- Use `make dev` for development with hot reload
- Use `make clean && make build` for production builds
- Test builds before deployment

## File Naming Conventions

- CSS components: Use Tailwind utility classes directly in templates and JavaScript
- Templates: `purpose.html` in `templates/`
- Go packages: lowercase, descriptive names
- JavaScript: Semantic function and component names in Preact style

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
3. Style with Tailwind utility classes in the template
4. Add Preact components to `workshop.js` if interactive features needed
5. Add route in `main.go`

### Styling Updates

1. Use Tailwind utility classes following the established design system
2. Leverage Tailwind's responsive utilities for mobile compatibility
3. Follow the black/white color scheme strictly
4. Use flat hover effects with `hover:bg-fg hover:text-bg`
5. Test across different screen sizes

### JavaScript Features

1. **ES6 Modules**: Use import/export for all new code with clear module boundaries
2. **Preact Components**: Workshop interface uses Preact functional components with hooks
3. **HTM Templates**: JSX-like syntax using htm template literals, no build step required
4. **Modern JavaScript**: async/await, destructuring, optional chaining, arrow functions
5. **Event Handling**: Preact event handlers with proper context binding
6. **Error Boundaries**: Comprehensive try/catch with user-friendly error messages
7. **Request Cancellation**: AbortController for cancelling in-flight requests
8. **Accessibility**: Focus management, keyboard navigation, ARIA attributes

### JavaScript Development Tasks

1. **Adding New Components**: Create Preact functional components in `workshop.js` using hooks
2. **Updating Existing Features**: Modify component logic while maintaining hook dependencies
3. **Event Handling**: Use Preact event handlers with proper state management
4. **Error Handling**: Implement try/catch with toast notifications for user-friendly feedback
5. **State Management**: Use useState, useEffect for component lifecycle and data flow

## Preact Workshop Architecture

The workshop interface is built as a single-page application using Preact with HTM (no build step required).

### Component Organization

- **`workshop.js`**: Complete Preact application with all components and state management
- **Components**: Toast, Loading, Modals, ArtworkCard, WorkshopApp (main component)
- **State Management**: React hooks (useState, useEffect, useRef) for component state
- **Event Handling**: Preact event handlers with proper context and state updates

### Architecture Patterns

- **Functional Components**: All components use function syntax with hooks
- **HTM Templates**: JSX-like syntax using template literals, no compilation needed
- **Component State**: Local state with useState, side effects with useEffect
- **Props-based Communication**: Parent-child communication through props and callbacks
- **Error Boundaries**: Try/catch with toast notifications for user feedback

### Modern Patterns

```javascript
// Preact functional component with hooks
const ArtworkCard = ({ artwork, onRegenerate, onConfigure }) => {
  const [isGenerating, setIsGenerating] = useState(false);

  const handleGenerate = async () => {
    setIsGenerating(true);
    try {
      await onRegenerate(artwork.id);
    } catch (error) {
      showToast(`Generation failed: ${error.message}`, "error");
    } finally {
      setIsGenerating(false);
    }
  };

  return html`
    <div class="border border-border ${isGenerating ? "opacity-60" : ""}">
      <button onClick=${handleGenerate} disabled=${isGenerating}>Generate</button>
    </div>
  `;
};
```

### Event Handling

- Use Preact event handlers (onClick, onInput, onSubmit)
- State updates trigger automatic re-rendering
- Proper context binding with arrow functions
- Form handling with controlled components

### Error Management

- Comprehensive error boundaries with try/catch
- User-friendly toast notifications for errors
- Console logging for debugging
- Graceful degradation for network failures

## Inspiration & Context

This project was directly inspired by [Simon Willison's LLM SVG benchmark](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/), where he tested 16 different AI models on their ability to generate SVG code for "a pelican riding a bicycle." Key insights from this experiment:

- **Creative Prompt Selection**: The pelican bicycle prompt was chosen because it's unlikely to exist in training data
- **Model Diversity**: Tested models from OpenAI, Anthropic, Google Gemini, Meta's Llama, and others
- **Quality Variations**: Different models produce significantly different SVG quality and styles
- **Accessibility**: Demonstrated that AI can create unique, creative SVG art beyond what's in training data

Pelican Art Gallery extends this concept by providing a user-friendly interface for exploring AI-generated SVG art with any creative prompt, not just pelicans on bicycles.
