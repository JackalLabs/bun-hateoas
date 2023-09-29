package hateoas

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/uptrace/bunrouter"
)

type Route struct {
	Type string
	Path string
}

type Group struct {
	Name     string
	Paths    PathSet
	BunGroup *bunrouter.Group
}

type HATEOAS struct {
	Groups []*Group
}

type GroupType interface {
	WithGroup(path string, fn func(group *bunrouter.Group))
}

// PathSet is a set of routes where no routes can be duplicated
type PathSet map[Route]bool

// Add puts the route into the set
func (p PathSet) Add(route Route) {
	p[route] = true
}

// List returns a slice of routes from the set
func (p PathSet) List() []Route {
	keys := make([]Route, len(p))

	i := 0
	for k := range p {
		keys[i] = k
		i++
	}

	return keys
}

// Append adds a new route to the group. Handler is usually GET or POST.
func (g *Group) Append(handlerType string, path string, handlerFunction bunrouter.HandlerFunc) {
	g.Paths.Add(Route{
		Type: handlerType,
		Path: path,
	})

	g.BunGroup.Handle(handlerType, path, handlerFunction)
}

// GET adds a new GET route to the group
func (g *Group) GET(path string, handlerFunction bunrouter.HandlerFunc) {
	g.Append("GET", path, handlerFunction)
}

// POST adds a new POST route to the group
func (g *Group) POST(path string, handlerFunction bunrouter.HandlerFunc) {
	g.Append("POST", path, handlerFunction)
}

// DELETE adds a new DELETE route to the group
func (g *Group) DELETE(path string, handlerFunction bunrouter.HandlerFunc) {
	g.Append("DELETE", path, handlerFunction)
}

type BunGroupFunc func(group *Group)

func (h *HATEOAS) WithGroup(g GroupType, path string, groupFunc BunGroupFunc) {
	g.WithGroup(path, func(group *bunrouter.Group) {
		newGroup := &Group{
			Name:     path,
			Paths:    make(PathSet),
			BunGroup: group,
		}
		h.Groups = append(h.Groups, newGroup)
		groupFunc(newGroup)
	})
}

func NewHATEOAS() *HATEOAS {
	h := HATEOAS{
		Groups: make([]*Group, 0),
	}

	return &h
}

func (g *Group) toString() string {
	s := strings.Builder{}

	for route := range g.Paths {
		s.WriteString(fmt.Sprintf("<li>%6s <a href=\"%s%s\">%s%s</a></li>", route.Type, g.Name, route.Path, g.Name, route.Path))
	}
	return s.String()
}

func (h *HATEOAS) Print() string {
	s := strings.Builder{}

	s.WriteString("<html><head><title>STRATUS REST API</title><style>*{white-space: pre;font-family:monospace;list-style-type: none;}</style></head><body><ul>")

	for _, group := range h.Groups {
		s.WriteString(group.toString())
	}

	s.WriteString("</ul></body></html>")

	return s.String()
}

func FriendlyTimestamp() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d\n",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
}

func processError(block string, caughtError error) {
	fmt.Printf("***** Error in block: %s *****\n", block)
	fmt.Printf("***** Stamp: %s *****\n", FriendlyTimestamp())
	fmt.Println(caughtError)
	fmt.Println("***** End Error Report *****")
}

func (h *HATEOAS) Handler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := w.Write([]byte(h.Print()))
		if err != nil {
			processError("WriteError for VersionHandler", err)
		}
		return nil
	}
}
