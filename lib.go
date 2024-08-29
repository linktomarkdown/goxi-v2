package goxi_v2

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LibLogic struct {
}
type RedisOptions redis.Options
type CacheClient struct {
	*redis.Client
}
type RegisterOptions struct {
	Name           string
	Type           string
	Path           string
	Command        string
	CommandStop    string
	Status         int
	RemoteHost     string
	RemotePort     int
	RemoteUser     string
	RemotePassword string
	RemoteKey      string
	ApiUrl         string
}

func NewLibLogic() *LibLogic {
	return &LibLogic{}
}

// GenerateOrderNo 生成订单号
func (u *LibLogic) GenerateOrderNo() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// GenerateName 生成名称
func (u *LibLogic) GenerateName(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

// Md5V 密码md5加密
func (u *LibLogic) Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRound 四舍五入
func (u *LibLogic) GetRound(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (u *LibLogic) NewRedisConnect(options *RedisOptions) (*CacheClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &CacheClient{rdb}, nil
}

// InitLogrus 初始化日志
func (u *LibLogic) InitLogrus() {
	// InitLogrus 初始化日志，写入文件
	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	// 设置日志输出,创建日志文件，如果文件存在则追加，不存在则创建
	logFile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("打开日志文件失败: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

// GetEnvInfo 获取环境变量
func (u *LibLogic) GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

// StringToFloat64 字符串转浮点数
func (u *LibLogic) StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// StringToInt 字符串转整数
func (u *LibLogic) StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// InArray 判断元素是否在数组中
func (u *LibLogic) InArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

// TryCatch 捕获异常
func (u *LibLogic) TryCatch(f func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			handler(err)
		}
	}()
	f()
}

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

// GenerateRandomPassword 生成随机密码
func (u *LibLogic) GenerateRandomPassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[r.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[r.Intn(len(numBytes))]
		}
	}
	return string(b)
}

// AddPrefix 统一给地址添加前缀
func (u *LibLogic) AddPrefix(path string, prefix string) string {
	if strings.HasPrefix(path, prefix) {
		return path
	}
	return prefix + path
}

func (u *LibLogic) inArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

// CopyFile 复制单个文件
func (u *LibLogic) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// CopyDir 复制整个文件夹
func (u *LibLogic) CopyDir(src string, dst string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !s.IsDir() {
		return &os.PathError{
			Op:   "read",
			Path: src,
			Err:  os.ErrInvalid,
		}
	}

	os.MkdirAll(dst, s.Mode())

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = u.CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = u.CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
