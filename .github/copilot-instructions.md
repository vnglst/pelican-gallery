# Pelican Art Gallery — Copilot Instructions

## Architecture

- Go backend (custom routing, HTML templates)
- Frontend: Preact (HTM, ES6 modules), Tailwind CSS v4 (utility-first, black/white only)
- SQLite database (artwork groups, artworks)

## Coding Conventions

- Go: idiomatic, explicit errors, thin handlers, business logic in packages
- CSS: Tailwind utility classes only, strict black/white theme
- JS: ES6 modules, Preact functional components, hooks, HTM templates

## Design Principles

- Minimalist, high-contrast, black/white palette
- Responsive, accessible, semantic HTML
- No color accents—use bold, shadow, transform for emphasis

## File Structure

- `main.go`, `internal/`, `templates/`, `static/`, `config/`, `bin/`

## Common Tasks

- Add features: Go handler → template → Tailwind CSS → JS module/component
- Use Tailwind for all styling
- Use Go templates for shared components (partials)

## Inspiration

- Based on Simon Willison’s LLM SVG benchmark (“pelican on a bicycle”)
