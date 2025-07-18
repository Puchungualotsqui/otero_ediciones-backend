package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BookDataBase struct {
	SimplifiedName string `json:"simplified_name"`
	Titulo         string `json:"titulo"`
	Autor          string `json:"autor"`
	NivelEducativo string `json:"nivel_educativo"`
	Materia        string `json:"materia"`
	Tipo           string `json:"tipo"`
	Idioma         string `json:"idioma"`
	Ilustraciones  string `json:"ilustraciones"`
	Genero         string `json:"genero"`
	Paginas        string `json:"paginas"`
	Tamano         string `json:"tamano"`
	DepositoLegal  string `json:"deposito_legal"`
	ISBN           string `json:"isbn"`
	Edad           string `json:"edad"`
	FichaDidactica string `json:"ficha_didactica"`
}

var PersistentInfo []BookDataBase

type BookResponse struct {
	Titulo         string `json:"titulo"`
	SimplifiedName string `json:"simplified_name"`
	TapaSmallURL   string `json:"tapa_small"`
}

type BookCompleteInfo struct {
	Titulo          string `json:"titulo"`
	TapaOriginalURL string `json:"tapa_original_url"`
	Autor           string `json:"autor"`
	Ilustraciones   string `json:"ilustraciones"`
	Materia         string `json:"materia"`
	NivelEducativo  string `json:"nivel_educativo"`
	Genero          string `json:"genero"`
	GuiaDidactica   string `json:"guia_didactica"`
	Tamano          string `json:"tamano"`
	Paginas         string `json:"paginas"`
	ISBN            string `json:"isbn"`
	DepositoLegal   string `json:"deposito_legal"`
	SinopsisURL     string `json:"descripcion"`
	Tipo            string `json:"tipo"`
	Idioma          string `json:"idioma"`
}

type MainPageInfo struct {
	Titulo string   `json:"titulo"`
	Libros []string `json:"libros"`
}

type HomeRowInfo struct {
	Titulo        string         `json:"titulo"`
	BookResponses []BookResponse `json:"book_responses"`
}

var MainPagePersistentInfo []MainPageInfo

func loadMainInfoJson(path string) ([]MainPageInfo, error) {
	var entries []MainPageInfo

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func getHomeRows() []HomeRowInfo {
	var homeRows []HomeRowInfo

	for _, entry := range MainPagePersistentInfo {
		var books []BookResponse

		for _, name := range entry.Libros {
			book := GetBookBySimplifiedName(name)
			if book != nil {
				books = append(books, BookResponse{
					Titulo:         book.Titulo,
					SimplifiedName: book.SimplifiedName,
					TapaSmallURL:   getTapaSmallURL(book.SimplifiedName),
				})
			}
		}

		homeRows = append(homeRows, HomeRowInfo{
			Titulo:        entry.Titulo,
			BookResponses: books,
		})
	}

	return homeRows
}

func booksHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rows := getHomeRows()

	err := json.NewEncoder(w).Encode(rows)
	if err != nil {
		http.Error(w, "Failed to encode home data", http.StatusInternalServerError)
	}
}

func loadBooksFromCSV(path string) ([]BookDataBase, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // tab instead of comma
	reader.Read()       // Skip header

	var books []BookDataBase
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		books = append(books, BookDataBase{
			SimplifiedName: record[0],
			Titulo:         record[1],
			Autor:          record[2],
			NivelEducativo: record[3],
			Materia:        record[4],
			Tipo:           record[5],
			Idioma:         record[6],
			Ilustraciones:  record[7],
			Genero:         record[8],
			Paginas:        record[9],
			Tamano:         record[10],
			DepositoLegal:  record[11],
			ISBN:           record[12],
			Edad:           record[13],
			FichaDidactica: record[14],
		})
	}

	return books, nil
}

func GetBookBySimplifiedName(name string) *BookDataBase {
	for _, book := range PersistentInfo {
		if book.SimplifiedName == name {
			return &book
		}
	}
	return nil // not found
}

func getTapaSmallURL(simplifiedName string) string {
	return fmt.Sprintf("https://otero-ediciones.s3.amazonaws.com/tapas/small/%s-tapa.jpg", simplifiedName)
}

func getTapaOriginalURL(simplifiedName string) string {
	return fmt.Sprintf("https://otero-ediciones.s3.amazonaws.com/tapas/originals/%s-tapa.jpg", simplifiedName)
}

