package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	//router := mux.NewRouter()

	//dossier static pour les images et css
	//router.PathPrefix("/static/").

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	var err error
	// Configure la connexion à la base de données
	db, err = sql.Open("mysql", "root:Lee7474//@tcp(127.0.0.1:3306)/base_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Vérifie la connexion à la base de données
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/signup", signupHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080/signup")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "signup.html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Hacher le mot de passe avant de le stocker dans la base de données
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Insérer le nouvel utilisateur dans la base de données
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Impossible de créer l'utilisateur", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Utilisateur %s créé avec succès!", username)

}
