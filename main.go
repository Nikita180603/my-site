package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

const personName = "Калугин Никита Анреевич"

type link struct {
	Path  string
	Label string
}

type page struct {
	id    string
	path  string
	file  string
	label string
	title string
}

type view struct {
	Name  string
	Title string
	Nav   []link
}

var pages = []page{
	{id: "home", path: "/", file: "index.html", label: "Анкета", title: personName},
	{id: "skills", path: "/skills", file: "skills.html", label: "Профессиональные навыки", title: "Профессиональные навыки — " + personName},
	{id: "qualities", path: "/qualities", file: "qualities.html", label: "Личные качества", title: "Личные качества — " + personName},
	{id: "info", path: "/info", file: "info.html", label: "Дополнительная информация", title: "Дополнительная информация — " + personName},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serve)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("сайт открыт:")
	log.Println("  на этом компьютере: http://localhost:8080")
	for _, ip := range localIPv4() {
		log.Printf("  в локальной сети:   http://%s:8080", ip)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	switch path {
	case "/index.html":
		http.Redirect(w, r, "/", http.StatusFound)
		return
	case "/skills.html":
		http.Redirect(w, r, "/skills", http.StatusFound)
		return
	case "/qualities.html":
		http.Redirect(w, r, "/qualities", http.StatusFound)
		return
	case "/info.html":
		http.Redirect(w, r, "/info", http.StatusFound)
		return
	case "/styles.css", "/photo.jpg":
		http.ServeFile(w, r, "."+path)
		return
	}

	for _, current := range pages {
		if current.path != path {
			continue
		}
		render(w, current)
		return
	}

	http.NotFound(w, r)
}

func render(w http.ResponseWriter, current page) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		current.file,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "layout", view{
		Name:  personName,
		Title: current.title,
		Nav:   navExcept(current.id),
	})
	if err != nil {
		log.Println(err)
	}
}

func navExcept(currentID string) []link {
	var links []link
	for _, item := range pages {
		if item.id == currentID {
			continue
		}
		links = append(links, link{Path: item.path, Label: item.label})
	}
	return links
}

func localIPv4() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	var ips []string
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		if ip := ipNet.IP.To4(); ip != nil {
			ips = append(ips, ip.String())
		}
	}
	return ips
}
