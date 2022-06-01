document.getElementsByClassName('triangle')[0].addEventListener('click', function () {
    if (document.cookie !="") {
        let Message = ""
        const popup = document.createElement("div")
        popup.setAttribute("id","interface_CP");
        popup.classList.add("interface_create_post")
        popup.innerHTML = 
        `
        <p class="CP_Information">Balance ton post :</p>
            <input class="CP_Message" name="SendPost" id="Message" type="text" placeholder="C'est ici ton blabla ;)">
            <input class="CP_Send" type="submit" value="&#10145" id="nextPost">
        <div class="CP_fermer" id="CP_close"> </div>
        <p class="ajout_photo">&#128193;</p>
        `
        
        const popupNext = document.createElement("div")
        popupNext.setAttribute("id","interface_CP_next");
        popupNext.classList.add('interface_create_post_Next')
        popupNext.innerHTML =
        `
        <div class="CP_fermer" id="CP_close_test"> </div>
        <form class="CP_form" method="POST" id="form_next" name="test">
            <input class="CP_Message" type="hidden" id="MessageValue" name="Message_Value">
            <p class="CP_Message_next" id="Value_Message">${Message}</p>
            <select name="Categorie" id="city-select">
                <option value=Technologie>Technologie</option>
                <option value=Debat>Debat</option>
            </select>
            <input class="CP_Send" type="submit" value="&#10145" id="publishPost">
        </form>
        `
        document.body.appendChild(popupNext)
        document.body.appendChild(popup)
    }else {
        alert("Veuillez vous connectez pour pouvoir poster ")
    }
    document.getElementById('CP_close').addEventListener('click', () => {
            document.getElementById("interface_CP").remove();
    })
    document.getElementById('nextPost').addEventListener('click', () => {
        Message = document.getElementById("Message").value
        document.getElementById("MessageValue").value = Message
        console.log("Value" + Message)
        document.getElementsByClassName("interface_create_post_Next")[0].style.zIndex = "7";
    });
    document.getElementById('form_next').addEventListener('submit', () => {
        console.log(document.getElementById("MessageValue").value)
    });
});


fetch("http://127.0.0.1:5555/Post")
.then(response => response.json())
.then(function (donnee) {
    console.log(donnee)
    donnee.forEach(element => {
        const posts = document.getElementById("InterfacePost")
        posts.innerHTML += `
        <div class="Post">
            <div class="Interface_User">
                <img src="https://pic.onlinewebfonts.com/svg/img_329115.png" class="profile_Post" width="35" height="35">
                <p class="UserName">${element.Username}</p>
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
        if (element.CommentaryPost != null) {
            element.CommentaryPost.forEach(commentary => {
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
                comm.appendChild(comment)
            });
        }
        document.body.appendChild(posts)
    });
})