func getSinopsisURL(simplifiedName string) string {
	return fmt.Sprintf("https://otero-ediciones.s3.amazonaws.com/sinopsis/%s.txt", simplifiedName)
}

func getBookFullInfo(simplifiedName string) BookCompleteInfo {
	var result BookCompleteInfo

	for _, book := range PersistentInfo {
		if book.SimplifiedName == simplifiedName {
			result.Titulo = book.Titulo
			result.TapaOriginalURL = getTapaOriginalURL(book.SimplifiedName)
			result.Autor = book.Autor
			result.Ilustraciones = book.Ilustraciones
			result.Materia = book.Materia
			result.NivelEducativo = book.NivelEducativo
			result.Genero = book.Genero
			result.GuiaDidactica = book.FichaDidactica
			result.Tamano = book.Tamano
			result.Paginas = book.Paginas
			result.ISBN = book.ISBN
			result.DepositoLegal = book.DepositoLegal
			result.Tipo = book.Tipo
			result.Idioma = book.Idioma

			if book.Tipo != "TEXTO_EDUCATIVO" {
				result.SinopsisURL = getSinopsisURL(simplifiedName)
			}
			break
		}
	}
	return result
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/catalogo/")
	if path == "" || strings.Contains(path, "/") {
		http.NotFound(w, r)
		return
	}

	simplifiedName := path
	book := getBookFullInfo(simplifiedName)

	if book.Titulo == "" {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func getBooksInfo(niveles, materias, tipos, idiomas []string, busqueda string, startIndex int16) []*BookResponse {
	var results []*BookResponse

	// Clamp start and end indexes to valid range
	if startIndex < 0 {
		startIndex = 0
	}

	var matched int16 = 0
	for _, book := range PersistentInfo {
		if !matchAny(book.NivelEducativo, niveles) {
			continue
		}
		if !matchAny(book.Materia, materias) {
			continue
		}
		if !matchAny(book.Tipo, tipos) {
			continue
		}
		if !matchAny(book.Idioma, idiomas) {
			continue
		}
		if busqueda != "" && !(strings.Contains(strings.ToLower(book.Titulo), strings.ToLower(busqueda)) ||
			strings.Contains(strings.ToLower(book.Autor), strings.ToLower(busqueda))) {
			continue
		}

		if matched < startIndex {
			matched++
			continue
		}

		if len(results) >= 9 {

			break
		}

		results = append(results, &BookResponse{
			Titulo:         book.Titulo,
			SimplifiedName: book.SimplifiedName,
			TapaSmallURL:   getTapaSmallURL(book.SimplifiedName),
		})
	}

	return results
}

func matchAny(value string, options []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))

	if len(options) == 0 {
		return true
	}
	for _, opt := range options {
		opt = strings.ToLower(strings.TrimSpace(opt))
		if opt == value {
			return true
		}
	}
	return false
}

func booksCatalogoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	nivelRaw := query.Get("nivel")
	var niveles []string
	if nivelRaw != "" {
		niveles = strings.Split(nivelRaw, ",")
	}

	materiaRaw := query.Get("materia")
	var materias []string
	if materiaRaw != "" {
		materias = strings.Split(materiaRaw, ",")
	}

	tipoRaw := query.Get("tipo")
	var tipos []string
	if tipoRaw != "" {
		tipos = strings.Split(tipoRaw, ",")
	}

	idiomaRaw := query.Get("idioma")
	var idiomas []string
	if idiomaRaw != "" {
		idiomas = strings.Split(idiomaRaw, ",")
	}

	busqueda := query.Get("busqueda")

	startIndex, _ := strconv.ParseInt(query.Get("startIndex"), 10, 16)

	// Fetch filtered books
	books := getBooksInfo(niveles, materias, tipos, idiomas, busqueda, int16(startIndex))

	json.NewEncoder(w).Encode(books)
}

func main() {
	var err error
	PersistentInfo, err = loadBooksFromCSV("books.tsv")
	if err != nil {
		log.Fatal("Failed to load books:", err)
	}
	log.Printf("Loaded %d books into memory.\n", len(PersistentInfo))

	MainPagePersistentInfo, err = loadMainInfoJson("frontpage_categories.json")
	if err != nil {
		log.Fatal("Failed to load frontpage_categories.json:", err)
	}
	log.Printf("Loaded %d frontpage categories.\n", len(MainPagePersistentInfo))

	http.HandleFunc("/home", booksHomeHandler)

	http.HandleFunc("/catalogo", booksCatalogoHandler)
	http.HandleFunc("/catalogo/", bookHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
