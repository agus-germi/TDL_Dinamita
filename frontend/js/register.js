document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Evita el envío del formulario

    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;
    const email = document.getElementById('registerEmail').value;

    fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password, email }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Error en el registro');
        }
        return response.json(); // Cambiar a JSON si el servidor responde con un objeto JSON
    })
    .then(data => {
        document.getElementById('registerMessage').innerText = data.message; // Suponiendo que el mensaje está en el objeto data
        localStorage.setItem('username', username); // Guardar en localStorage
        window.location.href = '../html/main.html'; // Redirigir a la página principal
    })
    .catch(error => {
        document.getElementById('registerMessage').innerText = error.message;
    });
});
