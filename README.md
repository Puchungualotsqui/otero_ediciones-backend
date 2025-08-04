# Otero Ediciones - Book Catalog Backend

This is the official **Golang backend** for the [Otero Ediciones](https://otero-ediciones.com) landing page. It reads book metadata from a `.tsv` file and serves filtered catalog results and full book details via a simple HTTP API. It also supports dynamic generation of image and synopsis URLs stored in AWS S3.

## Features

- Parses book data from a local `books.tsv` file
- Serves:
    - Home rows (`/home`)
    - Filtered catalog (`/catalogo`)
    - Book detail by ID (`/catalogo/{simplifiedName}`)
- Supports query-based filtering: level, subject, type, language, and search
- Provides pagination using `startIndex`
- AWS S3 integration for tapa and sinopsis URLs
- Lightweight, no database required

## Project Structure
```
.
├── books.tsv # Main book dataset (tab-separated)
├── frontpage_categories.json # Categories and rows shown on the landing page
├── main.go # Main application logic
├── Dockerfile # Optional for containerized deployment
└── README.md # This file
```


## Data Format

**books.tsv** must contain 15 tab-separated columns in this order:
```
simplified_name, titulo, autor, nivel_educativo, materia, tipo,
idioma, ilustraciones, genero, paginas, tamano,
deposito_legal, isbn, edad, ficha_didactica
```


**frontpage_categories.json** must contain:

```json
[
  {
    "titulo": "Novedades Primaria",
    "libros": ["el_hogar_de_los_pajaros", "asi_era_la_vida", ...]
  },
  ...
]
```

## Endpoints
### GET /home

Returns curated rows of books for the homepage.

### GET /catalogo

Returns up to 9 filtered books from the catalog.

Query Parameters:

    nivel (comma-separated)

    materia (comma-separated)

    tipo (comma-separated)

    idioma (comma-separated)

    busqueda (string)

    startIndex (integer)

Example:
```
GET /catalogo?nivel=Primaria&materia=Literatura&startIndex=0
```

## 📄 License

This project is licensed under the  
**Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International (CC BY-NC-ND 4.0)**

You may view the source, but you may not modify, redistribute, or use it commercially.

© 2025 Otero Ediciones  
[Read full license](https://creativecommons.org/licenses/by-nc-nd/4.0/)