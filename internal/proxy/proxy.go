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
	const originalBase = '%s';
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

		if (event.data.type === 'navigate') {
			if (window.location.pathname === event.data.navigateTo) {
				return;
			}

			window.location.pathname = event.data.navigateTo
		}
	});

	 // if triggerScrollEvent is true, send event on scroll
	if (%t) {
		document.addEventListener("scroll", (event) => {
			window.parent.postMessage({"sender": identifier, "msg": "hello there from your child", "type": "scroll", "scrollTop": document.documentElement.scrollTop}, "*")
			// outboundQueue.push({"sender": identifier, "msg": "hello there from your child", "type": "scroll", "scrollTop": document.documentElement.scrollTop});
		});

		document.addEventListener('DOMContentLoaded', function() {
			const allLinks = document.getElementsByTagName('a');
			allLinks.forEach(item => {
				console.log("before => ", item.href);
				console.log("starts with => ", item.href.startsWith("/"));
				if (item.href.startsWith("/")) {
					console.log("first if");
					item.setAttribute('href', 'https://preview-1-wpsr7vaf.fermyon.app' + item.href)
				}
	
				if (item.href.startsWith(originalBase)) {
					console.log("second if");
					item.setAttribute('href', item.href.replace(originalBase, 'https://preview-1-wpsr7vaf.fermyon.app'));
				}
	
				console.log("after ", item);
			});

			window.parent.postMessage({"sender": identifier, "msg": "hello there from your child", "type": "navigate", "navigateTo": window.location.pathname}, "*");
		}, false);
	}
	</script>
`

func AddPreviewFunctions(input []byte, previewBase string, triggerScrollEvent bool) ([]byte, error) {
	replaceWith := fmt.Sprintf("<head>%s", fmt.Sprintf(ProxyFunctions, uuid.New().String(), previewBase, triggerScrollEvent))
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

	d, err := AddPreviewFunctions(rawRespBody, previewBase, triggerScrollEvent)
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
