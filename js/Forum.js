document.getElementsByClassName('triangle')[0]
        .addEventListener('click', function () {
            const popup = document.createElement("div")
            popup.setAttribute("id","interface_CP");
            popup.classList.add("interface_create_post")
        popup.innerHTML = 
            `
            <p class="CP_Information">Balance ton post :</p>
            <input class="CP_Message" type="text" placeholder="C'est ici ton blabla ;)">
            <div class="CP_fermer" id="CP_close")> </div>
            <p class="ajout_photo">&#128193;</p>
            `
    document.body.appendChild(popup)
    document.getElementById('CP_close')
        .addEventListener('click', async () => {
            document.getElementById("interface_CP").remove();
    })
});

