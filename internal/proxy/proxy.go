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
	let incomingQueue = [];
	let outboundQueue = [];
	let autoScrolling = false;

	setInterval(function() {
		if (incomingQueue.length === 0) {
			return;
		}

		const lastEvent = incomingQueue[incomingQueue.length - 1];		
		document.documentElement.scrollTo({top: lastEvent.data.scrollTop, behavior: 'smooth'});

		incomingQueue = [];
	}, 100)

	setInterval(function() {
		if (outboundQueue.length === 0) {
			return;
		}

		const lastEvent = outboundQueue[outboundQueue.length - 1];
		window.parent.postMessage(lastEvent, "*")
		outboundQueue = [];
	}, 2000)

	window.addEventListener("message", (event) => {
		if (event.data.sender === identifier) {
			return;
		}

		if (event.data.type === 'scroll') {
			if (document.documentElement.scrollTop === event.data.scrollTop) {
				return;
			}

			incomingQueue.push(event)
		}
	});

	// if triggerScrollEvent is true, send event on scroll
	if (%t) {
		document.addEventListener("scroll", (event) => {
			window.parent.postMessage({"sender": identifier, "msg": "hello there from your child", "type": "scroll", "scrollTop": document.documentElement.scrollTop}, "*")
			// outboundQueue.push({"sender": identifier, "msg": "hello there from your child", "type": "scroll", "scrollTop": document.documentElement.scrollTop});
		});
	}

	</script>
`

func AddPreviewFunctions(input []byte, triggerScrollEvent bool) ([]byte, error) {
	replaceWith := fmt.Sprintf("<head>%s", fmt.Sprintf(ProxyFunctions, uuid.New().String(), triggerScrollEvent))
	return []byte(strings.Replace(string(input), "<head>", replaceWith, 1)), nil
}

func ServeHTTP(w http.ResponseWriter, r *http.Request, triggerScrollEvent bool) {
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

	d, err := AddPreviewFunctions(rawRespBody, triggerScrollEvent)
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
