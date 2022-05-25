document.getElementsByClassName('triangle')[0]
        .addEventListener('click', function () {
            const popup = document.createElement("div")
            popup.setAttribute("id","interface_CP");
            popup.classList.add("interface_create_post")
        popup.innerHTML = 
            `
            <p class="CP_Information">Balance ton post :</p>
            <form class="CP_form" method="POST">
                <input class="CP_Message" name="SendPost" type="text" placeholder="C'est ici ton blabla ;)">
                <input class="CP_Send" type="submit" value="&#10145">
            </form>
            <div class="CP_fermer" id="CP_close"> </div>
            <p class="ajout_photo">&#128193;</p>
            `
    document.body.appendChild(popup)
    document.getElementById('CP_close')
        .addEventListener('click', async () => {
            document.getElementById("interface_CP").remove();
    })
});

fetch("http://127.0.0.1:5555/donneesJson")
.then(response => response.json())
.then(function (donnee) {
    console.log(donnee)
    donnee.forEach(element => {
        const posts = document.getElementById("InterfacePost")
        posts.innerHTML += `
        <p class="Post">${element.TextPost}</p>
        `
        document.body.appendChild(posts)
    });
})