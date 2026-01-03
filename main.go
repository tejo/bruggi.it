package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/flosch/pongo2/v6"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml/v2"
)

// Data Structures

type IndexFile struct {
	Hero         SharedHeroSection        `toml:"hero"`
	Welcome      SharedWelcomeSection     `toml:"welcome"`
	Itineraries  SharedItinerariesSection `toml:"itineraries"`
	Contacts     SharedContacts           `toml:"contacts"`
	AugustEvents SharedAugustEvents       `toml:"august_events"`
	It           IndexLocale              `toml:"it"`
	En           IndexLocale              `toml:"en"`
}

type EventsFile struct {
	Enabled bool               `toml:"enabled"`
	It      AugustEventsLocale `toml:"it"`
	En      AugustEventsLocale `toml:"en"`
}

type SharedHeroSection struct {
	Images []string `toml:"images"`
}

type SharedWelcomeSection struct {
	Image string `toml:"image"`
}

type SharedAugustEvents struct {
	Enabled bool `toml:"enabled"`
}

type SharedItinerariesSection struct {
	HeroImage string `toml:"hero_image"`
}

type SharedContacts struct {
	Email   string `toml:"email"`
	Phone   string `toml:"phone"`
	Address string `toml:"address"`
}

type IndexLocale struct {
	Nav           NavLocale           `toml:"nav"`
	Hero          HeroLocale          `toml:"hero"`
	Welcome       WelcomeLocale       `toml:"welcome"`
	Sections      SectionTitles       `toml:"sections"`
	ItineraryPage ItineraryPageLocale `toml:"itinerary_page"`
	WebcamPage    WebcamPageLocale    `toml:"webcam_page"`
	ContactInfo   ContactInfoLocale   `toml:"contact_info"`
	AugustEvents  AugustEventsLocale  `toml:"august_events"`
	Footer        FooterLocale        `toml:"footer"`
}

type AugustEventsLocale struct {
	Title string      `toml:"title"`
	Items []EventItem `toml:"items"`
}

type EventItem struct {
	Name string `toml:"name"`
	Date string `toml:"date"`
	Time string `toml:"time"`
}

type ItineraryPageLocale struct {
	TrailDetails     string `toml:"trail_details"`
	Author           string `toml:"author"`
	Type             string `toml:"type"`
	TypeHiking       string `toml:"type_hiking"`
	TypeBiking       string `toml:"type_biking"`
	Duration         string `toml:"duration"`
	Distance         string `toml:"distance"`
	ElevationGain    string `toml:"elevation_gain"`
	DownloadGPX      string `toml:"download_gpx"`
	GPXNotAvailable  string `toml:"gpx_not_available"`
	Description      string `toml:"description"`
	Difficulty       string `toml:"difficulty"`
	DifficultyEasy   string `toml:"difficulty_easy"`
	DifficultyMedium string `toml:"difficulty_medium"`
	DifficultyHard   string `toml:"difficulty_hard"`
}

type ContactInfoLocale struct {
	Title        string `toml:"title"`
	Subtitle     string `toml:"subtitle"`
	EmailLabel   string `toml:"email_label"`
	PhoneLabel   string `toml:"phone_label"`
	AddressLabel string `toml:"address_label"`
	FormTitle    string `toml:"form_title"`
	FormName     string `toml:"form_name"`
	FormEmail    string `toml:"form_email"`
	FormMessage  string `toml:"form_message"`
	FormSubmit   string `toml:"form_submit"`
}

