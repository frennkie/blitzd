package serve

import (
	"fmt"
	"github.com/frennkie/blitzinfod/internal/data"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {

	htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon</title>
</head>
<body>
	<ul>	
		<li><a href="%s">API</a></li>
		<li><a href="/info/">Info Page</a></li>
	</ul>
	<br>

	<hr>
	%s
	<br>

	<hr>
	Request: 
	<pre>%s</pre>
</body>
</html>`

	values := []interface{}{data.APIv1, r.RemoteAddr, r.URL.Path}

	html := fmt.Sprintf(htmlRaw, values...)

	_, _ = fmt.Fprintf(w, "%s", html)
}
