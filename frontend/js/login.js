document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Evita el envío del formulario

    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Error en el inicio de sesión');
        }
        return response.text(); // Cambiar a JSON si el servidor responde con un objeto JSON
    })
    .then(data => {
        document.getElementById('loginMessage').innerText = data.message; // Suponiendo que el mensaje está en el objeto data
        localStorage.setItem('username', username); // Guardar en localStorage
        console.log('storage' + localStorage.getItem('username'));
        window.location.href = 'main.html'; // Redirigir a la página principal
    })
    .catch(error => {
        document.getElementById('loginMessage').innerText = error.message;
    });
});
