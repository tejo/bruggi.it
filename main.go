package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml/v2"
)

// Data Structures

type IndexFile struct {
	Hero    SharedHeroSection    `toml:"hero"`
	Welcome SharedWelcomeSection `toml:"welcome"`
	It      IndexLocale          `toml:"it"`
	En      IndexLocale          `toml:"en"`
}

type SharedHeroSection struct {
	Image string `toml:"image"`
}

type SharedWelcomeSection struct {
	Image string `toml:"image"`
}

type IndexLocale struct {
	Nav      NavLocale      `toml:"nav"`
	Hero     HeroLocale     `toml:"hero"`
	Welcome  WelcomeLocale  `toml:"welcome"`
	Sections SectionTitles  `toml:"sections"`
}

type NavLocale struct {
	Home        string `toml:"home"`
	Itineraries string `toml:"itineraries"`
	Webcam      string `toml:"webcam"`
	Gallery     string `toml:"gallery"`
	Contact     string `toml:"contact"`
}

type HeroLocale struct {
	Title    string `toml:"title"`
	Subtitle string `toml:"subtitle"`
	CTA      string `toml:"cta"`
}

type WelcomeLocale struct {
	Title       string `toml:"title"`
	Subtitle    string `toml:"subtitle"`
	Description string `toml:"description"`
	Altitude    string `toml:"altitude"`
	Founded     string `toml:"founded"`
	CTAHistory  string `toml:"cta_history"`
}

type SectionTitles struct {
	ItinerariesTitle    string `toml:"itineraries_title"`
	ItinerariesSubtitle string `toml:"itineraries_subtitle"`
	SeeAllItineraries   string `toml:"see_all_itineraries"`
	ReadMore            string `toml:"read_more"`
	FilterAll           string `toml:"filter_all"`
	FilterHiking        string `toml:"filter_hiking"`
	FilterBiking        string `toml:"filter_biking"`
	GalleryTitle        string `toml:"gallery_title"`
	GallerySubtitle     string `toml:"gallery_subtitle"`
	SeeAllGallery       string `toml:"see_all_gallery"`
}

type GalleryData struct {
	Images []GalleryImage `toml:"images"`
}

type GalleryImage struct {
	Url string `toml:"url"`
	Alt string `toml:"alt"`
}

type ItineraryFile struct {
	Slug          string          `toml:"slug"`
	Type          string          `toml:"type"`
	Image         string          `toml:"image"`
	GpxFile       string          `toml:"gpx_file"`
	Difficulty    string          `toml:"difficulty"`
	DistanceKM    float64         `toml:"distance_km"`
	Duration      string          `toml:"duration"`
	ElevationGain int             `toml:"elevation_gain"`
	It            ItineraryLocale `toml:"it"`
	En            ItineraryLocale `toml:"en"`
}

type ItineraryLocale struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	LongDesc    string   `toml:"long_description"`
	Tags        []string `toml:"tags"`
}

// Renderable Item for Templates
type RenderItinerary struct {
	Slug          string
	Type          string
	Image         string
	GpxFile       string
	Difficulty    string
	DistanceKM    float64
	Duration      string
	ElevationGain int
	Title         string
	Description   string
	LongDesc      string
	Tags          []string
}

// Helper struct to pass to templates, flattening the structure
type RenderIndex struct {
	Nav      RenderNav
	Hero     RenderHero
	Welcome  RenderWelcome
	Sections SectionTitles
}

type RenderNav struct {
	Home        string
	Itineraries string
	Webcam      string
	Gallery     string
	Contact     string
}

type RenderHero struct {
	Title    string
	Subtitle string
	CTA      string
	Image    string
}

type RenderWelcome struct {
	Title       string
	Subtitle    string
	Description string
	Image       string
	Altitude    string
	Founded     string
	CTAHistory  string
}

func main() {
	serveMode := flag.Bool("serve", false, "Watch for changes and serve the site")
	flag.Parse()

	if *serveMode {
		watchAndServe()
	} else {
		buildSite()
	}
}

func buildSite() {
	fmt.Println("Building site...")
	start := time.Now()

	// 1. Load Data
	indexData, err := loadIndex("content/index.toml")
	if err != nil {
		log.Printf("Error loading index: %v", err)
		return
	}

	galleryData, err := loadGallery("content/galleries.toml")
	if err != nil {
		log.Printf("Error loading gallery: %v", err)
		return
	}

	itineraries, err := loadItineraries("content/itineraries")
	if err != nil {
		log.Printf("Error loading itineraries: %v", err)
		return
	}

	// 2. Prepare Output Directory
	if err := os.RemoveAll("dist"); err != nil {
		log.Printf("Error clearing dist: %v", err)
		return
	}
	if err := os.MkdirAll("dist", 0755); err != nil {
		log.Printf("Error creating dist: %v", err)
		return
	}
	if err := os.MkdirAll("dist/en", 0755); err != nil {
		log.Printf("Error creating dist/en: %v", err)
		return
	}
	if err := os.MkdirAll("dist/static", 0755); err != nil {
		log.Printf("Error creating dist/static: %v", err)
		return
	}

	// Copy Static Files
	copyDir("static", "dist/static")

	// 3. Render Pages for IT (Default)
	renderLocale("it", "", indexData, *galleryData, itineraries)

	// 4. Render Pages for EN
	renderLocale("en", "/en", indexData, *galleryData, itineraries)

	fmt.Printf("Build complete in %v\n", time.Since(start))
}

