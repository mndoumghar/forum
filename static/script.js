document.addEventListener('DOMContentLoaded', function() {
    var createBtn = document.getElementById('createPostBtn');
    var modal = document.getElementById('createPostModal');
    var closeModalBtn = document.getElementById('closeModalBtn');

    // Open modal on create button click
    if (createBtn && modal) {
        createBtn.addEventListener('click', function() {
            modal.style.display = 'flex';
        });
    }

    // Close modal on close button click
    if (closeModalBtn && modal) {
        closeModalBtn.addEventListener('click', function() {
            modal.style.display = 'none';
        });
    }

    // Close modal if user clicks outside the modal content
    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Function to style buttons consistently
    function styleButton(button) {
        if (!button) return;
        button.style.background = 'rgba(255, 255, 255, 0.18)';
        button.style.color = '#7357f0';
        button.style.backdropFilter = 'blur(8px)';
        button.style.boxShadow = '0 4px 16px rgba(115,87,240,0.18)';
        button.style.border = '1px solid rgba(115,87,240,0.25)';
        button.style.borderRadius = '8px';
        button.style.cursor = 'pointer';  // Add pointer cursor for buttons/links
        button.style.transition = 'background 0.3s ease'; // Smooth hover effect if desired
    }

    // Style the logout, login, and create post buttons
    styleButton(document.querySelector('a[href="/logout"]'));
    styleButton(document.querySelector('a[href="/login"]'));
    styleButton(createBtn);
});





document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('.toggle-commenters-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const list = btn.nextElementSibling;
      if (!list) return;
      list.classList.toggle('hidden');
      btn.textContent = list.classList.contains('hidden') ? 'Show Commenters' : 'Hide Commenters';
    });
  });
});

document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('.commenters-list li').forEach(li => {
    li.textContent = li.textContent.replace(/[{}]/g, '');
  });
});