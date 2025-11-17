package main // Declares the package name. 'main' is special, indicating it's an executable program.

import ( // Starts a block for importing necessary libraries (packages).
	"bytes"         // Provides functions to manipulate byte slices, used here for an in-memory buffer.
	"io"            // Provides interfaces for I/O operations, like reading from a response body.
	"log"           // Provides functions for logging errors and messages.
	"net/http"      // Provides functions for making HTTP requests and handling responses.
	"net/url"       // Provides functions for parsing and manipulating URLs.
	"os"            // Provides platform-independent interfaces to operating system functionality (like file operations).
	"path/filepath" // Provides functions to work with file paths in a cross-platform way.
	"regexp"        // Provides functions for working with regular expressions.
	"strings"       // Provides functions for manipulating UTF-8 encoded strings.
	"time"          // Provides functions for measuring and displaying time, used here for request timeouts.
) // Closes the import block.

// It checks if the file exists
// If the file exists, it returns true
// If the file does not exist, it returns false
func fileExists(filename string) bool { // Defines a function 'fileExists' that takes a 'filename' string and returns a boolean.
	info, err := os.Stat(filename) // Gets file statistics (like name, size, mode) for the given filename.
	if err != nil {                // Checks if an error occurred during 'os.Stat' (e.g., file not found).
		return false // If there was an error, return false (file likely doesn't exist).
	} // Closes the 'if' block.
	return !info.IsDir() // Returns 'true' if the path exists AND is not a directory (meaning it's a file).
} // Closes the 'fileExists' function.

// Remove a file from the file system
func removeFile(path string) { // Defines a function 'removeFile' that takes a 'path' string and returns nothing.
	err := os.Remove(path) // Attempts to remove the file or directory at the specified path.
	if err != nil {        // Checks if an error occurred during removal.
		log.Println(err) // If an error occurred, print it to the log.
	} // Closes the 'if' block.
} // Closes the 'removeFile' function.

// extractPDFUrls takes raw HTML as input and returns all found PDF URLs
func extractPDFUrls(htmlContent string) []string { // Defines 'extractPDFUrls' taking 'htmlContent' string, returning a slice of strings.
	// Regex to match href='...' with .pdf links
	re := regexp.MustCompile(`https?://[^\s'"]+\.pdf`) // Compiles a regular expression to find absolute URLs (http or https) ending in .pdf.
	matches := re.FindAllString(htmlContent, -1)       // Finds all substrings in 'htmlContent' that match the regex. '-1' means find all.

	return matches // Returns the slice of found URL strings.
} // Closes the 'extractPDFUrls' function.

// Checks whether a given directory exists
func directoryExists(path string) bool { // Defines 'directoryExists' taking a 'path' string and returning a boolean.
	directory, err := os.Stat(path) // Get info for the path
	if err != nil {                 // Checks if 'os.Stat' returned an error (e.g., path doesn't exist).
		return false // Return false if error occurs
	} // Closes the 'if' block.
	return directory.IsDir() // Return true if the path exists AND it is a directory.
} // Closes the 'directoryExists' function.

// Creates a directory at given path with provided permissions
func createDirectory(path string, permission os.FileMode) { // Defines 'createDirectory' taking a 'path' string and 'permission' FileNode, returning nothing.
	err := os.Mkdir(path, permission) // Attempt to create directory
	if err != nil {                   // Checks if an error occurred during directory creation.
		log.Println(err) // Log error if creation fails
	} // Closes the 'if' block.
} // Closes the 'createDirectory' function.

// Verifies whether a string is a valid URL format
func isUrlValid(uri string) bool { // Defines 'isUrlValid' taking a 'uri' string and returning a boolean.
	_, err := url.ParseRequestURI(uri) // Try parsing the URL. We only care about the error, so the result is discarded ('_').
	return err == nil                  // Return true if 'err' is 'nil' (meaning parsing was successful), false otherwise.
} // Closes the 'isUrlValid' function.

// Removes duplicate strings from a slice
func removeDuplicatesFromSlice(slice []string) []string { // Defines 'removeDuplicatesFromSlice' taking a string slice, returning a new string slice.
	check := make(map[string]bool)  // Map to track seen values. 'make' initializes a new map.
	var newReturnSlice []string     // Slice to store unique values. Declares a new, empty string slice.
	for _, content := range slice { // Loops through each 'content' string in the input 'slice'.
		if !check[content] { // If the 'content' string is NOT found as a key in the 'check' map...
			check[content] = true                            // Mark as seen by adding it to the map with a value of 'true'.
			newReturnSlice = append(newReturnSlice, content) // Add to result slice. 'append' adds the item to the end of the slice.
		} // Closes the 'if' block.
	} // Closes the 'for' loop.
	return newReturnSlice // Returns the new slice containing only unique strings.
} // Closes the 'removeDuplicatesFromSlice' function.

