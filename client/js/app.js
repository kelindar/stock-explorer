var emitter = emitter.connect({
    secure: true
}); 
var resKey = 'Ywim6VqDa_jJijb29_14RZmFnEUgpU5q';
var reqKey = 'ifHzn-TvyAom7vDMXcC0Er3fDQu36PZP';
var vue = new Vue({
    el: '#app',
    data: {
        symbol: 'TSLA',
        result: { }
    },
    methods: {
        query: function () {
	        // publish a message to the chat channel
	        console.log('emitter: publishing ');
	        emitter.publish({
                key: reqKey,
                channel: "quote-request",
                message: JSON.stringify({
                    symbol: this.$data.symbol, 
                    reply: getPersistentVisitorId()
                })
            });

            // remove the text
            this.$data.symbol = '';
        },
    }
});

emitter.on('connect', function(){
    // once we're connected, subscribe to the 'chat' channel
    console.log('emitter: connected');
    emitter.subscribe({
        key: resKey,
        channel: "quote-response/" + getPersistentVisitorId()
    });

})

// on every message, print it out
emitter.on('message', function(msg){

    // log that we've received a message
    console.log('emitter: received ' + msg.asString() );

});