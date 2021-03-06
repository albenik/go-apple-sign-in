package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"
	jsoniter "github.com/json-iterator/go"

	applesignin "github.com/albenik-go/apple-sign-in"
	signinkey "github.com/albenik-go/apple-sign-in/key"
)

const (
	state = "static_state"
	nonce = "static_nonce"
)

func main() {
	listen := flag.String("addr", ":8080", "Listen address")
	audience := flag.String("aud", "", "Audience")
	teamID := flag.String("team", "", "Team ID")
	clientID := flag.String("client", "", "Client ID")
	keyID := flag.String("key", "", "Key ID")
	keyFile := flag.String("keyfile", "", "Key file")
	redirectURL := flag.String("redirect", "", "Redirect URL")
	flag.Parse()

	key, err := signinkey.ReadPrivateFromPEMFile(*keyFile)
	if err != nil {
		panic(err)
	}

	appl := applesignin.New(*teamID, *clientID, *keyID, key,
		applesignin.WithJWTParser(jwt.NewParser(jwt.WithAudience(*audience))))
	appl.RedirectURL = *redirectURL

	h := handler{appl: appl}
	http.HandleFunc("/", h.root)
	http.HandleFunc("/callback", h.callback)
	http.HandleFunc("/echo", h.callbackEcho)
	http.HandleFunc("/validate", h.validate)

	if err := http.ListenAndServe(*listen, nil); err != nil {
		panic(err)
	}
}

type handler struct {
	appl *applesignin.Client
}

func (h *handler) root(w http.ResponseWriter, _ *http.Request) {
	u := h.appl.AuthURL(applesignin.ResponseModePost, []string{
		applesignin.ScopeEmail,
		applesignin.ScopeName,
	}, state, nonce)

	if err := rootTemplate.Execute(w, u); err != nil {
		panic(err)
	}
}

func (h *handler) callback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
	default:
		s := http.StatusBadRequest
		http.Error(w, http.StatusText(s), s)
		return //nolint:wsl
	}

	if s := r.FormValue("state"); s != state {
		http.Error(w, fmt.Sprintf("Invalid state %q", s), http.StatusBadRequest)
		return
	}

	fmt.Println(">>>", r.Form)

	result, err := h.appl.ValidateCode(r.FormValue("code"), nonce, applesignin.MaxExpiration)
	h.printResult(w, result, err)
}

func (h *handler) callbackEcho(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
	default:
		s := http.StatusBadRequest
		http.Error(w, http.StatusText(s), s)
		return //nolint:wsl
	}

	if s := r.FormValue("state"); s != state {
		http.Error(w, fmt.Sprintf("Invalid state %q", s), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	for k, v := range r.Form {
		fmt.Fprintln(w, k, "=", v)
	}
}

func (h *handler) validate(w http.ResponseWriter, r *http.Request) {
	result, err := h.appl.ValidateRefreshToken(r.FormValue("token"), applesignin.MaxExpiration)
	h.printResult(w, result, err)
}

func (h *handler) printResult(w io.Writer, result *applesignin.TokenResponse, err error) {
	if err != nil {
		if err = resultTemplate.Execute(w, err); err != nil {
			panic(err)
		}
		return //nolint:wsl
	}

	json, err := jsoniter.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	if err = resultTemplate.Execute(w, string(json)); err != nil {
		panic(err)
	}
}
