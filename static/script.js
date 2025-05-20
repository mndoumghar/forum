document.addEventListener('DOMContentLoaded', function() {
    var createBtn = document.getElementById('createPostBtn');
    var modal = document.getElementById('createPostModal');
    var closeModalBtn = document.getElementById('closeModalBtn');

    if (createBtn && modal) {
        createBtn.onclick = function() {
            modal.style.display = 'flex';
        };
    }
    if (closeModalBtn && modal) {
        closeModalBtn.onclick = function() {
            modal.style.display = 'none';
        };
    }
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    };

    // Style the logout button
    var logoutBtn = document.querySelector('a[href="/logout"]');
    if (logoutBtn) {
        logoutBtn.style.background = 'rgba(255, 255, 255, 0.18)';
        logoutBtn.style.color = '#7357f0';
        logoutBtn.style.backdropFilter = 'blur(8px)';
        logoutBtn.style.boxShadow = '0 4px 16px rgba(115,87,240,0.18)';
        logoutBtn.style.border = '1px solid rgba(115,87,240,0.25)';
        logoutBtn.style.borderRadius = '8px';
    }

    // Style the login button
    var loginBtn = document.querySelector('a[href="/login"]');
    if (loginBtn) {
        loginBtn.style.background = 'rgba(255, 255, 255, 0.18)';
        loginBtn.style.color = '#7357f0';
        loginBtn.style.backdropFilter = 'blur(8px)';
        loginBtn.style.boxShadow = '0 4px 16px rgba(115,87,240,0.18)';
        loginBtn.style.border = '1px solid rgba(115,87,240,0.25)';
        loginBtn.style.borderRadius = '8px';
    }

    // Change the color of the Create Post button
    if (createBtn) {
        createBtn.style.background = 'rgba(255, 255, 255, 0.18)';
        createBtn.style.color = '#7357f0';
        createBtn.style.backdropFilter = 'blur(8px)';
        createBtn.style.boxShadow = '0 4px 16px rgba(115,87,240,0.18)';
        createBtn.style.border = '1px solid rgba(115,87,240,0.25)';
        createBtn.style.borderRadius = '8px';
    }
});