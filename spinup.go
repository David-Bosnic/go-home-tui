package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type CalendarEvent struct {
	Kind             string    `json:"kind"`
	Etag             string    `json:"etag"`
	Summary          string    `json:"summary"`
	Description      string    `json:"description"`
	Updated          time.Time `json:"updated"`
	TimeZone         string    `json:"timeZone"`
	AccessRole       string    `json:"accessRole"`
	DefaultReminders []struct {
		Method  string `json:"method"`
		Minutes int    `json:"minutes"`
	} `json:"defaultReminders"`
	NextSyncToken string `json:"nextSyncToken"`
	Items         []struct {
		Kind        string    `json:"kind"`
		Etag        string    `json:"etag"`
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		HTMLLink    string    `json:"htmlLink"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
		Summary     string    `json:"summary"`
		Description string    `json:"description,omitempty"`
		Location    string    `json:"location,omitempty"`
		Creator     struct {
			Email string `json:"email"`
			Self  bool   `json:"self"`
		} `json:"creator"`
		Organizer struct {
			Email string `json:"email"`
			Self  bool   `json:"self"`
		} `json:"organizer"`
		Start struct {
			DateTime string `json:"dateTime"`
			TimeZone string `json:"timeZone"`
		} `json:"start"`
		End struct {
			DateTime string `json:"dateTime"`
			TimeZone string `json:"timeZone"`
		} `json:"end"`
		Transparency string `json:"transparency,omitempty"`
		Visibility   string `json:"visibility,omitempty"`
		ICalUID      string `json:"iCalUID"`
		Sequence     int    `json:"sequence"`
		Attendees    []struct {
			Email          string `json:"email"`
			Organizer      bool   `json:"organizer"`
			Self           bool   `json:"self"`
			ResponseStatus string `json:"responseStatus"`
		} `json:"attendees"`
		GuestsCanInviteOthers bool `json:"guestsCanInviteOthers,omitempty"`
		Reminders             struct {
			UseDefault bool `json:"useDefault"`
		} `json:"reminders"`
		Source struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"source"`
		EventType      string `json:"eventType"`
		ConferenceData struct {
			EntryPoints []struct {
				EntryPointType string `json:"entryPointType"`
				URI            string `json:"uri"`
				Label          string `json:"label"`
				MeetingCode    string `json:"meetingCode"`
			} `json:"entryPoints"`
			ConferenceSolution struct {
				Key struct {
					Type string `json:"type"`
				} `json:"key"`
				Name    string `json:"name"`
				IconURI string `json:"iconUri"`
			} `json:"conferenceSolution"`
			ConferenceID string `json:"conferenceId"`
		} `json:"conferenceData"`
	} `json:"items"`
}

type apiConfig struct {
	accessToken  string
	refreshToken string
	calendarID   string
	clientID     string
	clientSecret string
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localHost := "http://localhost:" + os.Getenv("PORT")
		w.Header().Set("Access-Control-Allow-Origin", localHost)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}

		next.ServeHTTP(w, r)
	})
}

func SpinUp() {
	godotenv.Load()
	var apiConf apiConfig
	apiConf.accessToken = "Bearer " + os.Getenv("ACCESS_TOKEN")
	apiConf.calendarID = os.Getenv("CALENDAR_ID")
	apiConf.refreshToken = os.Getenv("REFRESH_TOKEN")
	apiConf.clientID = os.Getenv("CLIENT_ID")
	apiConf.clientSecret = os.Getenv("CLIENT_SECRET")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /calendar/events", apiConf.handlerEventsGet)
	mux.HandleFunc("POST /calendar/events", apiConf.handlerEventsPost)
	mux.HandleFunc("GET /auth/callback", apiConf.handleOauthCallback)
	mux.HandleFunc("GET /auth/google", apiConf.startOauthFlow)
	mux.HandleFunc("POST /admin/refresh", apiConf.refreshAccessTokenPost)
	mux.HandleFunc("PATCH /calendar/events", apiConf.handlerEventsPatch)
	mux.HandleFunc("DELETE /calendar/events", apiConf.handlerEventsDelete)

	ServerMux := http.Server{
		Handler: corsMiddleware(mux),
		Addr:    ":8080",
	}

	fmt.Println("Running Server")
	go func() {
		err := ServerMux.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Spinning up server")
		}
	}()
}
