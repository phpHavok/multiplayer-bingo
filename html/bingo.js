$(function () {
    var URL_PREFIX = location.protocol + '//' + location.hostname +
        (location.port ? ':' + location.port : '');
    var g_username;
    var g_uid;

    var createBingoCard = function (username, board) {
        var $player = $('<div class="bingo-player col-md-6 col-sm-12">');
        var $card = $('<div class="bingo-card">');
        for (var i = 0; i < board.length; i++) {
            var $cell = $('<div class="bingo-cell">');
            $cell.append($('<div>').text(board[i].phrase));
            if (board[i].marked) {
                $cell.addClass('marked');
            }
            $card.append($cell);
        }
        $player.append('<h2>' + username + '</h2>');
        $player.append($card);
        return $player;
    };

    var refreshGame = function () {
        $.get(URL_PREFIX + '/game',
            {
                uid: g_uid
            },
            function (data) {
                if (data.hasOwnProperty('error')) {
                    alert(data.error);
                } else {
                    $('#topic span').text(data.topic);
                    $('#topic').show();
                    var $bingo = $('#bingo');
                    $bingo.empty();
                    // Add me first for consistency.
                    for (var i = 0; i < data.players.length; i++) {
                        var player = data.players[i];
                        if (player.username == g_username) {
                            var $card = createBingoCard(player.username, player.bingo_board);
                            $bingo.append($card);
                        }
                    }
                    // Now add everyone else.
                    for (var i = 0; i < data.players.length; i++) {
                        var player = data.players[i];
                        if (player.username != g_username) {
                            var $card = createBingoCard(player.username, player.bingo_board);
                            $bingo.append($card);
                        }
                    }
                }
            },
            'json'
        );
    };

    var scheduleRefresh = function () {
        refreshGame();
        setInterval(refreshGame, 3000);
    };

    // Join a game.
    $('#joinForm').submit(function (evt) {
        evt.preventDefault();
        username = $('#joinUsername').val();
        room = $('#joinRoom').val();
        $.get(URL_PREFIX + '/join',
            {
                username: username,
                room: room
            },
            function (data) {
                if (data.hasOwnProperty('error')) {
                    alert(data.error);
                } else {
                    g_uid = data.uid;
                    g_username = username;
                    $('#join').hide();
                    $('#uid span').text(g_uid);
                    $('#uid').show();
                    scheduleRefresh();
                }
            },
            'json'
        );
    });

    // Rejoin a game.
    $('#rejoinForm').submit(function (evt) {
        evt.preventDefault();
        g_username = $('#rejoinUsername').val();
        g_uid = $('#rejoinUID').val();
        $('#join').hide();
        $('#uid span').text(g_uid);
        $('#uid').show();
        scheduleRefresh();
    });

    // Bind to first bingo card and allow clicking.
    $('#bingo').delegate('.bingo-player:first-child .bingo-cell', 'click', function (evt) {
        var $cell = $(this);
        var cell = $(this).index();
        var marked = !$(this).hasClass('marked');
        $.get(URL_PREFIX + '/cell',
            {
                uid: g_uid,
                cell: cell,
                marked: marked
            },
            function (data) {
                if (data.hasOwnProperty('error')) {
                    alert(data.error);
                } else {
                    $cell.toggleClass('marked', data.marked);
                }
            },
            'json'
        );
    });
}); 
