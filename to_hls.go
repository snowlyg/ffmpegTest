package go_ffmpeg

/*
#cgo pkg-config:  libavformat  libavutil libavcodec libswscale libswresample libavdevice libavfilter
#include <hls/hls.c>
*/
import "C"
import "os"

type rtspTransport int

const (
	_TCP rtspTransport = iota
	_UDP
)

func (r rtspTransport) String() string {
	switch r {
	case _TCP:
		return "tcp"
	case _UDP:
		return "udp"
	default:
		return "tcp"
	}
}

type Hls struct {
	InFilename    string
	OutFilename   string
	RtspTransport rtspTransport
}

//	养鸡rtsp回放：rtsp://www.mym9.com/101065?from=2019-06-28/01:12:13
//	rtmp://58.200.131.2:1935/livetv/hunantv
//	inFilename := "rtsp://183.59.168.27/PLTV/88888905/224/3221227272/10000100000000060000000001030757_0.smil?icip=88888888"
//	outFilename := "D:/Env/nginx/html/hls/ffmpeg/test.m3u8"
//	rtspTransport := "tcp"

func (h *Hls) ToHls() {

	err := CreateFile(h.OutFilename)
	if err != nil {
		panic(err)
	}

	outFilename := h.OutFilename + "/out.m3u8"

	C.to_hls(C.CString(h.InFilename), C.CString(outFilename), C.CString(h.RtspTransport.String()))
}

// 调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

//  判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//func main() {
//	inFilename := "rtsp://www.mym9.com/101065?from=2019-06-28/01:12:13"
//	outFilename := "./hls_files"
//	rtspTransport := "tcp"
//
//	hls := Hls{
//		InFilename: "rtsp://www.mym9.com/101065?from=2019-06-28/01:12:13",
//		OutFilename: "./hls_files",
//	}
//
//	hls.ToHls()
//}
