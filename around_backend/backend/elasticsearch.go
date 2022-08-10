package backend

import (
	"around/constants"
	"around/util"
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

var (
	//ESBackend is the pointer, it points to -> ElasticsearchBackend
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
	client *elastic.Client
}

// initiate the elastic search service each time the service starts and create the index(database)
func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	//create new client with ES internal IP, username and password to log in
	client, err := elastic.NewClient(
		elastic.SetURL(config.Address),
		elastic.SetBasicAuth(config.Username, config.Password))
	if err != nil {
		panic(err)
	}

	//check if index exists, return exists is true or false
	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	//if index does not exist, then we create POST INDEX
	//keyword and text are both string, but keyword we have to 100% match, text we do not need to
	if !exists {
		mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" }, 
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		_, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	//if index exists
	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	//if index does not exist, we create USER INDEX
	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
		_, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")

	//{first client is myclient that created in line 22, second client is in line 16}
	//it is like constructor, use a new (client: client) to get elasticsearchbackend
	//then use pointer ESBackend to point it
	ESBackend = &ElasticsearchBackend{client: client}
}

// read the result from elastic search and return it
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	searchResult, err := backend.client.Search().
		//Below is like SQL, select * from index where...
		Index(index). //search in Index(database) index
		Query(query). //specify the query and add it
		Pretty(true).
		Do(context.Background()) //excute
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
	_, err := backend.client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
	_, err := backend.client.DeleteByQuery().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())

	return err
}
