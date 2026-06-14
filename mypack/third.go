package mypack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  .env dosyası bulunamadı veya yüklenemedi:", err)
	} else {
		fmt.Println("✅ .env dosyası yüklendi.")
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GETpearOperation() map[string]string {

	URL := getEnv("API_URL", "")
	if URL == "" {
		fmt.Println("URL does not identify in .ENV file")
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	token := getEnv("AUTH_TOKEN", "")
	origin := getEnv("ORIGIN", "")
	referer := getEnv("REFERER", "")

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", token)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", origin)
	req.Header.Set("referer", referer)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/141.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("⚠️ JSON parse hatası:", err)
		fmt.Println("Raw response body:", string(body))
		return nil
	}

	serviceMap := make(map[string]string)
	for _, item := range data {
		id := ""
		serviceID := ""

		
		if v, ok := item["id"].(float64); ok {
			id = fmt.Sprintf("%.0f", v)
		}

		
		if service, ok := item["service"].(map[string]interface{}); ok {
			if v, ok := service["id"].(float64); ok {
				serviceID = fmt.Sprintf("%.0f", v)
			}
		}

		if id != "" && serviceID != "" {
			serviceMap[id] = serviceID
			fmt.Println("Record ID:", id, "Service ID:", serviceID)
		}
	}

	return serviceMap
}

func GETJobQuotes(serviceMap map[string]string, log func(msg string)) {
	client := &http.Client{}
	token := getEnv("AUTH_TOKEN", "")
	serviceURL := getEnv("SERVICE_URL", "")
	origin := getEnv("ORIGIN", "")
	referer := getEnv("REFERER", "")

	for recordID, serviceID := range serviceMap {
		URL := fmt.Sprintf("%s%s?service_id=%s", serviceURL, recordID, serviceID)

		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			log("❌ Request hatası: " + err.Error())
			continue
		}

		req.Header.Set("accept", "application/json")
		req.Header.Set("authorization", "Bearer "+token)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/141.0.0.0 Safari/537.36")
		req.Header.Set("origin", origin)
		req.Header.Set("referer", referer)

		resp, err := client.Do(req)
		if err != nil {
			log("❌ HTTP hatası: " + err.Error())
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		log("🌐 Status: " + resp.Status + " | Record: " + recordID)

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log("⚠️ JSON parse hatası: " + err.Error())
			continue
		}

		if isLead, ok := data["is_lead_required"].(bool); ok {
			if leadPriceRaw, ok := data["lead_price"].(float64); ok {
				msg := fmt.Sprintf("📌 Record: %s | Lead Required: %v | Lead Price: %.2f", recordID, isLead, leadPriceRaw)
				log(msg)
				if leadPriceRaw > 120.00 {
					log("🚫 Fiyat yüksek, reddediliyor: " + recordID)
					DeclineJob(recordID, log)
				}
			} else {
				log("⚠️ Lead price alınamadı: " + recordID)
			}
		} else {
			log("⚠️ Lead bilgisi alınamadı: " + recordID)
		}
	}
}

func DeclineJob(ClientID string, log func(msg string)) {
	client := &http.Client{}
	token := getEnv("AUTH_TOKEN", "")
	decline := getEnv("DECLINE_URL", "")
	url := fmt.Sprintf("%s%s/nothanks", decline, ClientID)

	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	JobReasonID := r.Int31n(5) + 13
	payload := map[string]interface{}{
		"answer":  "",
		"id":      JobReasonID,
		"job_id":  ClientID,
		"user_id": getEnv("USER_ID", ""),
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log("❌ Decline request hatası: " + err.Error())
		return
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", token)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/141.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log("❌ Decline HTTP hatası: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log("✅ Decline Status: " + resp.Status + " | " + string(body))
}


func SayMe() {
	fmt.Println("Hello, World!")
	fmt.Println("mypack paketinden sayMe fonksiyonu çağrıldı.")
}
