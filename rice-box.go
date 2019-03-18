package main

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "status.css",
		FileModTime: time.Unix(1552939422, 0),

		Content: string("body {\n    background-color: black;\n    font-family: sans-serif;\n    color:LightSteelBlue;\n}\n\np {\n    font-size: 12px;\n}\n\n#lcd-regular {\n    /* The size of the LCD on shopping list */\n    width: 480px;\n    height: 320px;\n}\n\n#nodename {\n    margin: 3px;\n    font-size: 20px;\n    font-weight: bold;\n}\n\n#nodeversion {\n    margin: 3px;\n    font-size: 16px;\n    float: right;\n}\n\n.status-bar {\n    width: 100%;\n    height: 40px;\n}\n\n.status-box {\n    width: 50%;\n    height: 45%;\n    box-sizing: border-box;\n    margin: 0px;\n    padding: 0px;\n    float: left;\n}\n\n.bar-background {\n    fill: lightgrey;\n}\n\n.bar-load-now {\n    fill: blue;\n}\n\n.bar-load-avg5 {\n    fill: lightblue;\n}\n\n.bar-load-avg5 {\n    fill: grey;\n}\n\n\n\n#meter{\n    display: inline;\n    width: 100%;\n    height: 100%;\n    transform: rotate(180deg);\n    transform: scale(-0.9);\n}\n\n#wrapper-temp {\n    width: 28px;\n    height: 100px;\n    margin-left: 0px;\n    margin-right: 3px;\n}\n\n#wrapper-mem {\n    width: 50px;\n    height: 100px;\n    margin-left: 3px;\n    margin-right: 3px;\n}\n\n#wrapper-hdd {\n    width: 50px;\n    height: 100px;\n    margin-left: 3px;\n    margin-right: 3px;\n}\n\n#wrapper-load {\n    margin-top:37px;\n    width: 70px;\n    height: 45px;\n}\n\n.status-symbol {\n    float: left;\n    width: 85px;\n    height: 100%;\n    margin-right: 1em;\n    padding: 0px;\n    box-sizing: border-box;\n}\n.status-text {\n    font-family: monospace;\n    padding: 2px;\n    font-size: 8px;\n    text-align: center;\n}\n\nsvg {\n    width: 100%;\n    height: 100%;\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1552939422, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "status.css"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`static`, &embedded.EmbeddedBox{
		Name: `static`,
		Time: time.Unix(1552939422, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"status.css": file2,
		},
	})
}

