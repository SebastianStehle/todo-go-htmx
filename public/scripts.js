htmx.defineExtension('bs-validation', {
    onEvent: (name, evt) => {
        if (name !== "htmx:afterProcessNode") {
            return;
        }

        let form = evt.detail.elt;
        // Check if trigger attribute and submit event exists for the form
        if (form.hasAttribute('hx-trigger')) {
            return;
        }

        // Set trigger for custom event bs-send
        form.setAttribute('hx-trigger', 'bs-send');
        // And attach the event only once
        form.addEventListener('submit', (event) => {
            if (form.checkValidity()) {
                // Trigger custom event hx-trigger="bs-send"
                htmx.trigger(form, "bsSend");
            }

            // Focus the first invalid field
            let invalidField = form.querySelector(':invalid');
            if (invalidField) {
                invalidField.focus();
            }

            event.preventDefault();
            event.stopPropagation();

            form.classList.add('was-validated');
        }, false);
    }
});