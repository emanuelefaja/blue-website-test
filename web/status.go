package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Service represents a monitored service
type Service struct {
	Name           string
	URL            string
	HealthEndpoint string
}

// CheckResult represents a single health check result
type CheckResult struct {
	ServiceName string
	ServiceURL  string
	Status      string // "up" or "down"
	CheckedAt   time.Time
}

// DayData represents daily status data
type DayData struct {
	Date   string  `json:"date"`
	Status string  `json:"status"`
	Uptime float64 `json:"uptime"`
}

// ServiceStatus represents current status for API response
type ServiceStatus struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	LastChecked string `json:"lastChecked"`
}

// ServiceHistory represents historical data for API response
type ServiceHistory struct {
	Name   string    `json:"name"`
	Uptime float64   `json:"uptime"`
	Days   []DayData `json:"days"`
}

// Hardcoded list of services to monitor
var monitoredServices = []Service{
	{Name: "Website", URL: "https://blue.cc", HealthEndpoint: ""},
	{Name: "API", URL: "https://api.blue.cc", HealthEndpoint: "https://api.blue.cc/health"},
	{Name: "Web App", URL: "https://app.blue.cc", HealthEndpoint: ""},
	{Name: "Beta Web App", URL: "https://beta.app.blue.cc", HealthEndpoint: ""},
	{Name: "Beta API", URL: "https://beta.api.blue.cc", HealthEndpoint: "https://beta.api.blue.cc/health"},
	{Name: "Background Services", URL: "https://bg.blue.cc", HealthEndpoint: "https://bg.blue.cc/health"},
	{Name: "Forms", URL: "https://forms.blue.cc", HealthEndpoint: ""}, // No health endpoint
	{Name: "Files", URL: "https://files.blue.cc", HealthEndpoint: ""},
	{Name: "White Label Files", URL: "https://wl-files.onrender.com", HealthEndpoint: ""},
	{Name: "White Label Forms", URL: "https://wl-forms.onrender.com", HealthEndpoint: ""},
	{Name: "Realtime Engine", URL: "https://collab.blue.cc", HealthEndpoint: ""},
	{Name: "Search", URL: "https://search.blue.cc", HealthEndpoint: "https://search.blue.cc/health"},
}

// D1Client handles communication with Cloudflare D1
type D1Client struct {
	AccountID  string
	DatabaseID string
	APIKey     string
	BaseURL    string
}

// D1Request represents a D1 API request
type D1Request struct {
	SQL    string   `json:"sql"`
	Params []string `json:"params,omitempty"`
}

// D1Response represents a D1 API response
type D1Response struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Result []struct {
		Meta struct {
			ChangedDB      bool    `json:"changed_db"`
			Changes        int     `json:"changes"`
			Duration       float64 `json:"duration"`
			LastRowID      int     `json:"last_row_id"`
			RowsRead       int     `json:"rows_read"`
			RowsWritten    int     `json:"rows_written"`
			ServedByRegion string  `json:"served_by_region"`
		} `json:"meta"`
		Results []map[string]interface{} `json:"results"`
		Success bool                     `json:"success"`
	} `json:"result"`
}

// NewD1Client creates a new D1 client from environment variables
func NewD1Client() *D1Client {
	return &D1Client{
		AccountID:  os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		DatabaseID: os.Getenv("CLOUDFLARE_DATABASE_ID"),
		APIKey:     os.Getenv("CLOUDFLARE_API_KEY"),
		BaseURL:    "https://api.cloudflare.com/client/v4",
	}
}

