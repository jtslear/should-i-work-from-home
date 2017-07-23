package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func WriteLog(handle http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		end := time.Now()
		latency := end.Sub(start)
		size := writer.length
		statusCode := writer.status
		log.Printf("\"%s %s\" %d %d %v", request.Method, request.URL.Path, statusCode, size, latency)
	}
}

type Question struct {
	Statement string
}

func main() {
	var a = []Question{
		Question{
			Statement: "Please rate your Fatigue.",
		},
		Question{
			Statement: "Are you suffering from Hives?",
		},
		Question{
			Statement: "Are you suffering from Dizziness?",
		},
		Question{
			Statement: "Are you suffering from Irregular Bowel Movements?",
		},
		Question{
			Statement: "Are you suffering from Fogginess?",
		},
		Question{
			Statement: "Are you suffering from Bloating?",
		},
		Question{
			Statement: "Are you suffering from Stomach Cramps?",
		},
		Question{
			Statement: "Please rate your Fever.",
		},
		Question{
			Statement: "Please rate your Nausea.",
		},
		Question{
			Statement: "Are you suffering from Vomit?",
		},
		Question{
			Statement: "Are you suffering from Excess Gas?",
		},
		Question{
			Statement: "Are you suffering from Muscle Soreness?",
		},
		Question{
			Statement: "Are you suffering from Sinuses Symptoms?",
		},
		Question{
			Statement: "Are you suffering from Shortness of Breath?",
		},
	}

	type Reply struct {
		Average float32
		Response string
	}

	http.HandleFunc("/api/formaverage", func (w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var length,total float32
		length = float32(len(r.Form))
		for k, _ := range r.Form {
			f, _ := strconv.ParseFloat(r.FormValue(k), 32)
			total = float32(f) + total
		}
		var respond Reply
		respond.Average = total/length

		switch {
		case respond.Average >= 7.5 && respond.Average < 9.5:
			respond.Response = "Work from Home"
		case respond.Average >= 9.5:
			respond.Response = "Take PTO"
		default:
			respond.Response = "Go to work"
		}

		b, err := json.Marshal(respond)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "I fail making golang")
		} else {
			fmt.Fprintf(w, "%s", b)
		}
	})

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.html").ParseFiles("./tmpl/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "I fail making golang resolve template files: %s\n", err)
		} else {
			err := t.Execute(w, a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "I fail making golang parse templates: %s\n", err)
			}
		}
	})

	var port string = "8080"

	log.Printf("Listening at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, WriteLog(http.DefaultServeMux)))
}
