# Get

```go
func main() {
	urlString := "https://www.google.com/search"
	params := url.Values{}
	apiUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("HTTP地址解析错误:", err)
	}
	params.Set("q", "food") // 设置请求参数

	apiUrl.RawQuery = params.Encode()
	fmt.Println("请求地址为:", apiUrl.String())

	body, err := Get(apiUrl.String())
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
	}
	fmt.Println(string(body))
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

# POST

```go
func main() {
	urlString := "https://www.google.com/search"
	contentType := "application/x-www-from-urlencoded"
	data := "wd=golang"
	
	body, err := Post(urlString, contentType, data)
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
	}

	fmt.Println(string(body))
}

func Post(url,contentType,data string) ([]byte ,error){
	resp,err := http.Post(url,contentType,strings.NewReader(data))
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
	return b,nil
}
```

# header

```go
func clientReq() {
	urlString := "https://www.google.com/"

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println("New Request失败:", err)
	}
	req.Header.Add("headerKey", "headerValue")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP网络请求失败:", err)
		// return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("读取网络响应失败:", err)
	}
	fmt.Println(string(body))
}
```