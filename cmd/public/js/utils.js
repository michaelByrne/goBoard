Array.prototype.remove = function (from, to) {
    var rest = this.slice((to || from) + 1 || this.length);
    this.length = from < 0 ? this.length + from : from;
    return this.push.apply(this, rest);
};


function firstpost(type, num, ob) {
    var fp = $('#fp_' + num);
    fp.toggle();

    if (!fp.attr('loaded')) {
        fp.html("loading...");
        $.ajax(
            {
                url: '/' + type + '/firstpost/' + num + '/&ajax=true',
                cache: false,
                success: function (html) {
                    fp.attr('loaded', true);
                    fp.html(html);
                }
            });
    }
    if (fp.css('display') == 'block') $(ob).html('&laquo;');
    else
        $(ob).html('&raquo;');
};

function uncollapser(type, media, count) {
    $('#uncollapse_links').hide();
    $('#uncollapse_loading').show();
    media = media ? "&media=true" : "";
    data = $('.post:last')[0].id.split('_');
    id = data[1];

    var start = $('.post:first').next()[0].id.split('_')[3] - 1;
    var num = count;
    if (count !== null) {
        start -= count;
    } else {
        num = start;
        start = 0;
    }

    if (start < 0) {
        num += start;
        start = 0;
    }
    if (start == Math.min(start, num)) {
        $("#uncollapse_all").hide();
        $("#uncollapse_more_counter").html("all " + Math.min(start, num));
        if (start == 1) $("#uncollapse_some").html("show final post");
    } else {
        $("#uncollapse_counter").html(start);
        $("#uncollapse_more_counter").html(Math.min(start, num));
    }
    $.ajax(
        {
            url: '/' + type + '/view/' + id + '/-' + start + '/' + num + '/&ajax=true' + media,
            cache: false,
            success: function (html) {
                $('#uncollapse').after(html);
                if (start <= 0) $('#uncollapse').hide();
                $('#uncollapse_loading').hide();
                $('#uncollapse_links').show();
            }
        });
};

const loadThreadPosts = (type, ob, pid) => {
    if (ob) {
        if (!ob.save) ob.save = ob.innerHTML;
        ob.innerHTML = "loading...";
    }
    let data = $('.post:last')[0].id.split('_');
    let id = data[1];
    let lastPosition = parseInt(data[data.length - 1]);

    fetch('/thread/post/' + pid + '/' + lastPosition, {
        method: 'GET',
        cache: 'no-cache',
    }).then((res) => {
        if (!res.ok) {
            throw new Error(`Network response was not ok ${res.status}`);
        }
        return res.text()
    }).then(html => {
        updateDOMWithNewPost(html, pid, id, lastPosition)
        if (ob) ob.innerHTML = ob.save;
    }).catch((error) => {
        $('#response_form').html('<div class="error">Error: ' + error + '</div>');
    })
};

const loadMessagePosts = (type, ob, pid) => {
    if (ob) {
        if (!ob.save) ob.save = ob.innerHTML;
        ob.innerHTML = "loading...";
    }
    let data = $('.post:last')[0].id.split('_');
    let id = data[1];
    let lastPosition = parseInt(data[data.length - 1]);

    fetch('/message/post/' + pid + '/' + lastPosition, {
        method: 'GET',
        cache: 'no-cache',
    }).then((res) => {
        if (!res.ok) {
            throw new Error(`Network response was not ok ${res.status}`);
        }
        return res.text()
    }).then(html => {
        updateDOMWithNewPost(html, pid,id, lastPosition)
        if (ob) ob.innerHTML = ob.save;
    }).catch((error) => {
        $('#response_form').html('<div class="error">Error: ' + error + '</div>');
    })
};

const updateDOMWithNewPost = (html, pid, id, lastPosition) => {
    $('#view_' + id).append(html);
    $('.post:last').attr('id', 'view_' + id + '_' + pid + '_' + (lastPosition + 1));
    $('.post:last').find('.count').text('#' + (lastPosition + 1));
    $('textarea').val('')
}

const showPreview = (post, type) => {
    let data = $('.post:last')[0].id.split('_');
    let lastPosition = parseInt(data[data.length - 1]);

    fetch('/' + type + '/previewpost/' + lastPosition, {
        method: 'POST',
        body: post
    }).then((res) => {
        if (!res.ok) {
            throw new Error(`Network response was not ok ${res.status}`);
        }
        return res.text()
    }).then(html => {
        $('#response_form').html(html);
    })
}

