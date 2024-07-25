document.querySelector('#search-btn').onclick = () => {
    let searchForm = document.querySelector('.search-artists'); // Corrected variable name
    searchForm.classList.toggle('active'); // Corrected classList method
};

let text = document.getElementById('text');
window.addEventListener('scroll', () => {
    let value = window.scrollY;
text.style.marginTop = value * 2.5 + 'px';
console.log(value)
if (value >= 250) {
    text.style.display = "none" 
} else {
    text.style.display = "block"
}
});

let loadMoreBtn = document.querySelector('#load-more');
let currentItem = 24;
loadMoreBtn.onclick = () => {
    let boxes = [...document.querySelectorAll('.cards .card')];
    for (let i = currentItem; i < currentItem + 24 && i < boxes.length; i++){
        boxes[i].style.display = 'inline-block'
    }
    currentItem += 24;
};
