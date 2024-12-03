document.addEventListener('DOMContentLoaded', function() {
    // Set default date to today in the workout form
    const dateInput = document.getElementById('date');
    if (dateInput) {
        dateInput.valueAsDate = new Date();
    }

    // Handle workout completion toggle
    const toggleForms = document.querySelectorAll('.toggle-form');
    toggleForms.forEach(form => {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            fetch(form.action, {
                method: 'POST',
                body: new FormData(form)
            }).then(() => {
                window.location.reload();
            }).catch(error => {
                console.error('Error:', error);
            });
        });
    });
}); 