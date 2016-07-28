/**
 * Created by yanjunhui on 16/7/18.
 */

$(function () {
    var conn;
    var host = $("#hostAddress").html();
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }
    
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + host + "/ws");
        conn.onclose = function (evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        };
        conn.onmessage = function (evt) {
            appendLog($("<pre/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
});