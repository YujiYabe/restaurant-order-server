async function showErrorModal(error) {
    printDebug('error', error); // -----------------------------

	const modal = new bootstrap.Modal(document.getElementById('modal'));
	$('#modal_header_text').addClass('text-center').text('注文できませんでした');

	$('#modal_body').addClass('text-center h3').text('最新データを取得します');

	$('#cancel_button').hide();
	$('#modal_footer_product').hide();
	$('#confirm_button_group').hide();
	$('#modal_footer_normal').show();

	$('#confirm_button').off('click', async () => {
		location.reload();
	});

	$('#confirm_button').on('click', async () => {
		location.reload();
	});

	modal.show();

}