func watchAndServe() {
	// Initial build
	buildSite()

	// Watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Modified file:", event.Name)
					buildSite()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add directories to watch
	dirsToWatch := []string{"content", "content/itineraries", "templates", "static"}
	for _, dir := range dirsToWatch {
		err = watcher.Add(dir)
		if err != nil {
			log.Printf("Error watching %s: %v", dir, err) // Don't crash if dir doesn't exist yet
		}
	}
	// Also watch individual files in static/js/ etc if needed, but 'static' covers direct children.
	// Recursive watch is not built-in to fsnotify, but we have a flat structure mostly.
	// Add static/js explicitly if needed.
	watcher.Add("static/js")

	// Server
	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)

	log.Println("Serving on http://localhost:8080")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-done
}

func renderLocale(locale string, baseUrl string, indexData *IndexFile, galleryT GalleryData, rawItineraries []ItineraryFile) {
	// Pick the right locale data
	var l IndexLocale
	if locale == "it" {
		l = indexData.It
	} else {
		l = indexData.En
	}

	// Merge shared and localized
	renderIndex := RenderIndex{
		Nav: RenderNav{
			Home:        l.Nav.Home,
			Itineraries: l.Nav.Itineraries,
			Webcam:      l.Nav.Webcam,
			Gallery:     l.Nav.Gallery,
			Contact:     l.Nav.Contact,
		},
		Hero: RenderHero{
			Title:    l.Hero.Title,
			Subtitle: l.Hero.Subtitle,
			CTA:      l.Hero.CTA,
			Image:    indexData.Hero.Image,
		},
		Welcome: RenderWelcome{
			Title:       l.Welcome.Title,
			Subtitle:    l.Welcome.Subtitle,
			Description: l.Welcome.Description,
			Altitude:    l.Welcome.Altitude,
			Founded:     l.Welcome.Founded,
			CTAHistory:  l.Welcome.CTAHistory,
			Image:       indexData.Welcome.Image,
		},
		Sections: l.Sections,
	}

	// Prepare Itineraries for this locale
	var localItineraries []RenderItinerary
	for _, raw := range rawItineraries {
		l := raw.It
		if locale == "en" {
			l = raw.En
		}
		localItineraries = append(localItineraries, RenderItinerary{
			Slug:          raw.Slug,
			Type:          raw.Type,
			Image:         raw.Image,
			GpxFile:       raw.GpxFile,
			Difficulty:    raw.Difficulty,
			DistanceKM:    raw.DistanceKM,
			Duration:      raw.Duration,
			ElevationGain: raw.ElevationGain,
			Title:         l.Title,
			Description:   l.Description,
			LongDesc:      l.LongDesc,
			Tags:          l.Tags,
		})
	}

	// Render Index
	ctx := pongo2.Context{
		"locale":         locale,
		"base_url":       baseUrl,
		"alternate_url":  computeAlternateUrl(locale, "/"),
		"page_title":     renderIndex.Hero.Title,
		"t":              renderIndex, // We pass our flattened struct as 't'
		"gallery_images": galleryT.Images,
		"itineraries":    localItineraries,
	}

	tpl := pongo2.Must(pongo2.FromFile("templates/index.html"))
	outPath := "dist/index.html"
	if locale == "en" {
		outPath = "dist/en/index.html"
	}

	err := renderToFile(tpl, ctx, outPath)
	if err != nil {
		log.Panic(err)
	}

	// Render Galleries
	galleryCtx := pongo2.Context{
		"locale":         locale,
		"base_url":       baseUrl,
		"alternate_url":  computeAlternateUrl(locale, "/galleries.html"),
		"page_title":     renderIndex.Sections.GalleryTitle,
		"t":              renderIndex,
		"gallery_images": galleryT.Images,
	}
	galTpl := pongo2.Must(pongo2.FromFile("templates/gallery.html"))
	galOutPath := "dist/galleries.html"
	if locale == "en" {
		galOutPath = "dist/en/galleries.html"
	}
	if err := renderToFile(galTpl, galleryCtx, galOutPath); err != nil {
		log.Panic(err)
	}

	// Render Webcam
	webcamCtx := pongo2.Context{
		"locale":        locale,
		"base_url":      baseUrl,
		"alternate_url": computeAlternateUrl(locale, "/webcam.html"),
		"page_title":    "Bruggi Webcams",
		"t":             renderIndex,
	}
	webcamTpl := pongo2.Must(pongo2.FromFile("templates/webcam.html"))
	webcamOutPath := "dist/webcam.html"
	if locale == "en" {
		webcamOutPath = "dist/en/webcam.html"
	}
	if err := renderToFile(webcamTpl, webcamCtx, webcamOutPath); err != nil {
		log.Panic(err)
	}

	// Render Contacts
	contactsCtx := pongo2.Context{
		"locale":        locale,
		"base_url":      baseUrl,
		"alternate_url": computeAlternateUrl(locale, "/contacts.html"),
		"page_title":    renderIndex.Nav.Contact,
		"t":             renderIndex,
	}
	contactsTpl := pongo2.Must(pongo2.FromFile("templates/contacts.html"))
	contactsOutPath := "dist/contacts.html"
	if locale == "en" {
		contactsOutPath = "dist/en/contacts.html"
	}
	if err := renderToFile(contactsTpl, contactsCtx, contactsOutPath); err != nil {
		log.Panic(err)
	}

	// Render Itineraries List (All + Filtered)
	filters := []string{"all", "hiking", "biking"}
	listTpl := pongo2.Must(pongo2.FromFile("templates/itinerary_list.html"))

	for _, filter := range filters {
		var filteredIts []RenderItinerary
		if filter == "all" {
			filteredIts = localItineraries
		} else {
			for _, it := range localItineraries {
				if it.Type == filter {
					filteredIts = append(filteredIts, it)
				}
			}
		}

		var relativePath string
		if filter == "all" {
			relativePath = "/itineraries.html"
		} else {
			relativePath = "/itineraries/" + filter + ".html"
		}

		listCtx := pongo2.Context{
			"locale":         locale,
			"base_url":       baseUrl,
			"alternate_url":  computeAlternateUrl(locale, relativePath),
			"page_title":     renderIndex.Sections.ItinerariesTitle,
			"t":              renderIndex,
			"itineraries":    filteredIts,
			"current_filter": filter,
		}

		var listOutPath string
		if filter == "all" {
			if locale == "en" {
				listOutPath = "dist/en/itineraries.html"
			} else {
				listOutPath = "dist/itineraries.html"
			}
		} else {
			// Ensure itineraries dir exists
			itineraryDir := "dist/itineraries"
			if locale == "en" {
				itineraryDir = "dist/en/itineraries"
			}
			if err := os.MkdirAll(itineraryDir, 0755); err != nil {
				log.Panic(err)
			}
			if locale == "en" {
				listOutPath = filepath.Join("dist/en/itineraries", filter+".html")
			} else {
				listOutPath = filepath.Join("dist/itineraries", filter+".html")
			}
		}

		if err := renderToFile(listTpl, listCtx, listOutPath); err != nil {
			log.Panic(err)
		}
	}

	// Render Itinerary Details
	detailTpl := pongo2.Must(pongo2.FromFile("templates/itinerary_detail.html"))

	// Create itineraries dir if not exists (for root/itineraries/...)
	// For EN, it will be dist/en/itineraries/...
	itineraryOutDir := "dist/itineraries"
	if locale == "en" {
		itineraryOutDir = "dist/en/itineraries"
	}
	if err := os.MkdirAll(itineraryOutDir, 0755); err != nil {
		log.Panic(err)
	}

	for _, it := range localItineraries {
		relativePath := "/itineraries/" + it.Slug + ".html"
		detailCtx := pongo2.Context{
			"locale":        locale,
			"base_url":      baseUrl,
			"alternate_url": computeAlternateUrl(locale, relativePath),
			"page_title":    it.Title,
			"itinerary":     it,
			"t":             renderIndex, // Pass main translations if needed for header/footer
		}
		detailOutPath := filepath.Join(itineraryOutDir, it.Slug+".html")
		if err := renderToFile(detailTpl, detailCtx, detailOutPath); err != nil {
			log.Panic(err)
		}
	}
}

