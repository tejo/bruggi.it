# Bruggi - Mountain Village Website

Welcome to the source code for the "Bruggi - Villaggio di Montagna" website. This project is a custom Static Site Generator (SSG) built with Go, designed to create a fast, localized, and feature-rich website for a mountain village.

## 🚀 Features

-   **Fast Static Generation:** Builds HTML from TOML content and Pongo2 templates.
-   **Localization:** Native support for Italian (IT) and English (EN).
-   **Image Optimization:** Automated thumbnail generation and unused image cleanup.
-   **Interactive Maps:** Leaflet.js integration for visualizing GPX tracks.
-   **Webcam & Weather:** Real-time weather data (Open-Meteo) and webcam time-lapse player.
-   **Responsive Design:** Styled with Tailwind CSS for mobile and desktop.

## 🛠️ Getting Started

### Prerequisites

*   [Go](https://go.dev/) (version 1.25+ recommended)
*   [Make](https://www.gnu.org/software/make/) (optional, for build commands)

### Installation

Clone the repository:
```bash
git clone https://github.com/your-username/bruggi.git
cd bruggi
```

Install dependencies:
```bash
go mod tidy
```

### Running Locally

To start the development server with live reloading:

```bash
make serve
```
Access the site at `http://localhost:8080`.

## 🏗️ Build Commands

| Command | Description |
| :--- | :--- |
| `make build` | Builds the static site into the `dist/` directory. |
| `make serve` | Runs the generator in watch mode, serving at `localhost:8080`. |
| `make build-arm` | Compiles the binary for Raspberry Pi (Linux ARM64). |
| `make clean` | Removes the `dist/` and `bin/` directories. |

## 📷 Webcam Updates

This project includes a built-in tool to manage webcam images. To update the "live" view and archive the previous image:

```bash
go run main.go -update-webcam /path/to/your/new_image.jpg
```

This command will:
1.  Copy the new image to `static/webcam/current.jpg`.
2.  Save a timestamped copy in `static/webcam/`.
3.  Regenerate only the webcam HTML pages to include the new image in the time-lapse history.

## 📂 Project Structure

*   **`content/`**: Edit TOML files here to change text, add itineraries, or update gallery images.
*   **`templates/`**: Modify HTML templates to change the site layout.
*   **`static/`**: Place raw assets here. Images in `static/img` are auto-processed.
*   **`main.go`**: The source code for the generator.

## 📄 License

This project is open source. Feel free to use and modify it.
