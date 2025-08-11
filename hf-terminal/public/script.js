document.getElementById('input').addEventListener('keydown', async function(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        const input = this.value.trim();
        if (input === '') return;

        // Display user input
        const output = document.getElementById('output');
        output.innerHTML += `<div><span class="prompt">$</span> ${input}</div>`;

        // Clear input
        this.value = '';

        // Send request to backend
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
                output.innerHTML += `<div style="color: #f00;">Error: ${data.error}</div>`;
            }
        } catch (error) {
            output.innerHTML += `<div style="color: #f00;">Error: Failed to connect to server</div>`;
        }

        // Scroll to bottom
        output.scrollTop = output.scrollHeight;
    }
});