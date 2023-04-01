package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rajatjindal/preview-compare/internal/preview"
)

func (s *Server) CreateNewPreviewRequest(w http.ResponseWriter, r *http.Request) {
	inp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := &preview.PreviewRequest{}
	err = json.Unmarshal(inp, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	preview, err := s.store.CreatePreview(r.Context(), req)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(preview)
}

func (s *Server) ComparePreviewWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	preview, err := s.store.GetPreview(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(os.Stderr, fmt.Sprintf("preview: %#v", preview))
	var rawhtml = `<!doctype html>
<html>

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body>
    <div class="grid grid-cols-2 gap-1 mx-full border-2 border-blue-100 w-1/2 mx-auto my-auto h-1/2 mt-20">
		<div id="container-1" class="col-span-1 w-full border border-red-900 h-screen">
			<iframe id="frame-id-1" src="https://preview-1-wpsr7vaf.fermyon.app?base=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
		</div>
		<div id="container-1" class="col-span-1 w-full border border-red-900 h-screen">
			<iframe id="frame-id-1" src="https://preview-2-b6p5mwqe.fermyon.app?base=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
		</div>
	</div>

	<script>
	// setInterval(function() {
	// 	document.getElementById('frame-id-1').contentWindow.postMessage("hello there from your parent", "*");
	// }, 2000);
	
	</script>
</body>

</html>
`
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(fmt.Sprintf(rawhtml, preview.ThisBase, preview.ThatBase)))
}
