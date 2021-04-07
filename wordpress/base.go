package wordpress

import (
	"fmt"
	"github.com/gosimple/slug"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"time"
)


type WpCat struct {
	Description string `json:"description"`
	Filter string `json:"filter"`
	Name string `json:"name"`
	Parent int `json:"parent"`
	Slug string `json:"slug"`
	Taxonomy string `json:"taxonomy"`
	TermGroup int `json:"term_group"`
	TermId int `json:"term_id"`
	TermTaxonomyId int `json:"term_taxonomy_id"`
}

type WpPost struct {
	Id int
	Title string
	Content string
	Date time.Time
	Link string
	Slug string
	Parent int
	Terms []WpCat
}

type Base struct {
	client *Client
	cnf []interface{}
	err error
}

type WpImage struct {
	Id int
	Url string
	UrlMedium string
}


func isNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}

func randStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func toInt(value string) int {
	var integer int = 0
	if value != "" {
		integer, _ = strconv.Atoi(value)
	}
	return integer
}

func (w *Base) Connect(url string, username string, password string, blogId int) *Client {

	resp, _ := http.PostForm(url + `/wp-admin/conn.php`, nil)

	if resp == nil || resp.StatusCode != 200 {
		resp, _ = http.PostForm(url + `/xmlrpc2.php`, nil)
		if resp == nil || resp.StatusCode != 200 {
			resp, _ = http.PostForm(url + `/xmlrpc.php`, nil)
			if resp == nil || resp.StatusCode != 200 {
				return nil
			}else{
				url += `/xmlrpc.php`
			}
		}else{
			url += `/xmlrpc2.php`
		}
	}else{
		url += `/wp-admin/conn.php`
	}

	c := NewClient(url, UserInfo{
		username,
		password,
	})
	w.client = c
	w.cnf = []interface{}{
		blogId, username, password,
	}
	return c
}

func (w *Base) GetError() error {
	return w.err
}

func (w *Base) PrepareCat(cat map[string]interface{}) WpCat {
	parentId, _ := strconv.Atoi(cat["parent"].(string))
	termGroup, _ := strconv.Atoi(cat["term_group"].(string))
	termId, _ := strconv.Atoi(cat["term_id"].(string))
	termTaxonomyId, _ := strconv.Atoi(cat["term_taxonomy_id"].(string))
	var description string
	if cat["description"] != nil {
		description = cat["description"].(string)
	}
	return WpCat{
		Description:    description,
		Filter:         cat["filter"].(string),
		Name:           cat["name"].(string),
		Parent:         parentId,
		Slug:           cat["slug"].(string),
		Taxonomy:       cat["taxonomy"].(string),
		TermGroup:      termGroup,
		TermId:         termId,
		TermTaxonomyId: termTaxonomyId,
	}
}

func (w *Base) PreparePost(post map[string]interface{}) WpPost {
	parent, _ := strconv.Atoi(post["post_parent"].(string))
	var cats []WpCat
	terms := post["terms"].([]interface{})
	if len(terms) > 0 {
		for _, item := range terms {
			cat := item.(map[string]interface{})
			cats = append(cats, w.PrepareCat(cat))
		}
	}
	id, _ := strconv.Atoi(post["post_id"].(string))

	wpPost := WpPost{
		Id: id,
		Date: post["post_date"].(time.Time),
		Parent: parent,
		Terms: cats,
	}
	if !isNil(post["post_content"]){
		wpPost.Content = post["post_content"].(string)
	}
	if !isNil(post["post_title"]){
		wpPost.Title = post["post_title"].(string)
	}
	if !isNil(post["post_name"]){
		wpPost.Slug = post["post_name"].(string)
	}
	if !isNil(post["link"]){
		wpPost.Link = post["link"].(string)
	}
	return wpPost
}

func (w *Base) GetCats() []WpCat {
	var result interface{}
	result, err := w.client.Call(`wp.getTerms`, struct {
		BlogId int `xml:"blogId"`
		Username string `xml`
	}{}append(
		w.cnf, "category",
	))
	if err != nil {
		w.err = err
		log.Println("Wordpress.GetCats.HasError", err)
	}
	var cats []WpCat
	if result != nil {
		res := result.([]interface{})
		if len(res) > 0 {
			for _, item := range res {
				cat := item.(map[string]interface{})
				cats = append(cats, w.PrepareCat(cat))
			}
		}
	}
	return cats
}


func (w *Base) NewTerm(name string, taxonomy string, slug string, description string, parentId int) int {
	params := map[string]string{
		"name": name,
		"taxonomy": taxonomy,
	}

	if slug != "" {
		params["slug"] = slug
	}

	if description != "" {
		params["description"] = description
	}

	if parentId > 0 {
		params["parent"] = strconv.Itoa(parentId)
	}

	result, err := w.client.Call(`wp.newTerm`, append(
		w.cnf, params,
	))
	if err != nil {
		w.err = err
		log.Println("Wordpress.NewTerm.HasError", err)
		return 0
	}

	fmt.Println(result)

	return 23
}

func (w *Base) GetPost(id int) WpPost {
	return WpPost{}
}

func (w *Base) EditPost(id int, title string, content string) bool {
	return true
}

func (w *Base) NewPost(title string, content string, catId int, photoId int) int {

	return 12
}

func (w *Base) CheckConn() bool {
	return w.client != nil
}

func (w *Base) UploadFile(url string, postId int, bytes *[]byte, encoded bool) (WpImage, error) {


	return WpImage{}, nil
}

func (w *Base) CatIdByName(name string) int {
	var catId int

	// Загружаем список категорий
	cats := w.GetCats()

	// Создавать ли категорию
	create := true

	// Пробегаем по всем категориям
	if len(cats) > 0 {
		for _, cat := range cats {
			// Проверка существования категории
			if cat.Name == name {
				catId = cat.TermId
				create = false
				break
			}
		}
	}

	// Создаём категорию
	if create {
		catId = w.NewTerm(name, "category", slug.Make(name), "", 0)
	}

	return catId
}
