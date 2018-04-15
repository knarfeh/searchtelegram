package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/knarfeh/searchtelegram/server/domain"

	"github.com/labstack/echo"
	tb "github.com/tucnak/telebot"
	elastic "gopkg.in/olivere/elastic.v5"
	// "gopkg.in/telegram-bot-api.v4"
	// "github.com/goinggo/mapstructure"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// OKResponse ...
var OKResponse = map[string]interface{}{"ok": true}

// Bind attaches api routes
func (api *API) Bind(group *echo.Group) {
	group.POST("/tgbot", api.TgBotWebhook)
	group.PUT("/v1/tg", api.UpdateTgResource)
	group.POST("/v1/tg", api.CreateTgResource)
	group.GET("/v1/tg/:tgID", api.GetTgResource)
	group.GET("/v1/tg/:tgID/exist", api.CheckTgResourceExist)
	group.GET("/v1/search", api.SearchTgResource)
	group.GET("/v1/tags", api.GetTgTags)
}

// GetTgTags ...
func (api *API) GetTgTags(c echo.Context) error {
	// TODO: Add cache
	app := c.Get("app").(*App)

	agg := elastic.NewTermsAggregation().Field("tags.tgid.keyword")
	search := app.ESClient.Client.Search().Index("telegram").Type("resource")
	allTags, _ := search.Aggregation("allTags", agg).Do(context.TODO())
	instance := domain.NewTgTagBuckets()
	json.Unmarshal(*allTags.Aggregations["allTags"], instance)
	result := make(map[string]interface{})
	result["results"] = instance.Buckets
	result["total"] = len(instance.Buckets)
	return c.JSON(http.StatusOK, result)
}

// SearchTgResource ...
func (api *API) SearchTgResource(c echo.Context) error {
	app := c.Get("app").(*App)

	queryString := c.FormValue("q")
	tags := c.FormValue("tags")
	page := c.FormValue("page")
	pageSize := c.FormValue("page_size")

	size := 10
	from := 0
	if pageSize != "" {
		size, _ = strconv.Atoi(pageSize)
	}
	if page != "" {
		intPage, _ := strconv.Atoi(page)
		from = (intPage - 1) * size
	}
	if tags == "" {
		tags = "people group channel"
	} else {
		tags = strings.Replace(tags, ",", " ", -1)
	}
	if queryString == "" {
		queryString = "*"
	}

	app.Engine.Logger.Infof("query: %s, tags: %s, page: %s, pageSize: %s", queryString, tags, page, pageSize)
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Filter(elastic.NewQueryStringQuery(queryString))
	boolQuery = boolQuery.Must(elastic.NewMoreLikeThisQuery().LikeText(tags).Field(
		"tags.tgid",
	).MinDocFreq(0).MinTermFreq(0))
	search := app.ESClient.Client.Search().Index("telegram").Type("resource").From(from).Size(size)
	searchResult, err := search.Query(boolQuery).Do(context.TODO())
	if err != nil {
		panic(err)
	}

	values := make([]*domain.TgResource, 0, searchResult.Hits.TotalHits)
	for _, hit := range searchResult.Hits.Hits {
		instance := domain.NewTgResource()
		json.Unmarshal(*hit.Source, instance)
		values = append(values, instance)
	}
	result := make(map[string]interface{})
	result["total"] = searchResult.TotalHits()
	result["from"] = from
	result["size"] = size
	result["results"] = values
	return c.JSON(http.StatusOK, result)
}

// CheckTgResourceExist ...
func (api *API) CheckTgResourceExist(c echo.Context) error {
	return c.JSON(404, "TODO")
}

// GetTgResource ...
func (api *API) GetTgResource(c echo.Context) error {
	app := c.Get("app").(*App)

	tgid := c.Param("tgid")
	fmt.Printf("resource tgid: %s", tgid)

	resourceResult, err := app.ESClient.Client.Get().Index("telegram").Type("resource").Id(tgid).Do(context.TODO())
	if err != nil {
		e, _ := err.(*elastic.Error)
		if e.Status == 404 {
			message := tgid + " is not exist"
			errorItem := make(map[string]string)
			errorItem["code"] = "resource_not_exist"
			errorItem["message"] = message
			errorItem["source"] = "10000"
			c.JSON(404, errorItem)
		}
	}

	tgResource := domain.NewTgResource()
	json.Unmarshal(*resourceResult.Source, tgResource)
	return c.JSON(http.StatusOK, tgResource)
}

// CreateTgResource ...
func (api *API) CreateTgResource(c echo.Context) error {
	// TODO: check already exist
	app := c.Get("app").(*App)
	tgResource := domain.NewTgResource()
	fmt.Println("Create tg resource with: ", *tgResource)
	if err := c.Bind(tgResource); err != nil {
		return err
	}
	if err := c.Validate(tgResource); err != nil {
		return err
	}

	tgResouceString, _ := json.Marshal(tgResource)
	app.Engine.Logger.Infof("Create tg resource: %s", tgResouceString)
	err := app.RedisClient.Client.Publish("st_submit", string(1)).Err()
	app.RedisClient.Client.LPush("st_submit_list", string(tgResouceString))
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, OKResponse)
}

// UpdateTgResource ... patch may be a better choice
func (api *API) UpdateTgResource(c echo.Context) error {
	app := c.Get("app").(*App)
	tgResource := domain.NewTgResource()
	if err := c.Bind(tgResource); err != nil {
		return err
	}
	if err := c.Validate(tgResource); err != nil {
		return err
	}

	// TODO: get type from telegram api
	updateResult, err := app.ESClient.Client.Update().Index("telegram").Type("resource").Id(tgResource.TgID).Doc(tgResource).Do(context.TODO())
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	app.Engine.Logger.Printf("New version of tgResource %q is now %d", updateResult.Id, updateResult.Version)
	return c.JSON(http.StatusOK, updateResult)
}

// TgBotWebhook ...
func (api *API) TgBotWebhook(c echo.Context) error {
	app := c.Get("app").(*App)

	update := &tb.Update{}
	if err := c.Bind(update); err != nil {
		return err
	}

	messageString, _ := json.Marshal(update)
	app.Engine.Logger.Printf("updateString: %s\n", messageString)
	app.TgBot.incommingUpdate(update, app)
	return c.JSON(http.StatusOK, OKResponse)

}
