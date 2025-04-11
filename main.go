package main

import (
	"fmt"
)

func main() {
	_, err := Get_vimeo_video("955118381")
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return
	}

	// websocket_url, err := Transfer_from_vimeo(url)
	// if err != nil {
	//   fmt.Println("Error transferring video:", err)
	//   return
	// }
	// fmt.Println("Websocket: " + websocket_url)
}
