<script>
export default {
    props : ['userId','username']
}
</script>

<template>
    <div class="userItem">
        <!-- Sarebbero le card degli utenti -->
        <!-- Nome -->
        <h3 class="username">{{ username }}</h3>
        <!-- Pulsante a destra -->
        <button class="Button" @click="$emit('new-conversation', userId)"> ✍🏼 </button>
    </div>
</template>



<style scoped>
.userItem{
    width: 60%;
    display: flex;
    justify-content: space-between;
    background-color: white;
    border-radius: 10px;
    border: 2px solid rgb(4, 103, 133);
    padding: 10px;
    margin-bottom: 8px;
}
.username{
    font-family: "SF Pro";
    color: rgb(4, 66, 85);
}
.Button{
    text-align: center;
	font-family: "SF Pro";
	color: rgb(4, 66, 85); 
	background-color: #bad4e9;
	border: 2px solid rgb(4, 103, 133);
	border-radius: 1rem;
    width: 70px;
	box-shadow: 1px 1px 0 rgb(4, 103, 133);
}
.Button:hover{
	color: white;
	background-color: rgb(10, 146, 187);
}
.Button:active{
	box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
</style>