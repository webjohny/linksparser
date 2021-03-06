package services

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CountryList struct{
	XMLName xml.Name `xml:"country-list"`
	Header struct{
		Name string `xml:"name"`
		FullName string `xml:"fullname"`
		English string `xml:"english"`
		Alpha2 string `xml:"alpha2"`
		Alpha3 string `xml:"alpha3"`
		Iso string `xml:"iso"`
		Location string `xml:"location"`
		LocationPrecise string `xml:"location-precise"`
	} `xml:"header"`
	Country []struct{
		Name string `xml:"name"`
		FullName string `xml:"fullname"`
		English string `xml:"english"`
		Alpha2 string `xml:"alpha2"`
		Alpha3 string `xml:"alpha3"`
		Iso string `xml:"iso"`
		Location string `xml:"location"`
		LocationPrecise string `xml:"location-precise"`
	} `xml:"country"`
}

func GetCountryList () (*CountryList, error) {
	var result *CountryList

	if xmlBytes, err := GetXML("https://www.artlebedev.ru/country-list/xml"); err != nil {
		log.Printf("Failed to get XML: %v", err)
		return nil, err
	} else {
		err = xml.Unmarshal(xmlBytes, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func SetTmpl(tmpl string, params map[string]string) string{
	for k, v := range params {
		tmpl = strings.ReplaceAll(tmpl, "{" + k + "}", v)
	}
	return tmpl
}

func GetXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func IsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear

}

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\":"\\\\", "'":`\'`, "\\0":"\\\\0", "\n":"\\n", "\r":"\\r", `"`:`\"`, "\x1a":"\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func ArrayRand(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(arr)
	return strings.Trim(arr[n], " ")
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ErrorHandler(err error) {
	if err != nil {
		log.Println(err)
	}
}

func SentenceSplit(str string) []string {
	return PregSplit(str, `([.?!])\s+`)
}

func PregSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes) + 1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func YoutubeEmbed(str string) string {
	if strings.Contains(str, "youtube.com/watch?v=") {
		str = strings.Replace(str, "youtube.com/watch?v=", "youtube.com/embed/", 1)
	} else if strings.Contains(str, "youtu.be/") {
		str = strings.Replace(str, "youtu.be/", "youtube.com/embed/", 1)
	} else if strings.Contains(str, "/watch?v=") {
		str = strings.Replace(str, "/watch?v=", "youtube.com/embed/", 1)
	}
	str = strings.Replace(str, "&", "?", 1)
	if !strings.Contains(str, "https") {
		return `https://` + str
	}
	return str

}


func Format(str string) string {
	reg := regexp.MustCompile(`(?m)<div[^<>]*><a[^<>]*>More items...</a></div>`)
	str = reg.ReplaceAllString(str, ``)

	reg = regexp.MustCompile(`(?m)<div[^<>]*>•</div>`)
	str = reg.ReplaceAllString(str, ``)

	reg = regexp.MustCompile(`(?m)<div[^<>]*>[JFMASOND][a-z]{2}\s\d{1,2},\s\d{4}</div>`)
	str = reg.ReplaceAllString(str, ``)

	reg = regexp.MustCompile(`(?m)<span[^<>]*>[JFMASOND][a-z]{2}\s\d{1,2},\s\d{4}</span>`)
	str = reg.ReplaceAllString(str, ``)

	headingMatch := PregMatch(`(?m)<div[^<>]*role="heading"><b>(?P<title>.+)</b></div>`, str)
	heading := headingMatch["title"]
	heading = StripTags(heading)

	str, _ = sanitize.HTMLAllowing(str, []string{
		"table", "thead", "tbody", "tr", "td", "th",
		"ol", "ul", "li",
		"dl", "dt", "dd",
		"p", "br"})

	str = strings.Replace(str, "...", ".", -1)

	str = strings.Replace(str, heading, "<h3>" + heading + "</h3>", 1)

	return str
}

func StripTags(html string) string {
	paaReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(paaReader)
	if err != nil {
		log.Println("Utils.StripTags.HasError", err)
	}

	return doc.Text()
}

func RandBool() bool {
	return rand.Float32() < 0.5
}

func ParseFormCollection(r *http.Request, typeName string) map[string]string {
	result := make(map[string]string)
	if err := r.ParseForm(); err != nil {
		log.Println("Utils.ParseFormCollection.HasError", err)
	}
	for key, values := range r.Form {
		re := regexp.MustCompile(typeName + "\\[(.+)\\]")
		matches := re.FindStringSubmatch(key)

		if len(matches) >= 2 {
			result[matches[1]] = values[0]
		}
	}
	return result
}

func ToInt(value string) int {
	var integer int = 0
	if value != "" {
		integer, _ = strconv.Atoi(value)
	}
	return integer
}

func PregMatch(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}


// Counters - work with mutex
type Counters struct {
	mx sync.Mutex
	m map[string]int
}

func (c *Counters) Load(key string) (int, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *Counters) Store(key string, value int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func NewCounters() *Counters {
	return &Counters{
		m: make(map[string]int),
	}
}