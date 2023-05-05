package api

import (
	"context"
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

	preview, err := s.createPreview(r.Context(), req)
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
	previewReq, err := s.store.GetPreview(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rendered, err := s.renderForCompare(r.Context(), previewReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.Write(rendered)
}

func (s *Server) ComparePreview(w http.ResponseWriter, r *http.Request) {
	url1 := r.URL.Query().Get("url1")
	url2 := r.URL.Query().Get("url2")

	req := &preview.PreviewRequest{
		ThisBase: url1,
		ThatBase: url2,
	}

	previewReq, err := s.createPreview(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rendered, err := s.renderForCompare(r.Context(), previewReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.Write(rendered)
}

// HELPER FUNCTIONS
func (s *Server) renderForCompare(ctx context.Context, previewReq *preview.PreviewRequest) ([]byte, error) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("preview: %#v", previewReq))
	var _ = `
<!doctype html>
<html>

<head>
	<script>
		let incomingQueue = [];
		setInterval(function() {
			if (incomingQueue.length === 0) {
				return;
			}

			const lastEvent = incomingQueue[incomingQueue.length - 1];
			document.getElementById('frame-id-2').contentWindow.postMessage(lastEvent.data, "*");
		
			incomingQueue = [];
		}, 1000)

		window.addEventListener("message", (event) => {
			// incomingQueue.push(event)
			document.getElementById('frame-id-2').contentWindow.postMessage(event.data, "*");
		});
	</script>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body>
    <div class="grid grid-cols-2 gap-1 mx-full border-2 w-full mx-auto my-auto h-screen">
		<div id="container-1" class="col-span-1 w-full h-screen resize-x">
			<iframe id="frame-id-1" src="https://preview-1.fermyon.app?previewBase=%s&previewId=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
		</div>
		<div id="container-1" class="col-span-1 w-full h-screen resize-x">
			<iframe id="frame-id-2" src="https://preview-2.fermyon.app?previewBase=%s&previewId=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
		</div>
	</div>
</body>

</html>
`

	return preview.Render(previewReq)
}

func (s *Server) createPreview(ctx context.Context, inp *preview.PreviewRequest) (*preview.PreviewRequest, error) {
	preview, err := s.store.CreatePreview(ctx, inp)
	if err != nil {
		return nil, err
	}

	return preview, nil
}