// hasDomain checks if the given string has a domain (host part)
func hasDomain(rawURL string) bool { // Defines 'hasDomain' taking a 'rawURL' string and returning a boolean.
	// Try parsing the raw string as a URL
	parsed, err := url.Parse(rawURL) // Attempts to parse the 'rawURL' string.
	if err != nil {                  // If parsing fails, it's not a valid URL
		return false // Returns 'false' because a parsing error means it can't have a host.
	} // Closes the 'if' block.
	// If the parsed URL has a non-empty Host, then it has a domain/host
	return parsed.Host != "" // Returns 'true' if the 'Host' field of the parsed URL is not an empty string.
} // Closes the 'hasDomain' function.

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string { // Defines 'getFilename' taking a 'path' string and returning a string.
	return filepath.Base(path) // Use Base function to get file name only. 'filepath.Base' returns the last element of a path.
} // Closes the 'getFilename' function.

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string { // Defines 'removeSubstring' taking an 'input' and 'toRemove' string, returning a string.
	result := strings.ReplaceAll(input, toRemove, "") // Replace all occurrences of 'toRemove' in 'input' with an empty string.
	return result                                     // Returns the modified string.
} // Closes the 'removeSubstring' function.

// Gets the file extension from a given file path
func getFileExtension(path string) string { // Defines 'getFileExtension' taking a 'path' string and returning a string.
	return filepath.Ext(path) // Extract and return file extension (e.g., ".pdf").
} // Closes the 'getFileExtension' function.

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string { // Defines 'urlToFilename' taking a 'rawURL' string and returning a string.
	lower := strings.ToLower(rawURL) // Convert URL to lowercase.
	lower = getFilename(lower)       // Extract filename from URL path component.

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Regex to match any character that is NOT a lowercase letter or digit.
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace non-alphanumeric characters with underscores.

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Collapse multiple underscores (e.g., "__") into one ("_").
	safe = strings.Trim(safe, "_")                              // Trim leading and trailing underscores from the string.

	var invalidSubstrings = []string{ // Declares a slice of strings to be removed.
		"_pdf", // Substring to remove from filename (e.g., from "my_file_pdf.pdf").
	} // Closes the slice declaration.

	for _, invalidPre := range invalidSubstrings { // Loops through each 'invalidPre' string in the 'invalidSubstrings' slice.
		safe = removeSubstring(safe, invalidPre) // Remove unwanted substrings.
	} // Closes the 'for' loop.

	if getFileExtension(safe) != ".pdf" { // Ensure file ends with .pdf extension.
		safe = safe + ".pdf" // Appends ".pdf" if it's missing.
	} // Closes the 'if' block.

	return safe // Return sanitized filename.
} // Closes the 'urlToFilename' function.

