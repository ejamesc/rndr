package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/unrolled/render"
)

const noTemplateMsg = "template '%s' does not exist"

var (
	switcher       HandlerSwitch
	rndr           *render.Render
	staticFilePath string
	templatesPath  string
	portNum        string
)

// HandlerSwitch is a switch that stores the right http handler.
// We use this to set return the right handler:
// "static" for the router that handles static files, and
// "templates" for the router that handles unrolled/render templates
type HandlerSwitch struct {
	static   http.Handler
	template http.Handler
}

// ServeHTTP makes HandlerSwitch conform to the http handler interface.
// In this case we serve static files if we detect /static in the url, else we serve rendered templates
func (hs HandlerSwitch) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.EscapedPath(), "/static") {
		switcher.static.ServeHTTP(w, req)
	} else {
		switcher.template.ServeHTTP(w, req)
	}
}

func init() {
	flag.StringVar(&staticFilePath, "s", "", "Static Directory Path. The full path of the static assets directory (required)")
	flag.StringVar(&templatesPath, "t", "", "Template Directory Path. The full path of the templates directory. All templates should be have the .html format (required)")
	flag.StringVar(&portNum, "p", "8095", "The port on which the server should start listening to")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `rndr is a webapp for quickly prototyping unrolled/render templates
Usage:

		rndr [options]
		
The options are:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	switcher = HandlerSwitch{}
	rndr = render.New(render.Options{
		Directory:     templatesPath,
		Extensions:    []string{".html"},
		Layout:        "base",
		IsDevelopment: true,
	})
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
	}

	ok := checkDirectoriesExist()
	if !ok {
		os.Exit(1)
	}

	setupRouters()

	fmt.Fprintln(os.Stderr, "~~|| RNDR // https://github.com/ejamesc/rndr ||~~")
	log.Printf("watching templates at: %s", templatesPath)
	log.Println("listening on port " + portNum)
	http.ListenAndServe(":"+portNum, switcher)
}

// TemplateServer grabs the request path, and renders the template file that corresponds to that path.
func TemplateServer(w http.ResponseWriter, req *http.Request) {
	path := req.URL.EscapedPath()
	if path == "/" {
		path = "/index"
	}
	if strings.Contains(path, "favicon.ico") {
		rndr.Data(w, http.StatusNotFound, []byte{})
		return
	}

	path = removeExt(path)
	tPath := getTemplateName(strings.Split(path, "/"))
	fpath := filepath.Join(templatesPath, tPath) + ".html"

	if !isFileExist(fpath) {
		log.Printf(noTemplateMsg, fpath)
		rndr.Text(w, http.StatusInternalServerError, fmt.Sprintf(noTemplateMsg, fpath))
		return
	}

	err := rndr.HTML(w, http.StatusOK, tPath, struct{}{})
	if err != nil {
		log.Println(fmt.Sprintf("problem rendering template '%s', got err %s", tPath, err))
	}
	log.Println(fmt.Sprintf("rendered template '%s'", tPath))
}

func setupRouters() {
	staticRouter, templateRouter := http.NewServeMux(), http.NewServeMux()
	templateRouter.HandleFunc("/", TemplateServer)
	staticRouter.Handle("/static/", http.FileServer(http.Dir(staticFilePath)))

	switcher.static = staticRouter
	switcher.template = templateRouter
}

func checkDirectoriesExist() bool {
	tExist := isFileExist(templatesPath)
	sExist := isFileExist(staticFilePath)
	if !tExist {
		fmt.Fprintf(os.Stderr, "templates dir does not exist: %s\n", templatesPath)
	}
	if !sExist {
		fmt.Fprintf(os.Stderr, "static dir does not exist: %s\n", staticFilePath)
	}
	return sExist || tExist
}

// Takes in '/some/dir/' and returns 'some/dir'
func getTemplateName(inArr []string) string {
	if len(inArr) < 1 {
		return ""
	}
	return strings.TrimRight(strings.Join(inArr, "/")[1:], "/")
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func removeExt(in string) string {
	if ext := filepath.Ext(in); ext != "" {
		in = strings.TrimRight(in, ext)
	}
	return in
}
