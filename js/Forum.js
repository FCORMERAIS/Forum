document.getElementsByClassName('triangle')[0]
        .addEventListener('click', function () {
            const popup = document.createElement("div")
            popup.setAttribute("id","interface_CP");
            popup.classList.add("interface_create_post")
        popup.innerHTML = 
            `
            <p class="CP_Information">Balance ton post :</p>
            <form class="CP_form" method="POST" id="form">
                <input class="CP_Message" name="SendPost" id="Message type="text" placeholder="C'est ici ton blabla ;)">
                <input class="CP_Send" type="submit" value="&#10145" id="nextPost">
            </form>
            <div class="CP_fermer" id="CP_close"> </div>
            <p class="ajout_photo">&#128193;</p>
            `
            document.body.appendChild(popup)
            const popupNext = document.createElement("div")
            popupNext.setAttribute("id","interface_CP_next");
            popupNext.classList.add('interface_create_post_Next')
            popupNext.innerHTML =
             `
            <div class="CP_fermer" id="CP_close_test"> </div>
            <form class="CP_form" method="POST" id="form_next">
            <input class="CP_Message" type="hidden" id="MessageValue">
            <input type="text">
            <input class="CP_Send" type="submit" value="&#10145" id="nextPost_test">
            </form>
            `
            document.body.appendChild(popupNext)
        document.getElementById('CP_close')
            .addEventListener('click', () => {
                document.getElementById("interface_CP").remove();
            })
        let Message = document.getElementById("Message").value
        const form  = document.getElementById('form');
        form.addEventListener('submit', (event) => {
                document.getElementsByClassName("interface_create_post_Next")[0].style.zIndex = "7";
                event.preventDefault();
        });
        const nextForm = document.getElementById('form_next');
        nextForm.addEventListener('submit', () => {
            document.getElementById("MessageValue").value = Message 
        });
    });




fetch("http://127.0.0.1:5555/donneesJson")
.then(response => response.json())
.then(function (donnee) {
    console.log(donnee)
    donnee.forEach(element => {
        const posts = document.getElementById("InterfacePost")
        posts.innerHTML += `
        <div class="Post">
            <div class="Interface_User">
                <img src="https://pic.onlinewebfonts.com/svg/img_329115.png" class="profile_Post" width="35" height="35">
                <p class="UserName">pseudo</p>
            </div>
            <p class="Message">${element.TextPost}</p>
            <form class="Barre_dinteraction" method="POST">
                <button class="ButtonLD" name="Like" type="submit" value='Like : ${element.ID_User_Post}'>&#x1F44D;</button>
                <button class="ButtonLD" name="Dislike" type="submit" value='Dislike : ${element.ID_User_Post}'>&#128078;</button>
            </form>
        </div>
        `
        document.body.appendChild(posts)
    });
})