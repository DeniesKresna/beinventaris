package Controllers

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var SessionId uint

// Result ..
type Result struct {
	CurrentPage  int         `json:"current_page"`
	Data         interface{} `json:"data"`
	FirstPageURL string      `json:"first_page_url"`
	From         int         `json:"from"`
	LastPage     int         `json:"last_page"`
	LastPageURL  string      `json:"last_page_url"`
	NextPageURL  string      `json:"next_page_url"`
	Path         string      `json:"path"`
	PerPage      int         `json:"per_page"`
	PrevPageURL  string      `json:"prev_page_url"`
	To           int64       `json:"to"`
	Total        int64       `json:"total"`
}

// Paginator ..
type Paginator interface {
	Paginate(db *gorm.DB) (Result, *gorm.DB)
}

// Config ..
//
// If you don't provide app url it will fetch the APP_URL from environment
type PConfig struct {
	Page    int
	Sort    string
	PerPage int
	AppURL  string
	Path    string
}

// Paginate ..
func (c *PConfig) Paginate(db *gorm.DB, any interface{}, count int64) (Result, *gorm.DB) {
	var r Result

	offset := (c.Page - 1) * c.PerPage
	d := db.Offset(offset).Limit(c.PerPage)

	if c.Sort != "" {
		d.Order(c.Sort)
	}

	db.Find(any)

	var lastIndex int64 = int64(c.PerPage) * int64(c.Page)
	if lastIndex > count {
		lastIndex = count
	}

	r.CurrentPage = c.Page
	r.NextPageURL = c.GetPageURL(c.Page + 1)
	r.FirstPageURL = c.GetPageURL(1)
	r.PrevPageURL = c.PreviousPageURL()
	r.PerPage = c.PerPage
	r.Path = c.Path
	r.To = lastIndex
	r.From = offset + 1
	r.Total = count
	r.Data = any
	r.LastPageURL = c.GetPageURL(r.GetLastPage())
	r.LastPage = r.GetLastPage()

	return r, d
}

// GetLastPage ..
func (r *Result) GetLastPage() int {
	return int(math.Ceil(float64(r.Total) / float64(r.PerPage)))
}

func (c *PConfig) GetPageURL(page int) string {
	return fmt.Sprintf("%s%s?page=%d&per_page=%d", c.GetAppURL(), c.Path, page, c.PerPage)
}

// PreviousPageURL ..
func (c *PConfig) PreviousPageURL() string {
	pageNumber := 1

	if c.Page > 1 {
		pageNumber = c.Page - 1
	}

	return c.GetPageURL(pageNumber)
}

// GetAppURL ..
func (c *PConfig) GetAppURL() string {
	if c.AppURL == "" {
		return os.Getenv("APP_URL")
	}

	return c.AppURL
}

func Debug(d interface{}) {
	str := fmt.Sprintf("%v", d)
	fmt.Println(str)
}

func SetSessionId(c *gin.Context) {
	sessionid, _ := c.Get("sessionId")
	SessionId = sessionid.(uint)
}

// Inject Struct Field Value from field in source struct to destination struct
func InjectStruct(sourceSctruct interface{}, destinationStruct interface{}) {
	var srcReflectValue = reflect.ValueOf(sourceSctruct)
	if srcReflectValue.Kind() == reflect.Ptr {
		srcReflectValue = srcReflectValue.Elem()
	}

	var srcReflectType = srcReflectValue.Type()

	var dstReflectValue = reflect.ValueOf(destinationStruct)
	if dstReflectValue.Kind() == reflect.Ptr {
		dstReflectValue = dstReflectValue.Elem()
	}

	var dstReflectType = dstReflectValue.Type()

	for i := 0; i < srcReflectType.NumField(); i++ {
		srcField := srcReflectType.Field(i).Name
		_, ok := dstReflectType.FieldByName(srcField)
		if ok {
			dstReflectValue.FieldByName(srcField).Set(srcReflectValue.FieldByName(srcField))
		}
	}
}

// Filter data on from all model field by the search param
func FilterModel(search string, sourceModel interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var srcReflectValue = reflect.ValueOf(sourceModel)
		if srcReflectValue.Kind() == reflect.Ptr {
			srcReflectValue = srcReflectValue.Elem()
		}

		var srcReflectType = srcReflectValue.Type()

		for i := 0; i < srcReflectValue.NumField(); i++ {
			fieldType := srcReflectType.Field(i).Type
			if (fieldType.Kind() >= 3 && fieldType.Kind() <= 12) || fieldType.Kind() == 24 {
				srcField := ToSnakeCase(srcReflectType.Field(i).Name)

				if i == 0 {
					db.Where(srcField+" LIKE ?", "%"+search+"%")
				} else {
					db.Or(srcField+" LIKE ?", "%"+search+"%")
				}
			}
		}

		return db
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func DownloadDocuments(c *gin.Context) {
	mediaFile := c.Query("path")
	f, err := os.Open("./" + mediaFile)
	if err != nil {
		Response.Json(c, 500, err)
		return
	}
	defer f.Close()
	var filename = strings.Split(mediaFile, "/")
	fmt.Println(filename[len(filename)-1])
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename[len(filename)-1])

	io.Copy(c.Writer, f)
}

func Copy(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
