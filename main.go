package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/pelletier/go-toml/v2"
)

// Data Structures

type IndexData struct {
	It IndexLocale `toml:"it"`
	En IndexLocale `toml:"en"`
}

type IndexLocale struct {
	Hero     HeroSection    `toml:"hero"`
	Welcome  WelcomeSection `toml:"welcome"`
	Sections SectionTitles  `toml:"sections"`
}

type HeroSection struct {
	Title    string `toml:"title"`
	Subtitle string `toml:"subtitle"`
	CTA      string `toml:"cta"`
	Image    string `toml:"image"`
}

type WelcomeSection struct {
	Title       string `toml:"title"`
	Subtitle    string `toml:"subtitle"`
	Description string `toml:"description"`
	Image       string `toml:"image"`
	Altitude    string `toml:"altitude"`
	Founded     string `toml:"founded"`
	CTAHistory  string `toml:"cta_history"`
}

type SectionTitles struct {
	ItinerariesTitle    string `toml:"itineraries_title"`
	ItinerariesSubtitle string `toml:"itineraries_subtitle"`
	SeeAllItineraries   string `toml:"see_all_itineraries"`
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
	Image         string          `toml:"image"`
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
	Image         string
	Difficulty    string
	DistanceKM    float64
	Duration      string
	ElevationGain int
	Title         string
	Description   string
	LongDesc      string
	Tags          []string
}

func main() {
	fmt.Println("Starting Bruggi Static Site Generator...")

	// 1. Load Data
	indexData, err := loadIndex("content/index.toml")
	if err != nil {
		panic(err)
	}

	galleryData, err := loadGallery("content/galleries.toml")
	if err != nil {
		panic(err)
	}

	itineraries, err := loadItineraries("content/itineraries")
	if err != nil {
		panic(err)
	}

	// 2. Prepare Output Directory
	if err := os.RemoveAll("dist"); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("dist", 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("dist/en", 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("dist/static", 0755); err != nil {
		panic(err)
	}

	// Copy Static Files
	copyDir("static", "dist/static")

	// 3. Render Pages for IT (Default)
	fmt.Println("Rendering IT locale...")
	renderLocale("it", "", indexData.It, *galleryData, itineraries)

	// 4. Render Pages for EN
	fmt.Println("Rendering EN locale...")
	renderLocale("en", "/en", indexData.En, *galleryData, itineraries)

	fmt.Println("Build Complete!")
}

func renderLocale(locale string, baseUrl string, indexT IndexLocale, galleryT GalleryData, rawItineraries []ItineraryFile) {
	// Prepare Itineraries for this locale
	var localItineraries []RenderItinerary
	for _, raw := range rawItineraries {
		l := raw.It
		if locale == "en" {
			l = raw.En
		}
		localItineraries = append(localItineraries, RenderItinerary{
			Slug:          raw.Slug,
			Image:         raw.Image,
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
		"page_title":     indexT.Hero.Title,
		"t":              indexT,
		"gallery_images": galleryT.Images, // Take first 4 for homepage?
		"itineraries":    localItineraries, // Maybe limit number in template or here
	}

	tpl := pongo2.Must(pongo2.FromFile("templates/index.html"))
	outPath := "dist/index.html"
	if locale == "en" {
		outPath = "dist/en/index.html"
	}
	
	err := renderToFile(tpl, ctx, outPath)
	if err != nil {
		panic(err)
	}

	// Render Galleries
	galleryCtx := pongo2.Context{
		"locale":         locale,
		"base_url":       baseUrl,
		"page_title":     indexT.Sections.GalleryTitle,
		"t":              indexT,
		"gallery_images": galleryT.Images,
	}
	galTpl := pongo2.Must(pongo2.FromFile("templates/gallery.html"))
	galOutPath := "dist/galleries.html"
	if locale == "en" {
		galOutPath = "dist/en/galleries.html"
	}
	if err := renderToFile(galTpl, galleryCtx, galOutPath); err != nil {
		panic(err)
	}

	// Render Webcam
	webcamCtx := pongo2.Context{
		"locale":     locale,
		"base_url":   baseUrl,
		"page_title": "Bruggi Webcams",
		"t":          indexT,
	}
	webcamTpl := pongo2.Must(pongo2.FromFile("templates/webcam.html"))
	webcamOutPath := "dist/webcam.html"
	if locale == "en" {
		webcamOutPath = "dist/en/webcam.html"
	}
	if err := renderToFile(webcamTpl, webcamCtx, webcamOutPath); err != nil {
		panic(err)
	}

	// Render Contacts
	contactsCtx := pongo2.Context{
		"locale":     locale,
		"base_url":   baseUrl,
		"page_title": "Contattaci", // TODO: Translate
		"t":          indexT,
	}
	contactsTpl := pongo2.Must(pongo2.FromFile("templates/contacts.html"))
	contactsOutPath := "dist/contacts.html"
	if locale == "en" {
		contactsOutPath = "dist/en/contacts.html"
	}
	if err := renderToFile(contactsTpl, contactsCtx, contactsOutPath); err != nil {
		panic(err)
	}

	// Render Itineraries List

	// Render Itineraries List
	listCtx := pongo2.Context{
		"locale":      locale,
		"base_url":    baseUrl,
		"page_title":  indexT.Sections.ItinerariesTitle, // Or a specific string
		"t":           indexT,
		"itineraries": localItineraries,
	}
	listTpl := pongo2.Must(pongo2.FromFile("templates/itinerary_list.html"))
	listOutPath := "dist/itineraries.html"
	if locale == "en" {
		listOutPath = "dist/en/itineraries.html"
	}
	if err := renderToFile(listTpl, listCtx, listOutPath); err != nil {
		panic(err)
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
		panic(err)
	}

	for _, it := range localItineraries {
		detailCtx := pongo2.Context{
			"locale":     locale,
			"base_url":   baseUrl,
			"page_title": it.Title,
			"itinerary":  it,
			"t":          indexT, // Pass main translations if needed for header/footer
		}
		detailOutPath := filepath.Join(itineraryOutDir, it.Slug+".html")
		if err := renderToFile(detailTpl, detailCtx, detailOutPath); err != nil {
			panic(err)
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

func loadIndex(path string) (*IndexData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data IndexData
	if err := toml.Unmarshal(b, &data); err != nil {
		return nil, err
	}
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
