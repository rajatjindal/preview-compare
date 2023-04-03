package preview

import (
	"fmt"
)

var rawTemplate1 = `
<!doctype html>
<html>

<head>
	<script>
		let incomingQueue = [];
		setInterval(function () {
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

	<script>
		function toggleLeftNav() {
			console.log('toggleLeftNav');
			document.getElementById('leftNav').classList.toggle('hidden');
			document.getElementById('bodyGrid').classList.toggle('w-10/12');
			document.getElementById('bodyGrid').classList.toggle('w-full');
		}

		function linkClick(link) {
			console.log('linkClick ', link);
			toggleLeftNav();
		}
	</script>
</head>

<body>
	<div class="h-12 py-3 bg-gray-900">
		<div class="my-auto text-indigo-400 flex" onclick="toggleLeftNav()">
			<div class="ml-2">
				<svg width="24" height="24" fill="none" viewBox="0 0 24 24">
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M4.75 5.75H19.25"></path>
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M4.75 18.25H19.25"></path>
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M4.75 12H19.25"></path>
				</svg>
			</div>
			<div class="ml-2 -mt-0.5 text-xl text-indigo-400 font-bold">Preview App</div>
		</div>
	</div>

	<div class="mx-full w-full mx-auto my-auto h-screen flex">
		<div id="leftNav" class="w-2/12 z-10 hidden col-span-1 flex-cols bg-gray-900 text-white text-xs">`

var changeLink = `
			<div class="mt-5 ml-3 text-sm font-medium flex" onclick="linkClick('%s')">
				<svg width="24" height="24" fill="none" viewBox="0 0 24 24">
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M7.75 19.25H16.25C17.3546 19.25 18.25 18.3546 18.25 17.25V9L14 4.75H7.75C6.64543 4.75 5.75 5.64543 5.75 6.75V17.25C5.75 18.3546 6.64543 19.25 7.75 19.25Z" />
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M18 9.25H13.75V5" />
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M9.75 15.25H14.25" />
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
						d="M9.75 12.25H14.25" />
				</svg>
				<div class="mt-1">%s</div>
			</div>
`

var rawTemplate2 = `
		</div>
		<div id="bodyGrid" class="w-full grid grid-cols-8 gap-1">
			<div id="container-1" class="col-span-4 w-full h-screen resize-x">
				<iframe id="frame-id-1" src="https://preview-1-wpsr7vaf.fermyon.app?previewBase=%s&previewId=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
			</div>
			<div id="container-1" class="col-span-4 w-full h-screen resize-x">
				<iframe id="frame-id-2" src="https://preview-2-b6p5mwqe.fermyon.app?previewBase=%s&previewId=%s" frameborder="0" style="width: 100%%; height: 100%%;"></iframe>
			</div>
		</div>
	</div>
</body>
</html>`

func Render(req *PreviewRequest) ([]byte, error) {
	changeLinks := ""
	for _, change := range req.Changes {
		changeLinks += fmt.Sprintf(changeLink, change.Link, change.LinkText)
	}
	lastPart := fmt.Sprintf(rawTemplate2, req.ThisBase, req.Id, req.ThatBase, req.Id)
	return []byte(rawTemplate1 + changeLinks + lastPart), nil
}
