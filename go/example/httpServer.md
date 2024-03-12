# 服务端代码
```go
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.FormValue("text"))
		if r.Method == "GET" && r.ParseForm() == nil {
			fmt.Fprintln(w, r.FormValue("text"))
			fmt.Fprintf(w, "\nGet 请求")
		}
		if r.Method == "POST" && r.ParseForm() == nil {
			fmt.Fprintln(w, r.FormValue("text"))
			fmt.Fprintf(w, "\nPost 请求")
		}

	})
    // 服务器配置
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
	fmt.Println("服务器启动完成:")
}

```

# 客户端代码
```go
func main() {
	// 普通请求
	// urlRoot := "http://localhost:8080"
	// b, err := Get(urlRoot)
	// if err != nil {
	// 	fmt.Println("出错了")
	// 	return
	// }
	// fmt.Println(string(b))

	// 带参数请求
	// urlApi := "http://localhost:8080/echo"
	// params := url.Values{}
	// params.Set("text", "food") // 设置请求参数
	// apiUrl, err := url.Parse(urlApi)
	// if err != nil {
	// 	fmt.Println("HTTP地址解析错误:", err)
	// }
	// apiUrl.RawQuery = params.Encode()
	// b1, err1 := Get(apiUrl.String())
	// if err1 != nil {
	// 	fmt.Println("出错了")
	// 	return
	// }
	// fmt.Println(string(b1))

	// POST

	apiUrl1 := "http://localhost:8080/echo?text=123"
	b2, err2 := Post(apiUrl1, "", "")
	if err2 != nil {
		fmt.Println("出错了")
		return
	}
	fmt.Println(string(b2))

}

func Post(url, contentType, data string) ([]byte, error) {
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
		return nil, err
	}
	return b, nil
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("读取网络响应失败:", err)
		return nil, err
	}
	return body, nil
}
```

