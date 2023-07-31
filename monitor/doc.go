// func main() {
// 	dirPath := "/path/to/directory" // Replace with your desired directory path.
// 	MaxSizeBytes := int64(100 * 1024 * 1024) // 100 MB
// 	interval := 10 * time.Second // Adjust the monitoring interval as needed.
//
// 	// Create a context with cancel to control the monitoring process.
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel() // Call cancel function when main exits to release resources.
//
// 	// Initialize the DirectoryMonitor
// 	monitor := NewDirectoryMonitor(dirPath, MaxSizeBytes, interval)
// 	// Start the monitoring process.
// 	go monitor.Start(ctx)
//
// 	// Run the main program for a while.
// 	time.Sleep(60 * time.Second)
//
// 	// Stop the monitoring process by canceling the context.
// 	cancel()
//
// 	// Wait for the monitoring process to exit gracefully.
// 	time.Sleep(1 * time.Second)
// 	fmt.Println("Main program finished.")
// }

package monitor
