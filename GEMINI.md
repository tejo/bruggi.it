# Bruggi - Static Website Project

## Project Overview

"Bruggi - Villaggio di Montagna" is a static website for a fictional or real mountain village named Bruggi. It showcases the village's attractions, itineraries, galleries, and contact information.

**Key Technologies:**
*   **Language:** Go (Golang)
*   **Generator:** Custom Static Site Generator (SSG) in `main.go`
*   **Templating:** [Pongo2](https://github.com/flosch/pongo2) (Django/Jinja2-like syntax)
*   **Configuration:** TOML files (`content/`) for data and localization
*   **Styling:** Tailwind CSS (via local script in `static/js/tailwindcss.js`) and custom CSS.
*   **Maps:** Leaflet.js with OpenTopoMap tiles.
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
*   `static/`: Static assets copied to `dist/` during build.
    *   `css/`: Stylesheets (including `fonts.css` and `leaflet.css`).
    *   `fonts/`: Local font files.
    *   `gpx/`: GPX tracks for itineraries.
    *   `img/`: Images for the site.
    *   `js/`: Client-side scripts (`main.js`, `leaflet.js`, `leaflet-gpx.js`).
*   `dist/`: The generated output directory (Git ignored recommended).
*   `legacy_html/`: Unused HTML files from the previous static version.

## Key Features

### Itineraries
*   **Filtering:** Itineraries can be filtered by type (`hiking` or `biking`). This is implemented via **static page generation** (e.g., `dist/itineraries/hiking.html`).
*   **Details:** Each itinerary page supports:
    *   **Interactive Map:** Leaflet map visualizing the GPX track.
    *   **GPX Download:** Button to download the associated `.gpx` file.
    *   **YouTube Video:** Embedded video below the map.
    *   **Gallery:** Grid of additional images specific to the itinerary.
*   **Data Structure:** Defined in TOML files with fields for `type`, `gpx_file`, `youtube_video_id`, and `gallery`.

### Localization
The site supports multiple locales (currently `it` and `en`).
*   **Italian (Default):** Generated at `dist/*.html`
*   **English:** Generated at `dist/en/*.html`
*   **Switcher:** The language switcher intelligently links to the *current page* in the alternate language.

### Assets
*   **Fonts:** Served locally from `static/fonts/` (no Google Fonts CDN).
*   **Paths:** Asset paths in TOML (e.g., `img/photo.jpg`) are relative to `static/`. The generator automatically prepends `/static/` during build.

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