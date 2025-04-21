package main

import (
	"os"
	"io"
    "fmt"
    "net/http"
	"image_recognize_server/img_recognize_sdk"
	"strings"
	"image_recognize_server/mongo"
	"image_recognize_server/tool"
	"image_recognize_server/define"
	"encoding/json"
)

var accessToken string;

type ResponseForReqTimes struct {
    UseTimes int32 `json:"UseTimes"`
    LimitTimes int32 `json:"LimitTimes"`
}

func postReqTimes(w http.ResponseWriter, r *http.Request) {
    // 解析表单数据
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // 获取 device_id 字段值
    deviceID := r.FormValue("device_id")
    if deviceID == "" {
        http.Error(w, "device_id not found", http.StatusBadRequest)
        return
    }

    // 打印 device_id
    fmt.Printf("Received device_id: %s\n", deviceID)

	sn := tool.GetSn(deviceID);

    // 创建响应数据
    response := ResponseForReqTimes{
        UseTimes: sn, // 示例值
        LimitTimes: define.LIMIT_TIMES_FOR_RECOGNIZE, // 示例值
    }

    // 将响应数据编码为 JSON 并发送
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// // 处理 POST 请求
// func postRecognize(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("postRecognize")
//     // 将请求传递给 forwardRequest 函数
//     resp, err := img_recognize_sdk.ImageRecognition(accessToken, r)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer resp.Body.Close()
//
//     // 将第三方服务器的响应头复制到客户端响应中
//     for key, values := range resp.Header {
//         for _, value := range values {
//             w.Header().Add(key, value)
//         }
//     }
//
//     // 将第三方服务器的响应状态码复制到客户端响应中
//     w.WriteHeader(resp.StatusCode)
//
//     // 将第三方服务器的响应 body 复制到客户端响应中
//     if _, err := io.Copy(w, resp.Body); err != nil {
//         fmt.Println("Error copying response body:", err)
//     }
// }

// 处理 POST 请求
func postRecognize(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postRecognize")
	// 读取并解析表单数据
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // 获取字段值
    deviceID := r.FormValue("device_id")
	// 删除 device_id 字段
    r.Form.Del("device_id")
	curSn := tool.GetSn(deviceID)
	fmt.Println("postRecognize, deviceID: ", deviceID, " curSn: ", curSn)
	if curSn >= define.LIMIT_TIMES_FOR_RECOGNIZE {
		fmt.Println("postRecognize, overflow times, curSn: ", curSn)
        http.Error(w, "overflow times", http.StatusBadRequest)
		return
	}
	tool.IncSn(deviceID)
    // 将请求传递给 forwardRequest 函数
    img_recognize_sdk.ImageRecognition(accessToken, w, r)
}


func loadToken() {
	// 打开文件
    file, err := os.Open("token.txt")
    if err != nil {
        fmt.Println("无法打开文件: %v", err)
    }
    defer file.Close()

    // 读取文件内容
    content, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("无法读取文件内容: %v", err)
    }

    // 将内容转换为字符串
    accessToken = strings.TrimRight(string(content), "\n")
}

func connectMongo() {
	// 连接 MongoDB
    mongo.Connect(define.DB_HOST);
}

func main() {
	connectMongo()
	loadToken()
    http.HandleFunc("/api/recognize", postRecognize)
    http.HandleFunc("/api/req_times", postReqTimes)
	fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
