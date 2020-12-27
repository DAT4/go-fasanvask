package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	APP_KEY = "golangcode.com"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/token", TokenHandler).Methods("GET")
	r.Handle("/book", AuthMiddleware(http.HandlerFunc(BookingHandler))).Methods("POST")
	r.Handle("/get", AuthMiddleware(http.HandlerFunc(GetBookingHandler))).Methods("GET")
	r.Handle("/delete", AuthMiddleware(http.HandlerFunc(DeleteBookingHandler))).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()
	user := Resident{
		Room:     r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}
	err := user.authenticate()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error":"invalid_credentials"}`)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Id,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

func AuthMiddleware(next http.Handler) http.Handler {
	if len(APP_KEY) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(APP_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	user, err := getCurrentUser(r)
	if err != nil {
		msg := fmt.Sprintf(`{"status":"%v"}`, err)
		io.WriteString(w, msg)
		return
	}

	var booking Booking
	_ = json.NewDecoder(r.Body).Decode(&booking)
	booking.User = *user

	addBooking(booking)

	io.WriteString(w, `{"status":"ok"}`)
}

func getCurrentUser(r *http.Request) (*Resident, error){
	tokenStr := strings.Split(r.Header.Get("Authorization")," ")[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexcpected signing method: %v", token.Header["alg"])
		}

		return []byte(APP_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["user"]
		user, err := getUser(id.(string))
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return nil, err
	}
}

func GetBookingHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getBookings())
}

func DeleteBookingHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Add("Content-Type", "application/json")
	var booking Booking
	_ = json.NewDecoder(r.Body).Decode(&booking)
	user, err := getCurrentUser(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	if booking.User.Id == user.Id {
		removeBooking(booking)
	} else {
		io.WriteString(w, `{"error":"You are not permitted to do that..."}`)
		return
	}
	io.WriteString(w, `{"status":"ok"}`)

}