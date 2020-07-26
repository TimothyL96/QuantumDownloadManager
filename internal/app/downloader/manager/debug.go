package manager

import "fmt"

// DebugUrl prints download URL details.
func (d *Download) DebugUrl() {
	// Debug print URL
	fmt.Println("URL DEBUG:")
	fmt.Println("URL scheme:", d.downloadURL.Scheme)
	fmt.Println("URL host:", d.downloadURL.Host)
	fmt.Println("URL Path:", d.downloadURL.Path)
	fmt.Println()
}

// DebugHeader prints download header details.
func (d *Download) DebugHeader() {
	// Debug print response header
	fmt.Println("HEADER DEBUG:")
	fmt.Println("Is download initialized:", d.isDownloadInitialized)
	fmt.Println("Response content length:", d.response.ContentLength)
	fmt.Println("Response headers:", d.response.Header)
	fmt.Println("Response status code:", d.response.StatusCode)
	fmt.Println()
}
