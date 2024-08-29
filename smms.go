package goxi_v2

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"strconv"
)

type SmmSLogic struct {
	BaseURI string
	Token   string
	Key     string
	Proxy   string
}

func NewSmmSLogic(key string, proxy string) *SmmSLogic {
	if proxy == "" {
		proxy = "http://localhost:9999"
	}
	return &SmmSLogic{
		BaseURI: "https://sm.ms/api/v2/",
		Token:   "",
		Key:     key,
		Proxy:   proxy,
	}
}

type TokenResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
	RequestId string `json:"RequestId"`
}

// GetToken User - Get API-Token
func (l *SmmSLogic) GetToken(username string, password string) (string, error) {
	logrus.Info("初始化登录SMMS")
	client := resty.New()
	client.SetProxy(l.Proxy)
	r, err := client.R().
		SetHeader("Authorization", l.Key).
		SetFormData(map[string]string{
			"username": username,
			"password": password,
		}).
		Post(l.BaseURI + "token")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var resp TokenResponse
	if err := json.Unmarshal(r.Body(), &resp); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	// 提取 token
	token := resp.Data.Token
	fmt.Println("Token:", token)
	l.Token = token
	return token, nil
}

type ProfileResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Username      string `json:"username"`
		Email         string `json:"email"`
		Role          string `json:"role"`
		GroupExpire   string `json:"group_expire"`
		EmailVerified int    `json:"email_verified"`
		DiskUsage     string `json:"disk_usage"`
		DiskLimit     string `json:"disk_limit"`
		DiskUsageRaw  int    `json:"disk_usage_raw"`
		DiskLimitRaw  int64  `json:"disk_limit_raw"`
	} `json:"data"`
	RequestId string `json:"RequestId"`
}

// GetProfile User - Get User Profile
func (l *SmmSLogic) GetProfile() (*ProfileResponse, error) {
	logrus.Info("获取SMMS Profile")
	client := resty.New()
	client.SetProxy(l.Proxy)
	r, err := client.R().
		SetHeader("Authorization", l.Key).
		Post(l.BaseURI + "profile")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var resp ProfileResponse
	if err := json.Unmarshal(r.Body(), &resp); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &resp, nil
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		FileId    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		Url       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	} `json:"data"`
	RequestId string `json:"RequestId"`
}

// Upload Image - Upload Image
func (l *SmmSLogic) Upload(name string, file io.Reader) (*UploadResponse, error) {
	logrus.Info("上传SMMS图片:", name)
	client := resty.New()
	client.SetProxy(l.Proxy)
	r, err := client.R().
		SetHeader("Authorization", l.Key).
		SetFileReader("smfile", name, file).
		Post(l.BaseURI + "upload")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Response Status Code:", r.StatusCode())
	fmt.Println("Response Status:", r.Status())
	fmt.Println("Response:", r)
	var resp UploadResponse
	if err := json.Unmarshal(r.Body(), &resp); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &resp, nil
}

type HistoryResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		FileId    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		CreatedAt string `json:"created_at"`
		Url       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	} `json:"data"`
	RequestId string `json:"RequestId"`
}

// History Image - Upload History
func (l *SmmSLogic) History(page int) (*HistoryResponse, error) {
	logrus.Info("获取SMMS上传历史")
	var response HistoryResponse
	NewLibLogic().TryCatch(func() {
		client := resty.New()
		client.SetProxy(l.Proxy)
		r, err := client.R().
			SetHeader("Authorization", l.Key).
			SetFormData(map[string]string{
				"format": "json",
				"page":   strconv.Itoa(page),
			}).
			Get(l.BaseURI + "upload_history")
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(r.Body(), &response); err != nil {
			fmt.Println("Error:", err)
		}
	}, func(err interface{}) {
		fmt.Println(err)
	})
	return &response, nil
}

type DeleteResponse struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"RequestId"`
}

// Delete Image - Image Deletion
func (l *SmmSLogic) Delete(hash string) (*DeleteResponse, error) {
	logrus.Info("删除SMMS图片")
	client := resty.New()
	client.SetProxy(l.Proxy)
	r, err := client.R().
		SetHeader("Authorization", l.Key).
		SetFormData(map[string]string{
			"format": "json",
			"hash":   hash,
		}).
		Get(l.BaseURI + "delete/" + hash)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var resp DeleteResponse
	if err := json.Unmarshal(r.Body(), &resp); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &resp, nil
}
