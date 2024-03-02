package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

// Assuming your testing function is within a *_test.go file

func TestLoginHandler(t *testing.T) {

	connectDB()

	// Set up a mock HTTP server
	router := setupRouter() // Assuming you have a setupRouter function

	// Prepare login form values
	form := url.Values{}
	form.Set("email", "ataytoleuov05@gmail.com")
	form.Set("password", "12345678")

	// Prepare login payload
	loginPayload := strings.NewReader(form.Encode())

	// Create a login request
	loginRequest, err := http.NewRequest("POST", "/login", loginPayload)
	assert.NoError(t, err)

	// Set the Content-Type header to simulate form data
	loginRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a recorder to capture the response
	loginResponse := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(loginResponse, loginRequest)

	// Assert the response status code
	assert.Equal(t, http.StatusSeeOther, loginResponse.Code)

	t.Log(loginResponse.Body.String())
}
func TestRegisterHandler(t *testing.T) {
	// Set up a mock HTTP server
	router := setupRouter() // Assuming you have a setupRouter function

	// Prepare registration form values
	form := url.Values{}
	form.Set("email", "test@example.com")
	form.Set("username", "testuser")
	form.Set("password", "testpassword")

	// Prepare registration payload
	registerPayload := strings.NewReader(form.Encode())

	// Create a registration request
	registerRequest, err := http.NewRequest("POST", "/register", registerPayload)
	assert.NoError(t, err)

	// Set the Content-Type header to simulate form data
	registerRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a recorder to capture the response
	registerResponse := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(registerResponse, registerRequest)

	// Assert the response status code
	assert.Equal(t, http.StatusSeeOther, registerResponse.Code)
	t.Log(registerResponse.Body.String())
}
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Replace these with your actual registration, verification, and login handlers
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/registration/verify", confirmVerificationCodeHandler)
	router.HandleFunc("/login", loginHandler)

	return router
}
func TestGenerateVerificationCode(t *testing.T) {
	// Set a fixed seed for the random number generator to ensure reproducibility
	randSeed := time.Now().UnixNano()
	rand.Seed(randSeed)

	// Call the function to generate a verification code
	verificationCode := generateVerificationCode()

	// Assert that the generated code is a 6-digit string
	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), verificationCode)
}

// Add other test functions for additional scenarios if needed

func startWebDriver(t *testing.T) (*selenium.Service, selenium.WebDriver) {

	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("C:/Users/abyla/OneDrive/Рабочий стол/chromedriver.exe", 4444, opts...)
	if err != nil {
		t.Fatalf("Failed to start ChromeDriver service: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string][]string{
			"args": []string{
				"--headless",
			},
		},
	}

	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 4444))
	if err != nil {
		t.Fatalf("Failed to create new WebDriver session: %v", err)
	}

	return service, webDriver
}

func stopWebDriver(t *testing.T, service *selenium.Service, webDriver selenium.WebDriver) {

	if webDriver != nil {
		if err := webDriver.Quit(); err != nil {
			t.Fatalf("Failed to quit WebDriver: %v", err)
		}
	}

	if service != nil {
		if err := service.Stop(); err != nil {
			t.Fatalf("Failed to stop WebDriver service: %v", err)
		}
	}
}

func TestWebPageInteractions(t *testing.T) {
	service, webDriver := startWebDriver(t)
	defer stopWebDriver(t, service, webDriver)

	if err := webDriver.Get("http://localhost:8080/products"); err != nil {
		t.Fatalf("Failed to navigate to the web page: %v", err)
	}

	// Example interaction: Find an element by CSS selector and check its text content
	element, err := webDriver.FindElement(selenium.ByCSSSelector, ".nav-link")
	if err != nil {
		t.Fatalf("Failed to find the element: %v", err)
	}

	// Get the text content of the element
	text, err := element.Text()
	if err != nil {
		t.Fatalf("Failed to get text content: %v", err)
	}

	// Verify that the text content is as expected
	expectedText := "Your Expected Text"
	if text != expectedText {
		t.Fatalf("Unexpected text content. Got %s, expected %s", text, expectedText)
	}

	// Add more interactions as needed

	// Wait for a few seconds to observe the changes
	time.Sleep(5 * time.Second)
}
