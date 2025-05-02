<script>
export default {
	data: function() {
		return {
			errormsg: null, //variabili globali 
            username: ""
		}
	},
	methods: {
		async doLogin() { //async se una funzione deve fare una richiesta al backend
			this.errormsg = null;
			//console.log(this.username)
			try {
				let response = await this.$axios.post("/login", {name: this.username}); //path dell'endpoint
                if (response) {
                    let auth = response.headers.authorization;   //let definisce una variabile 
				    localStorage.setItem("auth",auth)           //lo salvo nella cronologia (per prenderlo si usa getItem)

                    let identifier = response.data.identifier;  //setto il responso .data["identifier"]
                    localStorage.setItem("id",identifier)

                    //vado in un'altra pagina
                    this.$router.push("/home") //path della prossima pagina
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
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div class="subcontent"> 
			<form id="loginForm" @submit.prevent="doLogin">
			<p> Accedi o Crea un nuovo utente </p>
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
	background-color: rgb(233, 239, 243);
}
.subcontent{
	margin-top: 10rem;
	margin-bottom: 10rem;
	align-content: center;
	text-align: center;
}
h1 {
	margin-top: 2rem;
	text-align: center;
	color: rgb(4, 103, 133);
	text-shadow: 2px 2px 5px rgba(12, 92, 197, 0.3);
}
hr {
	width: 70%;
	margin: auto;
}

.inputButton {
	text-align: center;
	font-family: "SF Pro";
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 0.5rem;
	padding: 0.875rem;
	width: 100%;
	margin-top: 1rem;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}

.input {
  font-family: "SF Pro";
  width: 100%;
  padding: 0.875rem;
  /*font-size: 1rem;*/
  border: 2px solid rgb(4, 103, 133);
  border-radius: 0.5rem;
  box-shadow: 1px 3px 0 rgb(4, 103, 133);
  outline: none;
  transition: ease 0.25s;
}
.input:focus {
  box-shadow: 2.5px 5px 0 rgb(4, 103, 133);
}
</style>
