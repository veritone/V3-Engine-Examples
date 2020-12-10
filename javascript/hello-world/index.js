'use strict';

const KeywordExtraction = require("./keyword-extraction.js"),
    express = require('express'),
    multer = require('multer'),
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
        let buffer = req.file.buffer.toString();
        let output = KeywordExtraction.getOutput( buffer, null );
        return res.status(200).send( output );
    } catch (error) {
        return res.status(500).send(error);
    }
});
