package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Dane autora i port aplikacji
const autor = "Wiktoria Matacz"
const port = "8080"

// Mapa: pełna nazwa kraju => kod ISO
var kraje = map[string]string{
	"Polska": "PL",
	"Niemcy": "DE",
	"USA":    "US",
}

// Miasta dla każdego kraju (wg pełnej nazwy kraju)
var lokalizacje = map[string][]string{
	"Polska": {"Warszawa", "Kraków", "Gdańsk"},
	"Niemcy": {"Berlin", "Hamburg"},
	"USA":    {"Chicago", "San Francisco", "Los Angeles"},
}

// Klucz API pobrany z zmiennej środowiskowej
var kluczAPI = os.Getenv("OPENWEATHER_API_KEY")

// Struktura danych pogodowych odbieranych z API
type WynikPogody struct {
	Main struct {
		Temperatura float64 `json:"temp"`
		Wilgotnosc  int     `json:"humidity"`
	} `json:"main"`
	Pogoda []struct {
		Stan  string `json:"main"`
		Opis  string `json:"description"`
		Ikona string `json:"icon"`
	} `json:"weather"`
}

// Funkcja główna
func main() {
	if kluczAPI == "" {
		log.Fatal("Brak klucza API! Ustaw zmienną środowiskową OPENWEATHER_API_KEY.")
	}

	log.Printf("Data uruchomienia: %s", time.Now().Format(time.RFC3339))
	log.Printf("Autor: %s", autor)
	log.Printf("Numer portu: %s", port)

	http.HandleFunc("/", obslugaStartowa)
	http.HandleFunc("/pogoda", obslugaPogody)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Strona główna z formularzem wyboru kraju i miasta
func obslugaStartowa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Serializacja mapy lokalizacji do JSON
	lokalizacjeJSON, err := json.Marshal(lokalizacje)
	if err != nil {
		http.Error(w, "Błąd danych", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="pl">
		<head>
			<meta charset="UTF-8">
			<title>Aplikacja Pogodowa</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #e0f2ff;
					color: #003f5c;
					text-align: center;
				}
				form {
					margin-top: 30px;
				}
				.kafelek {
					margin: 30px auto;
					background: #cceeff;
					padding: 20px;
					border-radius: 10px;
					width: 250px;
					box-shadow: 2px 2px 6px rgba(0,0,0,0.2);
				}
				.kafelek img {
					width: 100px;
					height: 100px;
				}
			</style>
		</head>
		<body>
			<h1>Aplikacja Pogodowa</h1>
			<form action="/pogoda" method="get">
				Kraj:
				<select name="kraj" id="kraj" onchange="aktualizujMiasta()">
	`)

	// Generowanie opcji krajów
	for kraj := range lokalizacje {
		fmt.Fprintf(w, "<option value='%s'>%s</option>", kraj, kraj)
	}

	fmt.Fprint(w, `
				</select><br>
				Miasto:
				<select name="miasto" id="miasto"></select><br>
				<input type="submit" value="Sprawdź pogodę">
			</form>

			<script>
				const lokalizacje = `+string(lokalizacjeJSON)+`;

				function aktualizujMiasta() {
					const kraj = document.getElementById("kraj").value;
					const miastaSelect = document.getElementById("miasto");
					miastaSelect.innerHTML = "";

					if (lokalizacje[kraj]) {
						lokalizacje[kraj].forEach(miasto => {
							const option = document.createElement("option");
							option.value = miasto;
							option.text = miasto;
							miastaSelect.appendChild(option);
						});
					}
				}

				//Wywołane przy załadowaniu strony, by ustawić miasta dla domyślnego kraju
				window.onload = aktualizujMiasta;
			</script>
		</body>
		</html>
	`)
}

// Obsługa zapytań pogodowych
func obslugaPogody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	krajPelnaNazwa := r.URL.Query().Get("kraj")
	miasto := r.URL.Query().Get("miasto")

	if krajPelnaNazwa == "" || miasto == "" {
		http.Error(w, "Brakuje kraju lub miasta w zapytaniu", http.StatusBadRequest)
		return
	}

	kodKraju, ok := kraje[krajPelnaNazwa]
	if !ok {
		http.Error(w, "Nieznany kraj", http.StatusBadRequest)
		return
	}

	// Budowanie URL zapytania do API pogodowego
	q := fmt.Sprintf("%s,%s", miasto, kodKraju)
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=pl",
		url.QueryEscape(q), kluczAPI)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Błąd połączenia z API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("API zwróciło błąd: %s", body), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	// Dekodowanie odpowiedzi JSON
	var wynik WynikPogody
	if err := json.NewDecoder(resp.Body).Decode(&wynik); err != nil {
		http.Error(w, "Błąd dekodowania danych", http.StatusInternalServerError)
		return
	}

	// Generowanie HTML wyniku
	fmt.Fprint(w, `
	<!DOCTYPE html>
	<html lang="pl">
	<head>
		<meta charset="UTF-8">
		<title>Wynik pogody</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #e0f2ff;
				color: #003f5c;
				text-align: center;
			}
			.kafelek {
				margin: 30px auto;
				background: #cceeff;
				padding: 20px;
				border-radius: 10px;
				width: 250px;
				box-shadow: 2px 2px 6px rgba(0,0,0,0.2);
			}
			.kafelek img {
				width: 100px;
				height: 100px;
			}
			a {
				display: block;
				margin-top: 20px;
			}
		</style>
	</head>
	<body>
	`)

	fmt.Fprintf(w, `<div class="kafelek">
		<h2>%s, %s</h2>
		<img src="https://openweathermap.org/img/wn/%s@2x.png" alt="ikona pogodowa">
		<p><strong>Temperatura:</strong> %.1f°C</p>
		<p><strong>Wilgotność:</strong> %d%%</p>
		<p><strong>Stan:</strong> %s (%s)</p>
	</div>`,
		miasto, krajPelnaNazwa,
		wynik.Pogoda[0].Ikona,
		wynik.Main.Temperatura,
		wynik.Main.Wilgotnosc,
		wynik.Pogoda[0].Stan,
		wynik.Pogoda[0].Opis,
	)

	fmt.Fprint(w, `<a href="/">Wróć</a></body></html>`)
}
