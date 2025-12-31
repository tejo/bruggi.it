# Bruggi - Static Website Project

## Project Overview

"Bruggi - Villaggio di Montagna" is a static website for a fictional or real mountain village named Bruggi. It showcases the village's attractions, itineraries, galleries, and contact information.

**Key Technologies:**
*   **HTML5:** Semantic markup for structure.
*   **Tailwind CSS:** Used via CDN for styling, including a custom in-browser configuration for colors, fonts, and dark mode.
*   **Vanilla JavaScript:** `main.js` handles client-side dynamic behavior, primarily the injection of shared layout components (header/footer).
*   **External Assets:** Fonts (Google Fonts), Icons (Material Symbols), and Images are hosted externally.

## Directory Structure

*   `index.html`: The main landing page.
*   `itineraries.html`, `galleries.html`, `contacts.html`: Main content pages.
*   `itinerary.html`: Detail view for a specific itinerary.
*   `webcam.html`: Page likely for viewing a live webcam feed.
*   `layout.html`: Contains the source code for the shared Header and Footer.
*   `main.js`: The core script that orchestrates the layout injection and interactive elements (mobile menu, active link highlighting).

## Building and Running

Since this is a static site without a build process (no `package.json`, `npm`, or bundlers):

1.  **Run:** Open any `.html` file (e.g., `index.html`) directly in a web browser.
    *   *Note:* Because `main.js` uses `fetch()` to load `layout.html`, you **must** serve the files via a local web server to avoid CORS errors (browsers often block `fetch` on `file://` protocol).
    *   **Recommended:** Use a simple HTTP server.
        *   Python: `python3 -m http.server`
        *   Node.js: `npx serve .` or `npx http-server .`
        *   VS Code: "Live Server" extension.

## Development Conventions

### Shared Layout Pattern
This project avoids code duplication for the header and footer using a client-side injection technique:
1.  **Source:** `layout.html` defines elements with IDs `main-header` and `main-footer`.
2.  **Placeholders:** Content pages (like `index.html`) include `<div id="layout-header"></div>` and `<div id="layout-footer"></div>`.
3.  **Injection:** `main.js` fetches `layout.html`, parses the response, and replaces the placeholders with the actual content from the source.

### Styling
*   **Tailwind via CDN:** Styling is applied using utility classes. Configuration (colors, fonts) is defined in a `<script>` tag within the `<head>` of each page.
*   **Themes:** Supports Light and Dark modes via the `class="light"` (or dark) attribute on the `<html>` tag.

### Navigation
*   Links in the header should match the filenames (e.g., `itineraries.html`).
*   `main.js` automatically applies the `text-primary` class to the navigation link that matches the current URL path.