type WebcamPageLocale struct {
	Live            string `toml:"live"`
	HD              string `toml:"hd"`
	PanoramaTitle   string `toml:"panorama_title"`
	Location        string `toml:"location"`
	Share           string `toml:"share"`
	Snapshot        string `toml:"snapshot"`
	Timelapse       string `toml:"timelapse"`
	StatusOnline    string `toml:"status_online"`
	NextUpdate      string `toml:"next_update"`
	ReportIssue     string `toml:"report_issue"`
	ConditionsTitle string `toml:"conditions_title"`
	UpdatedAgo      string `toml:"updated_ago"`
	Temperature     string `toml:"temperature"`
	FeelsLike       string `toml:"feels_like"`
	Wind            string `toml:"wind"`
	Direction       string `toml:"direction"`
	Humidity        string `toml:"humidity"`
	Precip          string `toml:"precip"`
	Visibility      string `toml:"visibility"`
	VisRange        string `toml:"vis_range"`
	VisGood         string `toml:"vis_good"`
	VisPoor         string `toml:"vis_poor"`
	VisModerate     string `toml:"vis_moderate"`
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

type FooterLocale struct {
	Motto         string `toml:"motto"`
	ExploreTitle  string `toml:"explore_title"`
	ContactsTitle string `toml:"contacts_title"`
	Copyright     string `toml:"copyright"`
}

type GalleryData struct {
	Images []GalleryImage `toml:"images"`
}

type GalleryImage struct {
	Url       string `toml:"url"`
	Alt       string `toml:"alt"`
	Author    string `toml:"author"` // Instagram handle
	Thumbnail string // Populated during load
}

type ItineraryFile struct {
	Slug             string          `toml:"slug"`
	Type             string          `toml:"type"`
	Image            string          `toml:"image"`
	GpxFile          string          `toml:"gpx_file"`
	YoutubeVideoID   string          `toml:"youtube_video_id"`
	Gallery          []string        `toml:"gallery"`
	ProcessedGallery []GalleryImage  `toml:"-"`
	Difficulty       string          `toml:"difficulty"`
	DistanceKM       float64         `toml:"distance_km"`
	Duration         string          `toml:"duration"`
	ElevationGain    int             `toml:"elevation_gain"`
	Author           string          `toml:"author"` // Instagram handle
	It               ItineraryLocale `toml:"it"`
	En               ItineraryLocale `toml:"en"`
}

type ItineraryLocale struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	LongDesc    string   `toml:"long_description"`
	Tags        []string `toml:"tags"`
}

// Renderable Item for Templates
type RenderItinerary struct {
	Slug           string
	Type           string
	Image          string
	GpxFile        string
	YoutubeVideoID string
	Gallery        []GalleryImage
	Difficulty     string
	DistanceKM     float64
	Duration       string
	ElevationGain  int
	Author         string
	Title          string
	Description    string
	LongDesc       string
	Tags           []string
}

// Helper struct to pass to templates, flattening the structure
type RenderIndex struct {
	Nav           RenderNav
	Hero          RenderHero
	Welcome       RenderWelcome
	Itineraries   SharedItinerariesSection
	Sections      SectionTitles
	ItineraryPage ItineraryPageLocale
	WebcamPage    RenderWebcamPage
	Contacts      SharedContacts
	ContactInfo   ContactInfoLocale
	AugustEvents  RenderAugustEvents
	Footer        FooterLocale
}

type RenderAugustEvents struct {
	Enabled bool
	Title   string
	Items   []EventItem
}

type RenderWebcamPage struct {
	Live            string
	HD              string
	PanoramaTitle   string
	Location        string
	Share           string
	Snapshot        string
	Timelapse       string
	StatusOnline    string
	NextUpdate      string
	ReportIssue     string
	ConditionsTitle string
	UpdatedAgo      string
	Temperature     string
	FeelsLike       string
	Wind            string
	Direction       string
	Humidity        string
	Precip          string
	Visibility      string
	VisRange        string
	VisGood         string
	VisPoor         string
	VisModerate     string
	Images          []string
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
	Images   []string
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
	webcamUpdate := flag.String("update-webcam", "", "Path to new webcam image to add")
	flag.Parse()

	if *webcamUpdate != "" {
		handleWebcamUpdate(*webcamUpdate)
	} else if *serveMode {
		watchAndServe()
	} else {
		buildSite()
	}
}

