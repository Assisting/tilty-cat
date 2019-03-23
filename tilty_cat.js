const express = require('express');
const app = express();
const fs = require('fs');
const gm = require('gm');

const port = 3000;

var filename = '/cat.png';

var getCrocMode = function (req, res) {
	if(filename == '/cat.png') {
		res.send('false');
	} else {
		res.send('true');
	}
}


var setCrocMode = function (req, res) {
	if(Number.parseInt(req.params.bool) == 1) {
		filename = '/croc.png';
		res.status(200).end();
	} else {
		filename = '/cat.png';
		res.status(200).end();
	}
}

var sendIndex = function (req, res) {
	res.sendFile(__dirname+'/index.html')
}

var dogReq = function (req, res) {
	filename = '/dog.png';

	rotate(req, res);
}

var catReq = function (req, res) {
	filename = '/cat.png';

	rotate(req, res);
}

var rotate = function (req, res) {
	var degrees = Number.parseInt(req.params.degree);

	if(Number.isSafeInteger(degrees)) {
		degrees = degrees % 360;
	} else {
		degrees = 0;
	}

	res.header('Cache-Control', 'private, no-cache, no-store, must-revalidate');
	res.header('Content-Type', 'image/png');
	gm(__dirname+filename)
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

app.get('/tilty_cat/:degree.png', catReq);

app.get('/tilty_cat/croc_mode/', getCrocMode);

app.get('/tilty_cat/croc_mode/index.html', sendIndex);

app.post('/tilty_cat/croc_mode/:bool', setCrocMode);

// Is this scope creep
app.get('/tilty_dog/:degree.png', dogReq);

app.listen(port, () => console.log(`Now listening on port ${port}`))
