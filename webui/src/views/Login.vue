<script>
export default {
	data: function() {
		return {
			errormsg: null, //variabili globali 
            username: "",
			successmsg: null
		}
	},
	methods: {
		async doLogin() { //async se una funzione deve fare una richiesta al backend
			this.errormsg = null;
			//console.log(this.username)
			try {
				let response = await this.$axios.post("/login", {name: this.username}); //path dell'endpoint
				//console.log("Login response headers:", response.headers);
				//console.log("Login response data:", response.data);
                if (response.data) {
					//console.log("Response: ",response)
                    let auth = response.headers.authorization;   //let definisce una variabile 
				    localStorage.setItem("auth",auth)           //lo salvo nella cronologia (per prenderlo si usa getItem)
					//console.log("Auth: ",auth);
                    let identifier = response.data.identifier;  //setto il responso .data["identifier"]
                    localStorage.setItem("id",identifier)
					//console.log("Id: ",identifier);
					localStorage.setItem("username",this.username)

					this.successmsg = "Login effettuato con successo!";
					setTimeout(()=> {
						this.successmsg = null;
						//vado in un'altra pagina
						this.$router.push("/conversations") //path della prossima pagina
					},500);								//1 secondo di attesa prima del redirect
                }
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	},
	mounted() { //funzione che viene eseguita quando la schermata viene caricata su APP
	}
}
</script>

<template>
	<div class="content">
		<h1> WasaText </h1> 
		<hr>
		<ErrorMsg v-if="errormsg" :msg="errormsg"/>
		<SuccessMsg v-if="successmsg" :msg="successmsg" />
		<div class="subcontent"> 
			<form id="loginForm" @submit.prevent="doLogin">
			<p>Accedi o Crea un nuovo utente </p>
			<input class="input" type="text" id="username" v-model="username" placeholder="Inserisci il tuo username" required>
			<button class="inputButton" type="submit"> Accedi </button>
			</form>
		</div>
	</div>
</template>

<style scoped> 
/* scoped serve ad applicare lo stile solo a questa view*/
.content{
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}
.subcontent{
	margin-top: 4%;
	align-content: center;
	text-align: center;
}
h1 {
	margin-top: 5rem;
	text-align: center;
	color: rgb(4, 103, 133);
	text-shadow: 2px 2px 5px rgba(12, 92, 197, 0.3);
	font-size: 7rem;
}
hr {
	width: 50%;
	margin: auto;
}
p {
	font-size: 1.4rem;
	font-family: "SF Pro";
	color: rgb(4, 78, 100); 
}
.inputButton {
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 100%;
	margin-top: 1rem;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.inputButton:hover{
	background-color:rgb(10, 146, 187);
	color: white;
}
.inputButton:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}

.input {
  font-family: "SF Pro";
  width: 100%;
  padding: 0.875rem;
  border: 2px solid rgb(4, 103, 133);
  border-radius: 1rem;
  box-shadow: 1px 3px 0 rgb(4, 103, 133);
  outline: none;
  transition: ease 0.25s;
}
.input:focus {
  box-shadow: 2.5px 5px 0 rgb(4, 103, 133);
}
</style>
