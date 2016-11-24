package main

import (
	"flag"
	"github.com/blevesearch/bleve/analysis"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/huichen/sego"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	host            string
	ideographRegexp = regexp.MustCompile(`\p{Han}+`)
	r               = render.New()
	segmenter       *sego.Segmenter
	err             error
	router          = mux.NewRouter()
	dictFiles       string
	nested          = true
)

func init() {
	flag.StringVar(&host, "host", "localhost:3000", "The tokenizer server host.")
	flag.StringVar(&dictFiles, "dicts", "dict.txt", "dictionary file paths.")
	flag.Parse()
}

func appendToTokenStreams(stream analysis.TokenStream, segment *sego.Segment, start, pos int, nested, top bool) (analysis.TokenStream, int) {
	if nested && len(segment.Token().Segments()) > 0 {
		for _, one := range segment.Token().Segments() {
			stream, pos = appendToTokenStreams(stream, one, start+one.Start(), pos, nested, false)
		}
	}

	if top || !isFake(segment) {
		token := &analysis.Token{
			Term:     []byte(segment.Token().Text()),
			Start:    start,
			End:      start + segment.End() - segment.Start(),
			Position: pos,
			Type:     tokenType(segment.Token().Text()),
		}
		stream = append(stream, token)
		pos++
	}

	return stream, pos
}

func isFake(segment *sego.Segment) bool {
	return segment.Token().Frequency() == 1 && segment.Token().Pos() == "x"
}

func tokenType(s string) analysis.TokenType {
	if ideographRegexp.MatchString(s) {
		return analysis.Ideographic
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return analysis.Numeric
	}

	return analysis.AlphaNumeric
}

func main() {
	router = mux.NewRouter()

	if segmenter, err = getSegoSegmenter(dictFiles); err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/api/tokenizer/", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		data := req.Form.Get("data")
		stream := make(analysis.TokenStream, 0)
		pos := 1

		segments := segmenter.Segment([]byte(data))
		for _, segment := range segments {
			p := &segment
			stream, pos = appendToTokenStreams(stream, p, p.Start(), pos, nested, true)
		}
		r.JSON(w, http.StatusOK, stream)
	}).Methods("POST")

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)
	n.Run(host)
}
