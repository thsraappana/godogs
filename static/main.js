document.addEventListener('click', async (event) => {
    const fetchBreedsButton = document.getElementById('fetch-breeds-list');
    const fetchRandomImageButton = document.getElementById('fetch-random-image');
    const fetchDogBreedButton = document.getElementById('fetch-dog-breed-image');

    if (event.target.id === fetchRandomImageButton.id) {
        const response = await fetch('/fetch-random-image', { method: 'GET' });
        if (response.ok) {
            refreshPageContent(response)
        } else {
            // TODO: Handle errors
            console.log('something went wrong!')
        }
    }
    if (event.target.id === fetchBreedsButton.id) {
        const response = await fetch('/fetch-breeds', { method: 'GET' });
        if (response.ok) {
            refreshPageContent(response)
        } else {
            console.log('something went wrong!')
        }
    }
    if (event.target.id === fetchDogBreedButton.id) {
        const dogBreedTextInput = document.getElementById('dog-breed')

        if (!dogBreedTextInput.value.length) {
            alert("No dog breed given!")
            return;
        }
        const response = await fetch('/fetch-breed-image', { method: 'POST', body: JSON.stringify({ breed: dogBreedTextInput.value }) });
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


