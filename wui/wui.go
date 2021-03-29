package wui

import "net/http"

func ServeWUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//TODO: Decide if we want to inline or serve in a different method
	http.ServeFile(w, r, "./wui/index.html")
}
