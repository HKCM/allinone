package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	// 构建Alertmanager告警数据,以JSON格式
	// "status":"firing",
	alertData := `[
		{
		  "labels": {
			"alertname": "Cpu TEST golang4",
			"name": "test4",
			"instance": "example4",
			"hostname": "example4",
			"severity": "critical"
		  },
		  "annotations": {
			"node": "testNode4",
			"status": "95",
			"message": "test"
		  },
		  "StartsAt": "2023-10-13T16:30:00Z", 
		  "endsAt": "2023-10-14T17:54:00Z"
		}
	  ]`

	// Alertmanager的API端点URL
	alertmanagerURL := "http://localhost:9093/api/v1/alerts"

	// 发送HTTP POST请求
	resp, err := http.Post(alertmanagerURL, "application/json", bytes.NewBuffer([]byte(alertData)))

	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Alert sent successfully!")
	} else {
		fmt.Println("Failed to send alert. Status code:", resp.Status)
	}
}
