// package main

// import (
// 	"math/rand"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"os"
// 	"regexp"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/tebeka/selenium"
// 	"github.com/tebeka/selenium/chrome"
// )

// // Mock data and functions for testing

// func TestLoginHandler(t *testing.T) {
// 	// Assuming your connectDB() function initializes the database connection
// 	connectDB()

// 	// Set up a mock HTTP server
// 	router := setupRouter()

// 	// Prepare login form values
// 	form := url.Values{}
// 	form.Set("email", "ataytoleuov015@gmail.com")
// 	form.Set("password", "12345678")

// 	// Prepare login payload
// 	loginPayload := strings.NewReader(form.Encode())

// 	// Create a login request
// 	loginRequest, err := http.NewRequest("POST", "/login", loginPayload)
// 	assert.NoError(t, err)

// 	// Set the Content-Type header to simulate form data
// 	loginRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	// Create a recorder to capture the response
// 	loginResponse := httptest.NewRecorder()

// 	// Serve the request using the router
// 	router.ServeHTTP(loginResponse, loginRequest)

// 	// Assert the response status code
// 	assert.Equal(t, http.StatusSeeOther, loginResponse.Code)

// 	// Log the response body for debugging (optional)
// 	t.Log(loginResponse.Body.String())
// }

// func TestRegisterHandler(t *testing.T) {
// 	// Set up a mock HTTP server
// 	router := setupRouter()

// 	// Prepare registration form values
// 	form := url.Values{}
// 	form.Set("email", "test@example.com")
// 	form.Set("username", "testuser")
// 	form.Set("password", "testpassword")

// 	// Prepare registration payload
// 	registerPayload := strings.NewReader(form.Encode())

// 	// Create a registration request
// 	registerRequest, err := http.NewRequest("POST", "/register", registerPayload)
// 	assert.NoError(t, err)

// 	// Set the Content-Type header to simulate form data
// 	registerRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	// Create a recorder to capture the response
// 	registerResponse := httptest.NewRecorder()

// 	// Serve the request using the router
// 	router.ServeHTTP(registerResponse, registerRequest)

// 	// Assert the response status code
// 	assert.Equal(t, http.StatusSeeOther, registerResponse.Code)
// }

// func TestGenerateVerificationCode(t *testing.T) {
// 	// Set a fixed seed for the random number generator to ensure reproducibility
// 	randSeed := time.Now().UnixNano()
// 	rand.Seed(randSeed)

// 	// Call the function to generate a verification code
// 	verificationCode := generateVerificationCode()

// 	// Assert that the generated code is a 6-digit string
// 	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), verificationCode)
// }

// func setupRouter() *http.ServeMux {
// 	router := http.NewServeMux()

// 	// Replace these with your actual registration, verification, and login handlers
// 	router.HandleFunc("/register", registerHandler)
// 	router.HandleFunc("/registration/verify", confirmVerificationCodeHandler)
// 	router.HandleFunc("/login", loginHandler)

// 	return router
// }

// func TestFilterElements(t *testing.T) {

// 	// Start ChromeDriver service
// 	// Start ChromeDriver service with verbose logging
// 	service, err := selenium.NewChromeDriverService("C:/Users/abyla/OneDrive/Рабочий стол/chromedriver.exe", 4444, selenium.Output(os.Stderr))

// 	if err != nil {
// 		log.Fatal("Error:", err)
// 	}
// 	defer service.Stop()

// 	caps := selenium.Capabilities{}
// 	caps.AddChrome(chrome.Capabilities{Args: []string{
// 		"--headless-new", // comment out this line for testing
// 	}})

// 	// create a new remote client with the specified options
// 	driver, err := selenium.NewRemote(caps, "")

// 	if err != nil {
// 		log.Fatal("Error:", err)
// 	}
// 	log.Println("WebDriver started successfully")
// 	// maximize the current window to avoid responsive rendering
// 	err = driver.MaximizeWindow("")
// 	if err != nil {
// 		log.Fatal("Error:", err)
// 	}

// 	time.Sleep(15 * time.Second) // Adjust the sleep duration as needed
// 	err = driver.Get("https://minima-laptop-store-2bs3.onrender.com/productsTest")
// 	if err != nil {
// 		log.Fatal("Error navigating to the webpage:", err)
// 	} else {
// 		log.Println("Successfully navigated to the webpage")
// 	}

// 	// Locate the price range and sort filter elements
// 	minPriceRange := findElement(t, driver, selenium.ByID, "minPriceRange")
// 	maxPriceRange := findElement(t, driver, selenium.ByID, "maxPriceRange")
// 	sortFilter := findElement(t, driver, selenium.ByID, "sortFilter")
// 	applyFiltersButton := findElement(t, driver, selenium.ByID, "applyFilters")

// 	// Perform interactions
// 	err = minPriceRange.SendKeys("1000")
// 	if err != nil {
// 		log.Fatal("Error setting minPriceRange:", err)
// 	}

// 	err = maxPriceRange.SendKeys("5000")
// 	if err != nil {
// 		log.Fatal("Error setting maxPriceRange:", err)
// 	}

// 	err = sortFilter.SendKeys("desc") // Select descending order
// 	if err != nil {
// 		log.Fatal("Error setting sortFilter:", err)
// 	}

// 	// Click the Apply Filters button
// 	err = applyFiltersButton.Click()
// 	if err != nil {
// 		log.Fatal("Error clicking applyFiltersButton:", err)
// 	}

// 	// Wait for results to load (add explicit wait if needed)

// 	// Assertions to verify the results
// 	// For example, check if the displayed products match the applied filters
// 	// ...

// 	// Cleanup if needed
// }

// func findElement(t *testing.T, wd selenium.WebDriver, by, value string) selenium.WebElement {
// 	// Find the element by specified criteria
// 	element, err := wd.FindElement(by, value)
// 	if err != nil {
// 		t.Fatalf("Failed to find element by %s: %v", by, err)
// 	}
// 	return element
// }
