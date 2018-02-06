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
	elastic "gopkg.in/olivere/elastic.v5"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// OKResponse ...
var OKResponse = map[string]interface{}{"ok": true}

// Bind attaches api routes
func (api *API) Bind(group *echo.Group) {
	group.PUT("/v1/tg", api.UpdateTgResource)
	group.POST("/v1/tg", api.CreateTgResource)
	group.GET("/v1/tg/:name", api.GetTgResource)
	group.GET("/v1/tg/:name/exist", api.CheckTgResourceExist)
	group.GET("/v1/search", api.SearchTgResource)
	group.GET("/v1/tags", api.GetTgTags)
}

// GetTgTags ...
func (api *API) GetTgTags(c echo.Context) error {
	// TODO: Add cache
	app := c.Get("app").(*App)

	agg := elastic.NewTermsAggregation().Field("tags.name.keyword")
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
		"tags.name",
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

	name := c.Param("name")
	fmt.Printf("resource name: %s", name)

	resourceResult, err := app.ESClient.Client.Get().Index("telegram").Type("resource").Id(name).Do(context.TODO())
	if err != nil {
		e, _ := err.(*elastic.Error)
		if e.Status == 404 {
			message := name + " is not exist"
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
	app := c.Get("app").(*App)
	tgResource := domain.NewTgResource()
	if err := c.Bind(tgResource); err != nil {
		return err
	}
	if err := c.Validate(tgResource); err != nil {
		return err
	}

	// TODO: get type from telegram api
	_, err := app.ESClient.Client.Index().OpType("create").Index("telegram").Type("resource").Id(tgResource.Name).BodyJson(tgResource).Do(context.TODO())
	if err != nil {
		// Please make sure domain not exist
		e, _ := err.(*elastic.Error)
		if e.Status == 409 {
			errorItem := make(map[string]string)
			errorItem["code"] = "resource_already_exist"
			errorItem["message"] = e.Details.Reason
			errorItem["source"] = "10001"
			return c.JSON(e.Status, errorItem)
		}
		// Should not happen...
		panic(err)
	}

	// TODO, response
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
	updateResult, err := app.ESClient.Client.Update().Index("telegram").Type("resource").Id(tgResource.Name).Doc(tgResource).Do(context.TODO())
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	app.Engine.Logger.Infof("New version of tgResource %q is now %d", updateResult.Id, updateResult.Version)
	return c.JSON(http.StatusOK, updateResult)
}
