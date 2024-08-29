package goxi_v2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"mime/multipart"
	"net/http"
)

type WordPressLogic struct {
	URL           string
	BaseURI       string
	Authorization string
}

type Post struct {
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Status        string   `json:"status"`
	FeaturedMedia int      `json:"featured_media"`
	Categories    []int    `json:"categories"`
	Tags          []int    `json:"tags"`
	PostMeta      PostMeta `json:"post_meta"`
}

type PostMeta struct {
	CaoPrice         int                      `json:"cao_price"`
	CaoVipRate       int                      `json:"cao_vip_rate"`
	CaCloseNovicePay int                      `json:"cao_close_novip_pay"`
	CaoIsBoosVip     int                      `json:"cao_is_boosvip"`
	CaoExpireDay     int                      `json:"cao_expire_day"`
	CaoStatus        int                      `json:"cao_status"`
	CaoDownUrlNew    []map[string]interface{} `json:"cao_downurl_new"`
	CaoDemoUrl       string                   `json:"cao_demourl"`
	CaoDiyBtn        string                   `json:"cao_diy_btn"`
	CaoInfo          string                   `json:"cao_info"`
	CaPayum          int                      `json:"cao_paynum"`
}

func NewWordPressLogic(url string, secret string) *WordPressLogic {
	return &WordPressLogic{
		URL:           url,
		BaseURI:       url + "/wp-json/wp/v2",
		Authorization: "Basic " + base64.StdEncoding.EncodeToString([]byte(secret)),
	}
}

// GetCategories 获取分类
func (l *WordPressLogic) GetCategories() ([]map[string]interface{}, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Get(l.BaseURI + "/categories?per_page=100")
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetTags 获取标签
func (l *WordPressLogic) GetTags(page string, size string, search string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tags?page=%s&per_page=%s&search=%s", l.BaseURI, page, size, search)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Get(url)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	return data, err
}

// CreateTag 创建标签
func (l *WordPressLogic) CreateTag(name string) (map[string]interface{}, error) {
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("name", name)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, l.BaseURI+"/tags", payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", l.Authorization)
	req.Header.Add("Content-Type", "multipart/form-data;")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPosts 获取文章
func (l *WordPressLogic) GetPosts(page string, size string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/posts?page=%s&per_page=%s", l.BaseURI, page, size)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Get(url)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	return data, err
}

func (l *WordPressLogic) GetMedia(page string, size string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/media?page=%s&per_page=%s", l.BaseURI, page, size)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Get(url)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	return data, err
}

func (l *WordPressLogic) Publish(p Post) (map[string]interface{}, error) {
	postParams := map[string]interface{}{
		"title":          p.Title,
		"content":        p.Content,
		"status":         p.Status,
		"featured_media": p.FeaturedMedia,
		"categories":     p.Categories,
		"tags":           p.Tags,
	}
	postData, err := json.Marshal(postParams)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", l.BaseURI+"/posts", bytes.NewBuffer(postData))
	req.Header.Add("Authorization", l.Authorization)
	req.Header.Set("Content-Type", "application/json")
	cateResp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(cateResp.Body)
	body, _ := io.ReadAll(cateResp.Body)
	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		panic(err)
	}
	// 创建文章元数据
	data, err2 := l.PutMetaData(fmt.Sprintf("%v", respData["id"]), p.PostMeta)
	if err2 != nil {
		return nil, err2
	}
	// 合并data，respData
	for k, v := range data {
		respData[k] = v
	}
	return respData, err
}

// PutMetaData
// @params "a:1:{i:0;a:3:{s:4:\"name\";s:12:\"下载地址\";s:3:\"url\";s:47:\"https://pan.baidu.com/s/17VyHGbYbTU8wauc4tf49ZQ\";s:3:\"pwd\";s:4:\"ca6f\";}}"
func (l *WordPressLogic) PutMetaData(id string, t PostMeta) (map[string]interface{}, error) {
	postData, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	baseURI := fmt.Sprintf(l.URL+"/wp-json/xnf/v3/metadata/%s?Authorization=%s", id, l.Authorization)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", baseURI, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")
	cateResp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(cateResp.Body)
	body, _ := io.ReadAll(cateResp.Body)
	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	return respData, err
}

// FeaturedMedia 特色图片
func (l *WordPressLogic) FeaturedMedia(file *multipart.FileHeader) (map[string]interface{}, error) {
	client := &http.Client{}
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileBytes)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", l.BaseURI+"/media", body)
	req.Header.Add("Authorization", l.Authorization)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body2, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	err = json.Unmarshal(body2, &respData)
	return respData, err
}

// DeleteMedia 删除图片
func (l *WordPressLogic) DeleteMedia(id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/media/%s", l.BaseURI, id)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Delete(url)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	return data, err
}

// GetUsers 获取用户
func (l *WordPressLogic) GetUsers(page string, size string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/users?page=%s&per_page=%s", l.BaseURI, page, size)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", l.Authorization).
		Get(url)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	return data, err
}