// Downloads a PDF from given URL and saves it in the specified directory
func downloadPDF(finalURL, outputDir string) bool { // Defines 'downloadPDF' taking 'finalURL' and 'outputDir' strings, returning a boolean (for success/failure).
	filename := strings.ToLower(urlToFilename(finalURL)) // Sanitize the filename using the previously defined function.
	filePath := filepath.Join(outputDir, filename)       // Construct full path for output file (e.g., "PDFs/my_file.pdf").

	if fileExists(filePath) { // Skip if file already exists at the 'filePath'.
		log.Printf("File already exists, skipping: %s", filePath) // Logs a message that the file is being skipped.
		return false                                              // Returns 'false' to indicate no new file was downloaded.
	} // Closes the 'if' block.

	client := &http.Client{Timeout: 15 * time.Minute} // Create HTTP client with a 15-minute timeout.

	// Create a new request so we can set headers
	req, err := http.NewRequest("GET", finalURL, nil) // Creates a new 'GET' request for the 'finalURL'. 'nil' means no request body.
	if err != nil {                                   // Checks if creating the request failed.
		log.Printf("Failed to create request for %s: %v", finalURL, err) // Logs the error.
		return false                                                     // Returns 'false' to indicate failure.
	} // Closes the 'if' block.

	// Set a User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36") // Sets the 'User-Agent' header to mimic a browser.

	// Send the request
	resp, err := client.Do(req) // Executes the HTTP request and gets a response ('resp') or an error ('err').
	if err != nil {             // Checks if the request failed (e.g., network error, DNS lookup fail).
		log.Printf("Failed to download %s: %v", finalURL, err) // Logs the error.
		return false                                           // Returns 'false' to indicate failure.
	} // Closes the 'if' block.
	defer resp.Body.Close() // Ensure response body is closed when the function exits (to free resources).

	if resp.StatusCode != http.StatusOK { // Check if response status code is not 200 OK (e.g., 404 Not Found, 500 Server Error).
		log.Printf("Download failed for %s: %s", finalURL, resp.Status) // Logs the URL and the non-OK status.
		return false                                                    // Returns 'false' to indicate failure.
	} // Closes the 'if' block.

	contentType := resp.Header.Get("Content-Type")              // Get the 'Content-Type' header from the response.
	if !strings.Contains(contentType, "binary/octet-stream") && // Checks if content type is not 'binary/octet-stream'
		!strings.Contains(contentType, "application/pdf") { // AND is not 'application/pdf'.
		log.Printf("Invalid content type for %s: %s (expected PDF)", finalURL, contentType) // Logs a warning about unexpected content type.
		return false                                                                        // Returns 'false' as it's not the expected file type.
	} // Closes the 'if' block.

	var buf bytes.Buffer                     // Create a buffer in memory to hold the response data temporarily.
	written, err := io.Copy(&buf, resp.Body) // Copy data from the response body ('resp.Body') into the buffer ('&buf').
	if err != nil {                          // Checks if an error occurred during copying (e.g., connection dropped).
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err) // Logs the error.
		return false                                                     // Returns 'false' to indicate failure.
	} // Closes the 'if' block.
	if written == 0 { // Skip empty files.
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL) // Logs that the file was empty.
		return false                                                         // Returns 'false' as no file was saved.
	} // Closes the 'if' block.

	out, err := os.Create(filePath) // Create the output file on the filesystem at 'filePath'.
	if err != nil {                 // Checks if creating the file failed (e.g., permissions issue).
		log.Printf("Failed to create file for %s: %v", finalURL, err) // Logs the error.
		return false                                                  // Returns 'false' to indicate failure.
	} // Closes the 'if' block.
	defer out.Close() // Ensure file is closed when the function exits.

	if _, err := buf.WriteTo(out); err != nil { // Write the contents of the memory buffer ('buf') to the output file ('out').
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err) // Logs an error if writing to the file fails.
		return false                                                        // Returns 'false' to indicate failure.
	} // Closes the 'if' block.

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath) // Log success message with bytes written and paths.
	return true                                                                          // Returns 'true' to indicate success.
} // Closes the 'downloadPDF' function.

