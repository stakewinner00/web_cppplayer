function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
        for(var i=0;i < ca.length;i++) {
            var c = ca[i];
            while (c.charAt(0)==' ') c = c.substring(1,c.length);
                if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
        }
    return null;
}

function ChangeTheme(theme) {
    var d = new Date();
    d.setTime(d.getTime() + (365*24*60*60*1000));

    document.cookie = "tema=" + theme + "; expires=" + d.toUTCString();

    if (theme === 'dark') {
        $('.blue')
            .removeClass('blue')
            .addClass('teal');

        $('body')
            .addClass('blue-grey')
            .addClass('darken-4');

        $('.card')
            .addClass('blue-grey')
            .addClass('lighten-1');
    } else {
        $('.teal')
            .removeClass('teal')
            .addClass('blue');

        $('body')
            .removeClass('blue-grey')
            .removeClass('darken-4');

        $('.card')
            .removeClass('blue-grey')
            .removeClass('lighten-1');
    }
}

$(document).ready(function() {
	$('#btn-play-pause').click(function(e) {
		e.preventDefault();

		if ($($('#btn-play-pause').find('.material-icons')).text() == 'pause') {
			$.ajax({
				url: '/pause',
				method: 'GET',
				async: true,
				dataType: 'text',
				success: function(data) {
					$($('#btn-play-pause').find('.material-icons')).text('play_arrow');
				}
			});
		} else {
			$.ajax({
				url: '/pause',
				method: 'GET',
				async: true,
				dataType: 'text',
				success: function(data) {
					$($('#btn-play-pause').find('.material-icons')).text('pause');
				}
			});            
		}
	});

	$('#btn-next').click(function(e) {
		e.preventDefault();
		$.ajax({
			url: '/next',
			method: 'GET',
			async: true,
			dataType: 'text',
			success: function(data) {
				// Función para obtener datos de la canción y setearlos en la página
			}
		});
	});

	$('#btn-prev').click(function(e) {
		e.preventDefault();
		$.ajax({
			url: '/prev',
			method: 'GET',
			async: true,
			dataType: 'text',
			success: function(data) {
				// Función para obtener datos de la canción y setearlos en la página
			}
		});
	});

	$('#btn-stop').click(function(e) {
		e.preventDefault();
		$.ajax({
			url: '/stop',
			method: 'GET',
			async: true,
			dataType: 'text',
			success: function(data) {}
		});
	});

	$('#valor-volumen').change(function() {
		var vol = $(this).val();

		$.ajax({
			url: '/setvolume/' + vol,
			method: 'GET',
			async: false,
			dataType: 'text',
			success: function(data) {}
		});
	});

    $(".button-collapse").sideNav();
    
    
    /** Datos de precarga **/
    var tema = 'light';
    if (getCookie('tema') != null) {
        tema = getCookie('tema');
    }

    ChangeTheme(tema);

    $('.change-theme').click(function(e){
        e.preventDefault();
        ChangeTheme($(this).attr('data-theme'));
    });

    // Prevcarga Volumen
    $.ajax({
        url: '/setvolume/' + vol,
        method: 'GET',
        async: false,
        dataType: 'text',
        success: function(data) {
            var vol = 100;
            if (data != null) {
                vol = data;
            }
            
            $('#valor-volumen').val(vol);
        }
    });
});
