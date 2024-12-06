document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    const name = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;
    const email = document.getElementById('registerEmail').value;
    const registerMessage = document.getElementById('registerMessage');

    if (!name) {
        registerMessage.innerText = 'El nombre es obligatorio.';
        return;
    }

    if (!password) {
        registerMessage.innerText = 'La contraseña es obligatoria.';
        return;
    } else {
        if (password.length < 8 || password.length > 15) {
            registerMessage.innerText = 'La contraseña debe tener entre 8 y 15 caracteres.';
            return;
        }
    }

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email)) {
        registerMessage.innerText = 'El correo electrónico no tiene un formato válido.';
        return;
    }

    fetch('http://localhost:8080/api/v1/auth/signup', { 
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, password, email }), 
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Error en el registro');
        }
        return response.json();
    })
    .then(data => {
        document.getElementById('registerMessage').innerText = data.message; 
        localStorage.setItem('username', name); 
        window.location.href = '/static/index.html'; 
    })
    .catch(error => {
        document.getElementById('registerMessage').innerText = error.message;
    });
});
