package goxi_v2

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	mu sync.Mutex // 用于保证并发安全
)

// 处理视频工具方法

type VideoLogic struct {
	mediaDir string
}

func NewVideoLogic(mediaDir string) *VideoLogic {
	if mediaDir == "" {
		mediaDir = "media"
	}
	return &VideoLogic{
		mediaDir: mediaDir,
	}
}

// GenerateSegments 生成视频分段
func (v *VideoLogic) GenerateSegments(videoPath string, videoName string) error {
	mu.Lock()
	defer mu.Unlock()

	// 创建目录用于存储生成的M3U8文件和TS分段
	allFilePath := filepath.Join(v.mediaDir, videoName)
	_ = os.MkdirAll(allFilePath, os.ModePerm)
	//1.将原视频整体转码为 ts 格式
	//ffmpeg -y -i Test.mp4  -vcodec copy -acodec copy -vbsf h264_mp4toannexb EncodeTest.ts
	//2. ts 文件切片并生成索引
	//ffmpeg -i EncodeTest.ts -c copy -map 0 -f segment -segment_list Index.m3u8 -segment_time 10 TestSeg_%3d.ts

	// 1. 将原视频整体转码为 ts 格式
	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath,
		"-vcodec", "copy", "-acodec", "copy", "-vbsf", "h264_mp4toannexb",
		filepath.Join(allFilePath, videoName+".ts"), "-ss", "00:00:01", "-vframes", "1", "-f", "image2", "-y",
		filepath.Join(allFilePath, "cover.jpg"))

	err := cmd.Run()
	if err != nil {
		logrus.Error("Failed to generate TS:", err)
	}
	// 2. ts 文件切片并生成索引
	cmd = exec.Command("ffmpeg",
		"-i", filepath.Join(allFilePath, videoName+".ts"),
		"-c", "copy",
		"-map", "0",
		"-f", "segment",
		"-segment_list", filepath.Join(allFilePath, "playlist.m3u8"),
		"-segment_time", "10",
		filepath.Join(allFilePath, "segment_%d.ts"))
	err = cmd.Run()
	if err != nil {
		logrus.Error("Failed to generate segments:", err)
	}

	return nil
}

// CheckSegmentsExist 检查视频分段是否存在
func (v *VideoLogic) CheckSegmentsExist(videoName string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()
	segmentsPath := filepath.Join(v.mediaDir, videoName)
	_, err := os.Stat(segmentsPath)
	if err == nil {
		return true, nil // 分段文件已存在
	}
	if os.IsNotExist(err) {
		return false, nil // 分段文件不存在
	}
	return false, err // 其他错误
}
