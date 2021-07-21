package api

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"poc/internal/model"
	"poc/internal/util"
)

var (
	tpl          *template.Template
	sessionStore = sessions.NewCookieStore([]byte("okta-custom-login-session-store"))
	state        = generateState()
	nonce        = "NonceNotSetYet"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func RegisterAuthApi(router *mux.Router) {
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/authorization-code/callback", authorizationCodeHandler)
}

func generateState() string {
	// Generate a random byte array for state parameter
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache")
	nonce, _ := util.GenerateNonce()

	issuerParts, _ := url.Parse(viper.GetString("Issuer"))
	baseUrl := issuerParts.Scheme + "://" + issuerParts.Hostname()

	data := model.AuthData{
		Profile:         getProfileData(r),
		IsAuthenticated: isAuthenticated(r),
		BaseUrl:         baseUrl,
		RedirectUri:     viper.GetString("RedirectUrl"),
		ClientId:        viper.GetString("ClientId"),
		Issuer:          viper.GetString("Issuer"),
		State:           state,
		Nonce:           nonce,
	}

	if err := tpl.ExecuteTemplate(w, "login.gohtml", data); err != nil {
		panic(err)
	}
}

func authorizationCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check the state that was returned in the query string is the same as the above state
	if r.URL.Query().Get("state") != state {
		_, err := fmt.Fprintln(w, "The state was not as expected")
		if err != nil {
			panic(err)
		}
		return
	}
	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		fmt.Fprintln(w, "The code was not returned or is not accessible")
		return
	}

	exchange := exchangeCode(r.URL.Query().Get("code"), r)

	if exchange.Error != "" {
		fmt.Println(exchange.Error)
		fmt.Println(exchange.ErrorDescription)
		return
	}

	session, err := sessionStore.Get(r, "okta-custom-login-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, verificationError := verifyToken(exchange.IdToken)

	if verificationError != nil {
		fmt.Println(verificationError)
	}

	if verificationError == nil {
		session.Values["id_token"] = exchange.IdToken
		session.Values["access_token"] = exchange.AccessToken

		session.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func exchangeCode(code string, r *http.Request) model.Exchange {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(viper.GetString("ClientId") + ":" + viper.GetString("ClientSecret")))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Set("code", code)
	q.Add("redirect_uri", viper.GetString("RedirectUrl"))

	url := viper.GetString("Issuer") + "/oauth2/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(body)
	log.Println(bodyStr)

	defer resp.Body.Close()
	var exchange model.Exchange
	if err := json.Unmarshal(body, &exchange); err != nil {
		panic(err)
	}

	return exchange
}

func getProfileData(r *http.Request) map[string]string {
	m := make(map[string]string)

	session, err := sessionStore.Get(r, "okta-custom-login-session-store")

	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return m
	}

	reqUrl := viper.GetString("Issuer") + "/v1/userinfo"

	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+session.Values["access_token"].(string))
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(body)
	log.Println(bodyStr)

	defer resp.Body.Close()
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}

	return m
}

func isAuthenticated(r *http.Request) bool {
	session, err := sessionStore.Get(r, "okta-custom-login-session-store")

	if err != nil || session.Values["id_token"] == nil || session.Values["id_token"] == "" {
		return false
	}

	return true
}

func verifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = viper.GetString("ClientId")
	jv := verifier.JwtVerifier{
		Issuer:           viper.GetString("Issuer"),
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}