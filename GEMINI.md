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

## Directory Structure

*   `main.go`: The core generator logic.
*   `content/`: TOML data files defining the site's content.
    *   `index.toml`: Homepage content.
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

## Building and Running

1.  **Generate the Site:**
    Run the Go program to build the static files into the `dist/` directory.
    ```bash
    go run main.go
    ```

2.  **Serve the Site:**
    Serve the `dist/` directory using any static file server.
    *   **Python:** `python3 -m http.server -d dist`
    *   **Node.js:** `npx serve dist`
    *   **Go:** `go run github.com/jessvdk/go-static@latest -d dist`

## Localization
The site supports multiple locales (currently `it` and `en`).
*   **Italian (Default):** Generated at `dist/*.html`
*   **English:** Generated at `dist/en/*.html`

Data in TOML files is structured with `[it]` and `[en]` sections for translation.