var https = require('https');
var fs = require('fs');

var options = {
    hostname:'192.168.101.220',
    port:443,
    path:'/',
    method:'GET',
    pfx:fs.readFileSync('./client.pfx'),
    passphrase:'Aa123456',
    agent:false
};

options.agent = new https.Agent(options);
var req = https.request(options,function(res){
console.log("statusCode: ", res.statusCode);
  console.log("headers: ", res.headers);
    res.setEncoding('utf-8');
    res.on('data',function(d){
        console.log(d);
    })
});

req.end();

req.on('error',function(e){
    console.log(e);
});
