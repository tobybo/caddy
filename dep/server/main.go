package main

import (
	"path/filepath"
	"io"
	"os"
	"strings"

    "encoding/json"
    "fmt"
    "net/http"
	"server/player_mng"
	"server/mongo"
	"server/define"
)

const (
	E_OK = 0
	E_ERROR = 1
	E_EXISTS = 2
	E_NOT_EXISTS = 3
	E_PWD_WRONG = 4
)

// 响应结构体
type Response struct {
	Code int `json:"code"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

// 处理 GET 请求
// func greetHandler(w http.ResponseWriter, r *http.Request) {
//     response := Response{Message: "Hello, welcome to the Go API!"}
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

// 处理 POST 请求
func postRegister(w http.ResponseWriter, r *http.Request) {
    fmt.Println("postRegister")
    var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data);
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    response := Response{Message: "Data received!"}
	fmt.Println(data)
	if !player_mng.CheckPlayerExists(data["name"].(string)) {
		player_mng.CreatePlayer(data["name"].(string), data["pwd"].(string))
		response.Code = E_OK
	} else {
		fmt.Println("Player already exists")
		response.Message = "Player already exists"
		response.Code = E_EXISTS
	}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postLogin")
    var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data);
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    response := Response{Message: "Data received!"}
	fmt.Println(data)
	ply, exists := player_mng.GetPlayer(data["name"].(string))
	if !exists {
		response.Code = E_NOT_EXISTS
		fmt.Println("Player is not exists")
		response.Message = "Player is not exists"
	} else {
		if ply.Pwd != data["pwd"].(string) {
			response.Code = E_PWD_WRONG
			fmt.Println("Password is not correct")
			response.Message = "Password is not correct"
		} else {
			response.Code = E_OK
            response.Data = map[string]interface{}{"pid": ply.Pid, "name": ply.Name, "permission": ply.Permission}
		}
	}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

type ImageInfo struct {
	ImagePath string `json:"image_path"`
	Description string `json:"description"`
}

func fetchImages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchImages")
	imagesDir := "../myhtml/images"
	var images []ImageInfo

	// 读取目录下的所有文件夹
	entries, err := os.ReadDir(imagesDir)
	if err != nil {
        fmt.Println("fetchImages, failed 1, error:", err)
		http.Error(w, "Failed to read images directory", http.StatusInternalServerError)
		return
	}

	for _, entry := range entries {
        if entry.IsDir() {
            folderPath := filepath.Join(imagesDir, entry.Name())
            var description string

            // 查找描述文件
            files, err := os.ReadDir(folderPath)
            if err != nil {
                continue // 如果读取文件夹失败，跳过这个文件夹
            }

            for _, file := range files {
                if strings.HasSuffix(file.Name(), ".txt") {
                    // 读取描述文件内容
                    descPath := filepath.Join(folderPath, file.Name())
                    descriptionBytes, err := os.ReadFile(descPath)
                    if err == nil {
                        description = string(descriptionBytes)
                    }
                }
            }

            // 查找图片文件
            var imagePath string
            imageFormats := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"} // 支持的图片格式

            for _, file := range files {
                for _, format := range imageFormats {
                    if strings.HasSuffix(file.Name(), format) {
                        imagePath = filepath.Join("./images", entry.Name(), file.Name()) // 返回相对路径
                        break
                    }
                }
                if imagePath != "" {
                    break // 找到第一个有效的图片文件后停止查找
                }
            }

            // 如果找到了图片路径，添加到列表中
            if imagePath != "" {
                images = append(images, ImageInfo{
                    ImagePath:   imagePath,
                    Description: description,
                })
            }
        }
    }

	// for _, entry := range entries {
	//     if entry.IsDir() {
	//         descriptionPath := filepath.Join(imagesDir, entry.Name(), "desc.txt") // 假设描述文件名为 description.txt
    //
	//         // 读取描述文件内容
	//         description, err := os.ReadFile(descriptionPath)
	//         if err != nil {
	//             fmt.Println("fetchImages, failed, error:", err)
	//             continue // 如果描述文件不存在，跳过这个文件夹
	//         }
    //
	//         images = append(images, ImageInfo{
	//             ImagePath: filepath.Join("./images", entry.Name(), "bx比心.png"), // 返回相对路径
	//             Description: string(description),
	//         })
    //         fmt.Println(images)
	//     }
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploadHandler")
    if r.Method != http.MethodPost {
        fmt.Println("uploadHandler, invalid method")
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // 解析表单数据
    err := r.ParseMultipartForm(10 << 20) // 限制最大上传文件为 10MB
    if err != nil {
        fmt.Println("uploadHandler, parse form error:", err)
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

	// 获取文件
    file, fileHeader, err := r.FormFile("image")
    if err != nil {
        fmt.Println("uploadHandler, get file error:", err)
        http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 获取文件扩展名
    ext := filepath.Ext(fileHeader.Filename)

    title := r.FormValue("title") // 获取描述字段

	dir := "../myhtml/images/" + title // 你可以根据需要修改路径
    err = os.MkdirAll(dir, os.ModePerm) // 创建多级目录
    if err != nil {
        fmt.Println("uploadHandler, mkdir error:", err)
        http.Error(w, "Unable to create directories", http.StatusInternalServerError)
        return
    }

    // 创建文件
    filePath := dir + "/" + title + ext // 组合路径
    outFile, err := os.Create(filePath)
    if err != nil {
        fmt.Println("uploadHandler, create file error:", err)
        http.Error(w, "Unable to create file", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()

    // 将上传的文件内容复制到新文件
    _, err = io.Copy(outFile, file)
    if err != nil {
        fmt.Println("uploadHandler, copy file error:", err)
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
    }

    // 获取描述信息
    description := r.FormValue("description") // 获取描述字段

    // 处理描述信息（例如，保存到文件）
    err = os.WriteFile(filepath.Join(dir, "desc.txt"), []byte(description), 0644) // 保存描述信息
    if err != nil {
        fmt.Println("uploadHandler, save description error:", err)
        http.Error(w, "Unable to save description", http.StatusInternalServerError)
        return
    }

    // 返回成功消息
    fmt.Fprintf(w, "File and description uploaded successfully!")
}

func main() {
	mongo.Connect(define.DB_HOST);
	player_mng.LoadPlys();
	player_mng.ViewPlys();
	// http.HandleFunc("/api/greet", greetHandler)
	http.HandleFunc("/api/register", postRegister)
	http.HandleFunc("/api/login", postLogin)
	http.HandleFunc("/images", fetchImages)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
