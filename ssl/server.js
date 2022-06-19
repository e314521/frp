var https = require('https');
var fs = require('fs');

var options = {
    key: fs.readFileSync('./server-key.pem'),ca: [fs.readFileSync('./ca-cert.pem')],cert: fs.readFileSync('./server-cert.pem')
};
 

https.createServer(options, function(req,res){
    res.writeHead(200, {
        "content-type":"text/plain"
    });
    res.write("hello nodejs\n");
    res.end();
}).listen(8080);
