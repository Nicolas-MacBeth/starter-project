function apiCall(){
    // Value of ingredients, put into an array, with all space around each item trimmed
    const ingredients = document.getElementById('ingredients').value.split(',').map(item => item.toLowerCase().trim());

    // Alert the user that there isn't even 1 correct input
    if (ingredients.every(item => item === '')){
        alert('You must input something other than empty spaces');
        return;
    }

    // HTTP POST setup
    var http = new XMLHttpRequest();
    var url = '/findfood';
    var params = JSON.stringify({
        IngredientsList: ingredients,
    });

    http.open('POST', url, true);
    http.setRequestHeader('Content-type', 'application/json');

    // Function called when state changes, with http codes
    http.onreadystatechange = function() {
        if(http.readyState === 4 && http.status === 200) {
            const payload = JSON.parse(http.responseText);
            populateView(payload);
        } 
        else if (http.readyState === 4 && http.status === 404){
            notFoundView();
        }
        else if (http.readyState === 4 && http.status !== 200){
            console.error(http.responseText)
        }
    }

    http.send(params);

    // Reset the view for new results
    document.getElementById('results').innerHTML = '';
}

// Handle 'enter' keypress
function enterKey(event){
    if (event.keyCode === 13){
        event.preventDefault();
        document.getElementById('submit').click();
    }    
}

function populateView(payload) {
    if (!payload.ListOfResults){
        return;
    }

    let previous = null;
    const parentNode = document.getElementById('results');

    // Populate results on the HTML
    payload.ListOfResults.forEach(element => {
        // Add a header if it's a "new" ingredient
        if(element.Ingredient !== previous){
            previous = element.Ingredient;
            let newHeader = document.createElement('H2');
            newHeader.innerHTML = element.Ingredient;
            parentNode.appendChild(newHeader);
        }

        // add the result
        let newResult = document.createTextNode(`Vendor ID: ${element.VendorID}, Price ${element.Price}$/kg, Inventory: ${element.Inventory} units, Vendor name: ${element.VendorName}`);
        parentNode.appendChild(newResult);
        parentNode.appendChild(document.createElement('br'));
    });
}

function notFoundView(){
    const parentNode = document.getElementById('results');
    let errorMessageNode = document.createElement('H2');
    errorMessageNode.innerHTML = "No vendors found for those ingredients";
    errorMessageNode.style.color = 'red';
    parentNode.appendChild(errorMessageNode);
}