const captureNewThreadPostSubmit = (event) => {
    event.preventDefault();
    fetch(event.target.action, {
        method: 'POST',
        body: new URLSearchParams(new FormData(event.target)) // event.target is the form
    }).then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json(); // or response.text() or whatever the server sends
    }).then(({post_id}) => {
        $('#response_form').html('');
        loadThreadPosts('thread', false, post_id)
        $('.submit').attr('disabled', false);
    }).catch((error) => {
        // TODO handle error
    });
}

const captureNewMessagePostSubmit = (event) => {
    event.preventDefault();
    fetch(event.target.action, {
        method: 'POST',
        body: new URLSearchParams(new FormData(event.target)) // event.target is the form
    }).then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json(); // or response.text() or whatever the server sends
    }).then(({post_id}) => {
        $('#response_form').html('');
        loadMessagePosts('message', false, post_id)
        $('.submit').attr('disabled', false);
    }).catch((error) => {
        // TODO handle error
    });
}

const captureNewMessageSubmit = (event) => {
    event.preventDefault();
    fetch(event.target.action, {
        method: 'POST',
        body: new URLSearchParams(new FormData(event.target)) // event.target is the form
    }).then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json(); // or response.text() or whatever the server sends
    }).then(({id}) => {
        window.location.href = '/message/list/1';
        $('.submit').attr('disabled', false);
    }).catch((error) => {
        console.log(error)
        // TODO handle error
    });
}

const doCollapse = (maxViewablePosts, minPostsToHide) => {
    let post = $('.post')
    if (post.length > maxViewablePosts) {
        $(post.slice(0, minPostsToHide).hide());
        $('#uncollapse').show();
    }
}

const uncollapse = () => {
    $('.post').show();
    $('#uncollapse').hide();
}

function quote_post(id) {
    const info = jQuery.trim($('#post_info_' + id).text());
    const body = jQuery.trim($('#post_body_' + id).text());
    $('#body').val($('#body').val() + '[quote]' + info + '\n' + body + '[/quote]\n\n');
}

function toggle_ignore_thread(id) {
    var status;
    if ($('#ignorecmd').html() == "ignore") status = "ignoring...";
    if ($('#ignorecmd').html() == "unignore") status = "unignoring...";
    if (!status) return;
    $('#ignorecmd').html(status);

    $.ajax(
        {
            url: '/thread/toggleignore/' + id + '/',
            cache: false,
            success: function (html) {
                $('#ignorecmd').html(jQuery.trim(html));
            }
        });
}

function toggle_favorite(id) {
    var status;
    if ($('#fcmd').html() == "add") status = "adding...";
    if ($('#fcmd').html() == "remove") status = "removing...";
    if (!status) return;
    $('#fcmd').html(status);

    $.ajax(
        {
            url: '/thread/togglefavorite/' + id + '/',
            cache: false,
            success: function (html) {
                $('#fcmd').html(jQuery.trim(html));
            }
        });
}

function undot(id) {
    var status;
    $('#undot').html('undotting...');

    $.ajax(
        {
            url: '/thread/undot/' + id + '/',
            cache: false,
            success: function (html) {
                $('#undot').html(jQuery.trim(html));
            }
        });
}

// Member Add Box
function catch_enter(e) {
    var code = e ? e.keyCode : event.keyCod;
    if (code == 13) {
        check_member();
        return false;
    } else
        return true;
}

function check_member() {
    $('#notice').remove();
    if (!$('#_recipients').val()) return;
    $('#add').val('Adding..');

    members = $('#_recipients').val();
    $('#_recipients').val('');
    $.post("/message/addmember", {names: members}, function (data) {
        // data = jQuery.trim(data);
        // if (data) {
        //     data = data.split(',');
        //     for (i = 0; i < data.length; i = i + 2) add_member(data[i], data[i + 1]);
        // }
        data.members.forEach(({id, name}) => {
            add_member(id, name)
        })
        $('#add').val('Add');
    });
}

function add_member(id, name) {
    var mm = $('#message_members').val();
    if (!mm) mm = [];
    else
        mm = mm.split(',');
    if (jQuery.inArray(id, mm) === -1) {
        members = $('#message_members').val();
        if ($('#m').html() == '-') $('#m').empty();
        $('#m').append('<span id="m' + id + '"><sup><a style="color: white" href="javascript:;" onclick="remove_member(\'' + id + '\')">x</a></sup>&nbsp;' + name + '&nbsp;&nbsp;&nbsp;</span> ');
        $('#message_members').val((members ? members + ',' : '') + id);
    }
}

function remove_member(id) {
    var mm = $('#message_members').val().split(',');
    if (jQuery.inArray(id, mm) != -1) {
        mm.remove(jQuery.inArray(id, mm));
        $('#message_members').val(mm.join(','));
        $('#m' + id).remove();
        if (!mm.length) $('#m').html('-');
    }
}
