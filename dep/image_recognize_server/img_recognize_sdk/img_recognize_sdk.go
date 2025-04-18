package img_recognize_sdk

import (
    "fmt"
    "net/http"
)

// func ImageRecognition(accessToken string, r *http.Request) (*http.Response, error) {
//     requestURL := "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general"
//
//     // 设置请求参数
//     requestURL = requestURL + "?access_token=" + accessToken
//     headers := map[string]string{
//         "Content-Type": "application/x-www-form-urlencoded",
//     }
//
//     // 发送 POST 请求
//     client := &http.Client{}
//     req, err := http.NewRequest("POST", requestURL, r.Body)
//     if err != nil {
//         fmt.Println("Error creating request:", err)
//         return nil, err
//     }
//
//     for key, value := range headers {
//         req.Header.Set(key, value)
//     }
//
//     return client.Do(req)
// }

func ImageRecognition(accessToken string, w http.ResponseWriter, r *http.Request) {
    requestURL := "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general"

    // 设置请求参数
    requestURL = requestURL + "?access_token=" + accessToken
    headers := map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    }

    resp, err := http.PostForm(requestURL, r.Form)
    if err != nil {
        http.Error(w, "Failed to send request to third party", http.StatusInternalServerError)
        return
    }

    // 将第三方服务器的响应状态码复制到客户端响应中
    w.WriteHeader(resp.StatusCode)

    // 将第三方服务器的响应 body 复制到客户端响应中
    if _, err := io.Copy(w, resp.Body); err != nil {
        fmt.Println("Error copying response body:", err)
    }
}

// import (
//     "bytes"
//     "fmt"
//     "io"
//     "net/http"
// )

// func imageRecognition(filePath string, accessToken string) {
//     requestURL := "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general"
//
//     // 读取本地文件
//     fileData, err := os.ReadFile(filePath)
//     if err != nil {
//         fmt.Println("Error reading file:", err)
//         return
//     }
//
//     // Base64 编码
//     imgBase64 := base64.StdEncoding.EncodeToString(fileData)
//
//     // 使用 url.QueryUnescape 进行解码
//     urlString := url.QueryEscape(imgBase64)
//
//     // 设置请求参数
//     params := "image=" + urlString
//     requestURL = requestURL + "?access_token=" + accessToken
//     headers := map[string]string{
//         "Content-Type": "application/x-www-form-urlencoded",
//     }
//
//     // 发送 POST 请求
//     client := &http.Client{}
//     req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(params)))
//     if err != nil {
//         fmt.Println("Error creating request:", err)
//         return
//     }
//
//     for key, value := range headers {
//         req.Header.Set(key, value)
//     }
//
//     resp, err := client.Do(req)
//     if err != nil {
//         fmt.Println("Error making request:", err)
//         return
//     }
//     defer resp.Body.Close()
//
//     // 读取响应
//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         fmt.Println("Error reading response:", err)
//         return
//     }
//
//     fmt.Println("Response:", string(body))
// }
//
// func main() {
//     base64str := "./Snipaste_2025-04-09_19-56-43.png"
//     accessToken := "24.c6a8580200ca25b7fb39d99dc226afb9.2592000.1746696929.282335-42280315"
//     imageRecognition(base64str, accessToken)
// }
//
