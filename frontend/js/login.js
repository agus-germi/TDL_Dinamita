document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Evita el envío del formulario

    const email = document.getElementById('loginEmail').value; // Cambiar loginUsername a loginEmail
    const password = document.getElementById('loginPassword').value;

    if (!email) {
        document.getElementById('loginMessage').innerText = 'El email es obligatorio.';
        return;
    }
    
    if (!password) {
        document.getElementById('loginMessage').innerText = 'La contraseña es obligatoria.';
        return;
    }

    fetch('http://localhost:8080/users/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }), 
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Error en el inicio de sesión');
        }
        return response.json(); // Asegúrate de que el backend responde con JSON
    })
    .then(data => {
        document.getElementById('loginMessage').innerText = "Inicio de sesión exitoso";
        localStorage.setItem('username', data.name); // Guardar nombre de usuario en localStorage
        console.log('Usuario en storage: ' + localStorage.getItem('username'));
        window.location.href = '/static/html/main.html'; // Redirigir a la página principal
    })
    .catch(error => {
        document.getElementById('loginMessage').innerText = error.message;
    });
});
