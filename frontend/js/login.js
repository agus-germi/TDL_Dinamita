document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    if (!email) {
        document.getElementById('loginMessage').innerText = 'El email es obligatorio.';
        return;
    }

    if (!password) {
        document.getElementById('loginMessage').innerText = 'La contraseña es obligatoria.';
        return;
    }

    fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => {
                throw new Error(err.message || 'Error en el inicio de sesión');
            });
        }
        return response.json();
    })
    .then(data => {
        localStorage.setItem('token', data.token);

        const decodedToken = jwt_decode(data.token);
        const userId = decodedToken.user_id || decodedToken.id;
        const name = decodedToken.name;
        console.log('ID del usuario: ' + userId);

        localStorage.setItem('userId', userId);
        localStorage.setItem('userName', name);

        document.getElementById('loginMessage').innerText = "Inicio de sesión exitoso";
        window.location.href = '/static/html/main.html';
    })
    .catch(error => {
        document.getElementById('loginMessage').innerText = error.message;
    });
});
