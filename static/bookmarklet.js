const arefollower_bookmarklet = type => {
    let api = "http://are.takashi_trap.trap.show/api/";
    switch (type) {
        case 1:
            api += "24ago.jsonp";
            break;
        case 2:
            api += "yesterday.jsonp?hour=6";
            break;
        default:
            api += "yesterday.jsonp";
            break;
    }
    $.ajax({
        type: "GET",
        url: api,
        dataType: "jsonp",
        jsonpCallback: "callback",
        success: json => json.forEach(v => $(".contentBody li.item[data-id=" + v + "]").remove())
    });
};