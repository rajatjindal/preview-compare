package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
)

//https://gist.github.com/raedatoui/b33fac34fb24ae5ecaabd5f7b3b67e0c

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
			<iframe id="frame-id-1" src="http://localhost:3000/first?base=https://rajatjindal.com" frameborder="0" style="width: 100%; height: 100%;"></iframe>
		</div>
		<div id="container-1" class="col-span-1 w-full border border-red-900 h-screen">
			<iframe id="frame-id-1" src="http://localhost:3000/first?base=https://quicksilver1183.com" frameborder="0" style="width: 100%; height: 100%;"></iframe>
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
	fmt.Fprint(os.Stderr, fmt.Sprintf("%s\n", r.URL.Path))
	if r.URL.Path != "/compare" {
		original, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		after, err := url.Parse("https://rajatjindal.com")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		after.Path, _ = strings.CutPrefix(original.Path, "/first")

		req, err := http.NewRequest(r.Method, after.String(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := spinhttp.Send(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			fmt.Fprint(os.Stderr, fmt.Sprintf("%#v\n", resp.Header))
			writeResponse(w, resp)
			return
		}

		fmt.Fprint(os.Stderr, fmt.Sprintf("%#v\n", resp.Header))
		d, err := inject(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		for k, v := range resp.Header {
			for _, vi := range v {
				w.Header().Add(k, vi)
			}
		}

		w.Write(d)
		return
	}

	if r.URL.Path == "/second" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}

	if r.URL.Path == "/compare" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawhtml))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("default route"))
}

func init() {
	spinhttp.Handle(handler)
}

func main() {}

var fns = `
console.log('injected'); 
window.addEventListener("message", (event) => {
	console.log("hello from child");
	console.log(event);
});

console.log('injected'); 
document.body.scrollTop=600; 
document.documentElement.scrollTop=600;
`

func inject(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	x := strings.Replace(string(b), "<script>", fmt.Sprintf("<script>%s", fns), 1)
	return []byte(x), nil
}
