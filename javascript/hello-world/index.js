'use strict';

const KeywordExtraction = require("./keyword-extraction.js"),
    express = require('express'),
    multer = require('multer'),
    fetch = require('node-fetch-commonjs'),
    upload = multer({storage: multer.memoryStorage()}),
    app = express();

let chunkUpload = upload.single('chunk');
let server = app.listen( 8080 );
server.setTimeout( 10 * 60 * 1000 );

// READY WEBHOOK
app.get('/ready', (req, res) => {
    res.status(200).send('OK');
});

// PROCESS WEBHOOK
app.post('/process', chunkUpload, async (req, res)=>{
    try {
        const cacheUri = req.body.cacheURI;
        console.log(`Reading from ${cacheUri}`);

        const response = await fetch(cacheUri);
        const buffer = await response.text();

        let output = KeywordExtraction.getOutput( buffer, null );
        return res.status(200).send( output );
    } catch (error) {
        console.log(error);
        return res.status(500).send(error);
    }
});
