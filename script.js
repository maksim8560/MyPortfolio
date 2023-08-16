let cursor = document.querySelector('#cursor');
document.addEventListener('mousemove', (e) => {
    cursor.style.top = e.pageY + 'px';
    cursor.style.left = e.pageX + 'px';
})

let links = document.querySelectorAll('ul li a');
links.forEach(link => {
    let letters = link.innerText.split('');
    let newHTML = letters.map((letter, i) => `<span style="transition-delay:${i * 40}ms;">${letter}</span>`).join('');
    link.innerHTML = newHTML;
});
