package server

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Credentials struct {
	DB_NAME                string
	DB_USER                string
	DB_PASSWORD            string
	ADMIN_LOGIN            string
	ADMIN_PASSWORD         string
	DB_TABLE_WITH_ARTICLES string
}

type Article struct {
	ID        int    `gorm:"column:id"`
	Header    string `gorm:"column:header"`
	Body      string `gorm:"column:body"`
	TimeStamp int64  `gorm:"column:time"`
}

type MarkdownArticle struct {
	ID             int
	MarkdownBody   template.HTML
	MarkdownHeader template.HTML
}

// отображение статьи по отдельности
func ArticlePageHandler(w http.ResponseWriter, r *http.Request) {
	adminCredentialsData, err := Readadmin_credentials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't Readadmin_credentials %v</p></html>", err)
		return
	}

	id, err := getArticleID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't getArticleID %v</p></html>", err)
		return
	}

	articl, err := getArticleByID(DB, id, &adminCredentialsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't find article with this id</p></html>")
		return
	}
	var markdownBody bytes.Buffer
	var markdownHeader bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	if err := md.Convert([]byte(articl.Header), &markdownHeader); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Failed to convert markdown %v</p></html>", err)
		return
	}
	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()), // Включает рендеринг "опасного" HTML (если нужно)
	)
	if err := md.Convert([]byte(articl.Body), &markdownBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Failed to convert markdown %v</p></html>", err)
		return
	}

	data := struct {
		MarkdownHeader template.HTML
		MarkdownBody   template.HTML
	}{
		MarkdownHeader: template.HTML(markdownHeader.String()),
		MarkdownBody:   template.HTML(markdownBody.String()), // Безопасное включение HTML
	}
	tmplPath := filepath.Join("templates", "article.html")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't paarse article.html</p></html>")
		return
	}
	t.Execute(w, data)

}

func getArticleID(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("articleid")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return -1, fmt.Errorf("invalid 'id'")
	}
	return id, nil
}
func ArticlesPageHandler(w http.ResponseWriter, r *http.Request) {
	adminCredentialsData, err := Readadmin_credentials()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read admin credentials: %v", err), http.StatusInternalServerError)
		return
	}

	// db, err := ConnectToDB(&adminCredentialsData)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to connect to DB: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	page, err := getPage(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get page: %v", err), http.StatusInternalServerError)
		return
	}

	limit := 3
	articles, err := getSomeArticles((page-1)*limit, limit, DB, &adminCredentialsData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get articles: %v", err), http.StatusInternalServerError)
		return
	}

	var markdownArticles []MarkdownArticle
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	for _, article := range articles {
		var markdownBody, markdownHeader bytes.Buffer
		var markdownArticle MarkdownArticle

		if err := md.Convert([]byte(article.Header), &markdownHeader); err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert markdown header: %v", err), http.StatusInternalServerError)
			return
		}

		if err := md.Convert([]byte(article.Body), &markdownBody); err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert markdown body: %v", err), http.StatusInternalServerError)
			return
		}
		markdownArticle.ID = article.ID
		markdownArticle.MarkdownBody = template.HTML(markdownBody.String())
		markdownArticle.MarkdownHeader = template.HTML(markdownHeader.String())
		markdownArticles = append(markdownArticles, markdownArticle)
	}

	data := struct {
		MarkdownArticles []MarkdownArticle
		Page             int
	}{
		MarkdownArticles: markdownArticles,
		Page:             page,
	}

	var tmplPath string
	if len(articles) == 0 {
		tmplPath = filepath.Join("templates", "page_not_exist.html")
	} else {
		tmplPath = filepath.Join("templates", "articles.html")
	}

	t, err := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"minus": func(a, b int) int { return a - b },
		"plus":  func(a, b int) int { return a + b },
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't parse template %s: %v", tmplPath, err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Can't execute template: %v", err), http.StatusInternalServerError)
		return
	}
}

// из ссылки http://127.0.0.1:8888/article?page=2 получает номер страницы
func getPage(r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return -1, fmt.Errorf("invalid 'page'")
	}
	return page, nil
}

// возвращает срез статей из БД начиная с fromPos в заданном количестве numberOfArticles
func getSomeArticles(fromPos int, numberOfArticles int, db *gorm.DB, adminCredentialsData *Credentials) ([]Article, error) {
	var articles []Article
	err := db.Table(adminCredentialsData.DB_TABLE_WITH_ARTICLES).
		Offset(fromPos).
		Limit(numberOfArticles).
		Find(&articles).Error
	if err != nil {
		return []Article{}, err
	}
	return articles, nil
}

func getArticleByID(db *gorm.DB, id int, cred *Credentials) (Article, error) {
	var article Article
	result := db.Table(cred.DB_TABLE_WITH_ARTICLES).First(&article, id)
	if result.Error != nil {
		return Article{}, result.Error
	}
	return article, nil
}

