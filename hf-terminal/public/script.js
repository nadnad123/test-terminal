document.getElementById('input').addEventListener('keydown', async function(event) {
    // Check for Enter key without Shift
    if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault(); // Prevent newline in textarea

        const input = this.value.trim();
        if (input === '') {
            return; // Ignore empty input
        }

        // Display user input in output
        const output = document.getElementById('output');
        output.innerHTML += `<div><span class="prompt">$</span> ${input}</div>`;

        // Clear textarea
        this.value = '';

        // Disable textarea during request to prevent multiple submissions
        this.disabled = true;

        try {
            const response = await fetch('/api/completion', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    messages: [{ role: 'user', content: input }],
                }),
            });

            const data = await response.json();
            if (data.response) {
                output.innerHTML += `<div>${data.response}</div>`;
            } else {
                output.innerHTML += `<div style="color: #f00;">Error: ${data.error || 'Unknown error'}</div>`;
            }
        } catch (error) {
            console.error('Fetch error:', error); // Log error for debugging
            output.innerHTML += `<div style="color: #f00;">Error: Failed to connect to server</div>`;
        } finally {
            // Re-enable textarea
            this.disabled = false;
            this.focus(); // Restore focus to textarea
        }

        // Scroll to bottom of output
        output.scrollTop = output.scrollHeight;
    }
});

// Ensure textarea is focused on page load
document.getElementById('input').focus();
