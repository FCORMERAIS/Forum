console.log("JS Loaded")
const url = "127.0.0.1:5555"
var inputForm = document.getElementById("inputForm")
var register = document.getElementById("RegisterBtn")
var connect = document.getElementById("connectBtn")
const registerPush = false
const connectedPush = false
// inputForm.addEventListener("submit", (e)=>{
//     //prevent auto submission
//     e.preventDefault()
//     const formdata = new FormData(inputForm)
//     fetch(url,{
//         method:"POST",
//         body:formdata,
//     }).then(
//         response => response.text()
//     ).then(
//         (data) => {console.log(data);document.getElementById("serverMessageBox").innerHTML=data}
//     ).catch(
//         error => console.error(error)
//     )
// })

register.addEventListener("click",function() {
    if (registerPush == false) {
        const popup = document.createElement("div")
        popup.classList.add("popup")
        popup.innerHTML = `
        <div class="popup-content">
            <h1 style="margin-top:17%; margin-left:5%;">S'enregistrer : </h1>
            <form id="registerForm">
                </br>
                <input class ="formulaire" type="email" name="username" placeholder="Email">
                </br>
                </br>
                <input class ="formulaire" type="text" name="email" placeholder="Username">
                </br>
                </br>
                <input class ="formulaire" type="password" name="password" placeholder="Password">
                </br>
                </br>
                </br>
                <input style="color:#000000" class ="formulaire" type="submit" value="Register">
            </form>
        </div>
        `
        document.body.appendChild(popup)
    }
    if (registerPush ==false){registerPush=true}
})

connect.addEventListener("click",function() {
    if (connectedPush == false) {
        const popup = document.createElement("div")
        popup.classList.add("popup")
        popup.innerHTML = `
            <div class="popup-content">
            <h1 style="margin-top:10%; margin-left:5%;">Connexion : </h1>
            <form id="registerForm">
                </br>
                <input class ="formulaire" type="email" name="email" placeholder="Email">
                </br>
                </br>
                <input class ="formulaire" type="password" name="password" placeholder="Password">
                </br>
                </br>
                <input style="color:#000000" class ="formulaire" type="submit" value="se connecter">
            </form>
        </div>
        `
        document.body.appendChild(popup)
    }
    if (connectedPush ==false){connectedPush=true}
})