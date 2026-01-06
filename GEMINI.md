# Bruggi - Static Website Project

## Project Overview

"Bruggi - Villaggio di Montagna" is a static website for a mountain village named Bruggi. It showcases the village's attractions, itineraries, galleries, webcams, and contact information.

**Key Technologies:**
*   **Language:** Go (Golang)
*   **Generator:** Custom Static Site Generator (SSG) in `main.go`
*   **Templating:** [Pongo2](https://github.com/flosch/pongo2) (Django/Jinja2-like syntax)
*   **Image Processing:** [imaging](https://github.com/disintegration/imaging) for resizing and thumbnail generation.
*   **Configuration:** TOML files (`content/`) for data and localization.
*   **Styling:** Tailwind CSS (via local script in `static/js/tailwindcss.js`) and custom CSS.
*   **Maps:** Leaflet.js with OpenTopoMap tiles.
*   **Weather:** Real-time data via Open-Meteo API.
*   **CI/CD:** GitHub Actions for automated builds and deployment to the repository.
*   **Output:** Static HTML files generated in `dist/`.
*   **Watcher:** `fsnotify` for auto-rebuilding during development.

## Directory Structure

*   `main.go`: The core generator logic.
*   `Makefile`: Build automation commands.
*   `content/`: TOML data files defining the site's content.
    *   `index.toml`: Homepage content, navigation, and webcam localization.
    *   `history.toml`: Localized content for the history and legends page.
    *   `august_events.toml`: List of events for the village festival.
    *   `galleries/`: Multiple gallery definitions (e.g., `main.toml`, `nature.toml`).
    *   `itineraries/*.toml`: Individual itinerary definitions.
*   `templates/`: Pongo2 HTML templates.
    *   `base.html`: Shared layout (Header/Footer).
    *   `index.html`: Homepage template.
    *   `history.html`: History and legends page.
    *   `itinerary_list.html`: List of itineraries.
    *   `itinerary_detail.html`: Detail view for a single itinerary.
    *   `gallery_list.html`, `gallery_detail.html`: Multi-gallery templates.
    *   `webcam.html`, `contacts.html`: Other page templates.
*   `static/`: Static assets copied to `dist/` during build.
    *   `css/`: Stylesheets (`fonts.css`, `leaflet.css`, `lightbox.css`).
    *   `fonts/`: Local font files.
    *   `gpx/`: GPX tracks for itineraries.
    *   `img/`: High-resolution images for the site.
    *   `thumbs/`: Auto-generated thumbnails (do not edit manually).
    *   `webcam/`: Webcam history images.
    *   `js/`: Client-side scripts (`main.js`, `leaflet.js`, `lightbox.js`, `slideshow.js`, etc.).
*   `dist/`: The generated output directory (tracked in Git for deployment).

## Key Features

### Image Management
*   **Auto-Thumbnails:** The generator automatically creates optimized thumbnails for images referenced in TOML files, storing them in `static/thumbs/`.
*   **Slideshow:** A hero slideshow on the homepage (via `slideshow.js`) cycles through multiple images with crossfade effects.
*   **Lightbox:** A custom JS/CSS lightbox allows users to view high-resolution images by clicking on thumbnails in galleries and itineraries.
*   **Cleanup:** The build process automatically removes unused images and thumbnails while strictly preserving all referenced assets (including history hero images and contact photos).

### History & Legends
*   **Content:** Dedicated page showcasing the village's history and folk stories, managed via `history.toml`.
*   **Formatting:** Supports basic Markdown-like syntax (bold and italics) within TOML strings, processed during generation.
*   **UX:** Uses an accordion-style interface (HTML `<details>`) for better readability.

### Events
*   **Dynamic Visibility:** Events for August can be enabled or disabled globally via `august_events.toml`.
*   **Integration:** When enabled, events are automatically displayed on the homepage.

### Webcam & Weather
*   **Live View:** Displays the latest timestamped image from `static/webcam/`.
*   **Static Timestamps:** The timestamp is extracted directly from the filename (e.g., `YYYY-MM-DD_HH-MM-SS.jpg`).
*   **Improved Time-lapse:** 
    *   Plays in reverse (Newest -> Oldest) starting from the latest image.
    *   On-demand preloading with a visual loader inside the button.
    *   Synchronized timestamps that update for each frame during playback.
*   **Real-time Weather:** Fetches live temperature, wind, and visibility data for Bruggi using the Open-Meteo API.

### Itineraries
*   **Filtering:** Static pages generated for `hiking` and `biking` types.
*   **Automatic Stats:** The generator automatically parses GPX files to calculate elevation gain and distance.
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
    Generate the `dist/` folder and perform image cleanup.
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



4.  **Deployment:**
    The project uses GitHub Actions to automatically rebuild and commit the `dist/` folder whenever changes are pushed to `main`.
