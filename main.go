package main

import (
	"fmt"
	"io"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
)

var frame = `<!doctype html>
<html>

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body class="w-full h-screen">
	<iframe style="width: 100%%; height: 100%%" id="%s" src="%s" frameborder="0"></iframe>
</body>

</html>
`

var html = `<!doctype html>
<html>

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body>
    <div class="grid grid-cols-2 gap-1 mx-full h-screen border-2 border-blue-100">
        <div class="col-span-1 h-screen w-full border border-red-900">
            <iframe id="frame-id-1" src="http://localhost:3000/first" frameborder="0" style="width: 100%; height: 100%"></iframe>
        </div>
        <div class="col-span-1 h-full w-full">
            <iframe id="frame-id-2" src="http://localhost:3000/second" frameborder="0" style="width: 100%; height: 100%"></iframe>
        </div>
    </div>
</body>

</html>
`

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	for k, v := range resp.Header {
		for _, vi := range v {
			w.Header().Add(k, vi)
		}
	}
	io.Copy(w, resp.Body)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/first" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(frame, "id-1", "https://fermyon.com")))
		return
	}

	if r.URL.Path == "/second" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(frame, "id-2", "https://fermyon.com")))
		return
	}

	if r.URL.Path == "/compare" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("default route"))
}

func init() {
	spinhttp.Handle(handler)
}

func main() {}
