var emitter = emitter.connect({
    secure: true
}); 
var resKey = 'Ywim6VqDa_jJijb29_14RZmFnEUgpU5q';
var reqKey = 'ifHzn-TvyAom7vDMXcC0Er3fDQu36PZP';
var vue = new Vue({
    el: '#app',
    data: {
        symbol: 'AAPL',
        result: new Object()
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
    },

    filters: {
        currencyDisplay: {
            read: function(val) {
                if(val === undefined) return ''
                return '$' + val.toFixed(2);
            }
        },

        percentageDisplay: {
            read: function(val) {
                if(val === undefined) return ''
                return val.toFixed(2) + '%';
            }
        },

        millionsDisplay: {
            read: function(val) {
                if(val === undefined) return ''
                return '$' + val.toFixed(0) + 'M';
            }
        },

        dateDisplay: {
            read: function(val) {
                if(val === undefined) return ''
                return formatDate(val);
            }
        },
    },  
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
    var data = msg.asObject();
    console.log('emitter: received ', data);

    // sort financials 
    if (data.Financials){
        data.Financials = data.Financials.sort(byDate)
    }

    // sort financials 
    if (data.DividendHistory != null){
        data.DividendHistory = data.DividendHistory.sort(byDate)
    }

    // bind the result to the screen
    vue.$data.result = data;

    // make sure to empty the div
    if (data.DividendHistory === null){
        document.getElementById('dividends-chart').innerHTML = "";
        return;
    }



    labels = []
    series = []
    data.DividendHistory.forEach(function(d){
        labels.push(formatDate(d.Date))
        series.push(d.Value)
    });

    // apply the chart
    new Chartist.Line('#dividends-chart', {
        labels: labels,
        series: [series]
        }, {
            fullWidth: true,
            axisX: {
                showGrid: false,
                labelInterpolationFnc: function(value, index) {
                    return index % 2 === 0 ? value : null;
                }
            }
    });
});

function formatDate(d){
    return d.substring(0, d.indexOf('-', 5))
}

function byDate(a, b) {
    d1 = new Date(formatDate(a.Date));
    d2 = new Date(formatDate(b.Date));
    return d1 - d2;
}