func handleWebcamUpdate(srcPath string) {
	fmt.Printf("Updating webcam with image: %s\n", srcPath)

	// 1. Prepare Paths
	webcamDir := "static/webcam"
	if err := os.MkdirAll(webcamDir, 0755); err != nil {
		log.Fatalf("Error creating webcam dir: %v", err)
	}

	distWebcamDir := "dist/static/webcam"
	// Ensure dist exists (if not, we might be running this without a previous build,
	// but we try to support it)
	if err := os.MkdirAll(distWebcamDir, 0755); err != nil {
		log.Fatalf("Error creating dist webcam dir: %v", err)
	}

	// 2. Generate Filenames
	currentName := "current.jpg"

	now := time.Now()
	timestampName := fmt.Sprintf("%s.jpg", now.Format("2006-01-02_15-04-05"))

	// 3. Copy files to static/webcam (Source of Truth)
	if err := copyFile(srcPath, filepath.Join(webcamDir, currentName)); err != nil {
		log.Fatalf("Error updating current.jpg in static: %v", err)
	}
	if err := copyFile(srcPath, filepath.Join(webcamDir, timestampName)); err != nil {
		log.Fatalf("Error adding timestamped image in static: %v", err)
	}

	// 4. Copy files to dist/static/webcam (Served Content)
	if err := copyFile(srcPath, filepath.Join(distWebcamDir, currentName)); err != nil {
		log.Fatalf("Error updating current.jpg in dist: %v", err)
	}
	if err := copyFile(srcPath, filepath.Join(distWebcamDir, timestampName)); err != nil {
		log.Fatalf("Error adding timestamped image in dist: %v", err)
	}

	// 5. Update Pages
	indexData, err := loadIndex("content/index.toml")
	if err != nil {
		log.Fatalf("Error loading index: %v", err)
	}
	eventsData, err := loadEvents("content/august_events.toml")
	if err != nil {
		log.Fatalf("Error loading events: %v", err)
	}

	updateWebcamPages(indexData, eventsData)
	fmt.Println("Webcam update complete.")
}

