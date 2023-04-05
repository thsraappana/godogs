const fetchBreedsButton = document.getElementById('fetch-breeds-list');
const fetchRandomImageButton = document.getElementById('fetch-random-image');

document.addEventListener('click', async (event) => {
    if (event.target.id === fetchRandomImageButton.id) {
        const response = await fetch('/handle-random-image-button-click', { method: 'POST' });
        if (response.ok) {
            refreshPageContent(response)
        } else {
            // TODO: Handle errors
            console.log('something went wrong!')
        }
    }
    if (event.target.id === fetchBreedsButton.id) {
        const response = await fetch('/handle-fetch-breeds-button-click', { method: 'POST' });
        if (response.ok) {
            refreshPageContent(response)
        } else {
            console.log('something went wrong!')
        }
    }
});

async function refreshPageContent(response) {
    const pageContent = await response.text();
    const parser = new DOMParser();
    const newDocument = parser.parseFromString(pageContent, 'text/html');
    document.documentElement.innerHTML = newDocument.documentElement.innerHTML;
}


