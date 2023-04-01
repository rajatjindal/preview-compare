package proxy

import (
	"fmt"
	"strings"
)

const fns = `
	console.log('proxying'); 
	window.addEventListener("message", (event) => {
		console.log("hello from child");
		console.log(event);
	});

	console.log('proxied'); 
	document.body.scrollTop=600; 
	document.documentElement.scrollTop=600;
`

func AddPreviewFunctions(input []byte) ([]byte, error) {
	return []byte(strings.Replace(string(input), "<script>", fmt.Sprintf("<script>%s", fns), 1)), nil
}
