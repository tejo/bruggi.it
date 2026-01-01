# Bruggi - Static Website Project

## Project Overview

"Bruggi - Villaggio di Montagna" is a static website for a fictional or real mountain village named Bruggi. It showcases the village's attractions, itineraries, galleries, webcams, and contact information.

**Key Technologies:**
*   **Language:** Go (Golang)
*   **Generator:** Custom Static Site Generator (SSG) in `main.go`
*   **Templating:** [Pongo2](https://github.com/flosch/pongo2) (Django/Jinja2-like syntax)
*   **Image Processing:** [imaging](https://github.com/disintegration/imaging) for resizing and thumbnail generation.
*   **Configuration:** TOML files (`content/`) for data and localization.
*   **Styling:** Tailwind CSS (via local script in `static/js/tailwindcss.js`) and custom CSS.
*   **Maps:** Leaflet.js with OpenTopoMap tiles.
*   **Weather:** Real-time data via Open-Meteo API.
*   **Output:** Static HTML files generated in `dist/`.
*   **Watcher:** `fsnotify` for auto-rebuilding during development.

## Directory Structure

*   `main.go`: The core generator logic.
*   `Makefile`: Build automation commands.
*   `content/`: TOML data files defining the site's content.
    *   `index.toml`: Homepage content, navigation, and webcam localization.
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
    *   `css/`: Stylesheets (`fonts.css`, `leaflet.css`, `lightbox.css`).
    *   `fonts/`: Local font files.
    *   `gpx/`: GPX tracks for itineraries.
    *   `img/`: High-resolution images for the site.
    *   `thumbs/`: Auto-generated thumbnails (do not edit manually).
    *   `webcam/`: Webcam history images.
    *   `js/`: Client-side scripts (`main.js`, `leaflet.js`, `lightbox.js`, etc.).
*   `dist/`: The generated output directory (Git ignored).

## Key Features

### Image Management
*   **Auto-Thumbnails:** The generator automatically creates optimized thumbnails for images referenced in TOML files, storing them in `static/thumbs/`.
*   **Lightbox:** A custom JS/CSS lightbox allows users to view high-resolution images by clicking on thumbnails in galleries and itineraries.
*   **Cleanup:** The build process automatically removes unused images and thumbnails from the `static` folder to keep the project clean.

### Webcam & Weather
*   **Live View:** Displays the latest image from `static/webcam/current.jpg`.
*   **Time-lapse:** A client-side player cycles through historical images stored in `static/webcam/`.
*   **Real-time Weather:** Fetches live temperature, wind, and visibility data for Bruggi (lat/lon: 44.71143, 9.18697) using the Open-Meteo API.
*   **Update Tool:** A dedicated flag `-update-webcam` allows easy updating of the current view and history without a full site rebuild.

### Itineraries
*   **Filtering:** Static pages generated for `hiking` and `biking` types.
*   **Details:** Includes interactive Leaflet maps (GPX tracks), elevation profiles, YouTube embeds, and photo galleries.

### Localization
*   **Languages:** Italian (`dist/*.html`) and English (`dist/en/*.html`).
*   **Smart Switching:** Language switcher links preserve the current page context.

## Building and Running

1.  **Development Mode:**
    Build, watch for changes, and serve at `http://localhost:8080`.
    ```bash
    make serve
    # OR
    go run main.go -serve
    ```

2.  **Build Static Site:**
    Generate the `dist/` folder.
    ```bash
    make build
    # OR
    go run main.go
    ```

3.  **Build for Raspberry Pi (ARM64):**
    Generate a binary for ARM64 Linux in `bin/`.
    ```bash
    make build-arm
    ```

4.  **Update Webcam:**
    Add a new webcam image (updates `current.jpg`, adds a timestamped copy, and refreshes the webcam page).
    ```bash
    go run main.go -update-webcam /path/to/new/image.jpg
    # OR using the binary
    ./bin/bruggi -update-webcam /path/to/new/image.jpg
    ```
