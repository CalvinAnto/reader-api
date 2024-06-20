package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Manga struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Title_JP    string `json:"title_jp"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Thumbnail   string `json:"thumbnail"`
	Pages       []Page `json:"pages"`
}

type Page struct {
	PageNo    int    `json:"page_no"`
	Url       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connStr := os.Getenv("CONNSTR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/manga", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Get("/manga/:id", func(c *fiber.Ctx) error {
		return getMangaByIdHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {

	var manga Manga
	var mangas []Manga

	rows, err := db.Query("SELECT id, title, description, cover, thumbnail FROM manga")
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&manga.ID, &manga.Title, &manga.Description, &manga.Cover, &manga.Thumbnail)
		mangas = append(mangas, manga)
	}
	// return c.Render("index", fiber.Map{
	// 	"Mangas": mangas,
	// })

	return c.JSON(mangas)
}

func getMangaByIdHandler(c *fiber.Ctx, db *sql.DB) error {

	mangaId, _ := strconv.Atoi(c.Params("id"))
	var manga Manga

	query := "SELECT id, title, description, cover, thumbnail FROM manga WHERE id = $1"
	rows, err := db.Query(query, mangaId)
	defer rows.Close()

	if err != nil {
		log.Println("A")
		log.Fatalln(err)
		c.JSON("An error occured")
	}

	rows.Next()

	rows.Scan(&manga.ID, &manga.Title, &manga.Description, &manga.Cover, &manga.Thumbnail)

	rows.Close()

	rows, err = db.Query("SELECT page_no, url, thumbnail FROM page WHERE manga_id = $1", mangaId)

	defer rows.Close()

	if err != nil {
		log.Println("B")
		log.Fatalln(err)
		c.JSON("An error occured")
	}

	var page Page
	var pages []Page

	for rows.Next() {
		rows.Scan(&page.PageNo, &page.Url, &page.Thumbnail)
		pages = append(pages, page)
	}

	manga.Pages = pages

	return c.JSON(manga)

}
