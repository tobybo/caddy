package main

import (
	"os"
	"io"
    "fmt"
    "net/http"
	"image_recognize_server/img_recognize_sdk"
	"strings"
)

var accessToken string;

// 响应结构体
type Response struct {
	Code int `json:"code"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

// 处理 POST 请求
func postRecognize(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postRecognize")
    // 将请求传递给 forwardRequest 函数
    resp, err := img_recognize_sdk.ImageRecognition(accessToken, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // 将第三方服务器的响应头复制到客户端响应中
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // 将第三方服务器的响应状态码复制到客户端响应中
    w.WriteHeader(resp.StatusCode)

    // 将第三方服务器的响应 body 复制到客户端响应中
    if _, err := io.Copy(w, resp.Body); err != nil {
        fmt.Println("Error copying response body:", err)
    }
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

// // 处理 POST 请求
// func postRecognize(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("postRegister")
//     var data map[string]interface{}
//     err := json.NewDecoder(r.Body).Decode(&data);
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
//     response := Response{Message: "Data received!"}
//     fmt.Println(data)
//
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

func main() {
	loadToken()
    http.HandleFunc("/api/recognize", postRecognize)
	fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
