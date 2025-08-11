// Handle button click
document.getElementById('send-button').addEventListener('click', async function() {
    const inputElement = document.getElementById('input');
    const input = inputElement.value.trim();
    if (input === '') {
        return; // Ignore empty input
    }

    // Display user input in output
    const output = document.getElementById('output');
    output.innerHTML += `<div><span class="prompt">$</span> ${input}</div>`;

    // Clear textarea
    inputElement.value = '';

    // Disable button and textarea during request
    this.disabled = true;
    inputElement.disabled = true;

    try {
        console.log('Sending request to /api/completion with:', input); // Debug log
        const response = await fetch('/api/completion', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                messages: [{ role: 'user', content: input }],
            }),
        });

        console.log('Response status:', response.status); // Debug log
        const data = await response.json();
        console.log('API Response:', data); // Debug log

        if (data.response) {
            output.innerHTML += `<div>${data.response}</div>`;
        } else {
            output.innerHTML += `<div style="color: #f00;">Error: ${data.error || 'Unknown error'}</div>`;
        }
    } catch (error) {
        console.error('Fetch error:', error); // Debug log
        output.innerHTML += `<div style="color: #f00;">Error: Failed to connect to server - ${error.message}</div>`;
    } finally {
        // Re-enable button and textarea
        this.disabled = false;
        inputElement.disabled = false;
        inputElement.focus(); // Restore focus to textarea
    }

    // Scroll to bottom of output
    output.scrollTop = output.scrollHeight;
});

// Ensure textarea is focused on page load
document.getElementById('input').focus();