func init() {

	// define files
	file4 := &embedded.EmbeddedFile{
		Filename:    "info.tmpl",
		FileModTime: time.Unix(1552941957, 0),

		Content: string("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\" />\n    <title>RaspiBlitz Status</title>\n    <link rel=\"stylesheet\" href=\"../static/status.css\">\n</head>\n\n<body>\n\n<div id=\"lcd-regular\">\n    <div class=\"status-bar\">\n        <span id=\"nodename\">Status ⚡ Mockup</span>\n        <span id=\"nodeversion\">RaspiBlitz v1.0 (Up: {{ .Uptime.Value }})</span>\n    </div>\n    <div class=\"status-box\">\n        <div id=\"wrapper-temp\" class=\"status-symbol\">\n            <svg id=\"bar-temp\">\n                <rect id=\"bar-temp-background\"\n                      x=\"8\" y=\"2\"\n                      rx=\"5\" ry=\"5\"\n                      width=\"14\" height=\"90\"\n                      fill=\"#404040\" />\n                <circle r=\"12\" cx=\"15\" cy=\"85\" fill=\"#404040\" />\n\n                <rect\n                        x=\"11\" y=\"4\"\n                        rx=\"3\" ry=\"3\"\n                        width=\"8\" height=\"85\"\n                        fill=\"black\" />\n                <circle r=\"9\" cx=\"15\" cy=\"85\" fill=\"black\" />\n\n                <rect id=\"bar-temp-display\"\n                      x=\"12\" y=\"40\"\n                      rx=\"0\" ry=\"0\"\n                      width=\"6\" height=\"45\"\n                      fill=\"green\" />\n                <circle r=\"8\" cx=\"15\" cy=\"85\" fill=\"green\" />\n            </svg>\n            <div class=\"status-text\">56°C CPU</div>\n        </div>\n\n        <div id=\"wrapper-mem\" class=\"status-symbol\">\n            <svg id=\"bar-mem\">\n                <defs>\n                    <g id=\"mem-foot\">\n                        <rect x=\"4\" y=\"0\" width=\"6\" height=\"10\" fill=\"#404040\" />\n                        <rect x=\"0\" y=\"3\" width=\"8\" height=\"4\"  fill=\"#404040\" />\n                    </g>\n                </defs>\n\n                <rect id=\"bar-mem-background\"\n                      x=\"8\" y=\"0\"\n                      rx=\"3\" ry=\"3\"\n                      width=\"35\" height=\"100\"\n                      stroke=\"dimgray\"\n                      fill=\"#202020\" />\n\n                <rect id=\"bar-mem-background-free\"\n                      x=\"16\" y=\"9\"\n                      width=\"19\" height=\"33\"\n                      fill=\"darkgray\" />\n                <rect id=\"bar-mem-background-used\"\n                      x=\"16\" y=\"43\"\n                      width=\"19\" height=\"47\"\n                      fill=\"orange\" />\n\n                <use x=\"0\" y=\"5\" xlink:href=\"#mem-foot\" />\n                <use x=\"0\" y=\"25\" xlink:href=\"#mem-foot\" />\n                <use x=\"0\" y=\"45\" xlink:href=\"#mem-foot\" />\n                <use x=\"0\" y=\"65\" xlink:href=\"#mem-foot\" />\n                <use x=\"0\" y=\"85\" xlink:href=\"#mem-foot\" />\n\n                <use x=\"-51\" y=\"5\" xlink:href=\"#mem-foot\" transform=\"scale(-1,1)\" />\n                <use x=\"-51\" y=\"25\" xlink:href=\"#mem-foot\" transform=\"scale(-1,1)\" />\n                <use x=\"-51\" y=\"45\" xlink:href=\"#mem-foot\" transform=\"scale(-1,1)\" />\n                <use x=\"-51\" y=\"65\" xlink:href=\"#mem-foot\" transform=\"scale(-1,1)\" />\n                <use x=\"-51\" y=\"85\" xlink:href=\"#mem-foot\" transform=\"scale(-1,1)\" />\n            </svg>\n            <div class=\"status-text\">386/927 MB free</div>\n        </div>\n\n        <div id=\"wrapper-hdd\" class=\"status-symbol\">\n            <svg>\n                <rect x=\"0\" y=\"0\" rx=\"5\" ry=\"5\" width=\"50\" height=\"80\" fill=\"#404040\" />\n                <path d=\"M 10 78 Q 10 95 30 90\" stroke=\"#404040\" stroke-width=\"3\" fill=\"transparent\" />\n                <rect x=\"25\" y=\"84\" rx=\"3\" ry=\"3\" width=\"15\" height=\"12\" fill=\"#404040\" />\n                <rect x=\"36\" y=\"85\" width=\"12\" height=\"10\" fill=\"grey\" stroke=\"#404040\" />\n                <rect x=\"31\" y=\"83\" width=\"10\" height=\"14\" fill=\"#404040\" />\n                <rect x=\"43\" y=\"88\" width=\"3\"  height=\"1\"  fill=\"#404040\" />\n                <rect x=\"43\" y=\"91\" width=\"3\"  height=\"1\"  fill=\"#404040\" />\n                <circle r=\"14\" cx=\"25\" cy=\"30\" fill=\"#404040\" stroke=\"dimgray\" stroke-width=\"15\" />\n                <circle r=\"14\" cx=\"25\" cy=\"30\" fill=\"transparent\" stroke=\"darkorange\" stroke-width=\"15\" stroke-dasharray=\"35 96\" />\n            </svg>\n            <div class=\"status-text\">602/990 GB free</div>\n        </div>\n\n        <div id=\"wrapper-load\" class=\"status-symbol\">\n            <svg>\n                <circle r=\"25\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"#404040\" stroke-width=\"4\" />\n                <circle r=\"25\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"red\"\n                        stroke-width=\"2\" stroke-dasharray=\"55 250\" transform=\"rotate(273 35 30)\" />\n                <circle r=\"20\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"lightgrey\"\n                        stroke-width=\"3\" stroke-dasharray=\"1 10\" transform=\"rotate(50 35 30)\" />\n                <circle r=\"22\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"black\"\n                        stroke-width=\"28\" stroke-dasharray=\"45 200\" transform=\"rotate(32 35 30)\" />\n                <circle r=\"22\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"black\"\n                        stroke-width=\"28\" stroke-dasharray=\"45 200\" transform=\"rotate(32 35 30)\" />\n                <circle r=\"14\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"#505050\"\n                        stroke-width=\"17\" stroke-dasharray=\"10 500\" transform=\"rotate(200 35 30)\" />\n                <circle r=\"7\" cx=\"35\" cy=\"30\" fill=\"transparent\" stroke=\"#404040\"\n                        stroke-width=\"17\" stroke-dasharray=\"7 500\" transform=\"rotate(190 35 30)\" />\n                <polygon points=\"32,30 35,5 38,30\" fill=\"orange\" transform=\"rotate(-25 35 30)\"/>\n                <circle r=\"3\" cx=\"35\" cy=\"30\" fill=\"white\" stroke=\"orange\" stroke-width=\"4\" />\n            </svg>\n            <pre class=\"status-text\">CPU load\n  1 min∅ 2.7\n  5 min∅ 1.6\n 15 min∅ 1.9</pre>\n        </div>\n    </div>\n\n    <div class=\"status-box\">\n        <p>LND 0.5.2-beta wallet 212784 sat</p>\n        <p>6/7 Channels 359780 sat 12 peers</p>\n        <p>Idea: lines for each channel like Shango</p>\n        <p>colors for channel states: pending, open, active, closing</p>\n    </div>\n    <div class=\"status-box\">\n        <p>bitcoin v0.17.0.1 mainnet Sync OK (100%)</p>\n        <p>▼2.8GiB ▲1002.8MiB</p>\n        <p>37 connections</p>\n\n    </div>\n    <div class=\"status-box\">\n        <p>ssh admin@192.168.178.42</p>\n        <p>web admin http://192.168.178.42:3000</p>\n        <p>DynDNS myblitz.ignorelist.com</p>\n    </div>\n</div>\n\n\n</body>\n\n</html>\n"),
	}

	// define dirs
	dir3 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1552941957, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file4, // "info.tmpl"

		},
	}

	// link ChildDirs
	dir3.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1552941957, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir3,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"info.tmpl": file4,
		},
	})
}
