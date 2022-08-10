package service

import (
	"mime/multipart"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

// support user-based search.
func SearchPostsByUser(user string) ([]model.Post, error) {
	query := elastic.NewTermQuery("user", user)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

// support keywords-based search.
func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	//here it does not need mewtermquery, because keywords can be part of the content, not like user 100% match
	query := elastic.NewMatchQuery("message", keywords)

	//there are possible more than 1 keyword, relationship we use AND
	query.Operator("AND")

	//if no keyword provided, return all the contents
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}

	//packge.variable pointer.function()
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

// get previous search result and transfer return result to []model.Post
func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post
	var posts []model.Post

	//reflect.TypeOf() == java.instanceOf(), to check each result if it is post-type result, if yes go loop
	//because ES is non-relational databse not like MYSQL
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		//item.() == casting, it is to cast item type to post type
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}

func SavePost(post *model.Post, file multipart.File) error {
	//create the media file link with SaveToGCS return result
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}

	//update the post with media link added
	post.Url = medialink

	//save updated post: id, url, message, user type to the ES
	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}

func DeletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
