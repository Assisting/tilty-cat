const express = require('express');
const app = express();
const fs = require('fs');
const gm = require('gm');

const port = 80;

var rotate = function (req, res) {
	var degrees = Number.parseInt(req.params.degree);

	if(Number.isSafeInteger(degrees)) {
		degrees = degrees % 360;
	} else {
		degrees = 0;
	}

	res.header('Cache-Control', 'private, no-cache, no-store, must-revalidate');
	res.header('Content-Type', 'image/png');
	gm(__dirname+'/cat.png')
	.rotate('transparent', degrees)
	.resize(128, 128)
	.trim()
	.resize(128, 128)
	.stream(function streamOut (err, stdout, stderr) {
            if (err) return console.dir(arguments)
			console.log("Created: " + arguments[3])
			stdout.pipe(res);
		}
	); // End stream()
}

app.get('/tilty_cat/:degree.png', rotate);

app.listen(port, () => console.log(`Now listening on port ${port}`))
