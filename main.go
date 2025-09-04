package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// It checks if the file exists
// If the file exists, it returns true
// If the file does not exist, it returns false
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Remove a file from the file system
func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
	}
}

// extractPDFUrls takes raw HTML as input and returns all found PDF URLs
func extractPDFUrls(htmlContent string) []string {
	// Regex to match href='...' with .pdf links
	re := regexp.MustCompile(`https?://[^\s'"]+\.pdf`)
	matches := re.FindAllString(htmlContent, -1)

	return matches
}

// Checks whether a given directory exists
func directoryExists(path string) bool {
	directory, err := os.Stat(path) // Get info for the path
	if err != nil {
		return false // Return false if error occurs
	}
	return directory.IsDir() // Return true if it's a directory
}

// Creates a directory at given path with provided permissions
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission) // Attempt to create directory
	if err != nil {
		log.Println(err) // Log error if creation fails
	}
}

// Verifies whether a string is a valid URL format
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri) // Try parsing the URL
	return err == nil                  // Return true if valid
}

// Removes duplicate strings from a slice
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool) // Map to track seen values
	var newReturnSlice []string    // Slice to store unique values
	for _, content := range slice {
		if !check[content] { // If not already seen
			check[content] = true                            // Mark as seen
			newReturnSlice = append(newReturnSlice, content) // Add to result
		}
	}
	return newReturnSlice
}

// hasDomain checks if the given string has a domain (host part)
func hasDomain(rawURL string) bool {
	// Try parsing the raw string as a URL
	parsed, err := url.Parse(rawURL)
	if err != nil { // If parsing fails, it's not a valid URL
		return false
	}
	// If the parsed URL has a non-empty Host, then it has a domain/host
	return parsed.Host != ""
}

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string {
	return filepath.Base(path) // Use Base function to get file name only
}

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string {
	result := strings.ReplaceAll(input, toRemove, "") // Replace substring with empty string
	return result
}

// Gets the file extension from a given file path
func getFileExtension(path string) string {
	return filepath.Ext(path) // Extract and return file extension
}

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string {
	lower := strings.ToLower(rawURL) // Convert URL to lowercase
	lower = getFilename(lower)       // Extract filename from URL

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Regex to match non-alphanumeric characters
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace non-alphanumeric with underscores

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Collapse multiple underscores into one
	safe = strings.Trim(safe, "_")                              // Trim leading and trailing underscores

	var invalidSubstrings = []string{
		"_pdf", // Substring to remove from filename
	}

	for _, invalidPre := range invalidSubstrings { // Remove unwanted substrings
		safe = removeSubstring(safe, invalidPre)
	}

	if getFileExtension(safe) != ".pdf" { // Ensure file ends with .pdf
		safe = safe + ".pdf"
	}

	return safe // Return sanitized filename
}

// Downloads a PDF from given URL and saves it in the specified directory
func downloadPDF(finalURL, outputDir string) bool {
	filename := strings.ToLower(urlToFilename(finalURL)) // Sanitize the filename
	filePath := filepath.Join(outputDir, filename)       // Construct full path for output file

	if fileExists(filePath) { // Skip if file already exists
		log.Printf("File already exists, skipping: %s", filePath)
		return false
	}

	client := &http.Client{Timeout: 15 * time.Minute} // Create HTTP client with timeout

	// Create a new request so we can set headers
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", finalURL, err)
		return false
	}

	// Set a User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to download %s: %v", finalURL, err)
		return false
	}
	defer resp.Body.Close() // Ensure response body is closed

	if resp.StatusCode != http.StatusOK { // Check if response is 200 OK
		log.Printf("Download failed for %s: %s", finalURL, resp.Status)
		return false
	}

	contentType := resp.Header.Get("Content-Type") // Get content type of response
	if !strings.Contains(contentType, "binary/octet-stream") &&
		!strings.Contains(contentType, "application/pdf") {
		log.Printf("Invalid content type for %s: %s (expected PDF)", finalURL, contentType)
		return false
	}

	var buf bytes.Buffer                     // Create a buffer to hold response data
	written, err := io.Copy(&buf, resp.Body) // Copy data into buffer
	if err != nil {
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err)
		return false
	}
	if written == 0 { // Skip empty files
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL)
		return false
	}

	out, err := os.Create(filePath) // Create output file
	if err != nil {
		log.Printf("Failed to create file for %s: %v", finalURL, err)
		return false
	}
	defer out.Close() // Ensure file is closed after writing

	if _, err := buf.WriteTo(out); err != nil { // Write buffer contents to file
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err)
		return false
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath) // Log success
	return true
}