func updateWebcamPages(indexData *IndexFile, eventsData *EventsFile) {
	// Re-render ONLY webcam.html for IT and EN

	webcamImages, err := loadWebcamImages("static/webcam")
	if err != nil {
		log.Printf("Error loading webcam images: %v", err)
	}

	render := func(locale string, baseUrl string, outPath string) {
		renderIndex := createRenderIndex(locale, indexData, eventsData)
		renderIndex.WebcamPage.Images = webcamImages

		ctx := pongo2.Context{
			"locale":        locale,
			"base_url":      baseUrl,
			"alternate_url": computeAlternateUrl(locale, "/webcam.html"),
			"page_title":    "Bruggi Webcams",
			"t":             renderIndex,
		}

		// Ensure output dir exists
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			log.Panic(err)
		}

		tpl := pongo2.Must(pongo2.FromFile("templates/webcam.html"))
		if err := renderToFile(tpl, ctx, outPath); err != nil {
			log.Panic(err)
		}
	}

	render("it", "", "dist/webcam.html")
	render("en", "/en", "dist/en/webcam.html")
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

	eventsData, err := loadEvents("content/august_events.toml")
	if err != nil {
		log.Printf("Error loading events: %v", err)
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
	renderLocale("it", "", indexData, eventsData, *galleryData, itineraries)

	// 4. Render Pages for EN
	renderLocale("en", "/en", indexData, eventsData, *galleryData, itineraries)

	// 5. Cleanup Unused Images
	// usedImages := collectUsedImages(indexData, galleryData, itineraries)
	// if err := cleanupImages(usedImages); err != nil {
	// 	log.Printf("Error cleaning up images: %v", err)
	// }

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

func createRenderIndex(locale string, indexData *IndexFile, eventsData *EventsFile) RenderIndex {
	var l IndexLocale
	var el AugustEventsLocale
	if locale == "it" {
		l = indexData.It
		el = eventsData.It
	} else {
		l = indexData.En
		el = eventsData.En
	}

	return RenderIndex{
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
			Images:   indexData.Hero.Images,
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
		Itineraries:   indexData.Itineraries,
		Sections:      l.Sections,
		ItineraryPage: l.ItineraryPage,
		Contacts:      indexData.Contacts,
		ContactInfo:   l.ContactInfo,
		AugustEvents: RenderAugustEvents{
			Enabled: eventsData.Enabled,
			Title:   el.Title,
			Items:   el.Items,
		},
		WebcamPage: RenderWebcamPage{
			Live:            l.WebcamPage.Live,
			HD:              l.WebcamPage.HD,
			PanoramaTitle:   l.WebcamPage.PanoramaTitle,
			Location:        l.WebcamPage.Location,
			Share:           l.WebcamPage.Share,
			Snapshot:        l.WebcamPage.Snapshot,
			Timelapse:       l.WebcamPage.Timelapse,
			StatusOnline:    l.WebcamPage.StatusOnline,
			NextUpdate:      l.WebcamPage.NextUpdate,
			ReportIssue:     l.WebcamPage.ReportIssue,
			ConditionsTitle: l.WebcamPage.ConditionsTitle,
			UpdatedAgo:      l.WebcamPage.UpdatedAgo,
			Temperature:     l.WebcamPage.Temperature,
			FeelsLike:       l.WebcamPage.FeelsLike,
			Wind:            l.WebcamPage.Wind,
			Direction:       l.WebcamPage.Direction,
			Humidity:        l.WebcamPage.Humidity,
			Precip:          l.WebcamPage.Precip,
			Visibility:      l.WebcamPage.Visibility,
			VisRange:        l.WebcamPage.VisRange,
			VisGood:         l.WebcamPage.VisGood,
			VisPoor:         l.WebcamPage.VisPoor,
			VisModerate:     l.WebcamPage.VisModerate,
		},
		Footer: FooterLocale{
			Motto:         l.Footer.Motto,
			ExploreTitle:  l.Footer.ExploreTitle,
			ContactsTitle: l.Footer.ContactsTitle,
			Copyright:     strings.ReplaceAll(l.Footer.Copyright, "{year}", fmt.Sprintf("%d", time.Now().Year())),
		},
	}
}

func renderLocale(locale string, baseUrl string, indexData *IndexFile, eventsData *EventsFile, galleryT GalleryData, rawItineraries []ItineraryFile) {
	// Merge shared and localized
	renderIndex := createRenderIndex(locale, indexData, eventsData)

	// Prepare Itineraries for this locale
	var localItineraries []RenderItinerary
	for _, raw := range rawItineraries {
		// Filter: Only include itineraries with a GPX file
		if raw.GpxFile == "" {
			continue
		}

		l := raw.It
		if locale == "en" {
			l = raw.En
		}
		localItineraries = append(localItineraries, RenderItinerary{
			Slug:           raw.Slug,
			Type:           raw.Type,
			Image:          raw.Image,
			GpxFile:        raw.GpxFile,
			YoutubeVideoID: raw.YoutubeVideoID,
			Gallery:        raw.ProcessedGallery, // Use processed gallery
			Difficulty:     raw.Difficulty,
			DistanceKM:     raw.DistanceKM,
			Duration:       raw.Duration,
			ElevationGain:  raw.ElevationGain,
			Author:         raw.Author,
			Title:          l.Title,
			Description:    l.Description,
			LongDesc:       l.LongDesc,
			Tags:           l.Tags,
		})
	}

	webcamImages, err := loadWebcamImages("static/webcam")
	if err != nil {
		log.Printf("Error loading webcam images: %v", err)
	}

	// Limit gallery images for the index page to 8
	indexGalleryImages := galleryT.Images
	if len(indexGalleryImages) > 8 {
		indexGalleryImages = indexGalleryImages[:8]
	}

	// Render Index
	ctx := pongo2.Context{
		"locale":         locale,
		"base_url":       baseUrl,
		"alternate_url":  computeAlternateUrl(locale, "/"),
		"page_title":     renderIndex.Hero.Title,
		"t":              renderIndex, // We pass our flattened struct as 't'
		"gallery_images": indexGalleryImages,
		"itineraries":    localItineraries,
	}

	// Update WebcamPage with loaded images
	renderIndex.WebcamPage.Images = webcamImages

	tpl := pongo2.Must(pongo2.FromFile("templates/index.html"))
	outPath := "dist/index.html"
	if locale == "en" {
		outPath = "dist/en/index.html"
	}

	err = renderToFile(tpl, ctx, outPath)
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
	for i := range data.Hero.Images {
		data.Hero.Images[i] = "/static/" + data.Hero.Images[i]
	}
	data.Welcome.Image = "/static/" + data.Welcome.Image
	if data.Itineraries.HeroImage != "" {
		data.Itineraries.HeroImage = "/static/" + data.Itineraries.HeroImage
	}
	return &data, nil
}

func loadEvents(path string) (*EventsFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data EventsFile
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
	for i := range data.Images {
		url, thumb, err := processImage(data.Images[i].Url)
		if err != nil {
			log.Printf("Warning: processing image %s failed: %v", data.Images[i].Url, err)
			data.Images[i].Url = "/static/" + data.Images[i].Url
			data.Images[i].Thumbnail = data.Images[i].Url // Fallback
		} else {
			data.Images[i].Url = url
			data.Images[i].Thumbnail = thumb
		}
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
				// Calculate Elevation Gain
				// Handle both relative path from content/itineraries or absolute-ish path
				// The GpxFile string usually comes as "gpx/foo.gpx" or "/static/gpx/foo.gpx"
				// We need the filesystem path: "static/gpx/foo.gpx"

				cleanPath := it.GpxFile
				if strings.HasPrefix(cleanPath, "/static/") {
					cleanPath = strings.TrimPrefix(cleanPath, "/static/")
				}
				fsPath := filepath.Join("static", cleanPath)

				gain, dist, err := processGpx(fsPath)
				if err != nil {
					log.Printf("Warning: failed to process GPX %s: %v", fsPath, err)
				} else {
					it.ElevationGain = gain
					it.DistanceKM = dist
				}

				it.GpxFile = "/static/" + cleanPath
			}

			// Process Gallery
			it.ProcessedGallery = make([]GalleryImage, len(it.Gallery))
			for i, rawPath := range it.Gallery {
				url, thumb, err := processImage(rawPath)
				if err != nil {
					log.Printf("Warning: processing itinerary image %s failed: %v", rawPath, err)
					it.ProcessedGallery[i] = GalleryImage{Url: "/static/" + rawPath, Thumbnail: "/static/" + rawPath}
				} else {
					it.ProcessedGallery[i] = GalleryImage{Url: url, Thumbnail: thumb}
				}
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

func collectUsedImages(index *IndexFile, gallery *GalleryData, itineraries []ItineraryFile) map[string]bool {
	used := make(map[string]bool)

	// Helper to add path
	add := func(p string) {
		if p == "" {
			return
		}
		// p is like "/static/img/foo.jpg" or "img/hero.jpg" (from index.toml)
		clean := p
		if strings.HasPrefix(clean, "/") {
			clean = strings.TrimPrefix(clean, "/")
		} else {
			// If it doesn't have /static prefix, it might be from TOML relative to static/
			if !strings.HasPrefix(clean, "static/") {
				clean = "static/" + clean
			}
		}
		used[clean] = true
	}

	for _, img := range index.Hero.Images {
		add(img)
	}
	add(index.Welcome.Image)

	for _, img := range gallery.Images {
		add(img.Url)
		add(img.Thumbnail)
	}

	for _, it := range itineraries {
		add(it.Image)
		for _, img := range it.ProcessedGallery {
			add(img.Url)
			add(img.Thumbnail)
		}
	}

	return used
}

func cleanupImages(usedImages map[string]bool) error {
	// 1. Scan static/img
	err := filepath.WalkDir("static/img", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// path is like "static/img/foo.jpg"
			if !usedImages[path] {
				log.Printf("Removing unused image: %s", path)
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 2. Scan static/thumbs/img
	thumbsDir := "static/thumbs/img"
	if _, err := os.Stat(thumbsDir); os.IsNotExist(err) {
		return nil
	}

	err = filepath.WalkDir(thumbsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// path is like "static/thumbs/img/foo.jpg"
			if !usedImages[path] {
				log.Printf("Removing unused thumbnail: %s", path)
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}

// processImage ensures a thumbnail exists for the given image and returns the web paths for original and thumbnail.
func processImage(rawPath string) (originalWeb string, thumbWeb string, err error) {
	// Clean rawPath
	cleanPath := rawPath
	if strings.HasPrefix(cleanPath, "/static/") {
		cleanPath = strings.TrimPrefix(cleanPath, "/static/")
	} else if strings.HasPrefix(cleanPath, "static/") {
		cleanPath = strings.TrimPrefix(cleanPath, "static/")
	}
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	srcPath := filepath.Join("static", cleanPath)
	thumbPath := filepath.Join("static", "thumbs", cleanPath)

	// Check if source exists
	info, err := os.Stat(srcPath)
	if err != nil {
		return "", "", fmt.Errorf("source image not found: %w", err)
	}

	// Check if thumb exists and is newer
	thumbInfo, err := os.Stat(thumbPath)
	if err == nil && thumbInfo.ModTime().After(info.ModTime()) {
		// Thumb is up to date
		return "/static/" + cleanPath, "/static/thumbs/" + cleanPath, nil
	}

	// Create thumb dir
	if err := os.MkdirAll(filepath.Dir(thumbPath), 0755); err != nil {
		return "", "", fmt.Errorf("failed to create thumb dir: %w", err)
	}

	// Generate
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open image: %w", err)
	}

	// Resize to width 600, preserving aspect ratio
	thumbImg := imaging.Resize(srcImg, 600, 0, imaging.Lanczos)

	if err := imaging.Save(thumbImg, thumbPath); err != nil {
		return "", "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return "/static/" + cleanPath, "/static/thumbs/" + cleanPath, nil
}

func loadWebcamImages(dir string) ([]string, error) {
	var images []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return images, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".jpg") {
			if entry.Name() != "current.jpg" {
				images = append(images, "/static/webcam/"+entry.Name())
			}
		}
	}
	return images, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// GPX Parsing Structures

type Gpx struct {
	Trk []Trk `xml:"trk"`
}

type Trk struct {
	TrkSeg []TrkSeg `xml:"trkseg"`
}

type TrkSeg struct {
	TrkPt []TrkPt `xml:"trkpt"`
}

type TrkPt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele float64 `xml:"ele"`
}

func processGpx(path string) (elevationGain int, distanceKm float64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	var gpx Gpx
	if err := xml.NewDecoder(f).Decode(&gpx); err != nil {
		return 0, 0, err
	}

	var gain float64
	var dist float64
	var prevEle float64
	var prevLat, prevLon float64
	first := true

	for _, trk := range gpx.Trk {
		for _, seg := range trk.TrkSeg {
			for _, pt := range seg.TrkPt {
				if first {
					prevEle = pt.Ele
					prevLat = pt.Lat
					prevLon = pt.Lon
					first = false
					continue
				}

				// Elevation Gain
				diff := pt.Ele - prevEle
				if diff > 0 {
					gain += diff
				}
				prevEle = pt.Ele

				// Distance
				dist += haversine(prevLat, prevLon, pt.Lat, pt.Lon)
				prevLat = pt.Lat
				prevLon = pt.Lon
			}
		}
	}

	return int(math.Round(gain)), math.Round((dist/1000)*100) / 100, nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
