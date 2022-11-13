package sktch

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
)

var defaultWidth = 210
var defaultHeight = 297

type Inputs struct {
	Width  int
	Height int
	SVG    template.HTML
}

type SketchFunc func(inputs Inputs, canvas *svg.SVG) error

type ServerOpt interface {
	Apply(s *Server) error
}

type Server struct {
	name  string
	addr  string
	out   io.Writer
	srv   *http.Server
	mux   *http.ServeMux
	tmpls *template.Template
}

func NewServer(opts ...ServerOpt) (*Server, error) {
	s, err := defaultServer()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err := opt.Apply(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) AddSketch(path string, f SketchFunc) error {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		inputs := Inputs{
			Width:  defaultWidth,
			Height: defaultHeight,
		} // TODO build inputs

		b := &bytes.Buffer{}

		c := svg.New(b)
		c.StartviewUnit(inputs.Width, inputs.Height, "mm", 0, 0, inputs.Width, inputs.Height)

		err := f(inputs, c)
		if err != nil {
			panic(err) //todo remove this
		}

		c.End()

		inputs.SVG = template.HTML(b.String())

		s.tmpls.ExecuteTemplate(w, "index.html.tmpl", inputs)
	})

	// s.mux.HandleFunc(path+"svg.svg", func(w http.ResponseWriter, _ *http.Request) {
	// 	inputs := Inputs{
	// 		Width:  640,
	// 		Height: 480,
	// 	} // TODO build inputs
	//
	// 	c := svg.New(w)
	// 	c.StartviewUnit(inputs.Width, inputs.Height, "mm", 0, 0, inputs.Width, inputs.Height)
	// 	c.Line(0, 0, inputs.Width, inputs.Height, `style="stroke:black;fill-opacity:0.0;"`)
	// 	c.End()
	// })

	return nil
}

func (s *Server) ListenAndServe() error {
	s.srv.Addr = s.addr
	s.srv.Handler = s.mux
	fmt.Fprintf(s.out, "[%s] listening %s\n", s.name, s.addr)
	return s.srv.ListenAndServe()
}

func defaultServer() (*Server, error) {
	tmpl, err := template.ParseFiles("public/index.html.tmpl")
	if err != nil {
		return nil, err
	}

	return &Server{
		name:  "sktch",
		addr:  ":8181",
		out:   os.Stdout,
		srv:   &http.Server{},
		mux:   http.NewServeMux(),
		tmpls: tmpl,
	}, nil
}
