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
				url: '/play',
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
			url: '/volumen/' + vol,
			method: 'GET',
			async: false,
			dataType: 'text',
			success: function(data) {}
		});
	});
});