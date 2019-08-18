var messages = require('./hello_pb');
var services = require('./hello_grpc_pb');

const grpc = require('grpc');

function subscribe() {
    return new Promise((resolve, reject) => {
        let client = new services.GreeterClient('localhost:50051', grpc.credentials.createInsecure());
        let request = new messages.HelloRequest();
        request.setName(process.argv[0] || "world");
        var call = client.sayHello(request);
        call.on('data', function(msg) {
            console.log('data', msg);
        });
        call.on('end', function() {
            // The server has finished sending
            console.log('end');
            resolve();
        });
        call.on('error', function(e) {
            // An error has occurred and the stream has been closed.
            console.error('error', e);
            reject(e);
        });
        call.on('status', function(status) {
            // process status
            console.log('status', status);
        });
    });
}

async function main() {
    while (true) {
        try {
            await subscribe();
        } catch (e) {

        }
    }
}

main();
