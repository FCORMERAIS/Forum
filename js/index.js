console.log("JS Loaded")
const url = "127.0.0.1:5555"
var connected
var popupActivate = false

const header = document.getElementById("header")
if (document.cookie == "") {
    document.getElementById("connected").remove();
    const popup = document.createElement("div")
    popup.classList.add("connected")
    popup.innerHTML = `
            <button class ="button" id="connectBtn">Se connecter</button>
            <button class ="button" id="RegisterBtn">S'inscrire</button>
    `
    header.appendChild(popup)
    connected = false
}else {
    console.log("hello")
    document.getElementById("connected").remove();
    const popup = document.createElement("div")
    popup.classList.add("connected")
    popup.innerHTML = `
        <form method="POST">
        <button type="submit" class ="button" id="disconnectBtn">Se Déconnecter</button>
        </form>
    `
    header.appendChild(popup)
    console.log(document.cookie == null)
    connected = true
}
var register = document.getElementById("RegisterBtn")
var connect = document.getElementById("connectBtn")
if (connected == false){
    register.addEventListener("click",function() {
        console.log("hello")
        if (popupActivate) {
            document.getElementById("popup-content").remove();
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
        document.body.appendChild(popup)
        popupActivate = true
    })

    connect.addEventListener("click",function() {
        console.log("hello")
        if (popupActivate) {
            document.getElementById("popup-content").remove();
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
        document.body.appendChild(popup)
        popupActivate = true
    })
}else {
    var disconnect = document.getElementById("disconnectBtn")
    disconnect.addEventListener("click",function(){
        document.cookie = "UserSessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC;"
        document.cookie = "Value=; expires=Thu, 01 Jan 1970 00:00:00 UTC;"
    })
}