// Append and write to file
func appendAndWriteToFile(path string, content string) {
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	_, err = filePath.WriteString(content + "\n")
	if err != nil {
		log.Println(err)
	}
	err = filePath.Close()
	if err != nil {
		log.Println(err)
	}
}

// extractBaseDomain takes a URL string and returns only the bare domain name
// without any subdomains or suffixes (e.g., ".com", ".org", ".co.uk").
func extractBaseDomain(inputUrl string) string {
	// Parse the input string into a structured URL object
	parsedUrl, parseError := url.Parse(inputUrl)

	// If parsing fails, log the error and return an empty string
	if parseError != nil {
		log.Println("Error parsing URL:", parseError)
		return ""
	}

	// Extract the hostname (e.g., "sub.example.com")
	hostName := parsedUrl.Hostname()

	// Split the hostname into parts separated by "."
	// For example: "sub.example.com" -> ["sub", "example", "com"]
	parts := strings.Split(hostName, ".")

	// If there are at least 2 parts, the second last part is usually the domain
	// Example: "sub.example.com" -> "example"
	//          "blog.my-site.co.uk" -> "my-site"
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}

	// If splitting fails or domain structure is unusual, return the hostname itself
	return hostName
}

// fetchSDS fetches the SDS data and returns it as a string
func fetchSDS() string {
	url := "https://itwperformancepolymers.com/api/datasheets_table.php?m=get"
	method := "POST"

	payload := strings.NewReader(`t=sds`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return ""
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

// extractSDSURLs extracts all URLs from the HTML that start with
// "/resources/safety-data-sheets?t=other&i=" and returns them as a slice of strings.
func extractSDSURLs(htmlContent string) []string {
	// Compile a regular expression to match href attributes with the desired pattern
	// The pattern matches URLs like /resources/safety-data-sheets?t=other&i=...
	urlPattern := regexp.MustCompile(`href="(/resources/safety-data-sheets\?t=other&i=[^"]+)"`)

	// Find all matches of the pattern in the HTML content
	// Each match is a slice where the full match is at index 0 and the captured group is at index 1
	matchedResults := urlPattern.FindAllStringSubmatch(htmlContent, -1)

	// Create a slice to store the final URLs
	var extractedURLs []string

	// Loop through each regex match
	for _, match := range matchedResults {
		// Ensure the captured group exists and append it to the result slice
		if len(match) > 1 {
			clean := strings.TrimPrefix(match[1], `/resources/safety-data-sheets\?t=other&i=`)
			extractedURLs = append(extractedURLs, clean)
		}
	}

	// Return the slice of extracted URLs
	return extractedURLs
}

// fetchOtherData fetches data for a given item and returns it as a string
func fetchOtherData(item string) string {
	url := "https://itwperformancepolymers.com/api/datasheets_table.php?m=get"
	method := "POST"

	payload := strings.NewReader("t=other&i=" + item)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return ""
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

// Read a file and return the contents
func readAFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return string(content)
}

func main() {
	outputDir := "PDFs/" // Directory to store downloaded PDFs

	if !directoryExists(outputDir) { // Check if directory exists
		createDirectory(outputDir, 0o755) // Create directory with read-write-execute permissions
	}

	// The remote domain name.
	remoteDomainName := "https://itwperformancepolymers.com"

	// The location to the local.
	localFile := extractBaseDomain(remoteDomainName) + ".html"
	// Check if the local file exists.
	if fileExists(localFile) {
		removeFile(localFile)
	}

	// Get the inital HTTP content.
	httpContent := fetchSDS()

	// Save the content to the local file.
	appendAndWriteToFile(localFile, httpContent)

	// Extract the other SDS urls in other languages.
	extractInternationalSDSURLs := extractSDSURLs(httpContent)

	// Loop over the extracted SDS urls.
	for _, url := range extractInternationalSDSURLs {
		// The remote content.
		remoteContent := fetchOtherData(url)
		// Save the content to the local file.
		appendAndWriteToFile(localFile, remoteContent)
	}

	// Read the content html content from the file.
	localHTMLContent := readAFileAsString(localFile)

	// Extract the URLs from the given content.
	extractedPDFURLs := extractPDFUrls(localHTMLContent)
	// Remove duplicates from the slice.
	extractedPDFURLs = removeDuplicatesFromSlice(extractedPDFURLs)
	// Loop through all extracted PDF URLs
	for _, urls := range extractedPDFURLs {
		if !hasDomain(urls) {
			urls = remoteDomainName + urls

		}
		if isUrlValid(urls) { // Check if the final URL is valid
			downloadPDF(urls, outputDir) // Download the PDF
		}
	}
}
