package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/google/uuid"
)

const ProxyFunctions = `
	<script>
	const identifier = '%s';
	let eventsQueue = [];

	setInterval(function() {
		// console.log("identifier: ", identifier, " queue length: ", eventsQueue.length);
		if (eventsQueue.length === 0) {
			// console.log("no events in scroll queue")
			return;
		}

		// const lastEvent = eventsQueue[eventsQueue.length - 1];
		console.log("updating scrollTop for self:  ", identifier, " for event: ", lastEvent.data);
		document.documentElement.scrollTo({top: lastEvent.data.scrollTop, behavior: 'smooth'});

		eventsQueue = [];
		// console.log("identifier: ", identifier, " resetting queue : ", eventsQueue.length);
	}, 2000)

	console.log('proxying'); 
	window.addEventListener("message", (event) => {
		// console.log("hello from inside child with self-identifier ", identifier);
		// console.log("self-identifier", identifier, " event.data: ", event.data);
		if (event.data.sender === identifier) {
			// console.log("ignore msg from self");
			return;
		}

		if (event.data.type === 'scroll') {
			if (document.documentElement.scrollTop === event.data.scrollTop) {
				// console.log("already updated scrollTop");
				return;
			}

			eventsQueue.push(event)
		}
	});

	document.addEventListener("scroll", (event) => {
		// console.log("inside scroll event ", event);
		// console.log("document.documentElement.scrollTop", document.documentElement.scrollTop);
		window.parent.postMessage({"sender": identifier, "msg": "hello there from your child", "type": "scroll", "scrollTop": document.documentElement.scrollTop}, "*")
	});

	</script>
`

func AddPreviewFunctions(input []byte) ([]byte, error) {
	replaceWith := fmt.Sprintf("<head>%s", fmt.Sprintf(ProxyFunctions, uuid.New().String()))
	return []byte(strings.Replace(string(input), "<head>", replaceWith, 1)), nil
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	previewBase := getPreviewBase(r)
	original, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	after, err := url.Parse(previewBase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	after.Path = original.Path

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
	rawRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d, err := AddPreviewFunctions(rawRespBody)
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

	expiration := time.Now().Add(2 * time.Minute)
	cookie := http.Cookie{Name: "previewBase", Value: previewBase, Expires: expiration}
	http.SetCookie(w, &cookie)

	w.Write(d)
	return
}

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	for k, v := range resp.Header {
		for _, vi := range v {
			w.Header().Add(k, vi)
		}
	}
	io.Copy(w, resp.Body)
}

func getPreviewBase(r *http.Request) string {
	previewBase := r.URL.Query().Get("previewBase")

	if previewBase != "" {
		return previewBase
	}

	for _, cookie := range r.Cookies() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%#v", cookie))
		if cookie.Name == "previewBase" {
			return cookie.Value
		}
	}

	return ""
}