// Query executes a SQL query against D1
func (d *D1Client) Query(sql string, params ...string) (*D1Response, error) {
	url := fmt.Sprintf("%s/accounts/%s/d1/database/%s/query", d.BaseURL, d.AccountID, d.DatabaseID)

	reqBody := D1Request{
		SQL:    sql,
		Params: params,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d1Resp D1Response
	if err := json.NewDecoder(resp.Body).Decode(&d1Resp); err != nil {
		return nil, err
	}

	if !d1Resp.Success {
		if len(d1Resp.Errors) > 0 {
			return nil, fmt.Errorf("D1 error: %s", d1Resp.Errors[0].Message)
		}
		return nil, fmt.Errorf("D1 request failed")
	}

	return &d1Resp, nil
}

// HealthChecker manages health checking and caching
type HealthChecker struct {
	d1Client      *D1Client
	cache         *sync.Map // Thread-safe map for caching
	lastCheckTime *sync.Map // Last check time for each service
	mu            sync.RWMutex
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(d1Client *D1Client) *HealthChecker {
	return &HealthChecker{
		d1Client:      d1Client,
		cache:         &sync.Map{},
		lastCheckTime: &sync.Map{},
	}
}

// Initialize sets up the database and loads historical data
func (h *HealthChecker) Initialize() error {
	// Create tables
	if err := h.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Clean up old data
	if err := h.cleanupOldData(); err != nil {
		log.Printf("Failed to cleanup old data: %v", err)
		// Continue - not critical
	}

	// Load historical data into cache
	if err := h.loadHistoricalData(); err != nil {
		log.Printf("Failed to load historical data: %v", err)
		// Continue with empty cache
	}

	return nil
}

// createTables creates the necessary database tables
func (h *HealthChecker) createTables() error {
	schema := `
		CREATE TABLE IF NOT EXISTS service_checks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service_name TEXT NOT NULL,
			service_url TEXT NOT NULL,
			status TEXT NOT NULL,
			checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_service_checks_name_date 
		ON service_checks(service_name, checked_at);
		
		CREATE INDEX IF NOT EXISTS idx_service_checks_date 
		ON service_checks(checked_at);
	`

	_, err := h.d1Client.Query(schema)
	return err
}

// cleanupOldData removes data older than 90 days
func (h *HealthChecker) cleanupOldData() error {
	sql := "DELETE FROM service_checks WHERE checked_at < datetime('now', '-90 days')"
	_, err := h.d1Client.Query(sql)
	return err
}

// loadHistoricalData loads 90 days of data into cache
func (h *HealthChecker) loadHistoricalData() error {
	sql := `
		SELECT service_name, status, checked_at 
		FROM service_checks 
		WHERE checked_at >= datetime('now', '-90 days')
		ORDER BY checked_at DESC
	`

	resp, err := h.d1Client.Query(sql)
	if err != nil {
		return err
	}

	// Process results and populate cache
	for _, result := range resp.Result {
		for _, row := range result.Results {
			serviceName, _ := row["service_name"].(string)
			status, _ := row["status"].(string)
			checkedAtStr, _ := row["checked_at"].(string)
			
			checkedAt, _ := time.Parse(time.RFC3339, checkedAtStr)
			
			check := CheckResult{
				ServiceName: serviceName,
				Status:      status,
				CheckedAt:   checkedAt,
			}
			
			h.addToCache(check)
		}
	}

	return nil
}

// addToCache adds a check result to the cache
func (h *HealthChecker) addToCache(result CheckResult) {
	key := fmt.Sprintf("%s:%s", result.ServiceName, result.CheckedAt.Format("2006-01-02"))
	
	// Get existing checks for this service/day
	var checks []CheckResult
	if existing, ok := h.cache.Load(key); ok {
		checks = existing.([]CheckResult)
	}
	
	checks = append(checks, result)
	h.cache.Store(key, checks)
	
	// Update last check time
	h.lastCheckTime.Store(result.ServiceName, result.CheckedAt)
}

// CheckService performs a health check on a single service
func (h *HealthChecker) CheckService(service Service) CheckResult {
	// Try health endpoint first, fallback to simple ping
	endpoint := service.HealthEndpoint
	if endpoint == "" {
		endpoint = service.URL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(endpoint)

	status := "up"
	if err != nil {
		status = "down"
		log.Printf("Health check failed for %s: %v", service.Name, err)
	} else if resp != nil {
		// Special case for White Label Files - 404 means server is alive
		if service.Name == "White Label Files" && resp.StatusCode == 404 {
			status = "up"
		} else if resp.StatusCode >= 400 {
			status = "down"
		}
	}
	
	if resp != nil {
		resp.Body.Close()
	}

	result := CheckResult{
		ServiceName: service.Name,
		ServiceURL:  service.URL,
		Status:      status,
		CheckedAt:   time.Now().UTC(),
	}

	// Store result in D1
	h.storeResult(result)
	
	// Add to cache
	h.addToCache(result)

	return result
}

// storeResult saves a check result to D1
func (h *HealthChecker) storeResult(result CheckResult) {
	sql := `INSERT INTO service_checks (service_name, service_url, status, checked_at) 
			VALUES (?, ?, ?, ?)`

	params := []string{
		result.ServiceName,
		result.ServiceURL,
		result.Status,
		result.CheckedAt.Format(time.RFC3339),
	}

	if _, err := h.d1Client.Query(sql, params...); err != nil {
		log.Printf("Failed to store check result for %s: %v", result.ServiceName, err)
		// Continue operation - don't crash the checker
	}
}

// hasRecentCheck checks if any service has been checked within the last 5 minutes
func (h *HealthChecker) hasRecentCheck() bool {
	sql := `SELECT COUNT(*) as count FROM service_checks 
			WHERE checked_at >= datetime('now', '-5 minutes')`
	
	resp, err := h.d1Client.Query(sql)
	if err != nil {
		log.Printf("Failed to check for recent checks: %v", err)
		return false // If we can't check, assume no recent checks and proceed
	}
	
	for _, result := range resp.Result {
		for _, row := range result.Results {
			if count, ok := row["count"].(float64); ok {
				return count > 0
			}
		}
	}
	
	return false
}

// CheckAllServicesIfNeeded performs health checks only if none have been done in the last 5 minutes
func (h *HealthChecker) CheckAllServicesIfNeeded() {
	if h.hasRecentCheck() {
		log.Println("⏭️  Recent health checks found (within last 5 minutes), skipping new checks")
		return
	}
	
	log.Println("🏥 No recent health checks found, performing health checks on all services")
	h.CheckAllServices()
}

// CheckAllServices performs health checks on all services
func (h *HealthChecker) CheckAllServices() {
	var wg sync.WaitGroup
	for _, service := range monitoredServices {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			h.CheckService(svc)
		}(service)
	}
	wg.Wait()
}

// GetCurrentStatus returns the current status of all services
func (h *HealthChecker) GetCurrentStatus() []ServiceStatus {
	var statuses []ServiceStatus
	
	for _, service := range monitoredServices {
		// Get last check time
		var lastChecked string
		if checkTime, ok := h.lastCheckTime.Load(service.Name); ok {
			duration := time.Since(checkTime.(time.Time))
			lastChecked = formatDuration(duration)
		} else {
			lastChecked = "Never"
		}
		
		// Get current status from latest check
		status := h.getLatestStatus(service.Name)
		
		statuses = append(statuses, ServiceStatus{
			Name:        service.Name,
			Status:      status,
			LastChecked: lastChecked,
		})
	}
	
	return statuses
}

// getLatestStatus returns the most recent status for a service
func (h *HealthChecker) getLatestStatus(serviceName string) string {
	today := time.Now().UTC().Format("2006-01-02")
	key := fmt.Sprintf("%s:%s", serviceName, today)
	
	if checks, ok := h.cache.Load(key); ok {
		checkList := checks.([]CheckResult)
		if len(checkList) > 0 {
			// Return the most recent check
			return checkList[len(checkList)-1].Status
		}
	}
	
	return "unknown"
}

// GetHistoricalData returns 90-day history for all services
func (h *HealthChecker) GetHistoricalData() []ServiceHistory {
	var histories []ServiceHistory
	
	for _, service := range monitoredServices {
		days := make([]DayData, 0, 90)
		totalUptime := 0.0
		
		// Generate data for last 90 days
		for i := 89; i >= 0; i-- {
			date := time.Now().UTC().AddDate(0, 0, -i)
			dayData := h.getDayDataForService(service.Name, date)
			days = append(days, dayData)
			totalUptime += dayData.Uptime
		}
		
		// Calculate overall uptime
		overallUptime := totalUptime / 90.0
		
		histories = append(histories, ServiceHistory{
			Name:   service.Name,
			Uptime: overallUptime,
			Days:   days,
		})
	}
	
	return histories
}

// getDayDataForService returns status data for a specific service and day
func (h *HealthChecker) getDayDataForService(serviceName string, date time.Time) DayData {
	key := fmt.Sprintf("%s:%s", serviceName, date.Format("2006-01-02"))
	
	// Check cache for data
	if checks, ok := h.cache.Load(key); ok {
		checkList := checks.([]CheckResult)
		if len(checkList) > 0 {
			// Calculate uptime percentage
			upCount := 0
			for _, check := range checkList {
				if check.Status == "up" {
					upCount++
				}
			}
			
			uptime := (float64(upCount) / float64(len(checkList))) * 100
			status := "operational"
			if uptime < 99 {
				status = "degraded"
			}
			if uptime < 95 {
				status = "outage"
			}
			
			return DayData{
				Date:   date.Format("2006-01-02"),
				Status: status,
				Uptime: uptime,
			}
		}
	}
	
	// No data = assume operational
	return DayData{
		Date:   date.Format("2006-01-02"),
		Status: "operational",
		Uptime: 100.0,
	}
}

// formatDuration formats a duration as "X minutes ago"
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "Less than a minute ago"
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	return "More than a day ago"
}

// HTTP Handlers

// StatusPageHandler handles the main status page
func StatusPageHandler(checker *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The page will be rendered by the existing template system
		// This handler is just a placeholder for any server-side logic
		http.ServeFile(w, r, "pages/platform/status.html")
	}
}

// CurrentStatusAPIHandler returns current status as JSON
func CurrentStatusAPIHandler(checker *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		statuses := checker.GetCurrentStatus()
		
		response := map[string]interface{}{
			"services": statuses,
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode current status response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// HistoricalDataAPIHandler returns historical data as JSON
func HistoricalDataAPIHandler(checker *HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		histories := checker.GetHistoricalData()
		
		response := map[string]interface{}{
			"services": histories,
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode historical data response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}