func renderToFile(tpl *pongo2.Template, ctx pongo2.Context, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.ExecuteWriter(ctx, f)
}

func loadIndex(path string) (*IndexFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data IndexFile
	if err := toml.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	data.Hero.Image = "/static/" + data.Hero.Image
	data.Welcome.Image = "/static/" + data.Welcome.Image
	return &data, nil
}

func loadGallery(path string) (*GalleryData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data GalleryData
	if err := toml.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	for i := range data.Images {
		data.Images[i].Url = "/static/" + data.Images[i].Url
	}
	return &data, nil
}

func loadItineraries(dir string) ([]ItineraryFile, error) {
	var its []ItineraryFile
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".toml") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var it ItineraryFile
			if err := toml.Unmarshal(b, &it); err != nil {
				return err
			}
			it.Image = "/static/" + it.Image
			if it.GpxFile != "" {
				it.GpxFile = "/static/" + it.GpxFile
			}
			its = append(its, it)
		}
		return nil
	})
	return its, err
}

func copyDir(src string, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == src {
			return nil
		}
		rel, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dst, rel)
		
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}
		
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
}

func computeAlternateUrl(currentLocale string, relativePath string) string {
	if currentLocale == "it" {
		if relativePath == "/" {
			return "/en/"
		}
		return "/en" + relativePath
	} else {
		// current is en
		if relativePath == "/" {
			return "/"
		}
		return relativePath
	}
}