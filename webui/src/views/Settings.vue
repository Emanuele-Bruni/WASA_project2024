<script>
export default {
    data: function() {
        return {
			errormsg: null,
		}
    },
    methods: {
        logout(){
			localStorage.removeItem("auth")
			localStorage.removeItem("id")
			this.$router.push("/login") //torno al login
		},
        usersSection(){
			this.$router.push("/users")
		},
        conversationsSection(){
            this.$router.push("/conversations")
        }
    },
    mounted() {
    }
}
</script>
<template>
    <div class="navbar">
		<div class="container">
			<button type="button"class="btn" @click="conversationsSection"> Conversazioni </button>
			<button type="button"class="btn" @click="usersSection"> Utenti </button>
			<button type="button" class="currentButton"> Impostazioni </button>
			<button type="button" class="logoutButton" @click="logout"> Logout </button>
		</div> 
	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>


<style scoped>
.navbar{
	background-color: rgb(178, 207, 233);
	padding: 0.75rem 1rem;
	position: fixed;
	left: 0;
	top: 0;
	width: 100%;
	z-index: 1000;
	box-shadow: 0 5px 5px rgba(53, 81, 94, 0.5	) /* x y sfocatura colore(trasparenza)*/
}
.container{
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 0.2rem;
}
.logoutButton{	
	color: white;
	background-color: rgb(235, 85, 85);
	text-align: center;
	font-family: "SF Pro";
	padding: 0.875rem;
	border: 2px solid rgb(243, 32, 32);
	border-radius: 1rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(243, 32, 32);
}
.logoutButton:hover{
	background-color: rgb(241, 41, 41);
}
.logoutButton:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.btn{
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
.btn:hover{
	color: white;
	background-color: rgb(10, 146, 187);
}
.btn:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.currentButton{
	text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: rgb(199, 177, 177);
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
	padding: 0.875rem;
	width: 118px;
	box-shadow: 1px 3px 0 rgb(4, 103, 133);
}
</style>