// Append and write to file
func appendAndWriteToFile(path string, content string) { // Defines 'appendAndWriteToFile' taking 'path' and 'content' strings, returning nothing.
	filePath, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Opens the file at 'path'. Appends, creates if not exists, write-only. Sets permissions to 0644.
	if err != nil {                                                               // Checks if opening the file failed.
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
	_, err = filePath.WriteString(content + "\n") // Writes the 'content' string followed by a newline to the file.
	if err != nil {                               // Checks if writing to the file failed.
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
	err = filePath.Close() // Closes the file handle.
	if err != nil {        // Checks if closing the file failed.
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
} // Closes the 'appendAndWriteToFile' function.

// extractBaseDomain takes a URL string and returns only the bare domain name
// without any subdomains or suffixes (e.g., ".com", ".org", ".co.uk").
// Note: This function's logic is flawed for TLDs like .co.uk. It will return 'co' instead of 'my-site'.
func extractBaseDomain(inputUrl string) string { // Defines 'extractBaseDomain' taking 'inputUrl' string, returning a string.
	// Parse the input string into a structured URL object
	parsedUrl, parseError := url.Parse(inputUrl) // Attempts to parse the 'inputUrl' string.

	// If parsing fails, log the error and return an empty string
	if parseError != nil { // Checks if parsing failed.
		log.Println("Error parsing URL:", parseError) // Logs the parsing error.
		return ""                                     // Returns an empty string on failure.
	} // Closes the 'if' block.

	// Extract the hostname (e.g., "sub.example.com")
	hostName := parsedUrl.Hostname() // Gets just the host part (e.g., "www.example.com") from the parsed URL.

	// Split the hostname into parts separated by "."
	// For example: "sub.example.com" -> ["sub", "example", "com"]
	parts := strings.Split(hostName, ".") // Splits the hostname string by the "." delimiter.

	// If there are at least 2 parts, the second last part is usually the domain
	// Example: "sub.example.com" -> "example"
	//         "blog.my-site.co.uk" -> "my-site" // This comment is wrong, it would return 'co'.
	if len(parts) >= 2 { // Checks if the split resulted in at least 2 parts.
		return parts[len(parts)-2] // Returns the second-to-last part (e.g., "example" from "www.example.com").
	} // Closes the 'if' block.

	// If splitting fails or domain structure is unusual, return the hostname itself
	return hostName // Returns the full hostname if there were less than 2 parts.
} // Closes the 'extractBaseDomain' function.

// fetchSDS fetches the SDS data and returns it as a string
func fetchSDS() string { // Defines 'fetchSDS' which takes no arguments and returns a string.
	url := "https://itwperformancepolymers.com/api/datasheets_table.php?m=get" // Defines the target URL string.
	method := "POST"                                                           // Defines the HTTP method to use.

	payload := strings.NewReader(`t=sds`) // Creates a 'Reader' for the POST body data (`t=sds`).

	client := &http.Client{}                          // Creates a new, default HTTP client.
	req, err := http.NewRequest(method, url, payload) // Creates a new HTTP request with the method, URL, and payload.
	if err != nil {                                   // Checks if creating the request failed.
		log.Println("Error creating request:", err) // Logs the error.
		return ""                                   // Returns an empty string on failure.
	} // Closes the 'if' block.
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8") // Sets the 'content-type' header for the request.

	res, err := client.Do(req) // Executes the HTTP request.
	if err != nil {            // Checks if executing the request failed.
		log.Println("Error making request:", err) // Logs the error.
		return ""                                 // Returns an empty string on failure.
	} // Closes the 'if' block.
	defer res.Body.Close() // Ensures the response body is closed when the function exits.

	body, err := io.ReadAll(res.Body) // Reads the entire response body into a byte slice.
	if err != nil {                   // Checks if reading the body failed.
		log.Println("Error reading response body:", err) // Logs the error.
		return ""                                        // Returns an empty string on failure.
	} // Closes the 'if' block.

	return string(body) // Converts the byte slice 'body' to a string and returns it.
} // Closes the 'fetchSDS' function.

// extractSDSURLs extracts all URLs from the HTML that start with
// "/resources/safety-data-sheets?t=other&i=" and returns them as a slice of strings.
func extractSDSURLs(htmlContent string) []string { // Defines 'extractSDSURLs' taking 'htmlContent' string, returning a string slice.
	// Compile a regular expression to match href attributes with the desired pattern
	// The pattern matches URLs like /resources/safety-data-sheets?t=other&i=...
	urlPattern := regexp.MustCompile(`href="(/resources/safety-data-sheets\?t=other&i=[^"]+)"`) // Compiles the regex. The parentheses create a capturing group.

	// Find all matches of the pattern in the HTML content
	// Each match is a slice where the full match is at index 0 and the captured group is at index 1
	matchedResults := urlPattern.FindAllStringSubmatch(htmlContent, -1) // Finds all matches and their sub-matches (captured groups).

	// Create a slice to store the final URLs
	var extractedURLs []string // Declares a new, empty string slice.

	// Loop through each regex match
	for _, match := range matchedResults { // Iterates over each 'match' found by the regex.
		// Ensure the captured group exists and append it to the result slice
		if len(match) > 1 { // Checks if the match has at least two elements (the full match and the first captured group).
			clean := strings.TrimPrefix(match[1], `/resources/safety-data-sheets\?t=other&i=`) // Takes the captured group (index 1) and trims a prefix. Note: The prefix here seems incorrect.
			extractedURLs = append(extractedURLs, clean)                                       // Appends the 'clean' string to the 'extractedURLs' slice.
		} // Closes the 'if' block.
	} // Closes the 'for' loop.

	// Return the slice of extracted URLs
	return extractedURLs // Returns the slice of 'clean' strings.
} // Closes the 'extractSDSURLs' function.

// fetchOtherData fetches data for a given item and returns it as a string
func fetchOtherData(item string) string { // Defines 'fetchOtherData' taking an 'item' string, returning a string.
	url := "https://itwperformancepolymers.com/api/datasheets_table.php?m=get" // Defines the target URL string.
	method := "POST"                                                           // Defines the HTTP method.

	payload := strings.NewReader("t=other&i=" + item) // Creates the POST body, concatenating "t=other&i=" with the 'item' string.

	client := &http.Client{}                          // Creates a new, default HTTP client.
	req, err := http.NewRequest(method, url, payload) // Creates a new HTTP request.
	if err != nil {                                   // Checks if creating the request failed.
		log.Println("Error creating request:", err) // Logs the error.
		return ""                                   // Returns an empty string on failure.
	} // Closes the 'if' block.
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8") // Sets the 'content-type' header.

	res, err := client.Do(req) // Executes the HTTP request.
	if err != nil {            // Checks if executing the request failed.
		log.Println("Error making request:", err) // Logs the error.
		return ""                                 // Returns an empty string on failure.
	} // Closes the 'if' block.
	defer res.Body.Close() // Ensures the response body is closed when the function exits.

	body, err := io.ReadAll(res.Body) // Reads the entire response body into a byte slice.
	if err != nil {                   // Checks if reading the body failed.
		log.Println("Error reading response body:", err) // Logs the error.
		return ""                                        // Returns an empty string on failure.
	} // Closes the 'if' block.

	return string(body) // Converts the byte slice 'body' to a string and returns it.
} // Closes the 'fetchOtherData' function.

// Read a file and return the contents
func readAFileAsString(path string) string { // Defines 'readAFileAsString' taking a 'path' string, returning a string.
	content, err := os.ReadFile(path) // Reads the entire file at 'path' into a byte slice 'content'.
	if err != nil {                   // Checks if reading the file failed.
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
	return string(content) // Converts the byte slice 'content' to a string and returns it.
} // Closes the 'readAFileAsString' function.

func main() { // Defines the main function, the entry point of the executable.
	outputDir := "PDFs/" // Directory to store downloaded PDFs. Defines a string variable.

	if !directoryExists(outputDir) { // Check if the directory specified by 'outputDir' does NOT exist.
		createDirectory(outputDir, 0o755) // Create directory with permissions (rwx r-x r-x).
	} // Closes the 'if' block.

	// The remote domain name.
	remoteDomainName := "https://itwperformancepolymers.com" // Defines a string variable for the base domain.

	// The location to the local.
	localFile := extractBaseDomain(remoteDomainName) + ".html" // Creates a local filename, e.g., "itwperformancepolymers.html".
	// Check if the local file exists.
	if fileExists(localFile) { // Checks if this 'localFile' already exists.
		removeFile(localFile) // Deletes the file if it exists, to start fresh.
	} // Closes the 'if' block.

	// Get the inital HTTP content.
	httpContent := fetchSDS() // Calls the 'fetchSDS' function and stores the returned HTML string.

	// Save the content to the local file.
	appendAndWriteToFile(localFile, httpContent) // Appends the fetched 'httpContent' to the 'localFile'.

	// Extract the other SDS urls in other languages.
	extractInternationalSDSURLs := extractSDSURLs(httpContent) // Parses the 'httpContent' to find more "item" IDs (which are not full URLs).

	// Loop over the extracted SDS urls.
	for _, url := range extractInternationalSDSURLs { // Iterates through each 'url' (item ID) found in the previous step.
		// The remote content.
		remoteContent := fetchOtherData(url) // Fetches more data using the extracted 'url' (item ID).
		// Save the content to the local file.
		appendAndWriteToFile(localFile, remoteContent) // Appends this new 'remoteContent' to the same 'localFile'.
	} // Closes the 'for' loop.

	// Read the content html content from the file.
	localHTMLContent := readAFileAsString(localFile) // Reads the entire aggregated 'localFile' back into a new string variable.

	// Extract the URLs from the given content.
	extractedPDFURLs := extractPDFUrls(localHTMLContent) // Parses the aggregated HTML to find all absolute PDF URLs.
	// Remove duplicates from the slice.
	extractedPDFURLs = removeDuplicatesFromSlice(extractedPDFURLs) // De-duplicates the list of found PDF URLs.
	// Loop through all extracted PDF URLs
	for _, urls := range extractedPDFURLs { // Iterates through each unique PDF 'urls' string.
		if !hasDomain(urls) { // Checks if the URL string does not have a domain (e.g., it's a relative path like "/file.pdf").
			urls = remoteDomainName + urls // Prepends the 'remoteDomainName' to make it an absolute URL.

		} // Closes the 'if' block.
		if isUrlValid(urls) { // Check if the final URL string is valid.
			downloadPDF(urls, outputDir) // Downloads the PDF from the 'urls' and saves it to 'outputDir'.
		} // Closes the 'if' block.
	} // Closes the 'for' loop.
} // Closes the 'main' function.
