# Bruggi - Static Website Project

## Project Overview

"Bruggi - Villaggio di Montagna" is a static website for a fictional or real mountain village named Bruggi. It showcases the village's attractions, itineraries, galleries, and contact information.

**Key Technologies:**
*   **Language:** Go (Golang)
*   **Generator:** Custom Static Site Generator (SSG) in `main.go`
*   **Templating:** [Pongo2](https://github.com/flosch/pongo2) (Django/Jinja2-like syntax)
*   **Configuration:** TOML files (`content/`) for data and localization
*   **Styling:** Tailwind CSS (via CDN in templates)
*   **Output:** Static HTML files generated in `dist/`
*   **Watcher:** `fsnotify` for auto-rebuilding during development.

## Directory Structure

*   `main.go`: The core generator logic.
*   `content/`: TOML data files defining the site's content.
    *   `index.toml`: Homepage content and navigation.
    *   `galleries.toml`: Photo collection.
    *   `itineraries/*.toml`: Individual itinerary definitions.
*   `templates/`: Pongo2 HTML templates.
    *   `base.html`: Shared layout (Header/Footer).
    *   `index.html`: Homepage template.
    *   `itinerary_list.html`: List of itineraries.
    *   `itinerary_detail.html`: Detail view for a single itinerary.
    *   `gallery.html`: Photo gallery page.
    *   `webcam.html`, `contacts.html`: Other page templates.
*   `static/`: Static assets (JS, CSS, images) copied to `dist/` during build.
    *   `js/main.js`: Client-side scripts.
*   `dist/`: The generated output directory (Git ignored recommended).
*   `legacy_html/`: Unused HTML files from the previous static version.

## Building and Running

1.  **Generate the Site:**
    Run the Go program to build the static files into the `dist/` directory.
    ```bash
    go run main.go
    ```

2.  **Watch and Serve (Development Mode):**
    Build the site, watch for file changes, and serve at `http://localhost:8080`.
    ```bash
    go run main.go -serve
    ```

## Localization
The site supports multiple locales (currently `it` and `en`).
*   **Italian (Default):** Generated at `dist/*.html`
*   **English:** Generated at `dist/en/*.html`

Data in TOML files is structured with shared fields at the root and localized fields in `[it]` and `[en]` sections. The navigation bar is also fully localized via `index.toml`.
