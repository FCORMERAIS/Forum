document.getElementsByClassName('triangle')[0].addEventListener('click', async function () { // on ajoute un eventListener "click" au triangle 
    if (document.cookie !="") { // on vérifie si il a des cookies sur le serveur pour savoir si la pop-up peut-etre activé
        const popup = document.createElement("div")
        popup.setAttribute("id","interface_CP");
        popup.classList.add("interface_create_post")
        popup.innerHTML = // on ajoute a la pop-up son corps
        `
        <p class="CP_Information">Balance ton post :</p>
            <input class="CP_Message" name="SendPost" id="Message" type="text" placeholder="C'est ici ton blabla ;)">
            <input class="CP_Send" type="submit" value="&#10145" id="nextPost">
        <div class="CP_fermer" id="CP_close"> </div>
        <p class="ajout_photo">&#128193;</p>
        `
        let options =""
        await fetch("http://127.0.0.1:5555/JsonCategories") // on vas chercher les informations des Catégories pour qu'il puisse afficher les choix de categories possible
            .then(response => response.json())
            .then( function (categories) {
                categories.forEach(categorie => {
                    options +=`
                    <option value="${categorie.Name}">${categorie.Name}</option> 
                    `
                })
            })
        const popupNext = document.createElement("div")
        popupNext.setAttribute("id","interface_CP_next");
        popupNext.classList.add('interface_create_post_Next')
        popupNext.innerHTML =
        `
        <div class="CP_fermer" id="CP_close_next"> </div>
        <form class="CP_form" method="POST" id="form_next" name="test">
            <input class="CP_Message" type="hidden" id="MessageValue" name="Message_Value">
            <p class="CP_Message" id="Value_Message"></p>
            <select name="Categorie" id="categorie" class="Categories_CP" required>
            ${options}
            </select>
            <input class="CP_Send" type="submit" value="&#10145" id="publishPost">
        </form>
        `
        document.body.appendChild(popupNext) // on ajoute la pop-up pour choisir la categories
        document.body.appendChild(popup) // on ajoute la pop-up a notre html
    }else { // si il n'y a pas de cookie dans le serveur alors on envois une alerte a l'utilisateur indiquant qu'il doit se connecter
        alert("Veuillez vous connectez pour pouvoir poster ")
    }
    document.getElementById('CP_close').addEventListener('click', () => { // on ajoute un eventListener a la pop-up pour la fermer
            document.getElementById("interface_CP").remove();
            document.getElementById("interface_CP_next").remove()
    })
    document.getElementById('CP_close_next').addEventListener('click', () => { // on ajoute un event listener a la deuxieme pop-up de categories pour aussi la fermer
        document.getElementById("interface_CP_next").remove();
        document.getElementById("interface_CP").remove();
})
    document.getElementById('nextPost').addEventListener('click', () => { // on ajoute un eventListener au post pour quand il l'ajoute
        let Message = document.getElementById("Message").value
        document.getElementById("MessageValue").value = Message
        document.getElementById("Value_Message").innerHTML = Message
        document.getElementsByClassName("interface_create_post_Next")[0].style.zIndex = "7";
    });
    document.getElementById('form_next').addEventListener('submit', () => { // on envoie les données pour qu'il envois une requete de post
        console.log(document.getElementById("MessageValue").value)
    });
});


fetch("http://127.0.0.1:5555/Post") // on avs chercher tout les posts a afficher
.then(response => response.json())
.then(function (donnee) {
    if (donnee != null) { // si donnee est null on afficher rien du tout 
        donnee.forEach(element => {//sinon on creer un Post pour chaque Post qu'il y a dans donnée
            const posts = document.getElementById("InterfacePost")
            posts.innerHTML += `
            <div class="Post" id="${element.IdPost}">
                <div class="Post_top">
                    <div class="Interface_User">
                        <img src="https://pic.onlinewebfonts.com/svg/img_329115.png" class="profile_Post" width="35" height="35">
                        <p class="UserName">${element.Username}</p>
                    </div>
                    <div id="Delete_${element.IdPost}" >

                    </div>
                </div>
                <p class="Message">${element.TextPost}</p>
                <div class=likeDislike style="display:flex">
                    <form class="LikePost" method="POST">
                        <input type="hidden" name="Like" value="${element.IdPost}">
                        <button class="ButtonLD" type="submit">&#x1F44D; ${element.LikePost}</button>
                    </form>
                    <form class="DislikePost" method="POST">
                        <input type="hidden" name="Dislike" value="${element.IdPost}">
                        <button class="ButtonLD" type="submit">&#128078;${element.DislikePost}</button>
                    </form>
                </div>
                <form class="commentary" method="POST">
                    <input type="hidden" name="idPost" value="${element.IdPost}">
                    <input type="text" name="textCommentary"> 
                </form>
                <div class="commentary" id=${element.IdPost}>
                </div>
            </div>
            `
        if (element.CommentaryPost != null) { // on vérifie si il y a des commentaire dans les posts
            element.CommentaryPost.forEach(commentary => { // si il y en a on les ajoute tous un par un 
                const comm = document.getElementById(commentary.IdPost)
                const comment = document.createElement("div")
                comment.innerHTML =`
                <div class="comm" style="display:flex">
                    <p>${commentary.Username} : ${commentary.Text}</p>
                    <div class=likeDislike style="display:flex">
                        <form class="LikeComment" method="POST">
                            <input type="hidden" name="LikeComm" value="${commentary.IdCommentary}">
                            <button class="ButtonLD" type="submit">&#x1F44D; ${commentary.Like}</button>
                        </form>
                        <form class="DislikeComment" method="POST">
                            <input type="hidden" name="DislikeComm" value="${commentary.IdCommentary}">
                            <button class="ButtonLD" type="submit">&#128078;${commentary.Dislike}</button>
                        </form>
                    </div>
                </div>
                `
                comm.appendChild(comment) // on ajoute le commentaire a la liste de commentaire
            });
        }
        document.body.appendChild(posts)
        document.getElementById(element.IdPost).style.backgroundColor = element.CategorieColor;
        if (element.SamePersonWhithSession) { // on verifie si l'utilisateur est l'auteur des posts et si cest le cas il peut supprimer les posts
            document.getElementById(`Delete_${element.IdPost}`).classList.add("Inteface_Delete")
            document.getElementById(`Delete_${element.IdPost}`).innerHTML = `
            <form method="POST">
                <button class="corbeille" name="delete" value="${element.IdPost}" type="submit">&#128465;</button>
            </form>
            `
        }
        });
    }
})

fetch("http://127.0.0.1:5555/JsonCategories") // Permet de recuperer toutes les categories disponibles pour les afficher en bas a gauche en forme de bouton avec leur couleur respective
.then(response => response.json())
.then( function (categories) {
    categories.forEach(categorie => {
        const categorieAdd = document.getElementById("AffichageCategorie")
        categorieAdd.innerHTML += `
        <form method="post" class="button_categories">
                <input name="categorieForm" class="categorie" type="submit" id="ButonCategorie${categorie.Name}" value=${categorie.Name} />
        </form>
        `
        document.body.appendChild(categorieAdd)
        document.getElementById(`ButonCategorie${categorie.Name}`).style.backgroundColor = categorie.Color;
    });
});
