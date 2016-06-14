$(document).ready(function() {
    $('#btn-play-pause').click(function() {
        if ($($('#btn-play-pause').find('.material-icons')).text() == 'pause') {
            // Envía ajax de Pausar /pause
            $($('#btn-play-pause').find('.material-icons')).text('play_arrow')
        } else {
            // Envía ajax de Play /play
            $($('#btn-play-pause').find('.material-icons')).text('pause')
        }
    });    
});