// читает файл admin_credentials.txt и заполняет структуру Credentials для входа в БД и логин на сайт
func Readadmin_credentials() (Credentials, error) {
	file, err := os.Open("admin_credentials.txt")
	if err != nil {
		return Credentials{}, fmt.Errorf("can't open admin_credentials")
	}

	defer file.Close()
	var adminCredentialsData Credentials
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		line := scaner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "DB_NAME":
			adminCredentialsData.DB_NAME = parts[1]
		case "DB_USER":
			adminCredentialsData.DB_USER = parts[1]
		case "DB_PASSWORD":
			adminCredentialsData.DB_PASSWORD = parts[1]
		case "ADMIN_LOGIN":
			adminCredentialsData.ADMIN_LOGIN = parts[1]
		case "ADMIN_PASSWORD":
			adminCredentialsData.ADMIN_PASSWORD = parts[1]
		case "DB_TABLE_WITH_ARTICLES":
			adminCredentialsData.DB_TABLE_WITH_ARTICLES = parts[1]
		}
		if err = scaner.Err(); err != nil {
			return Credentials{}, fmt.Errorf("can't parse admin_credentials")
		}

	}
	return adminCredentialsData, nil
}

func ConnectToDB(creds *Credentials) error {
	sqlConnect := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		creds.DB_USER, creds.DB_PASSWORD, creds.DB_NAME)
	var err error
	DB, err = gorm.Open(postgres.Open(sqlConnect), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// отражает главную страницу
func MainPage(w http.ResponseWriter, r *http.Request) {
	cread, err := Readadmin_credentials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't read admin_credentials in mainPage</p></html>")
		return
	}

	articles, err := getSomeArticles(0, 3, DB, &cread)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't getSomeArticles mainPage</p></html>")
		return
	}
	var markdownArticles []MarkdownArticle
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	for _, article := range articles {
		var markdownBody, markdownHeader bytes.Buffer
		var markdownArticle MarkdownArticle

		if err := md.Convert([]byte(article.Header), &markdownHeader); err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert markdown header: %v", err), http.StatusInternalServerError)
			return
		}

		if err := md.Convert([]byte(article.Body), &markdownBody); err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert markdown body: %v", err), http.StatusInternalServerError)
			return
		}
		markdownArticle.ID = article.ID
		markdownArticle.MarkdownBody = template.HTML(markdownBody.String())
		markdownArticle.MarkdownHeader = template.HTML(markdownHeader.String())
		markdownArticles = append(markdownArticles, markdownArticle)
	}

	data := struct {
		MarkdownArticles []MarkdownArticle
	}{
		MarkdownArticles: markdownArticles,
	}
	tmplParsed, err := template.ParseFiles("templates/mainPage.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse html</p></html>")
		return
	}
	tmplParsed.Execute(w, data)
}

// Секретный ключ для подписания JWT токенов
var jwtKey = []byte("your_secret_key")

// Структура для хранения данных JWT токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// страница логина для создания поста
func Login(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err == nil {
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err == nil && tkn.Valid {
			// Если токен валиден, перенаправляем на страницу /admin/insert
			http.Redirect(w, r, "/admin/insert", http.StatusSeeOther)
			return
		}
	}

	cred, err := Readadmin_credentials()
	if r.Method == http.MethodPost {
		if r.FormValue("username") == cred.ADMIN_LOGIN && r.FormValue("password") == cred.ADMIN_PASSWORD {
			expirationTime := time.Now().Add(1 * time.Minute)
			claims := &Claims{
				Username: r.FormValue("username"),
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Could not create token")
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

			http.Redirect(w, r, "/admin/insert", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/admin/wrong_auth", http.StatusSeeOther)
			return
		}
	}

	tmpl := filepath.Join("templates", "login.html")
	tmplParsed, err := template.ParseFiles(tmpl)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse html</p></html>")
		return
	}
	tmplParsed.Execute(w, nil)

}

// middleware чтобы нельзя было перейти по сслыке /admin/insert и писать пост
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/admin/wrong_auth", http.StatusUnauthorized)
				return
			} else {
				http.Redirect(w, r, "/admin/wrong_auth", http.StatusUnauthorized)
				return
			}
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/admin/wrong_auth", http.StatusUnauthorized)
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Can't parse html</p></html>")
			return
		}

		next.ServeHTTP(w, r)

	})

}

// пишем пост
func WritePost(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		cred, err := Readadmin_credentials()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse make_raticle.html</p></html>")
			return
		}

		// DB, err := ConnectToDB(&cred)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse make_raticle.html</p></html>")
		// 	return
		// }
		fmt.Println("ok")
		article := Article{
			Header:    r.FormValue("header"),
			Body:      r.FormValue("content"),
			TimeStamp: time.Now().Unix(),
		}
		if err := DB.Table(cred.DB_TABLE_WITH_ARTICLES).Create(&article).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse make_raticle.html</p></html>")
			return
		}
		http.Redirect(w, r, "/admin/insert", http.StatusSeeOther)
		return
	}
	tmpl := filepath.Join("templates", "make_article.html")
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Can't parse make_raticle.html</p></html>")
		return
	}
	t.Execute(w, t)
}

// неверная аунтефикация
func WrongAuthPage(w http.ResponseWriter, r *http.Request) {
	templ := filepath.Join("templates", "wrong_auth.html")
	t, err := template.ParseFiles(templ)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><h1>500 Status Internal Server Error</h1><p>Wrong LOGIN data</p></html>")
		return
	}
	t.Execute(w, nil)
}
