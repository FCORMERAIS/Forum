document.getElementsByClassName('triangle')[0]
        .addEventListener('click', function () {
            const popup = document.createElement("div")
            popup.classList.add("interface_create_post")
        popup.innerHTML = 
            `
            <input class="CP_Message" type="text">
            `
    document.body.appendChild(popup)
    });