var connected
var popupActivate = false

const header = document.getElementById("header") // permet d'afficher le header sur la page d'acceuil
if (document.cookie == "") { // on verifie si il y a des cookies sur le serveur pour savoir si l'utilisateur est connecter 
    document.getElementById("connected").remove(); // on supprime les boutons de connections
    const popup = document.createElement("div")
    popup.classList.add("connected")
    popup.innerHTML = `
            <button class ="button" id="connectBtn">Se connecter</button>
            <button class ="button" id="RegisterBtn">S'inscrire</button>
    `
    header.appendChild(popup)// on ajoute la pop-up pour que l'utilisateur se connecte 
    connected = false // et on definie sa connection en false
}else { // a l'inverse si l'utilisateur est connecter
    document.getElementById("connected").remove();
    const popup = document.createElement("div")
    popup.classList.add("connected")
    popup.innerHTML = `
        <form method="POST">
        <button type="submit" class ="button" id="disconnectBtn">Se Déconnecter</button>
        </form>
    `
    header.appendChild(popup) // on ajoute un bouton qui lui permet de se deconnecter 
    connected = true // et on definie connected sur true
}
if (connected == false){ // si l'utilisateur
    var register = document.getElementById("RegisterBtn") // permet de recuperer la reference sur le bouton de register
    var connect = document.getElementById("connectBtn") // permet de recuperer la reference sur le bouton de connection
    register.addEventListener("click",function() { // si le bouton est appuyer dans ce cas la on ajoute une pop-up pour afficher les données que l'on a besoin pour l'inscription
        if (popupActivate) { // si une pop-up est activé
            document.getElementById("popup-content").remove(); // on l'enlève
        }
        const popup = document.createElement("div")
        popup.classList.add("popup")
        popup.innerHTML = `
        <div class="popup-content" id="popup-content">
            <h1 style="margin-top:17%; margin-left:5%;">S'enregistrer : </h1>
            <form id="registerForm" method="POST">
                </br>
                <input class ="formulaire" type="email" name="email" placeholder="Email">
                </br></br>
                <input class ="formulaire" type="text" name="username" placeholder="Username">
                </br></br>
                <input class ="formulaire" type="password" name="password" placeholder="Password">
                </br></br></br>
                <input style="color:#000000" class ="formulaire" type="submit" value="S'enregistrer">
                </br></br>
            </form>
        </div>
        `
        document.body.appendChild(popup) // on ajoute les données a body qui l'affiche
        popupActivate = true
    })

    connect.addEventListener("click",function() { // l'utilisateur appuie sur la pop-up pour se connecter
        if (popupActivate) { // si une pop-up est acitvé
            document.getElementById("popup-content").remove(); // on l'enlève
        }
        const popup = document.createElement("div")
        popup.classList.add("popup")
        popup.innerHTML = `
            <div class="popup-content" id="popup-content">
            <h1 style="margin-top:10%; margin-left:5%;">Connexion : </h1>
            <form id="registerForm" method="POST">
                </br>
                <input class ="formulaire" type="email" name="email2" placeholder="Email">
                </br></br>
                <input class ="formulaire" type="password" name="password2" placeholder="Password">
                </br></br>
                <input style="color:#000000" class ="formulaire" type="submit" value="se connecter">
                </br></br>
            </form>
        </div>
        `
        document.body.appendChild(popup) //on ajoute a body les informations utile pour la connections (email password)
        popupActivate = true
    })
}else { // sinon l'utilisateur ne peut que se deconnecter
    var disconnect = document.getElementById("disconnectBtn")
    disconnect.addEventListener("click",function(){ // donc on ajoute un bouton qui lorsque l'on appuie dessus supprime tous les cookies
        document.cookie = "UserSessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC;"
        document.cookie = "Value=; expires=Thu, 01 Jan 1970 00:00:00 UTC;"
    })
}