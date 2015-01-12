/* global $:false, alert */

"use strict"; // jshint ignore:line

function saveDeviceList() {
    var listText = $('#deviceListConfig').val();

    $.post('/api/savedevicelist', {text: encodeURIComponent(listText)}, null, "json")
        .done(function(data) {
            if (!data.success) {
                alert(data.error);
            }
        });
}

(function() {
    $('#saveDeviceListBtn').click(saveDeviceList);
    return;
})();